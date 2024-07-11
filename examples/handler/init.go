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

func InitServices() {
	rdb := NewRedis()
	wsSrv = service.New(rdb)
}

func NewRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         "172.16.2.7:6379",
		Password:     "", // no password set
		DB:           0,  // use default DB
		DialTimeout:  time.Millisecond * 1000,
		ReadTimeout:  time.Millisecond * 1000,
		WriteTimeout: time.Millisecond * 1000,
	})
	cmd := rdb.Ping(context.Background())
	status, err := cmd.Result()
	if strings.ToUpper(status) == "PONG" {
		fmt.Println("redis ok")
	} else {
		fmt.Println("redis fail", err)
	}
	return rdb
}

func Shutdown() {
	wsSrv.Shutdown()
}
