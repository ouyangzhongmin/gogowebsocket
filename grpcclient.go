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

func GrpcSendMsg(server serverInfo, msg *WSBody, broadcast int) error {

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", server.ServerIP, server.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Errorln("grpc连接失败", server.ServerIP, server.Port, err)
		return err
	}
	defer conn.Close()

	c := protobuf.NewWSServerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	//转发的消息的body都需要转换为string, protobuf的any类型不能完全适用任意类型的数据
	body, err := convertSendBody(msg.BodyType, msg.Body)
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
		Body:       body,
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

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", server.ServerIP, server.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Errorln("grpc连接失败", server.ServerIP, server.Port, err)
		return err
	}
	defer conn.Close()

	c := protobuf.NewWSServerClient(conn)
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

func convertSendBody(bodyType int, body interface{}) ([]byte, error) {
	if bodyType == BODY_TYPE_TEXT {
		return []byte(body.(string)), nil
	}
	if bodyType == BODY_TYPE_BYTES {
		return body.([]byte), nil
	}
	tmp, err := json.Marshal(body)
	//tmp, err := pbany(msg.Body)
	if err != nil {
		return nil, err
	}
	return tmp, nil
}
