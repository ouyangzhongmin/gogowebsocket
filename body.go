package gogowebsocket

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
)

const (
	BODY_TYPE_TEXT      = 1
	BODY_TYPE_JSON      = 0
	BODY_TYPE_BYTES     = 2 // 如果是这个格式，将会把WSBody包按包体格式如下通过字节流发送给前端,包体格式参考packBinaryMessage函数
	BODY_TYPE_HEARTBEAT = 10000
	PONG                = "pong"

	ID_ERROR     = -1
	ID_HEARTBEAT = -2
)

type Error struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type WSBody struct {
	ClientID   string      `json:"-"`
	Client     *Client     `json:"-"`
	ProtocolId int64       `json:"protocol_id"`
	BodyType   int         `json:"body_type"`
	Queue      int         `json:"queue"` //是否全局队列处理消息,提供前端api接口时需要标记清楚
	Body       interface{} `json:"body"`
}

// 压包 转换为字节流传输
// 包体格式如下： uint32(0-3): protocolId, uint16(4-5):bodyType, uint32(6-9): bodySize, 10-end: body
func (m *WSBody) packBinaryMessage() ([]byte, error) {
	// 转换协议ID为固定字节序
	var protocolId = m.ProtocolId
	headBytes := make([]byte, 10)
	binary.BigEndian.PutUint32(headBytes, uint32(protocolId))

	// 获取body数据
	body, ok := m.Body.([]byte)
	if !ok {
		return nil, fmt.Errorf("invalid body type")
	}
	bodySize := len(body)

	binary.BigEndian.PutUint16(headBytes[4:6], uint16(m.BodyType))
	// 填充body大小 (6-9)
	binary.BigEndian.PutUint32(headBytes[6:], uint32(bodySize))

	// 创建缓冲区
	bytes := make([]byte, len(headBytes)+bodySize)
	// 填充协议头
	copy(bytes[0:len(headBytes)], headBytes)
	// 填充body内容 (10-end)
	copy(bytes[len(headBytes):], body)

	return bytes, nil
}

// 解包
func (m *WSBody) unpackBinaryMessage(bytes []byte) error {
	// 转换协议ID为固定字节序
	if len(bytes) < 10 {
		return fmt.Errorf("invalid message bytes")
	}
	var protocolId = binary.BigEndian.Uint32(bytes)
	bodyType := binary.BigEndian.Uint16(bytes[4:6])
	//bodySize := binary.BigEndian.Uint32(bytes[6:10])
	body := bytes[10:]
	m.ProtocolId = int64(protocolId)
	m.BodyType = int(bodyType)
	if bodyType == BODY_TYPE_TEXT {
		m.Body = string(body)
	} else if bodyType == BODY_TYPE_JSON {
		var tmpBody map[string]interface{}
		err := json.Unmarshal(body, &tmpBody)
		m.Body = tmpBody
		return err
	} else {
		m.Body = body
	}
	return nil
}

func (m *WSBody) BodyToStruct(to interface{}) error {
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

func (m *WSBody) BodyToString() (string, error) {
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
