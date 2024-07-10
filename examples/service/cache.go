//Copyright The ZHIYUNCo.All rights reserved.
//Created by admin at2024/7/5.

package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

const CACHE_HASH_KEY = "zp_wsmap_%s"
const CACHE_DURATION = 60 * 60 * 24 * 20 //10天

type cache struct {
	r *redis.Client
}

func newCache(r *redis.Client) *cache {
	return &cache{r: r}
}

// 更新映射到缓存
func (c *cache) putUserClientID(userKey string, clientId string) error {
	if c.r == nil {
		return errors.New("redis未初始化")
	}
	ctx := context.Background()
	key := fmt.Sprintf(CACHE_HASH_KEY, userKey)
	cmd := c.r.ZAdd(ctx, key, &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: clientId,
	})
	if cmd.Err() != nil {
		return cmd.Err()
	}
	//设置列表的有效期，这里要注意可能有连接保持很多天不会断开的，实际业务再实际分析处理
	cmd2 := c.r.Expire(ctx, key, time.Second*CACHE_DURATION)
	if cmd2.Err() != nil {
		return cmd2.Err()
	}
	return nil
}

func (c *cache) removeUserClientID(userKey string, clientId string) error {
	if c.r == nil {
		return errors.New("redis未初始化")
	}
	ctx := context.Background()
	key := fmt.Sprintf(CACHE_HASH_KEY, userKey)
	cmd := c.r.ZRem(ctx, key, clientId)
	return cmd.Err()
}

// 返回全部的服务器列表
func (c *cache) getUserClientIDs(userKey string) ([]string, error) {
	if c.r == nil {
		return nil, errors.New("redis未初始化")
	}
	key := fmt.Sprintf(CACHE_HASH_KEY, userKey)
	cmd := c.r.ZRevRange(context.Background(), key, 0, -1)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}
	list := cmd.Val()
	if len(list) > 5 {
		//限制数量，防止过多
		list = list[:5]
	}
	return list, nil
}
