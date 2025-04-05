package model

import (
	"gorm.io/gorm"
	"time"
)

// MakeMigrate 迁移数据表，在没有数据表结构变更时候，建议注释不执行
// 只支持创建表、增加表中没有的字段和索引
// 为了保护数据，并不支持改变已有的字段类型或删除未被使用的字段
func MakeMigrate(db *gorm.DB) error {
	// 设置表关联
	// 用于显式地配置一个多对多关系，其中 UserAuth 和 Role 通过一个关联表 UserAuthRole 进行关联。
	db.SetupJoinTable(&UserAuth{}, "Roles", &UserAuthRole{})
	db.SetupJoinTable(&Role{}, "Menus", &RoleMenu{})
	db.SetupJoinTable(&Role{}, "Resources", &RoleResource{})
	db.SetupJoinTable(&Role{}, "Users", &UserAuthRole{})

	return db.AutoMigrate(
		&Article{},      // 文章
		&Category{},     // 分类
		&Tag{},          // 标签
		&Comment{},      // 评论
		&Message{},      // 消息
		&FriendLink{},   // 友链
		&Page{},         // 页面
		&Config{},       // 网站设置
		&OperationLog{}, // 操作日志
		&UserInfo{},     // 用户信息

		&UserAuth{},     // 用户验证
		&Role{},         // 角色
		&Menu{},         // 菜单
		&Resource{},     // 资源（接口）
		&UserAuthRole{}, // 用户-角色 关联
	)
}

type Model struct {
	ID        int       `json:"id" gorm:"primary_key;auto_increment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OptionVO struct {
	ID   int    `json:"value"`
	Name string `json:"name"`
}

// Count 根据 where 条件统计数据
func Count[T any](db *gorm.DB, data *T, where ...any) (int, error) {
	var total int64
	db = db.Model(data)
	if len(where) > 0 {
		db = db.Where(where[0], where[1:]...)
	}
	result := db.Count(&total)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(total), nil
}

// Paginate 分页函数
func Paginate(page, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		switch {
		case size > 100:
			size = 100
		case size <= 10:
			size = 10
		}

		offset := (page - 1) * size
		return db.Offset(offset).Limit(size)
	}
}

// List 数据列表
func List[T any](db *gorm.DB, data T, slt, order, query string, args ...any) (T, error) {
	db = db.Model(data).Select(slt).Order(order)
	if query != "" {
		db = db.Where(query, args...)
	}

	// 数据存储在 data 中
	result := db.Find(&data)
	if result.Error != nil {
		return data, result.Error
	}
	return data, nil
}
