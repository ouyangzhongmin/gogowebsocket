package gogowebsocket

import (
	"github.com/gorilla/websocket"
	"github.com/ouyangzhongmin/gogowebsocket/logger"
	"time"
)

const (
// Time allowed to write a message to the peer.
//writeWait = 30 * time.Second

// Time allowed to read the next pong message from the peer.
//pongWait = 60 * time.Second

// Send pings to peer with this period. Must be less than pongWait.
// pingPeriod = (pongWait * 9) / 10
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type UserInfo interface {
	GetClientID() string
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	ws *WS

	userInfo UserInfo

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan *WSBody

	//目前的心跳规则为：服务器需要在读取超时pongWait之前发送ping命令给前端，
	//前端主动发送msg_type=MSG_BODY_TYPE_HEARTBEAT的心跳消息，后端返回MSG_BODY_TYPE_HEARTBEAT的消息
	//前端心跳主要解决异常断开时前端可以主动发现连接断开
	clientPingTime time.Time
	connectedTs    int64 //连接的时间戳
}

func newClient(userinfo UserInfo, ws *WS, conn *websocket.Conn) *Client {
	return &Client{userInfo: userinfo, ws: ws, conn: conn, send: make(chan *WSBody, 256), connectedTs: time.Now().Unix()}
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.ws.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(c.ws.opts.maxMessageSize) // 消息体大小限制
	c.conn.SetReadDeadline(time.Now().Add(c.ws.opts.pongWaitTime))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(c.ws.opts.pongWaitTime)); return nil })
	c.clientPingTime = time.Now()
	for {
		msg := &WSBody{}
		//_, messageBytes, err := c.conn.ReadMessage()
		err := c.conn.ReadJSON(msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Errorf("error0: %v \n", err)
			}
			break
		}
		//messageBytes = bytes.TrimSpace(bytes.Replace(messageBytes, newline, space, -1))
		//var msg PBMessageBody
		//err = proto.Unmarshal(messageBytes, &msg)
		//if err != nil {
		//	log.Printf("error1: %v", err)
		//}
		msg.ClientID = c.GetClientId()
		msg.Client = c
		if msg.BodyType == BODY_TYPE_HEARTBEAT {
			//收到前端的心跳包
			c.clientPingTime = time.Now()
			err = c.Send(&WSBody{
				ProtocolId: ID_HEARTBEAT,
				BodyType:   BODY_TYPE_HEARTBEAT,
				Body:       PONG,
			})
			if err != nil {
				logger.Errorln("reponse pong err: %v", err)
			}
			continue
		}
		if msg.Queue != 0 {
			//消息是否需要通道队列执行
			c.ws.receiveQueue <- msg
		} else {
			//本携程内直接按顺序执行, 开发时需要考虑并发访问问题
			logger.Debugf("Receive msg: ClientID=%s, BodyType=%d, ProtocolId=%d \n", msg.ClientID, msg.BodyType, msg.ProtocolId)
			c.ws._doReceiveMessage(msg)
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	// 优化过多ticker的问题
	// ticker := time.NewTicker(pingPeriod)
	ticker := c.ws.timew.After(c.ws.opts.pingPeriod)
	defer func() {
		//ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(c.ws.opts.writeWaitTime))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.conn.WriteJSON(message)
			if err != nil {
				logger.Errorln("WriteJson err:", c.GetClientId(), err.Error())
			}
		//case <-ticker.C:
		case <-ticker:
			ticker = c.ws.timew.After(c.ws.opts.pingPeriod)
			//在读取时间内必须发送心跳包，否则会超时1006 error断开连接
			logger.Debugln("发送服务器心跳" + c.GetClientId())
			c.conn.SetWriteDeadline(time.Now().Add(c.ws.opts.writeWaitTime))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				//心跳发送失败了，则认为连接已断开
				return
			}
			// 更新缓存状态
			if c.ws.grpcServer != nil {
				c.ws.grpcServer.putClientToCache(c)
			}
		}
	}
}

func (c *Client) Send(msg *WSBody) error {
	select {
	case c.send <- msg:
	default:
		close(c.send)
		c.Close()
		c.ws.unregister <- c
	}
	return nil
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Client) GetUserInfo() UserInfo {
	return c.userInfo
}

func (c *Client) GetClientId() string {
	if c.userInfo != nil {
		return c.userInfo.GetClientID()
	}
	return "0"
}
