package service

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/ouyangzhongmin/gogowebsocket"
	"github.com/ouyangzhongmin/gogowebsocket/examples/model"
	"github.com/ouyangzhongmin/gogowebsocket/logger"
	"net/http"
)

type WSUserInfo struct {
	UserId     int
	Production string
	DeviceId   string
	UserName   string
	Avatar     string
	NickName   string
	Mail       string
	Mobile     string
	clientid   string
}

func (u *WSUserInfo) GetClientID() string {
	if u.clientid == "" {
		//为了支持多设备同时连接在线，这里把clientid里包含deviceId， 映射关系保存到redis内
		u.clientid = fmt.Sprintf("%s_%d_%s", u.Production, u.UserId, u.DeviceId)
	}
	return u.clientid
}

func getUserKey(productionExt string, userId int) string {
	return fmt.Sprintf("%s_%d", productionExt, userId)
}

type Service struct {
	ws    *gogowebsocket.WS
	cache *cache
}

func New(rdClient *redis.Client, rpcPort string) (s *Service) {
	s = &Service{}
	s.ws = gogowebsocket.New("pppp", rpcPort, rdClient)
	s.ws.RegisterHandler(s.messageHandler)
	s.ws.RegisterEventHandler(s.eventHandler)
	s.cache = newCache(rdClient)
	return
}

func (s *Service) ServeWS(u *WSUserInfo, w http.ResponseWriter, r *http.Request) error {
	return s.ws.ServeWs(u, w, r)
}

func (s *Service) Shutdown() {
	s.ws.Shutdown()
}

func (s *Service) messageHandler(ws *gogowebsocket.WS, msg *gogowebsocket.WSBody) {
	str, _ := msg.BodyToString()
	logger.Println("测试用ws广播：", msg.ProtocolId, str)
	err := ws.Broadcast(msg, nil)
	if err != nil {
		logger.Println("测试用ws转发err：", err)
		return
	}
}

func (s *Service) eventHandler(client *gogowebsocket.Client, event string) {
	if event == gogowebsocket.EVENT_REGISTER {
		//连接设备
		u := client.GetUserInfo().(*WSUserInfo)
		userKey := getUserKey(u.Production, u.UserId)
		clientId := u.GetClientID()
		err := s.cache.putUserClientID(userKey, clientId)
		if err != nil {
			logger.Errorln("cache.putUserClientID err::", err)
		}
	} else if event == gogowebsocket.EVENT_UNREGISTER {
		//断开了连接删除掉缓存的数据
		u := client.GetUserInfo().(*WSUserInfo)
		userKey := getUserKey(u.Production, u.UserId)
		clientId := u.GetClientID()
		err := s.cache.removeUserClientID(userKey, clientId)
		if err != nil {
			logger.Errorln("cache.removeUserClientID err::", err)
		}
	}
}

func (s *Service) getClientID(productionExt string, userId int) []string {
	userKey := getUserKey(productionExt, userId)
	//通过映射查找真实的clientId
	clientIds, err := s.cache.getUserClientIDs(userKey)
	if err != nil {
		logger.Errorln("cache.getUserClientIDs err::", clientIds)
		return []string{}
	}
	return clientIds
}

func (s *Service) sendUserData(data *model.WsUserData) {
	productionExt := data.ProductionExt
	logger.Println("通知ws客户端:", data.UserID, productionExt, data.WsProtocalID)
	toClientIds := s.getClientID(productionExt, data.UserID)
	if len(toClientIds) == 0 {
		return
	}
	var msg gogowebsocket.WSBody
	msg.ProtocolId = data.WsProtocalID
	msg.BodyType = gogowebsocket.BODY_TYPE_JSON
	msg.Body = *data
	err := s.ws.SendMore(toClientIds, &msg)
	if err != nil {
		logger.Errorln("通知ws客户端err:", err)
		return
	}
}
