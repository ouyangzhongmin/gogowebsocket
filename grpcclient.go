//Copyright The ZHIYUNCo.All rights reserved.
//Created by admin at2024/7/5.

package gogowebsocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ouyangzhongmin/gogowebsocket/logger"
	"github.com/ouyangzhongmin/gogowebsocket/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// link::https://github.com/grpc/grpc-go/blob/master/examples/helloworld/greeter_client/main.go
func GrpcSendMsg(server serverInfo, msg *WSBody, broadcast int) error {
	// Set up a connection to the server.
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", server.ServerIP, server.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Errorln("连接失败", server.ServerIP, server.Port, err)
		return err
	}
	defer conn.Close()

	c := protobuf.NewAccServerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	//转发的消息的body都需要转换为string, protobuf的any类型不能完全适用任意类型的数据
	bb, err := convertSendBody(msg.BodyType, msg.Body)
	if err != nil {
		logger.Errorln("body marshal err::", err)
		return err
	}
	req := protobuf.SendMsgReq{
		Clientid:   msg.ClientID,
		ProtocolId: msg.ProtocolId,
		BodyType:   int32(msg.BodyType),
		Queue:      int32(msg.Queue),
		Broadcast:  int32(broadcast),
		Body:       string(bb),
	}
	rsp, err := c.SendMsg(ctx, &req)
	if err != nil {
		logger.Errorln("rpc.SendMsg:", server.ServerIP+":"+server.Port, err)
		return err
	}

	if rsp.GetErrcode() != 0 {
		logger.Errorln("rpc.SendMsg err::", rsp.String())
		err = errors.New(fmt.Sprintf("发送消息失败 code:%d", rsp.GetErrcode()))
		return err
	}
	logger.Debugln("rpc.SendMsg 成功")
	return nil
}
func GrpcForceDisconnect(server serverInfo, clientId string) error {
	// Set up a connection to the server.
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", server.ServerIP, server.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Errorln("连接失败", server.ServerIP, server.Port, err)
		return err
	}
	defer conn.Close()

	c := protobuf.NewAccServerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	req := protobuf.ForceDisconnectReq{
		Clientid: clientId,
	}
	rsp, err := c.ForceDisconnect(ctx, &req)
	if err != nil {
		logger.Errorln("rpc.GrpcForceDisconnect:", server.ServerIP+":"+server.Port, err)
		return err
	}
	if rsp.GetErrcode() != 0 {
		logger.Errorln("rpc.GrpcForceDisconnect err::", rsp.String())
		err = errors.New(fmt.Sprintf("rpc.GrpcForceDisconnect code:%d", rsp.GetErrcode()))
		return err
	}
	logger.Debugln("rpc.GrpcForceDisconnect 成功")
	return nil
}

func convertSendBody(bodyType int, body interface{}) (string, error) {
	if bodyType == BODY_TYPE_TEXT {
		return body.(string), nil
	}
	if bodyType == BODY_TYPE_BYTES {
		return string(body.([]byte)), nil
	}
	bb, err := json.Marshal(body)
	//bb, err := pbany(msg.Body)
	if err != nil {
		return "", err
	}
	return string(bb), nil
}
