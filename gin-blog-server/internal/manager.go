package ginblog

import (
	"gin-blog-server/docs"
	"gin-blog-server/internal/handle"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	// 后端管理系统接口
	userAuthAPI handle.UserAuth // 用户账号
)

func RegisterHandlers(r *gin.Engine) {
	// Swagger 配置：设置 Swagger API 文档的基础路径为 /api
	// SwaggerInfo 是由 `swaggo/swag` 生成的文档信息配置结构体
	docs.SwaggerInfo.BasePath = "/api" // 设置 Swagger 文档的基础路径为 "/api"

	// 注册 Swagger UI 路由，`ginSwagger.WrapHandler` 会将 Swagger 文档与 UI 渲染结合
	// 这样用户可以通过访问 `/swagger` 路由来查看 API 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 调用 registerBaseHandler 注册其他基础的路由处理
	// 例如，可能是一些常见的基础 API 路由
	registerBaseHandler(r)
}

// 通用接口: 全部不需要 登录 + 鉴权
func registerBaseHandler(r *gin.Engine) {
	base := r.Group("/api")

	// TODO: 登录, 注册 记录日志
	base.POST("/login", userAuthAPI.Login)            // 登录
	base.POST("/register", userAuthAPI.Register)      // 注册
	base.GET("/email/verify", userAuthAPI.VerifyCode) // 邮箱验证
	base.GET("/logout", userAuthAPI.Logout)           // 退出登录
}
