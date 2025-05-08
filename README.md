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


js解析数据：
```javascript
conn.onmessage = function (evt) {
            if(evt.data instanceof Blob){
                // 字节流
                let reader = new FileReader()
                // 以ArrayBuffer形式读取Blob
                reader.readAsArrayBuffer(evt.data);
                reader.onload = e => {
                    const arrayBuffer = e.target.result;
                    // 验证数据长度
                    if (arrayBuffer.byteLength < 10) {
                        throw new Error("无效数据包：数据长度不足8字节");
                    }
                    const dataView = new DataView(arrayBuffer);
                    // 读取protocolId（大端序）
                    const protocolId = dataView.getUint32(0, false);
                    const bodyType = dataView.getUint16(4, false);
                    // 读取bodySize（大端序）
                    const bodySize = dataView.getUint32(6, false);
                    // 验证数据完整性
                    const totalLength = 10 + bodySize;
                    if (arrayBuffer.byteLength < totalLength) {
                        throw new Error(`数据包不完整，期望长度：${totalLength}，实际长度：${arrayBuffer.byteLength}`);
                    }
                    console.log("收到字节流", protocolId, bodyType, bodySize);
                    // 提取body数据
                    const bodyBuffer = arrayBuffer.slice(10, totalLength);
                    const decoder = new TextDecoder(); // 默认使用UTF-8解码
                    const bodytext = decoder.decode(bodyBuffer);
                    var item = document.createElement("div");
                    item.innerText = bodytext;
                    appendLog(item);
                }
                reader.onerror = function(e) {
                    console.error(e);
                };
            }else{
                let msgobj = JSON.parse(evt.data)
                console.log(msgobj)
                let bodytext = msgobj["body"]
                if (msgobj["body_type"] == BODY_TYPE_JSON){
                    bodytext = JSON.stringify(bodytext)
                }else if (msgobj["body_type"] == BODY_TYPE_BYTES){
                    return
                }
                var item = document.createElement("div");
                item.innerText = bodytext;
                appendLog(item);
            }
        };
```

js发送数据:
```javascript
if (msg.value=="binary"){
    msgobj.protocol_id=9100
    msgobj.body_type=BODY_TYPE_TEXT
    msgobj.body = "你好，我是前端数据"
    msgobj.body_format = 'binary'
}else{
    msgobj.protocol_id=1000
    msgobj.body_type=BODY_TYPE_JSON
    msgobj.body = {
        msg: msg.value,
    }
}
sendMessge(msgobj);

function sendMessge(msgobj) {
    if (msgobj.body_format == 'binary'){
        // 给服务器发送字节流数据
        const encoder = new TextEncoder();
        // 假设 BODY_TYPE_BYTES 是已定义的常量（示例值 0x01）
        const bodyType = BODY_TYPE_BYTES;
        // 1. 编码业务数据
        const databytes = encoder.encode(msgobj.body); // 这里只是简单的按字符串
        // 2. 计算各部分长度
        const headerSize = 10; // 协议头固定10字节
        const bodySize = databytes.length; // body部分的实际长度
        const totalSize = headerSize + bodySize; // 总长度
        // 3. 创建缓冲区
        const buf = new ArrayBuffer(totalSize);
        const dataView = new DataView(buf);
        const uint8Array = new Uint8Array(buf); // 用于处理字节流操作
        // 4. 写入协议头
        dataView.setUint32(0, msgobj.protocol_id, false);         // protocolId (大端序)
        dataView.setUint16(4, msgobj.body_type, false); // bodyType (大端序)
        dataView.setUint32(6, bodySize, false);     // bodySize (大端序)
        // 5. 写入body数据
        // 使用Uint8Array的set方法将业务数据复制到缓冲区
        uint8Array.set(databytes, 10); // 从第10字节开始覆盖
        // 6. 发送完整数据包
        console.log('发送的数据包结构：');
        console.log('protocolId:', 9100);
        console.log('bodyType:', BODY_TYPE_BYTES.toString(16));
        console.log('bodySize:', bodySize);
        console.log('body内容:', new TextDecoder().decode(databytes));
        conn.send(buf);
    }else{
        console.log("send::", JSON.stringify(msgobj));
        conn.send(JSON.stringify(msgobj));
    }
}
```