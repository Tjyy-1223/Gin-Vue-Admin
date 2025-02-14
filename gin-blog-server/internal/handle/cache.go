package handle

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

// redis context
var rctx = context.Background()

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
