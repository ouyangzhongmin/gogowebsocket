package gogowebsocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/ouyangzhongmin/gogowebsocket/logger"
	"github.com/ouyangzhongmin/gogowebsocket/protobuf"
	"google.golang.org/grpc"
)

type serverInfo struct {
	ServerIP string `json:"server_ip"`
	Port     string `json:"port"`
}

func (s serverInfo) toString() string {
	return fmt.Sprintf("%s-%s", s.ServerIP, s.Port)
}

func (s serverInfo) isLocal(s2 serverInfo) bool {
	return s.ServerIP == s2.ServerIP && s.Port == s2.Port
}

type server struct {
	protobuf.UnimplementedWSServerServer
	ws *WS
}

// 健康检查
func (s *server) CheckHealth(c context.Context, req *protobuf.CheckHealthReq) (rsp *protobuf.CheckHealthRsp, err error) {
	rsp = &protobuf.CheckHealthRsp{Ts: time.Now().Unix()}
	return
}

// 给用户发消息
func (s *server) SendMsg(c context.Context, req *protobuf.SendMsgReq) (rsp *protobuf.OkRsp, err error) {
	rsp = &protobuf.OkRsp{}
	logger.Debugln("收到rpc.SendMsg::", req.Clientid)
	body2, err := s.convertReceiveBody(int(req.BodyType), req.Body)
	if err != nil {
		return rsp, errors.New("convertBody err:" + err.Error())
	}
	msg := &WSBody{
		ClientID:   req.Clientid,
		Client:     nil,
		ProtocolId: req.ProtocolId,
		BodyType:   int(req.BodyType),
		Queue:      int(req.Queue),
		Body:       body2,
	}
	if req.Broadcast == 1 {
		err = s.ws.BroadcastLocal(msg, nil)
		if err != nil {
			rsp.Errcode = 1
			rsp.ErrMsg = err.Error()
		}
	} else {
		err = s.ws.SendLocal(req.Clientid, msg)
		if err != nil {
			rsp.Errcode = 1
			rsp.ErrMsg = err.Error()
		}
	}

	return
}

// 强制断开一个连接
func (s *server) ForceDisconnect(c context.Context, req *protobuf.ForceDisconnectReq) (rsp *protobuf.OkRsp, err error) {
	rsp = &protobuf.OkRsp{}
	if req.Clientid != "" {
		err = s.ws.ForceDisconnect(req.Clientid)
		if err != nil {
			rsp.Errcode = 1
			rsp.ErrMsg = err.Error()
		}
	}
	return
}

func (s *server) convertReceiveBody(bodyType int, body []byte) (interface{}, error) {
	if bodyType == BODY_TYPE_TEXT {
		return string(body), nil
	}
	if bodyType == BODY_TYPE_BYTES {
		return body, nil
	}
	var data interface{}
	err := json.Unmarshal(body, &data)
	return data, err
}

type grpcServer struct {
	serverInfo *serverInfo
	ws         *WS
	cache      *cache
	startTime  time.Time
	scheduler  *gocron.Scheduler
}

func newGrpcServer(ws *WS, cache *cache, port string) *grpcServer {
	serverIp := GetServerIp()
	sInfo := &serverInfo{
		ServerIP: serverIp,
		Port:     port,
	}
	return &grpcServer{
		serverInfo: sInfo,
		ws:         ws,
		cache:      cache,
		startTime:  time.Now(),
	}
}

// 开启grpc server
func (s *grpcServer) Start() {
	logger.Printf("grpc server starting: %s:%s\n", s.serverInfo.ServerIP, s.serverInfo.Port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.serverInfo.Port))
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}

	if s.serverInfo.Port == "0" {
		// 如果是随机port需要从监听里重新获取端口
		tmp := lis.Addr().String()
		tmp = strings.ReplaceAll(tmp, "[::]", "")
		tmpArr := strings.Split(tmp, ":")
		if len(tmpArr) != 2 {
			logger.Fatalf("grpc server addr 格式错误:%s", tmp)
		}
		s.serverInfo.Port = tmpArr[1]
	}
	sv := grpc.NewServer()
	protobuf.RegisterWSServerServer(sv, &server{ws: s.ws})

	logger.Printf("grpc server started: %s:%s\n", s.serverInfo.ServerIP, s.serverInfo.Port)
	//记录到缓存
	s.timerPushServerToCache()

	if err := sv.Serve(lis); err != nil {
		if s.cache != nil {
			s.cache.removeServerInfo(s.serverInfo)
		}
		logger.Fatalf("grpc start fatal: %v", err)
	}
}

func (s *grpcServer) Stop() {
	if s.cache != nil {
		s.cache.removeServerInfo(s.ws.grpcServer.serverInfo)
		if s.scheduler != nil {
			s.scheduler.Stop()
			s.scheduler.Clear()
			s.scheduler = nil
		}
	}
}

func (s *grpcServer) send(toClientId string, msg *WSBody) error {
	info, err := s.cache.getClientInfo(toClientId)
	if err != nil {
		return err
	}
	if info.isLocal(*s.serverInfo) {
		return errors.New("本服务器的消息无需转发:" + toClientId)
	}
	if info.isOnline() {
		err = GrpcSendMsg(info.serverInfo, msg, 0)
		if err != nil {
			return err
		}
	} else {
		return errors.New("缓存的设备连接已断开:" + info.toString())
	}
	return nil
}

// 全局广播
func (s *grpcServer) broadcast(msg *WSBody, ignoreIds map[string]bool) error {
	//广播给其他server
	servers, err := s.cache.getServerList()
	if err != nil {
		return err
	}
	for _, server := range servers {
		if server.isLocal(*s.serverInfo) {
			continue
		}
		if server.isOnline() {
			err = GrpcSendMsg(server.serverInfo, msg, 1)
			if err != nil {
				logger.Errorln("GrpcSendMsg error::", err)
				continue
			}
		} else {
			logger.Warnln("grpc服务器已离线:" + server.toString())
		}
	}
	return nil
}

func (s *grpcServer) disconnectRemoteIfNeed(clientId string) {
	cacheInfo, _ := s.cache.getClientInfo(clientId)
	if cacheInfo != nil && cacheInfo.isOnline() &&
		(cacheInfo.ServerIP != s.serverInfo.ServerIP || cacheInfo.Port != s.serverInfo.Port) {
		logger.Infoln("其他设备上有重复连接:", clientId)
		go GrpcForceDisconnect(cacheInfo.serverInfo, clientId)
	}
}

func (s *grpcServer) timerPushServerToCache() {
	if s.cache != nil {
		s.pushServerToCache()
		//开启定时器 每隔30S更新在线
		s.scheduler = gocron.NewScheduler(time.Local)
		s.scheduler.Every(30).Seconds().Do(s.pushServerToCache)
		s.scheduler.StartAsync()

		servers, err := s.cache.getServerList()
		if err != nil {
			logger.Errorln("cache.getServerList error:", err)
		}
		logger.Println("grpc server list :::", servers)
	}
}

func (s *grpcServer) pushServerToCache() {
	if s.cache != nil {
		err := s.cache.putServerInfo(s.serverInfo, s.startTime)
		if err != nil {
			logger.Errorln("pushServerToCache error:", err)
		}
	}
}

func (s *grpcServer) putClientToCache(c *Client) {
	if s.cache != nil {
		nowts := time.Now().Unix()
		info := cacheClientInfo{
			serverInfo: serverInfo{
				ServerIP: s.serverInfo.ServerIP,
				Port:     s.serverInfo.Port,
			},
			Ts:  nowts,         //最新在线的时间戳
			CTs: c.connectedTs, //连接的时间戳
		}
		err := s.cache.putClientInfo(c.GetClientId(), info)
		if err != nil {
			logger.Errorln("putClientToCache error::", err)
		}
	}
}

func (s *grpcServer) removeClientCache(clientId string) {
	if s.cache != nil {
		err := s.cache.removeClientInfo(clientId)
		if err != nil {
			logger.Errorln("removeClientCache error::", err)
		}
	}
}
