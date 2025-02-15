package handle

import (
	"context"
	"gin-blog-server/internal/global"
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
