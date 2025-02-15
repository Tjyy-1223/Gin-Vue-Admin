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
	"net/http"
	"strconv"
	"time"
)

type UserAuth struct{}

// LoginReq 发送的登陆请求
type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterReq struct {
	Username string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=4,max=20"`
}

// LoginVO 登陆信息返回给前端的数据
type LoginVO struct {
	model.UserInfo

	// 点赞 Set： 用于记录用户点赞过的文章，评论
	ArticleLikeSet []string `json:"article_like_set"`
	CommentLikeSet []string `json:"comment_like_set"`
	Token          string   `json:"token"`
}

// Login 完成登陆操作
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

// Register 完成注册功能
// 首先检查用户名是否存在，避免重复注册；其次吧用户输入的信息加密保存在 redis 中，等待验证
// 在以下情况下会出错：1-用户邮箱已经注册过；2-用户邮箱无效等原因导致邮件发送失败
// @Summary 注册
// @Description 注册
// @Tags UserAuth
// @Param form body RegisterReq true "注册"
// @Accept json
// @Produce json
// @Success 0 {object} string
// @Router /register [post]
func (*UserAuth) Register(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, global.ErrRequest, err)
		return
	}
	// 格式化用户名
	req.Username = utils.Format(req.Username)

	// 检查用户名是否存在，避免重复注册
	auth, err := model.GetUserAuthInfoByName(GetDB(c), req.Username)
	if err != nil {
		var flag = false
		if errors.Is(err, gorm.ErrRecordNotFound) {
			flag = true
		}
		if !flag {
			ReturnError(c, global.ErrDbOp, err)
			return
		}
	}

	// 用户名重复，不能正常进行注册
	if auth != nil {
		ReturnError(c, global.ErrUserExist, err)
		return
	}

	// 通过邮箱验证后才可以完成注册
	info := utils.GenEmailVerificationInfo(req.Username, req.Password)
	err = SetMailInfo(GetRDB(c), info, 15*time.Minute)
	if err != nil {
		ReturnError(c, global.ErrRedisOp, err)
		return
	}

	EmailData := utils.GetEmailData(req.Username, info)
	err = utils.SendEmail(req.Username, EmailData)
	if err != nil {
		ReturnError(c, global.ErrSendEmail, err)
		return
	}

	ReturnSuccess(c, nil)
}

// VerifyCode 邮箱验证
// 当用户点击邮箱中的链接时，会携带info（加密后的帐号密码）向这个接口发送请求。
// Verify会检查info是否存在redis中，若存在则认证成功，完成注册
// 会在以下方面出错： 1. 发送信息中没有info 2. info不存在redis中(已过期) 3. 创造新用户失败（数据库操作失败）
func (*UserAuth) VerifyCode(c *gin.Context) {
	var code string
	if code = c.Query("info"); code == "" {
		returnErrorPage(c)
		return
	}
	// 验证是否在 redis 数据库中
	ifExist, err := GetMailInfo(GetRDB(c), code)
	if err != nil {
		returnErrorPage(c)
		return
	}
	if !ifExist {
		returnErrorPage(c)
		return
	}

	err = DeleteMailInfo(GetRDB(c), code)
	if err != nil {
		returnErrorPage(c)
		return
	}

	// 从 code 中解析出来 用户名 和 密码
	username, password, err := utils.ParseEmailVerificationInfo(code)
	if err != nil {
		returnErrorPage(c)
		return
	}

	// 注册用户
	_, _, _, err = model.CreateNewUser(GetDB(c), username, password)
	if err != nil {
		returnErrorPage(c)
		return
	}

	// 注册成功，返回成功页面
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
        <!DOCTYPE html>
        <html lang="zh-CN">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>注册成功</title>
            <style>
                body {
                    font-family: Arial, sans-serif;
                    background-color: #f4f4f4;
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    height: 100vh;
                    margin: 0;
                }
                .container {
                    background-color: #fff;
                    padding: 20px;
                    border-radius: 8px;
                    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
                    text-align: center;
                }
                h1 {
                    color: #5cb85c;
                }
                p {
                    color: #333;
                }
            </style>
        </head>
        <body>
            <div class="container">
                <h1>注册成功</h1>
                <p>恭喜您，注册成功！</p>
            </div>
        </body>
        </html>
    `))
}

// c.Data 可以用来直接返回原始字节数据，而不是使用 Gin 中的 c.JSON、c.String 等方法。它特别适合于返回 非结构化数据，例如 HTML 页面、文本或文件。
func returnErrorPage(c *gin.Context) {
	c.Data(http.StatusInternalServerError, "text/html; charset=utf-8", []byte(`
        <!DOCTYPE html>
        <html lang="zh-CN">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>注册失败</title>
            <style>
                body {
                    font-family: Arial, sans-serif;
                    background-color: #f4f4f4;
                    display: flex;
                    justify-content: center;
                    align-items: center;
                    height: 100vh;
                    margin: 0;
                }
                .container {
                    background-color: #fff;
                    padding: 20px;
                    border-radius: 8px;
                    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
                    text-align: center;
                }
                h1 {
                    color: #d9534f;
                }
                p {
                    color: #333;
                }
            </style>
        </head>
        <body>
            <div class="container">
                <h1>注册失败</h1>
                <p>请重试。</p>
            </div>
        </body>
        </html>
    `))
}

// Logout 退出登录
// TODO：退出登录之后，应该将 jwt Token 失效就可以，session 的操作有些多余
// @Summary 退出登录
// @Description 退出登录
// @Tags UserAuth
// @Accept json
// @Produce json
// @Success 0 {object} string
// @Router /logout [get]
func (*UserAuth) Logout(c *gin.Context) {
	// 防止其他请求设置干扰
	c.Set(global.CTX_USER_AUTH, nil)

	// 已经退出登录
	auth, _ := CurrentUserAuth(c)
	if auth == nil {
		ReturnSuccess(c, nil)
		return
	}

	session := sessions.Default(c)
	session.Delete(global.CTX_USER_AUTH)
	session.Save()

	// 删除 Redis 中的在线状态
	rdb := GetRDB(c)
	onlineKey := global.ONLINE_USER + strconv.Itoa(auth.ID)
	rdb.Del(rctx, onlineKey)
	ReturnSuccess(c, nil)
}
