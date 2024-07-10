package main

import (
	"context"
	"github.com/ouyangzhongmin/gogowebsocket/examples/handler"
	"github.com/ouyangzhongmin/gogowebsocket/examples/router"
	"github.com/ouyangzhongmin/gogowebsocket/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger.Init(&logger.LogConf{
		FileName:     "",
		Level:        "debug",
		ReportCaller: true,
	})
	logger.Log.Println("init server , start gin ..")

	//初始化路由
	rt := router.InitRouter()
	addr := ":15800"
	srv := &http.Server{
		Addr:    addr,
		Handler: rt,
	}

	go func() {
		// service connections
		//gops工具接入,若需要远程访问，可配置 agent.Options{Addr: "0.0.0.0:6060"}，否则默认仅允许本地访问
		logger.Log.Println("server start listening at ", addr, "....")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Panicln("listen error ", err)
		}
	}()

	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Log.Fatal("Server Shutdown err:", err)
	}
	//通知ws断开所有连接
	handler.Shutdown()
	//catching ctx.Done(). timeout of 2 seconds.
	select {
	case <-ctx.Done():
		logger.Log.Println("timeout of 2 seconds.")
	}
	logger.Log.Println("Server exiting")
}