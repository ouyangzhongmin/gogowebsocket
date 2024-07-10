//Copyright The ZHIYUNCo.All rights reserved.
//Created by admin at2024/7/10.

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ouyangzhongmin/gogowebsocket/examples/handler"
	"github.com/ouyangzhongmin/gogowebsocket/logger"
	"net/http"
)

func InitRouter() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	//初始化业务层的故事
	handler.InitServices()

	//初始化路由
	initRouter(r)
	return r
}

func initRouter(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		v1.Any("/ws", Cors(nil), handler.OnWSHandler)
	}
}

func Cors(allowOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("origin")
		if allowOrigins != nil && len(allowOrigins) > 0 {
			hasFind := false
			for _, v := range allowOrigins {
				if v == origin {
					hasFind = true
					break
				}
			}
			if !hasFind {
				logger.Warnln("不支持的跨域请求：", origin)
				origin = ""
			}
		}
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			//c.Header("Access-Control-Allow-Headers", "*")
			acrheaders := c.GetHeader("access-control-request-Headers")
			if acrheaders != "" {
				c.Header("Access-Control-Allow-Headers", acrheaders)
			} else {
				c.Header("Access-Control-Allow-Headers", "*")
			}
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "*")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "60")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
			//放行所有OPTIONS方法
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(http.StatusNoContent)
				return
			}
		}

		c.Next()
	}
}
