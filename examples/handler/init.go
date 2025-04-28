package handler

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/ouyangzhongmin/gogowebsocket/examples/service"
	"strings"
	"time"
)

var (
	wsSrv *service.Service
)

func InitServices(rpcPort string) {
	rds := NewRedis()
	wsSrv = service.New(rds, rpcPort)
}

func NewRedis() *redis.Client {
	rds := redis.NewClient(&redis.Options{
		Addr:         "172.16.2.7:6379",
		Password:     "", // no password set
		DB:           0,  // use default DB
		DialTimeout:  time.Millisecond * 1000,
		ReadTimeout:  time.Millisecond * 1000,
		WriteTimeout: time.Millisecond * 1000,
	})
	cmd := rds.Ping(context.Background())
	status, err := cmd.Result()
	if strings.ToUpper(status) == "PONG" {
		fmt.Println("redis ok")
	} else {
		fmt.Println("redis fail", err)
	}
	return rds
}

func Shutdown() {
	wsSrv.Shutdown()
}
