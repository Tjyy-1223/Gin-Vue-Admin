package middleware

import (
	"fmt"
	"gin-blog-server/internal/global"
	"gin-blog-server/internal/handle"
	"gin-blog-server/internal/model"
	"gin-blog-server/internal/utils/jwt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"strings"
	"time"
)

// JWTAuth 基于 jwt 实现鉴权
// TODO: 如果存在 session, 则直接从 session 中获取用户信息
// 从 Authorization 中获取 token, 并解析 token 获取用户信息, 并设置到 session 中
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// FIXME: 前后台 session 混乱, 暂时无法将用户信息挂载在 gin context 缓存
		slog.Debug("[middleware-JWTAuth] user auth not exist, do jwt auth")

		db := c.MustGet(global.CTX_DB).(*gorm.DB)

		// 系统管理的资源需要进行用户鉴权，其他资源不需要鉴权
		//url, method := c.FullPath()[4:], c.Request.Method
		//resource, err := model.GetResource(db, url, method)
		//if err != nil {
		//	// 没有找到的资源，不需要鉴权，跳过后续的验证过程
		//	if errors.Is(err, gorm.ErrRecordNotFound) {
		//		slog.Debug("[middleware-JWTAuth] resource not exist, skip jwt auth")
		//		c.Set("skip_check", true)
		//		c.Next()
		//		c.Set("skip_check", false)
		//		return
		//	}
		//	handle.ReturnError(c, global.ErrDbOp, err)
		//	return
		//}

		// 匿名资源，不需要鉴权，跳过后续的验证过程
		//if resource.Anonymous {
		//	slog.Debug(fmt.Sprintf("[middleware-JWTAuth] resouce: %s %s is anonymous, skip jwt auth!", url, method))
		//	c.Set("skip_check", true)
		//	c.Next()
		//	c.Set("skip_check", false)
		//	return
		//}

		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			handle.ReturnError(c, global.ErrTokenNotExist, nil)
			return
		}

		// token 的正确格式: `Bearer [tokenString]`
		parts := strings.Split(authorization, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			handle.ReturnError(c, global.ErrTokenType, nil)
			return
		}

		claims, err := jwt.ParseToken(global.Conf.JWT.Secret, parts[1])
		if err != nil {
			handle.ReturnError(c, global.ErrTokenWrong, err)
			return
		}

		// 判断 token 已经过期
		if time.Now().Unix() > claims.ExpiresAt.Unix() {
			handle.ReturnError(c, global.ErrTokenRuntime, nil)
			return
		}

		// 获取用户信息
		user, err := model.GetUserAuthInfoById(db, claims.UserId)
		if err != nil {
			handle.ReturnError(c, global.ErrUserNotExist, err)
			return
		}

		// session 设置
		session := sessions.Default(c)
		session.Set(global.CTX_USER_AUTH, claims.UserId)
		session.Save()

		// gin context
		c.Set(global.CTX_USER_AUTH, user)
	}
}

// PermissionCheck 资源访问权限验证
func PermissionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetBool("skip_check") {
			c.Next()
			return
		}

		db := c.MustGet(global.CTX_DB).(*gorm.DB)
		auth, err := handle.CurrentUserAuth(c)
		if err != nil {
			handle.ReturnError(c, global.ErrUserNotExist, err)
			return
		}

		if auth.IsSuper {
			slog.Debug("[middleware-PermissionCheck]: super admin no need to check, pass!")
			c.Next()
			return
		}

		url := c.FullPath()[4:]
		method := c.Request.Method

		slog.Debug(fmt.Sprintf("[middleware-PermissionCheck] %v, %v, %v\n", auth.Username, url, method))
		for _, role := range auth.Roles {
			slog.Debug(fmt.Sprintf("[middleware-PermissionCheck] %v\n", role.Name))
			pass, err := model.CheckRoleAuth(db, role.ID, url, method)
			if err != nil {
				handle.ReturnError(c, global.ErrDbOp, err)
				return
			}
			if !pass {
				handle.ReturnError(c, global.ErrPermission, nil)
				return
			}
		}

		slog.Debug("[middleware-PermissionCheck]: pass")
		c.Next()
	}
}
