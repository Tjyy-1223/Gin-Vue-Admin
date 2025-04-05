package handle

import (
	"context"
	"encoding/json"
	"gin-blog-server/internal/global"
	"gin-blog-server/internal/model"
	"github.com/redis/go-redis/v9"
	"time"
)

// redis context
var rctx = context.Background()

// Config
// addConfigCache 将博客配置缓存到 Redis 中
func addConfigCache(rdb *redis.Client, config map[string]string) error {
	return rdb.HMSet(rctx, global.CONFIG, config).Err()
}

// removeConfigCache 删除 Redis 中博客配置缓存
func removeConfigCache(rdb *redis.Client) error {
	return rdb.Del(rctx, global.CONFIG).Err()
}

// 从 Redis 中获取博客配置缓存
// rdb.HGetAll 如果不存在 key, 不会返回 redis.Nil 错误, 而是返回空 map
func getConfigCache(rdb *redis.Client) (cache map[string]string, err error) {
	return rdb.HGetAll(rctx, global.CONFIG).Result()
}

// SetMailInfo 将邮箱信息存储到 rdb 中
func SetMailInfo(rdb *redis.Client, info string, expire time.Duration) error {
	return rdb.Set(rctx, info, true, expire).Err()
}

// GetMailInfo 检测 rdb 中是否存在邮箱信息
func GetMailInfo(rdb *redis.Client, info string) (bool, error) {
	return rdb.Get(rctx, info).Bool()
}

// DeleteMailInfo 从 rdb 中删除相关的邮箱信息
func DeleteMailInfo(rdb *redis.Client, info string) error {
	return rdb.Del(rctx, info).Err()
}

// addPageCache 将页面列表缓存到 Redis 中
func addPageCache(rdb *redis.Client, pages []model.Page) error {
	data, err := json.Marshal(pages)
	if err != nil {
		return err
	}
	return rdb.Set(rctx, global.PAGE, string(data), 0).Err()
}

// getPageCache 从 Redis 中获取页面列表缓存
// rdb.Get 如果不存在 key，会返回 redis.Nil 错误
func getPageCache(rdb *redis.Client) (cache []model.Page, err error) {
	s, err := rdb.Get(rctx, global.PAGE).Result()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(s), &cache); err != nil {
		return nil, err
	}
	return cache, nil
}

// 删除 Redis 中页面列表缓存
func removePageCache(rdb *redis.Client) error {
	return rdb.Del(rctx, global.PAGE).Err()
}
