package ginblog

import (
	"gin-blog-server/docs"
	"gin-blog-server/internal/handle"
	"gin-blog-server/internal/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	// 后端管理系统接口
	userAuthAPI handle.UserAuth // 用户账号
	blogInfoAPI handle.BlogInfo // 博客设置
	userAPI     handle.User     // 用户
	pageAPI     handle.Page     // 页面
	frontAPI    handle.Front    // 博客前台接口
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
	registerAdminHandler(r)
	registerBlogHandler(r)
}

// 通用接口: 全部不需要 登录 + 鉴权
func registerBaseHandler(r *gin.Engine) {
	base := r.Group("/api")

	// TODO: 登录, 注册 记录日志
	base.POST("/login", userAuthAPI.Login)            // 登录
	base.POST("/register", userAuthAPI.Register)      // 注册
	base.GET("/email/verify", userAuthAPI.VerifyCode) // 邮箱验证
	base.GET("/logout", userAuthAPI.Logout)           // 退出登录
	base.POST("/report", blogInfoAPI.Report)          // 上报信息
	base.GET("/config", blogInfoAPI.GetConfigMap)     // 获取配置
	base.PATCH("/config", blogInfoAPI.UpdateConfig)   // 更新配置
}

// 后台管理系统的接口: 全部需要 登录 + 鉴权
func registerAdminHandler(r *gin.Engine) {
	auth := r.Group("/api")

	// 注意使用中间件的顺序
	auth.Use(middleware.JWTAuth())
	auth.Use(middleware.PermissionCheck())
	auth.Use(middleware.OperationLog())
	auth.Use(middleware.ListenOnline())

	// 博客设置
	setting := auth.Group("/setting")
	{
		setting.GET("/about", blogInfoAPI.GetAbout)    // 获取关于我
		setting.PUT("/about", blogInfoAPI.UpdateAbout) // 编辑关于我
	}

	// 用户模块
	user := auth.Group("/user")
	{
		user.GET("/info", userAPI.GetInfo)          // 获取当前用户信息
		user.GET("/current", userAPI.UpdateCurrent) // 修改当前用户信息
	}

	// 资源模块
	//resource := auth.Group("/resource")
	//{
	//	resource.GET("/list", resourceAPI.GetTreeList) // 资源列表(树形)
	//}
}

// 博客前台相关接口：大部分不需要登陆，部分需要登陆
func registerBlogHandler(r *gin.Engine) {
	base := r.Group("/api/front")

	base.GET("/about", blogInfoAPI.GetAbout) // 获取关于我
	base.GET("/home", frontAPI.GetHomeInfo)  // 前台首页
	base.GET("/page", pageAPI.GetList)       // 前台页面

	// 需要登录才能进行的操作
	base.Use(middleware.JWTAuth())
	{
		base.GET("/user/info", userAPI.GetInfo)       // 根据 Token 获取用户信息
		base.PUT("/user/info", userAPI.UpdateCurrent) // 根据 Token 更新当前用户信息
	}
}
