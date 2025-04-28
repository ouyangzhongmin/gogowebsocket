package gogowebsocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/ouyangzhongmin/gogowebsocket/logger"
	"time"
)

const CACHE_HASH_KEY = "websocket_clients_%s"
const CACHE_SERVERS_KEY = "websocket_servers_%s"
const SERVER_OFFLINE_DURATION = 120
const CLIENT_OFFLINE_DURATION = 120
const SERVER_LIST_CACHE_TIME = 300

type cacheServerInfo struct {
	serverInfo
	StartTime string `json:"start_time"` //服务器启动时间
	Ts        int64  `json:"ts"`         //最新在线时间戳
	TimeStr   string `json:"time_str"`   //最新在线时间
}

// 判定是否在线
func (i *cacheServerInfo) isOnline() bool {
	nowTs := time.Now().Unix()
	//fmt.Println("server.isOnline::", nowTs, i.Ts)
	if nowTs-i.Ts >= SERVER_OFFLINE_DURATION {
		return false
	}
	return true
}

func (i *cacheServerInfo) toString() string {
	return fmt.Sprintf("%s-%s-%d-%s-%s", i.ServerIP, i.Port, i.Ts, i.TimeStr, i.StartTime)
}

type cacheClientInfo struct {
	serverInfo
	Ts  int64 `json:"ts"`  //最新在线时间戳
	CTs int64 `json:"cts"` //连接时间戳
}

// 判定是否在线
func (i *cacheClientInfo) isOnline() bool {
	nowTs := time.Now().Unix()
	if nowTs-i.Ts >= CLIENT_OFFLINE_DURATION {
		return false
	}
	return true
}

type cache struct {
	appId string
	r     *redis.Client
}

func newCache(appId string, r *redis.Client) *cache {
	return &cache{appId: appId, r: r}
}

func (c *cache) isAvailable() bool {
	return c.r != nil
}

// 更新服务器到缓存
func (c *cache) putServerInfo(info *serverInfo, startTime time.Time) error {
	if c.isAvailable() {
		tmp := &cacheServerInfo{
			serverInfo: *info,
			StartTime:  startTime.Format(time.RFC3339),  //服务器启动时间
			TimeStr:    time.Now().Format(time.RFC3339), //最新在线时间
			Ts:         time.Now().Unix(),               //最新在线时间戳
		}
		ctx := context.Background()
		key := fmt.Sprintf(CACHE_SERVERS_KEY, c.appId)
		str, err := json.Marshal(tmp)
		if err != nil {
			return err
		}
		cmd := c.r.HSet(ctx, key, info.toString(), str)
		if cmd.Err() != nil {
			return cmd.Err()
		}
		//设置列表的有效期，如果没有服务器更新了，则列表到期清空
		cmd2 := c.r.Expire(ctx, key, time.Second*SERVER_LIST_CACHE_TIME)
		if cmd2.Err() != nil {
			return cmd2.Err()
		}
	}
	return nil
}

// 判定服务器是否在线
func (c *cache) isServerOnline(info *serverInfo) (bool, error) {
	if !c.isAvailable() {
		return false, errors.New("redis未初始化")
	}
	key := fmt.Sprintf(CACHE_SERVERS_KEY, c.appId)
	cmd := c.r.HGet(context.Background(), key, info.toString())
	if cmd.Err() != nil {
		return false, cmd.Err()
	}
	str := cmd.Val()
	var cacheInfo cacheServerInfo
	err := json.Unmarshal([]byte(str), &cacheInfo)
	if err != nil {
		return false, err
	}
	return cacheInfo.isOnline(), nil
}

// 删除服务器上的缓存
func (c *cache) removeServerInfo(info *serverInfo) error {
	if c.isAvailable() {
		key := fmt.Sprintf(CACHE_SERVERS_KEY, c.appId)
		cmd := c.r.HDel(context.Background(), key, info.toString())
		if cmd.Err() != nil {
			return cmd.Err()
		}
	}
	return nil
}

// 返回全部的服务器列表
func (c *cache) getServerList() ([]cacheServerInfo, error) {
	if !c.isAvailable() {
		return []cacheServerInfo{}, nil
	}
	key := fmt.Sprintf(CACHE_SERVERS_KEY, c.appId)
	cmd := c.r.HGetAll(context.Background(), key)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	m := cmd.Val()
	list := make([]cacheServerInfo, 0)
	for _, v := range m {
		var tmp cacheServerInfo
		err := json.Unmarshal([]byte(v), &tmp)
		if err != nil {
			logger.Errorln("getServerInfos err::", err)
		} else {
			list = append(list, tmp)
		}
	}
	return list, nil
}

// 保存用户连接信息到redis
func (c *cache) putClientInfo(clientId string, info cacheClientInfo) error {
	if c.isAvailable() {
		key := fmt.Sprintf(CACHE_HASH_KEY, c.appId)
		jsonbytes, err := json.Marshal(&info)
		if err != nil {
			return err
		}
		cmd := c.r.HSet(context.Background(), key, clientId, string(jsonbytes))
		if cmd.Err() != nil {
			return cmd.Err()
		}
	}
	return nil
}

func (c *cache) getClientInfo(clientId string) (*cacheClientInfo, error) {
	if !c.isAvailable() {
		return nil, errors.New("redis未初始化")
	}
	key := fmt.Sprintf(CACHE_HASH_KEY, c.appId)
	cmd := c.r.HGet(context.Background(), key, clientId)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	str := cmd.Val()
	var info cacheClientInfo
	err := json.Unmarshal([]byte(str), &info)
	return &info, err
}

func (c *cache) removeClientInfo(clientId string) error {
	if c.isAvailable() {
		key := fmt.Sprintf(CACHE_HASH_KEY, c.appId)
		cmd := c.r.HDel(context.Background(), key, clientId)
		if cmd.Err() != nil {
			return cmd.Err()
		}
	}
	return nil
}
