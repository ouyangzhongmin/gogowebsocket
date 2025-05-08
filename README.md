# gogowebsocket
一个高并发的websocket库, 通过redis记录连接信息和服务器信息，grpc内部转发消息， 本库仅做参考

如果传输消息类型为:BODY_TYPE_BYTES ,则前端解析包体时按如下规则：“包体格式如下： uint32(0-3): protocolId, uint16(4-5):bodyType, uint32(6-9): bodySize, 10-end: body”。
如果前端需要向服务器传输字节流数据也需要按上面的包体格式。

其他消息类型可以直接按json解析, 服务器定义包体如下：
```
    type WSBody struct {
        ProtocolId int64       `json:"protocol_id"`
        BodyType   int         `json:"body_type"`
        Queue      int         `json:"queue"` 
        Body       interface{} `json:"body"`
    }
```


