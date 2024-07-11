const BODY_TYPE_HEARTBEAT = 10000
const BODY_TYPE_JSON      = 0
const BODY_TYPE_TEXT      = 1
const BODY_TYPE_BYTES     = 2
var conn;
var msg;
var log ;
var user = {}
function appendLog(item) {
    var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
    log.appendChild(item);
    if (doScroll) {
        log.scrollTop = log.scrollHeight - log.clientHeight;
    }
}

function onsubmit () {
    if (!conn) {
        return false;
    }
    if (!msg.value) {
        return false;
    }
    let msgobj = {}
    if (msg.value=="-2"){
        // 心跳
        msgobj.protocol_id=-2 //固定的协议
        msgobj.body_type=BODY_TYPE_HEARTBEAT
        msgobj.body = ""
    }else{
        msgobj.protocol_id=1000
        msgobj.body_type=BODY_TYPE_JSON
        msgobj.body = {
            uid:user.uid,
            name:user.name,
            production:user.production,
            deviceId: user.deviceId,
            msg: msg.value,
        }
    }
    conn.send(JSON.stringify(msgobj));
    msg.value = "";
    return false;
};

function createWebsocket(userInfo) {
    user = userInfo
    msg = document.getElementById("msg")
    log = document.getElementById("log")
    document.getElementById("form").onsubmit = onsubmit;
    if (window["WebSocket"]) {
        conn = new WebSocket("ws://localhost:15800/v1/ws?X-ZY-USERID="+user.uid+"&X-ZY-PRODUCTION-EXT="+user.production+"&X-ZY-DEVICEID="+user.deviceId);
        conn.onclose = function (evt) {
            var item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
            appendLog(item);
        };
        conn.onopen = function(evt){
            var item = document.createElement("div");
            item.innerHTML = "<b>ws connected.</b>"
            appendLog(item);
        };
        conn.onmessage = function (evt) {
            let msgobj = JSON.parse(evt.data)
            console.log(msgobj)
            let bodytext = msgobj["body"]
            if (msgobj["body_type"] == BODY_TYPE_JSON){
                bodytext = JSON.stringify(bodytext)
            }else if (msgobj["body_type"] == BODY_TYPE_BYTES){
                var decodeStr  = window.atob(bodytext)
                let blob = new Blob([decodeStr ]);
                let reader = new FileReader()
                reader.readAsText(blob, 'utf-8')
                reader.onload = e => {
                    let readerres = reader.result
                    bodytext = JSON.parse(readerres)
                    bodytext = JSON.stringify(bodytext)
                    var item = document.createElement("div");
                    item.innerText = bodytext;
                    appendLog(item);
                }
                return
            }
            var item = document.createElement("div");
            item.innerText = bodytext;
            appendLog(item);

        };
    } else {
        var item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        appendLog(item);
    }
}