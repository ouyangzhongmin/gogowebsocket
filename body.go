package gogowebsocket

import (
	"encoding/json"
	"errors"
)

const (
	BODY_TYPE_TEXT      = 1
	BODY_TYPE_JSON      = 0
	BODY_TYPE_BYTES     = 2 // 这个类型在grpc里转发不太适用，暂不支持
	BODY_TYPE_HEARTBEAT = 10000
	PONG                = "pong"

	ID_ERROR     = -1
	ID_HEARTBEAT = -2
)

type WSBody struct {
	ClientID   string      `json:"-"`
	Client     *Client     `json:"-"`
	ProtocolId int64       `json:"protocol_id"`
	BodyType   int         `json:"body_type"`
	Queue      int         `json:"queue"` //是否全局队列处理消息,提供前端api接口时需要标记清楚
	Body       interface{} `json:"body"`
}

func (m WSBody) BodyToStruct(to interface{}) error {
	if m.BodyType != BODY_TYPE_JSON {
		return errors.New("not suport not json msg")
	}
	tmp := m.Body.(map[string]interface{})
	bytes, err := json.Marshal(&tmp)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, to)
	return err
}

func (m WSBody) BodyToString() (string, error) {
	if m.BodyType != BODY_TYPE_JSON && m.BodyType != BODY_TYPE_TEXT {
		return "", errors.New("not suport not json msg")
	}
	if m.BodyType == BODY_TYPE_TEXT {
		return m.Body.(string), nil
	}

	bb, err := json.Marshal(m.Body)
	if err != nil {
		return "", err
	}
	return string(bb), nil
}

type Error struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
