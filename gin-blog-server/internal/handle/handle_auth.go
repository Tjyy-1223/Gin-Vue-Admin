package handle

import (
	"errors"
	"gin-blog-server/internal/global"
	"gin-blog-server/internal/model"
	"gin-blog-server/internal/utils"
	"gin-blog-server/internal/utils/jwt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"strconv"
)

type UserAuth struct{}

// LoginReq 发送的登陆请求
type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginVO 登陆信息返回给前端的数据
type LoginVO struct {
	model.UserInfo

	// 点赞 Set： 用于记录用户点赞过的文章，评论
	ArticleLikeSet []string `json:"article_like_set"`
	CommentLikeSet []string `json:"comment_like_set"`
	Token          string   `json:"token"`
}

// @Summary 登录
// @Description 登录
// @Tags UserAuth
// @Param form body LoginReq true "登录"
// @Accept json
// @Produce json
// @Success 0 {object} Response[LoginVO]
// @Router /login [post]
func (*UserAuth) Login(c *gin.Context) {
	// 创建 LoginReq 结构体用于绑定前端传来的 JSON 数据
	var req LoginReq

	// 绑定请求体 JSON 数据到 req 结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		// 如果绑定数据失败，返回错误信息
		ReturnError(c, global.ErrRequest, err)
		return
	}

	// 获取数据库和 Redis 客户端实例
	db := GetDB(c)
	rdb := GetRDB(c)

	// 查询数据库，获取用户的身份信息（UserAuth）
	userAuth, err := model.GetUserAuthInfoByName(db, req.Username)
	if err != nil {
		// 如果没有找到用户，返回用户不存在的错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ReturnError(c, global.ErrUserNotExist, nil)
			return
		}
		// 如果查询发生数据库操作错误，返回错误
		ReturnError(c, global.ErrDbOp, err)
		return
	}

	// 检查传入的密码与数据库中存储的密码是否匹配
	if !utils.BcryptCheck(req.Password, userAuth.Password) {
		// 如果密码不匹配，返回密码错误
		ReturnError(c, global.ErrPassword, nil)
		return
	}

	// 获取请求中的 IP 地址和 IP 来源信息
	// FIXME: 可能无法正确读取 IP 地址，这需要解决
	ipAddress := utils.IP.GetIpAddress(c)
	ipSource := utils.IP.GetIpSourceSimpleIdle(ipAddress)

	// 根据 UserAuth 中的 UserInfoId 查询用户信息
	userInfo, err := model.GetUserInfoById(db, userAuth.UserInfoId)
	if err != nil {
		// 如果没有找到用户信息，返回错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ReturnError(c, global.ErrUserNotExist, nil)
			return
		}
		// 数据库操作出错，返回数据库错误
		ReturnError(c, global.ErrDbOp, err)
		return
	}

	// 获取用户的角色 ID 列表
	roleIds, err := model.GetRoleIdsByUserId(db, userAuth.ID)
	if err != nil {
		// 获取角色信息出错，返回错误
		ReturnError(c, global.ErrDbOp, err)
		return
	}

	// 获取用户在 Redis 中的文章点赞记录
	articleLikeSet, err := rdb.SMembers(rctx, global.ARTICLE_USER_LIKE_SET+strconv.Itoa(userAuth.ID)).Result()
	if err != nil {
		// 获取文章点赞信息出错，返回 Redis 操作错误
		ReturnError(c, global.ErrRedisOp, err)
		return
	}

	// 获取用户在 Redis 中的评论点赞记录
	commentLikeSet, err := rdb.SMembers(rctx, global.COMMENT_USER_LIKE_SET+strconv.Itoa(userAuth.ID)).Result()
	if err != nil {
		// 获取评论点赞信息出错，返回 Redis 操作错误
		ReturnError(c, global.ErrRedisOp, err)
		return
	}

	// 登录信息验证通过后，生成 JWT Token
	// UUID 生成方法：可以使用 ip 地址、浏览器信息和操作系统信息来生成唯一标识符（具体实现可调整）
	// 这里使用 jwt.GenToken 来生成 Token
	conf := global.GetConfig().JWT
	token, err := jwt.GenToken(conf.Secret, conf.Issuer, int(conf.Expire), userAuth.ID, roleIds)
	if err != nil {
		// Token 生成失败，返回错误
		ReturnError(c, global.ErrTokenCreate, err)
		return
	}

	// 更新用户的登录信息，包括 IP 地址和上次登录时间
	err = model.UpdateUserLoginInfo(db, userAuth.ID, ipAddress, ipSource)
	if err != nil {
		// 更新登录信息失败，返回数据库操作错误
		ReturnError(c, global.ErrDbOp, err)
		return
	}

	// 登录成功，记录日志
	slog.Info("用户登录成功: " + userAuth.Username)

	// 使用 Gin 的 session 来存储用户的认证信息（UserAuth ID）
	session := sessions.Default(c)
	session.Set(global.CTX_USER_AUTH, userAuth.ID)
	session.Save() // 保存 session

	// 删除 Redis 中的用户离线状态标识
	offlineKey := global.OFFLINE_USER + strconv.Itoa(userAuth.ID)
	rdb.Del(rctx, offlineKey).Result()

	// 返回成功响应，携带用户信息、文章点赞记录、评论点赞记录和 JWT Token
	ReturnSuccess(c, LoginVO{
		UserInfo:       *userInfo,      // 返回用户信息
		ArticleLikeSet: articleLikeSet, // 返回用户的文章点赞记录
		CommentLikeSet: commentLikeSet, // 返回用户的评论点赞记录
		Token:          token,          // 返回生成的 JWT Token
	})
}
