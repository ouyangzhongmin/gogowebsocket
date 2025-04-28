package gogowebsocket

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/ouyangzhongmin/gogowebsocket/logger"
	"github.com/ouyangzhongmin/gogowebsocket/timingwheel"
	"net/http"
	"time"
)

const (
	EVENT_REGISTER   = "register"
	EVENT_UNREGISTER = "unregister"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type MessageHandler func(*WS, *WSBody)

type EventHandler func(*Client, string)

type WS struct {
	clientsMgr *clientsMgr
	//用于优化ticker过多
	timew *timingwheel.TimingWheel
	// 用于队列发送消息.
	receiveQueue chan *WSBody

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	shutdown chan struct{}

	grpcServer *grpcServer

	handlers     []MessageHandler
	eventHandler EventHandler
}

func New() *WS {
	hub := &WS{
		receiveQueue: make(chan *WSBody, 20),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		shutdown:     make(chan struct{}),
		clientsMgr:   newClientsMgr(),
		timew:        timingwheel.NewTimingWheel(time.Second, 100),
		handlers:     make([]MessageHandler, 0),
	}
	go hub.run()
	return hub
}

// 需要主动开启grpc服务, appId用于区分记录在缓存中的命名
func (ws *WS) StartGrpcServer(appId, grpcPort string, r *redis.Client) {
	if appId == "" {
		panic("appId is empty")
	}
	if grpcPort == "" {
		panic("grpcPort is empty")
	}
	c := newCache(appId, r)
	ws.grpcServer = newGrpcServer(ws, c, grpcPort)
	go ws.grpcServer.Start()

}

func (ws *WS) run() {
	for {
		select {
		case <-ws.shutdown:
			//退出程序, 断开所有连接
			logger.Println("退出断开所有连接")
			ws.timew.Stop()
			clients := ws.clientsMgr.getClients()
			for id := range clients {
				c := clients[id]
				c.Close()
				if ws.grpcServer != nil {
					ws.grpcServer.removeClientCache(id)
				}
				//close(client.send)
			}
			ws.clientsMgr.delAll()
			if ws.grpcServer != nil {
				ws.grpcServer.Stop()
			}
			return
		case client := <-ws.register:
			//注册新的客户端连接
			logger.Infoln("register client: ", client.GetClientId())
			//这里多设备连接需要负载均衡保持同一个ClientId路由到同一台服务器
			c := ws.clientsMgr.getClient(client.GetClientId())
			if c != nil {
				//有旧的连接还在，关闭掉
				logger.Infoln("断开重复的连接:", client.GetClientId())
				c.Close()
			}
			//检测缓存中记录的连接信息
			if ws.grpcServer != nil {
				ws.grpcServer.disconnectRemoteIfNeed(client.GetClientId())
			}
			// 保存连接
			ws.clientsMgr.addClient(client)
			//保存连接信息到缓存
			if ws.grpcServer != nil {
				ws.grpcServer.putClientToCache(client)
			}

			//触发外部回调函数
			go ws.postEventHandler(client, EVENT_REGISTER)
		case client := <-ws.unregister:
			//删除客户端连接
			logger.Infoln("unregister client: ", client.GetClientId())
			close(client.send)
			ws.clientsMgr.delClient(client.GetClientId())

			//删除缓存的记录
			if ws.grpcServer != nil {
				ws.grpcServer.removeClientCache(client.GetClientId())
			}

			//触发外部回调函数
			go ws.postEventHandler(client, EVENT_UNREGISTER)
		case msg := <-ws.receiveQueue:
			//通道队列处理消息
			logger.Debugf("Receive queue msg: ClientID=%s, BodyType=%d, ProtocolId=%d \n", msg.ClientID, msg.BodyType, msg.ProtocolId)
			ws._doReceiveMessage(msg)
		}
	}
}

// 处理消息转发
func (ws *WS) _doReceiveMessage(msg *WSBody) {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorln("_doReceiveMessage panic: ", err)
		}
	}()
	for i := range ws.handlers {
		ws.handlers[i](ws, msg)
	}
}

// serveWs handles websocket requests from the peer.
func (ws *WS) ServeWs(userinfo UserInfo, w http.ResponseWriter, r *http.Request) error {
	if userinfo == nil {
		return errors.New("userinfo must not null!")
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Errorln("Upgrade ws err:", err)
		return err
	}
	client := newClient(userinfo, ws, conn)
	client.ws.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
	return nil
}

// 处理消息回调
func (ws *WS) RegisterHandler(handler MessageHandler) {
	ws.handlers = append(ws.handlers, handler)
}

// 连接成功断开连接将会触发回调
func (ws *WS) RegisterEventHandler(handler EventHandler) {
	ws.eventHandler = handler
}

func (ws *WS) postEventHandler(client *Client, event string) {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorln("postEventHandler.panic err::", err)
		}
	}()
	if ws.eventHandler != nil {
		ws.eventHandler(client, event)
	}
}

// 本服务器内发送
func (ws *WS) SendLocal(toClientId string, msg *WSBody) error {
	if client := ws.clientsMgr.getClient(toClientId); client != nil {
		err := client.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

// 全局发送
func (ws *WS) Send(toClientId string, msg *WSBody) error {
	if client := ws.clientsMgr.getClient(toClientId); client != nil {
		return ws.SendLocal(toClientId, msg)
	} else if ws.grpcServer != nil {
		//不在本服务器内转发到其他服务器
		return ws.grpcServer.send(toClientId, msg)
	}
	return errors.New("未找到设备连接:" + toClientId)
}

// 多个连接发送数据
func (ws *WS) SendMore(toClients []string, msg *WSBody) error {
	if len(toClients) <= 0 {
		return errors.New("send can't sendTo <= 0！")
	}
	sucess := 0
	for _, clientId := range toClients {
		if err := ws.Send(clientId, msg); err != nil {
			logger.Warnln(fmt.Sprintf("Send err:%s", err.Error()))
			continue
		}
		sucess++
	}
	if sucess == 0 {
		return errors.New("全部发送失败")
	}
	return nil
}

// 本服务器广播
func (ws *WS) BroadcastLocal(msg *WSBody, ignoreIds map[string]bool) error {
	success := 0
	clients := ws.clientsMgr.getClients()
	for id := range clients {
		if ignoreIds != nil {
			if _, ok := ignoreIds[id]; ok {
				continue
			}
		}
		err := ws.Send(id, msg)
		if nil != err {
			logger.Errorln("Broadcast %s error: %v\n", id, err)
			continue
		}
		success++
	}
	if success == 0 {
		return errors.New("all send err!")
	}
	return nil
}

// 全局广播
func (ws *WS) Broadcast(msg *WSBody, ignoreIds map[string]bool) error {
	err := ws.BroadcastLocal(msg, ignoreIds)
	if err != nil {
		logger.Errorln("BroadcastLocal err::", err)
	}
	if ws.grpcServer != nil {
		err = ws.grpcServer.broadcast(msg, ignoreIds)
		if err != nil {
			return err
		}
	}
	return nil
}

// 给连接返回错误码
func (ws *WS) Error(toClientId string, errcode int, errmsg string) error {
	errMsg := &WSBody{
		ProtocolId: ID_ERROR,
		BodyType:   BODY_TYPE_JSON,
		Body: &Error{
			ErrCode: errcode,
			ErrMsg:  errmsg,
		},
	}
	logger.Errorln("send client Error:", toClientId, errcode, errmsg)
	return ws.Send(toClientId, errMsg)
}

// 强制断开连接
func (ws *WS) ForceDisconnect(clientId string) error {
	if client := ws.clientsMgr.getClient(clientId); client != nil {
		client.Close()
		logger.Println("其他服务器通知断开连接:", clientId)
	}
	return nil
}

// 需要优雅重启时，可以先调用本函数断开所有连接，
// 以防止优雅重启过程中再收到新数据
func (ws *WS) Shutdown() {
	logger.Println("Shutdown ws")
	ws.shutdown <- struct{}{}
}
