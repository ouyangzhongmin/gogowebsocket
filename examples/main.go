package main

import (
	"context"
	"flag"
	"github.com/ouyangzhongmin/gogowebsocket/examples/handler"
	"github.com/ouyangzhongmin/gogowebsocket/examples/router"
	"github.com/ouyangzhongmin/gogowebsocket/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var port = flag.String("p", "", "http port")
var rpcPort = flag.String("rp", "", "rpc port")

func main() {
	flag.Parse()
	logger.Init(&logger.LogConf{
		FileName:     "",
		Level:        "debug",
		ReportCaller: true,
	})
	logger.Log.Println("start gin ..")

	rport := *rpcPort
	if rport == "" {
		rport = "15200"
	}
	handler.InitServices(rport)
	rt := router.InitRouter()

	addr := *port
	if addr == "" {
		addr = ":15800"
	}
	srv := &http.Server{
		Addr:    addr,
		Handler: rt,
	}

	go func() {
		logger.Println("server start listening at ", addr, "....")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Panic("listen error ", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	//通知ws断开所有连接
	handler.Shutdown()
	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Fatal("Server Shutdown err:", err)
	}
	select {
	case <-ctx.Done():
		logger.Println("timeout of 2 seconds.")
	}
	logger.Println("Server exiting")
}
