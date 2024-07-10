//Copyright The ZHIYUNCo.All rights reserved.
//Created by admin at2024/7/5.

package gogowebsocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
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

// 给用户发消息
func (s *server) SendMsg(c context.Context, req *protobuf.SendMsgReq) (rsp *protobuf.SendMsgRsp, err error) {

	fmt.Println("grpc_request 给本机用户发消息", req.String())

	rsp = &protobuf.SendMsgRsp{}
	//body, err := anytop(req.Body)
	//if err != nil {
	//	return rsp, errors.New("anytop err:" + err.Error())
	//}
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
	logger.Debugln("grpc_response 给本机用户发消息", rsp.String())
	return
}

// 给用户发消息
func (s *server) ForceDisconnect(c context.Context, req *protobuf.ForceDisconnectReq) (rsp *protobuf.ForceDisconnectRsp, err error) {
	rsp = &protobuf.ForceDisconnectRsp{}
	if req.Clientid != "" {
		err = s.ws.ForceDisconnect(req.Clientid)
		if err != nil {
			rsp.Errcode = 1
			rsp.ErrMsg = err.Error()
		}
	}
	return
}

func (s *server) convertReceiveBody(bodyType int, body string) (interface{}, error) {
	if bodyType == BODY_TYPE_TEXT {
		return body, nil
	}
	if bodyType == BODY_TYPE_BYTES {
		return []byte(body), nil
	}
	var data interface{}
	err := json.Unmarshal([]byte(body), &data)
	//if bodyType == BODY_TYPE_TEXT {
	//	b, _ := json.Marshal(v)
	//	var txt string
	//}
	return data, err
}

type grpcServer struct {
	serverInfo *serverInfo
	Ws         *WS
	startTime  time.Time
	scheduler  *gocron.Scheduler
}

func newGrpcServer(ws *WS, serverInfo *serverInfo) *grpcServer {
	return &grpcServer{
		serverInfo: serverInfo,
		Ws:         ws,
		startTime:  time.Now(),
	}
}

// rpc server
func (s *grpcServer) Start() {
	fmt.Println("rpc server 启动", s.serverInfo.ServerIP, s.serverInfo.Port)
	lis, err := net.Listen("tcp", ":"+s.serverInfo.Port)
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
	sv := grpc.NewServer()
	protobuf.RegisterWSServerServer(sv, &server{ws: s.Ws})
	//记录到缓存
	s.pushToCache()
	if err := sv.Serve(lis); err != nil {
		s.Ws.cache.removeServerInfo(s.Ws.appId, s.serverInfo)
		logger.Fatalf("failed to serve: %v", err)
	}
}

func (s *grpcServer) Stop() {
	if s.scheduler != nil {
		s.scheduler.Stop()
		s.scheduler.Clear()
		s.scheduler = nil
	}
}

func (s *grpcServer) pushToCache() {
	s.doPushToCache()
	//开启定时器
	s.scheduler = gocron.NewScheduler(time.Local)
	s.scheduler.Every(30).Seconds().Do(s.doPushToCache)
	s.scheduler.StartAsync()

	servers, err := s.Ws.cache.getServerInfos(s.Ws.appId)
	if err != nil {
		logger.Errorln("getServerInfos err:", err)
	}
	logger.Println("grpc servers :::", servers)
}

func (s *grpcServer) doPushToCache() {
	err := s.Ws.cache.putServerInfo(s.Ws.appId, s.serverInfo, s.startTime)
	if err != nil {
		logger.Errorln("putGrpcServerInfo err:", err)
	}
}
