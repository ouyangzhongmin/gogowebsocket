//Copyright The ZHIYUNCo.All rights reserved.
//Created by admin at2022/9/20.

package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ouyangzhongmin/gogowebsocket/examples/service"
	"github.com/ouyangzhongmin/gogowebsocket/logger"
	"strconv"
)

func OnWSHandler(c *gin.Context) {
	production := c.GetHeader("X-ZY-PRODUCTION-EXT")
	deviceId := c.GetHeader("X-ZY-DEVICEID")
	uid, _ := strconv.ParseUint(c.GetHeader("X-ZY-USERID"), 10, 64)
	if production == "" {
		GinErrorWithMsgP(c, 1, "X-ZY-Production-Ext不允许为空")
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
