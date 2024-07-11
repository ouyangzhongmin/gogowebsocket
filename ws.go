package gogowebsocket

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/ouyangzhongmin/gogowebsocket/logger"
	"net/http"
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
	appId string
	// 连接的clients.
	clients map[string]*Client
	// 用于队列发送消息.
	receiveQueue chan *WSBody

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	shutdown chan struct{}

	grpcServer *grpcServer
	serverIp   string
	rpcPort    string
	cache      *cache

	handlers     []MessageHandler
	eventHandler EventHandler
}

func New(appId, rpcPort string, r *redis.Client) *WS {
	hub := &WS{
		appId:        appId,
		receiveQueue: make(chan *WSBody, 20),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		shutdown:     make(chan struct{}),
		clients:      make(map[string]*Client),
		handlers:     make([]MessageHandler, 0),
	}
	hub.cache = newCache(r)
	hub.serverIp = GetServerIp()
	hub.rpcPort = rpcPort
	hub.grpcServer = newGrpcServer(hub, &serverInfo{
		ServerIP: hub.serverIp,
		Port:     hub.rpcPort,
	})
	go hub.grpcServer.Start()
	go hub.run()
	return hub
}

func (ws *WS) run() {
	for {
		select {
		case <-ws.shutdown:
			//退出程序, 断开所有连接
			logger.Println("退出断开所有连接")
			ws.cache.removeServerInfo(ws.appId, ws.grpcServer.serverInfo)
			for id := range ws.clients {
				c := ws.clients[id]
				c.Close()
				c.removeCache()
				//close(client.send)
			}
			ws.clients = make(map[string]*Client)
			if ws.grpcServer != nil {
				ws.grpcServer.Stop()
			}
			return
		case client := <-ws.register:
			//注册新的客户端连接
			logger.Infoln("register client: ", client.GetClientId())
			//这里多设备连接需要负载均衡保持同一个ClientId路由到同一台服务器
			c := ws.clients[client.GetClientId()]
			if c != nil {
				//有旧的连接还在，关闭掉
				logger.Infoln("断开重复的连接:", client.GetClientId())
				c.Close()
			}
			//检测缓存中记录的连接信息
			cacheInfo, _ := ws.cache.getClientInfo(ws.appId, client.GetClientId())
			if cacheInfo != nil && cacheInfo.isOnline() && (cacheInfo.ServerIP != ws.serverIp || cacheInfo.Port != ws.rpcPort) {
				logger.Infoln("其他设备上有重复连接:", client.GetClientId())
				go GrpcForceDisconnect(cacheInfo.serverInfo, client.GetClientId())
			}
			// 保存连接
			ws.clients[client.GetClientId()] = client
			//保存连接信息到缓存
			client.putToCache()

			//触发外部回调函数
			go ws.postEventHandler(client, EVENT_REGISTER)
		case client := <-ws.unregister:
			//删除客户端连接
			logger.Infoln("unregister client: ", client.GetClientId())
			close(client.send)
			delete(ws.clients, client.GetClientId())

			//删除缓存的记录
			client.removeCache()

			//触发外部回调函数
			go ws.postEventHandler(client, EVENT_UNREGISTER)
		case msg := <-ws.receiveQueue:
			//全局的队列处理消息
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
	if client, ok := ws.clients[toClientId]; ok {
		err := client.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

// 全局发送
func (ws *WS) Send(toClientId string, msg *WSBody) error {
	if _, ok := ws.clients[toClientId]; ok {
		return ws.SendLocal(toClientId, msg)
	} else {
		//不在本服务器内转发到其他服务器
		info, err := ws.cache.getClientInfo(ws.appId, toClientId)
		if err != nil {
			return err
		}
		if info.isLocal(*ws.grpcServer.serverInfo) {
			return errors.New("本服务器的消息无需转发:" + toClientId)
		}
		if info.isOnline() {
			err = GrpcSendMsg(info.serverInfo, msg, 0)
			if err != nil {
				return err
			}
		} else {
			return errors.New("设备连接已断开:" + info.toString())
		}
	}
	return nil
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
	for id := range ws.clients {
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
	//广播给其他server
	servers, err := ws.cache.getServerInfos(ws.appId)
	if err == nil {
		for _, server := range servers {
			if server.isLocal(*ws.grpcServer.serverInfo) {
				continue
			}
			if server.isOnline() {
				err = GrpcSendMsg(server.serverInfo, msg, 1)
				if err != nil {
					logger.Errorln("GrpcSendMsg err::", err)
					continue
				}
			} else {
				logger.Warnln("rpc服务器已离线:" + server.toString())
			}
		}
	} else {
		logger.Errorln("getServerInfos err::", err)
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
	logger.Errorln("Reponse client Err:", toClientId, errcode, errmsg)
	return ws.Send(toClientId, errMsg)
}

// 强制断开连接
func (ws *WS) ForceDisconnect(clientId string) error {
	if client, ok := ws.clients[clientId]; ok {
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
