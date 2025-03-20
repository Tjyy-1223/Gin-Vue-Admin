package handle

import (
	"encoding/json"
	"gin-blog-server/internal/global"
	"gin-blog-server/internal/model"
	"gin-blog-server/internal/utils"
	"github.com/gin-gonic/gin"
	"sort"
	"strconv"
	"strings"
	"time"
)

type User struct{}

type UpdateCurrentUserReq struct {
	Nickname string `json:"nickname" binding:"required"`
	Avatar   string `json:"avatar"`
	Intro    string `json:"intro"`
	Website  string `json:"website"`
	Email    string `json:"email"`
}

type UserQuery struct {
	PageQuery
	LoginType int8   `form:"login_type"`
	Username  string `form:"username"`
	Nickname  string `form:"nickname"`
}

type UpdateUserReq struct {
	UserAuthId int    `json:"id"`
	Nickname   string `json:"nickname" binding:"required"`
	RoleIds    []int  `json:"role_ids"`
}

type UpdateUserDisableReq struct {
	UserAuthId int  `json:"id"`
	IsDisable  bool `json:"is_disable"`
}

type UpdateCurrentPasswordReq struct {
	NewPassword string `json:"new_password" binding:"required,min=4,max=20"`
	OldPassword string `json:"old_password" binding:"required,min=4,max=20"`
}

// GetInfo 根据 Token 获取用户信息
func (*User) GetInfo(c *gin.Context) {
	rdb := GetRDB(c)

	user, err := CurrentUserAuth(c)
	if err != nil {
		ReturnError(c, global.ErrTokenRuntime, err)
		return
	}

	userInfoVO := model.UserInfoVO{UserInfo: *user.UserInfo}
	userInfoVO.ArticleLikeSet, err = rdb.SMembers(rctx, global.ARTICLE_USER_LIKE_SET+strconv.Itoa(user.ID)).Result()
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}

	userInfoVO.CommentLikeSet, err = rdb.SMembers(rctx, global.COMMENT_USER_LIKE_SET+strconv.Itoa(user.ID)).Result()
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, userInfoVO)
}

// UpdateCurrent 更新当前用户信息, 不需要传 id, 从 Token 中解析出来
func (*User) UpdateCurrent(c *gin.Context) {
	var req UpdateCurrentUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, global.ErrRequest, err)
		return
	}

	auth, _ := CurrentUserAuth(c)
	err := model.UpdateUserInfo(GetDB(c), auth.UserInfoId, req.Nickname, req.Avatar, req.Intro, req.Website)
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, nil)
}

// GetList 获取当前存在的用户列表
func (*User) GetList(c *gin.Context) {
	var query UserQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, global.ErrRequest, err)
		return
	}

	list, count, err := model.GetUserList(GetDB(c), query.Page, query.Size, query.LoginType, query.Nickname, query.Username)
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, PageResult[model.UserAuth]{
		Size:  query.Size,
		Page:  query.Page,
		Total: count,
		List:  list,
	})
}

// Update 更新用户信息：主要是对用户名和用户角色进行更新
func (*User) Update(c *gin.Context) {
	var req UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, global.ErrRequest, err)
		return
	}

	if err := model.UpdateUserNicknameAndRole(GetDB(c), req.UserAuthId, req.Nickname, req.RoleIds); err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, nil)
}

// UpdateDisable 修改用户禁用状态
func (*User) UpdateDisable(c *gin.Context) {
	var req UpdateUserDisableReq

	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, global.ErrRequest, err)
		return
	}

	err := model.UpdateUserDisable(GetDB(c), req.UserAuthId, req.IsDisable)
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
	}

	ReturnSuccess(c, nil)
}

// UpdateCurrentPassword 修改当前用户密码：需要输入旧密码进行验证
func (*User) UpdateCurrentPassword(c *gin.Context) {
	// TODO: 前端密码是明文传输过来的
	var req UpdateCurrentPasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, global.ErrRequest, err)
		return
	}

	// 获取当前用户
	auth, _ := CurrentUserAuth(c)

	// 判断旧密码输入是否正确
	if !utils.BcryptCheck(req.OldPassword, auth.Password) {
		ReturnError(c, global.ErrOldPassword, nil)
		return
	}

	hashPassword, _ := utils.BcryptHash(req.NewPassword)
	err := model.UpdateUserPassword(GetDB(c), auth.ID, hashPassword)
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}

	// TODO: 修改完密码后，强制当前用户下线

	ReturnSuccess(c, nil)
}

// GetOnlineList 查询当前的在线用户，主要是 redis 操作；
func (*User) GetOnlineList(c *gin.Context) {
	keyword := c.Query("keyword")

	rdb := GetRDB(c)

	onlineList := make([]model.UserAuth, 0)
	// 查询redis中的键，模糊查询： "online:*"
	keys := rdb.Keys(rctx, global.ONLINE_USER+"*").Val()

	for _, key := range keys {
		var auth model.UserAuth
		val := rdb.Get(rctx, key).Val() // 从 redis 中获取对应的用户
		json.Unmarshal([]byte(val), &auth)

		// 如果关键词存在，但是该用户的用户名和名称不包含关键词，省略该用户
		if keyword != "" &&
			!strings.Contains(auth.Username, keyword) &&
			!strings.Contains(auth.UserInfo.Nickname, keyword) {
			continue
		}

		onlineList = append(onlineList, auth)
	}

	// 根据上次登录时间进行排序
	sort.Slice(onlineList, func(i, j int) bool {
		return onlineList[i].LastLoginTime.Unix() > onlineList[j].LastLoginTime.Unix()
	})

	ReturnSuccess(c, onlineList)
}

// ForceOffline 强制用户离线
func (*User) ForceOffline(c *gin.Context) {
	id := c.Param("id")
	uid, err := strconv.Atoi(id)
	if err != nil {
		ReturnError(c, global.ErrRequest, err)
		return
	}

	auth, err := CurrentUserAuth(c)
	if err != nil {
		ReturnError(c, global.ErrUserAuth, err)
		return
	}

	// 不能离线自己
	if auth.ID == uid {
		ReturnError(c, global.ErrForceOfflineSelf, nil)
		return
	}

	rdb := GetRDB(c)
	onlineKey := global.ONLINE_USER + strconv.Itoa(uid)
	offlineKey := global.OFFLINE_USER + strconv.Itoa(uid)

	rdb.Del(rctx, onlineKey)
	rdb.Set(rctx, offlineKey, auth, time.Hour)

	ReturnSuccess(c, "强制离线成功")
}
