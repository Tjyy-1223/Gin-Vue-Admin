package middleware

import (
	"errors"
	"gin-blog-server/internal/global"
	"gin-blog-server/internal/handle"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// CORS 是一个Gin中间件，用于处理跨域请求的配置
func CORS() gin.HandlerFunc {
	// 使用cors.New创建并返回一个跨域处理的中间件
	return cors.New(cors.Config{
		// AllowOrigins 配置了允许的源（域名），这里设置为 "*"，表示允许任何域名发起跨域请求
		AllowOrigins: []string{"*"},

		// AllowMethods 配置了允许跨域请求使用的HTTP方法。这里设置了常见的HTTP方法：PUT、POST、GET、DELETE、OPTIONS 和 PATCH。
		AllowMethods: []string{"PUT", "POST", "GET", "DELETE", "OPTIONS", "PATCH"},

		// AllowHeaders 配置了允许的请求头，表示在跨域请求中可以携带哪些请求头。
		AllowHeaders: []string{"Origin", "Authorization", "Content-Type", "X-Requested-With"},

		// ExposeHeaders 配置了允许客户端访问的响应头，这里暴露了 "Content-Type" 头。
		ExposeHeaders: []string{"Content-Type"},

		// AllowCredentials 是否允许跨域请求时携带用户凭证（如cookies）。设置为true表示允许。
		AllowCredentials: true,

		// AllowOriginFunc 是一个函数，允许你对源进行动态验证，决定是否允许该源进行跨域请求。此处设置为总是返回true，表示任何源都被允许。
		AllowOriginFunc: func(origin string) bool {
			return true
		},

		// MaxAge 设置了预检请求的缓存时间，单位为时间.Duration。这里设置为24小时，表示预检请求结果会被缓存24小时。
		MaxAge: 24 * time.Hour,
	})
}

// WithRedisDB 将 redis.Client 注入到 gin.Context
// handler 中通过 c.MustGet(g.CTX_RDB).(*redis.Client) 来使用
func WithRedisDB(rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(global.CTX_RDB, rdb)
		ctx.Next()
	}
}

// WithGormDB 将 gorm.DB 注入到 gin.Context
// handler 中通过 c.MustGet(g.CTX_DB).(*gorm.DB) 来使用
func WithGormDB(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(global.CTX_DB, db)
		ctx.Next()
	}
}

// WithCookieStore 基于 Cookie 的 Session 中间件，用于在 Gin 框架中创建基于 Cookie 存储的会话管理。
// name: 会话名称，用于标识和存储会话。
// secret: 用于加密 Cookie 的密钥，保证 Cookie 的安全性。
func WithCookieStore(name, secret string) gin.HandlerFunc {
	// 创建一个新的 Cookie 存储实例，使用传入的密钥（secret）进行加密。
	store := cookie.NewStore([]byte(secret))

	// 设置 Cookie 存储的选项
	store.Options(sessions.Options{
		Path:   "/", // 指定 Cookie 的有效路径，"/" 表示整个网站都有效
		MaxAge: 600, // 设置 Cookie 的最大有效期，单位是秒（这里是 600秒，即10分钟）
	})

	// 返回一个 Gin 中间件，使用 sessions.Sessions 方法将会话存储配置应用到 Gin 路由中。
	// `name` 是会话的名称，`store` 是用于存储会话数据的 Cookie 存储实例。
	return sessions.Sessions(name, store)
}

// Logger 日志记录
// Logger 中间件函数，用于记录每个请求的日志信息
func Logger() gin.HandlerFunc {
	// 返回一个 gin.HandlerFunc 类型的中间件函数
	return func(c *gin.Context) {
		// 记录请求开始时间
		start := time.Now()

		// 调用 c.Next() 执行后续的中间件和路由处理
		c.Next()

		// 计算请求处理的耗时
		cost := time.Since(start)

		// 使用 slog 记录日志信息
		slog.Info("[GIN]", // 日志的标识，方便区分是来自 Gin 中间件的日志
			// 请求的路径（URL）
			slog.String("path", c.Request.URL.Path),
			// 请求的查询参数
			slog.String("query", c.Request.URL.RawQuery),
			// HTTP 响应的状态码
			slog.Int("status", c.Writer.Status()),
			// 请求方法（GET、POST等）
			slog.String("method", c.Request.Method),
			// 客户端的 IP 地址
			slog.String("ip", c.ClientIP()),
			// 响应的大小（字节）
			slog.Int("size", c.Writer.Size()),
			// 请求处理的耗时
			slog.Duration("cost", cost),
			// 以下是注释掉的部分，可以根据需要开启：
			// 请求体内容
			// slog.String("body", c.Request.PostForm.Encode()),
			// 用户代理（User-Agent）
			// slog.String("user-agent", c.Request.UserAgent()),
			// 错误信息
			// slog.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
		)
	}
}

// Recovery 中间件用于捕获并处理应用程序中的 panic 错误
// stack 参数控制是否打印 panic 的 stack trace 信息
func Recovery(stack bool) gin.HandlerFunc {
	// 返回一个 gin.HandlerFunc 类型的中间件函数
	return func(c *gin.Context) {
		// 使用 defer 确保函数结束时执行此处的错误处理逻辑
		defer func() {
			// 使用 recover() 捕获 panic 异常，如果有 panic 发生，err 将不为 nil
			if err := recover(); err != nil {
				// 检查是否是网络连接错误（例如：Broken Pipe），这些错误不需要堆栈信息
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					var se *os.SyscallError
					// 判断是否是 broken pipe 错误
					if errors.As(ne, &se) {
						seStr := strings.ToLower(se.Error())
						if strings.Contains(seStr, "broken pipe") ||
							strings.Contains(seStr, "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// 如果发生 panic，返回给客户端通用错误信息
				handle.ReturnHttpResponse(c, http.StatusInternalServerError, global.FAIL, global.GetMsg(global.FAIL), err)

				// 打印出发生 panic 时的 HTTP 请求信息
				httpRequest, _ := httputil.DumpRequest(c.Request, false)

				// 如果是 broken pipe 错误，不需要打印堆栈信息，直接记录错误并中止请求
				if brokenPipe {
					// 记录错误日志
					slog.Error(c.Request.URL.Path,
						slog.Any("error", err),
						slog.String("request", string(httpRequest)),
					)
					// 如果连接已经断开，不能再写入响应状态
					_ = c.Error(err.(error))
					c.Abort()
					return
				}

				// 如果需要打印 stack trace（堆栈信息），记录带有堆栈信息的错误日志
				if stack {
					slog.Error("[Recovery from panic]",
						slog.Any("error", err),
						slog.String("request", string(httpRequest)),
						slog.String("stack", string(debug.Stack())), // 打印堆栈信息
					)
				} else {
					// 如果不需要打印堆栈信息，只记录错误日志
					slog.Error("[Recovery from panic]",
						slog.Any("error", err),
						slog.String("request", string(httpRequest)),
					)
				}
				// 中止请求并返回 500 错误（内部服务器错误）
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		// 执行后续中间件和路由处理
		c.Next()
	}
}
