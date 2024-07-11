package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ouyangzhongmin/gogowebsocket/examples/service"
	"github.com/ouyangzhongmin/gogowebsocket/logger"
	"strconv"
)

func OnWSHandler(c *gin.Context) {
	production := c.Query("production")
	deviceId := c.Query("deviceId")
	uid, _ := strconv.ParseUint(c.Query("uid"), 10, 64)
	if production == "" {
		GinErrorWithMsgP(c, 1, "production不允许为空")
		return
	}
	logger.Println("用户发起ws连接：", production, uid)
	err := wsSrv.ServeWS(&service.WSUserInfo{
		UserId:     int(uid),
		Production: production,
		DeviceId:   deviceId,
	}, c.Writer, c.Request)
	if err != nil {
		GinErrorWithMsgP(c, 2, err.Error())
		return
	}
}

func GinErrorWithMsgP(c *gin.Context, errcode int, errmsg string) {
	c.Set("errcode", errcode)
	c.Set("errmsg", errmsg)
	logger.Errorln("errcode:", errcode, ",errmsg:", errmsg)
}
