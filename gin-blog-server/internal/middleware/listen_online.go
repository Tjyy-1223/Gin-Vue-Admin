package middleware

import (
	"context"
	"fmt"
	"gin-blog-server/internal/global"
	"gin-blog-server/internal/handle"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

// ListenOnline 监听在线状态
// 每次请求时检查用户是否被强制下线，并更新用户的在线状态，确保用户的在线状态在一定时间内保持有效。
func ListenOnline() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		rdb := c.MustGet(global.CTX_RDB).(*redis.Client)

		auth, err := handle.CurrentUserAuth(c)
		if err != nil {
			handle.ReturnError(c, global.ErrDbOp, err)
			return
		}

		onlineKey := global.ONLINE_USER + strconv.Itoa(auth.ID)
		offlineKey := global.OFFLINE_USER + strconv.Itoa(auth.ID)

		// 判断当前用户是否被强制下线
		exists, err := rdb.Exists(ctx, offlineKey).Result()
		if err != nil {
			// 处理 Redis 操作错误
			handle.ReturnError(c, global.ErrRedisOp, err)
			c.Abort()
			return
		}
		if exists == 1 {
			fmt.Println("用户被强制下线")
			handle.ReturnError(c, global.ErrForceOffline, nil)
			c.Abort()
			return
		}

		// 每次发送请求会更新 Redis 中的在线状态: 重新计算 10 分钟
		err = rdb.Set(ctx, onlineKey, auth, 10*time.Minute).Err()
		if err != nil {
			handle.ReturnError(c, global.ErrRedisOp, err)
			c.Abort()
			return
		}

		c.Next()
	}
}
