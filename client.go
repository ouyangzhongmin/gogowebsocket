package gogowebsocket

import (
	"github.com/gorilla/websocket"
	"github.com/ouyangzhongmin/gogowebsocket/logger"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
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
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
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
			//消息是否需要全局队列执行
			c.ws.receiveQueue <- msg
		} else {
			//本携程内直接按顺序执行
			//如果业务不需要单用户顺序，开发者需要回调函数内部再开携程并发处理以提升执行效率
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
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		//由于websocket不支持并发写入, 通过通道来处理
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.conn.WriteJSON(message)
			if err != nil {
				logger.Errorln("WriteJson err:", c.GetClientId(), err.Error())
			}
		case <-ticker.C:
			//在读取时间内必须发送心跳包，否则会超时1006 error断开连接
			logger.Debugln("发送服务器心跳" + c.GetClientId())
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				//心跳发送失败了，则认为连接已断开
				return
			}
			c.putToCache() //更新缓存的在线状态
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

func (c *Client) putToCache() {
	nowts := time.Now().Unix()
	info := cacheClientInfo{
		serverInfo: serverInfo{
			ServerIP: c.ws.serverIp,
			Port:     c.ws.rpcPort,
		},
		Ts:  nowts,         //最新在线的时间戳
		CTs: c.connectedTs, //连接的时间戳
	}
	err := c.ws.cache.putClientInfo(c.ws.appId, c.GetClientId(), info)
	if err != nil {
		logger.Errorln("cache.putClientInfo err::", err)
	}
}

func (c *Client) removeCache() {
	err := c.ws.cache.removeClientInfo(c.ws.appId, c.GetClientId())
	if err != nil {
		logger.Errorln("cache.removeClientInfo err::", err)
	}
}
