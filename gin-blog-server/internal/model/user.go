package model

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

// UserInfo 代表用户的个人信息
type UserInfo struct {
	Model
	Email    string `json:"email" gorm:"type:varchar(30)"`                    // 用户的邮箱，最大长度30字符，保存用户的电子邮件地址
	Nickname string `json:"nickname" gorm:"unique;type:varchar(30);not null"` // 用户的昵称，唯一，最大长度30字符，不能为空
	Avatar   string `json:"avatar" gorm:"type:varchar(1024);not null"`        // 用户头像，最大长度1024字符，不能为空
	Intro    string `json:"intro" gorm:"type:varchar(255)"`                   // 用户个人简介，最大长度255字符，用于描述用户的个人信息或介绍
	Website  string `json:"website" gorm:"type:varchar(255)"`                 // 用户的个人网站链接，最大长度255字符，用于存储用户的官网、博客等链接
}

type UserInfoVO struct {
	UserInfo
	ArticleLikeSet []string `json:"article_like_set"`
	CommentLikeSet []string `json:"comment_like_set"`
}

// GetUserInfoById 根据用户的 ID 从数据库中查询用户信息
// 参数:
//
//	db - GORM 数据库连接对象，用于执行查询操作
//	id - 用户的 ID，作为查询条件
//
// 返回:
//   - *UserInfo：指向查询到的 UserInfo 对象的指针。如果未找到用户，返回空字段
//   - error：查询过程中发生的错误。如果没有错误，则返回 nil
func GetUserInfoById(db *gorm.DB, id int) (*UserInfo, error) {
	var userInfo UserInfo
	result := db.Model(&userInfo).Where("id", id).First(&userInfo)
	return &userInfo, result.Error
}

// GetUserAuthInfoByName 根据用户名查询用户认证信息
// 参数:
//
//	db - GORM 的数据库实例，用于执行查询
//	name - 用户名，用于模糊匹配查询用户认证信息
//
// 返回:
//   - 如果找到用户，返回用户认证信息和 nil 错误
//   - 如果未找到用户或查询出错，返回 nil 和错误信息
func GetUserAuthInfoByName(db *gorm.DB, name string) (*UserAuth, error) {
	var userAuth UserAuth

	// 使用 GORM 构建查询，进行模糊查询（LIKE）以根据用户名查找用户
	result := db.Model(&userAuth).Where("username LIKE ?", name).First(&userAuth)

	// 检查查询结果是否出错，如果是记录未找到的错误（ErrRecordNotFound），则返回 nil 和错误
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	// 返回找到的用户认证信息和可能的其他错误
	return &userAuth, result.Error
}

// UpdateUserLoginInfo 更新用户登录信息
func UpdateUserLoginInfo(db *gorm.DB, id int, ipAddress, ipSource string) error {
	now := time.Now()
	userAuth := UserAuth{
		IpAddress:     ipAddress,
		IpSource:      ipSource,
		LastLoginTime: &now,
	}
	result := db.Where("id", id).Updates(userAuth)
	return result.Error
}

// UpdateUserInfo 根据 user-id 更新用户信息
func UpdateUserInfo(db *gorm.DB, id int, nickname, avatar, intro, website string) error {
	userInfo := UserInfo{
		Model:    Model{ID: id},
		Nickname: nickname,
		Avatar:   avatar,
		Intro:    intro,
		Website:  website,
	}

	result := db.Select("nickname", "avatar", "intro", "website").Updates(userInfo)
	return result.Error
}

// GetUserList 获取当前存在的用户列表
func GetUserList(db *gorm.DB, page, size int, loginType int8, nickname, username string) (list []UserAuth, total int64, err error) {
	if loginType != 0 {
		db = db.Where("login_type = ?", loginType)
	}
	if username != "" {
		db = db.Where("username LIKE ?", "%"+username+"%")
	}

	result := db.Model(&UserAuth{}).
		Joins("LEFT JOIN user_info ON user_info.id = user_auth.user_info_id").
		Where("user_info.nickname LIKE ?", "%"+nickname+"%").
		Preload("UserInfo").
		Preload("Roles").
		Count(&total).
		Scopes(Paginate(page, size)).
		Find(&list)

	return list, total, result.Error
}

// UpdateUserNicknameAndRole 更新用户昵称及角色信息
func UpdateUserNicknameAndRole(db *gorm.DB, authId int, nickname string, roleIds []int) error {
	userAuth, err := GetUserAuthInfoById(db, authId)
	if err != nil {
		return err
	}

	userInfo := UserInfo{
		Model:    Model{ID: userAuth.UserInfoId},
		Nickname: nickname,
	}

	// 更新用户信息
	result := db.Model(&userInfo).Updates(userInfo)
	if result.Error != nil {
		return result.Error
	}

	// 判断，至少有一个角色，一个用户也可以有多个角色
	if len(roleIds) == 0 {
		return nil
	}

	// 更新用户角色, 清空原本的 user_role 关系, 添加新的关系
	result = db.Where(UserAuthRole{UserAuthId: userAuth.ID}).Delete(UserAuthRole{})
	if result.Error != nil {
		return result.Error
	}

	var userRoles []UserAuthRole
	for _, id := range roleIds {
		userRoles = append(userRoles, UserAuthRole{
			RoleId:     id,
			UserAuthId: userAuth.ID,
		})
	}
	result = db.Create(&userRoles)
	return result.Error
}

// UpdateUserDisable 更新用户的禁用信息
func UpdateUserDisable(db *gorm.DB, id int, isDisable bool) error {
	userAuth := UserAuth{
		Model:     Model{ID: id},
		IsDisable: isDisable,
	}

	result := db.Model(&userAuth).Select("is_disable").Updates(&userAuth)
	return result.Error
}

// UpdateUserPassword 修改当前用户的密码
func UpdateUserPassword(db *gorm.DB, id int, password string) error {
	userAuth := UserAuth{
		Model:    Model{ID: id},
		Password: password,
	}
	result := db.Model(&userAuth).Updates(userAuth)
	return result.Error
}
