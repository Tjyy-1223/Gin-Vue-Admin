package model

import (
	"encoding/json"
	"gin-blog-server/internal/utils"
	"gorm.io/gorm"
	"log/slog"
	"strconv"
	"time"
)

// UserAuth 代表用户认证信息
type UserAuth struct {
	Model
	Username      string     `gorm:"unique;type:varchar(50)" json:"username"`           // 用户名，唯一，最大长度为50
	Password      string     `gorm:"type:varchar(100)" json:"-"`                        // 密码，最大长度为100，不会被JSON序列化
	LoginType     int        `gorm:"type:tinyint(1);comment:登录类型" json:"login_type"`    // 登录类型，Tinyint 类型，表示不同的登录方式（例如：用户名/密码、第三方登录等）
	IpAddress     string     `gorm:"type:varchar(20);comment:登录IP地址" json:"ip_address"` // 登录IP地址，最大长度为20
	IpSource      string     `gorm:"type:varchar(50);comment:IP来源" json:"ip_source"`    // IP来源，最大长度为50
	LastLoginTime *time.Time `json:"last_login_time"`                                   // 上次登录时间，类型为指针，以便为null
	IsDisable     bool       `json:"is_disable"`                                        // 是否禁用，布尔值，表示该用户是否被禁用
	IsSuper       bool       `json:"is_super"`                                          // 是否超级管理员，布尔值，超级管理员只能由后台设置
	UserInfoId    int        `json:"user_info_id"`                                      // 关联的用户信息表ID
	UserInfo      *UserInfo  `json:"info"`                                              // 关联的用户信息
	Roles         []*Role    `json:"roles" gorm:"many2many:user_auth_role"`             // 用户角色，表示与角色的多对多关系
}

func (u *UserAuth) MarshalBinary() (data []byte, err error) {
	return json.Marshal(u)
}

// Role 代表系统中的角色
type Role struct {
	Model
	Name      string `gorm:"unique" json:"name"`  // 角色名称，唯一
	Label     string `gorm:"unique" json:"label"` // 角色标签，唯一
	IsDisable bool   `json:"is_disable"`          // 是否禁用该角色，布尔值

	Resources []Resource `json:"resources" gorm:"many2many:role_resource"` // 角色拥有的资源，表示与资源的多对多关系
	Menus     []Menu     `json:"menus" gorm:"many2many:role_menu"`         // 角色拥有的菜单，表示与菜单的多对多关系
	Users     []UserAuth `json:"users" gorm:"many2many:user_auth_role"`    // 角色关联的用户，表示与用户的多对多关系
}

// Resource 代表系统中的资源
type Resource struct {
	Model
	Name      string `gorm:"unique;type:varchar(50)" json:"name"`    // 资源名称，最大长度为50，唯一
	ParentId  int    `json:"parent_id"`                              // 父资源ID，表示资源之间的层级关系
	Url       string `gorm:"type:varchar(255)" json:"url"`           // 资源的URL地址
	Method    string `gorm:"type:varchar(10)" json:"request_method"` // 请求方法，例如 GET, POST, PUT, DELETE 等
	Anonymous bool   `json:"is_anonymous"`                           // 是否是匿名访问的资源，布尔值，表示该资源是否不需要登录

	Roles []*Role `json:"roles" gorm:"many2many:role_resource"` // 资源关联的角色，表示与角色的多对多关系
}

/*
菜单设计:

目录: catalogue === true
  - 如果是目录，作为单独项，不展开子菜单（例如 "首页", "个人中心"）
  - 如果不是目录，且 parent_id 为 0，则为一级菜单，可以展开子菜单（例如 "文章管理" 下有 "文章列表", "文章分类", "文章标签" 等子菜单）
  - 如果不是目录，且 parent_id 不为 0，则为二级菜单

隐藏: hidden
  - 隐藏则不显示在菜单栏中

外链: external, external_link
  - 如果是外链，如果设置为外链，则点击后会在新窗口打开
*/

// Menu 代表系统中的菜单
type Menu struct {
	Model
	ParentId     int    `json:"parent_id"`                                                  // 父菜单ID，用于标识菜单的层级关系
	Name         string `gorm:"uniqueIndex:idx_name_and_path;type:varchar(20)" json:"name"` // 菜单名称，唯一索引
	Path         string `gorm:"uniqueIndex:idx_name_and_path;type:varchar(50)" json:"path"` // 菜单的路由地址，唯一索引
	Component    string `gorm:"type:varchar(50)" json:"component"`                          // 菜单组件路径
	Icon         string `gorm:"type:varchar(50)" json:"icon"`                               // 菜单图标
	OrderNum     int8   `json:"order_num"`                                                  // 菜单排序
	Redirect     string `gorm:"type:varchar(50)" json:"redirect"`                           // 菜单重定向地址
	Catalogue    bool   `json:"is_catalogue"`                                               // 是否为目录，目录项不展开子菜单
	Hidden       bool   `json:"is_hidden"`                                                  // 是否隐藏该菜单，隐藏则不显示在菜单栏
	KeepAlive    bool   `json:"keep_alive"`                                                 // 是否缓存该菜单
	External     bool   `json:"is_external"`                                                // 是否为外链菜单
	ExternalLink string `gorm:"type:varchar(255)" json:"external_link"`                     // 外链地址

	Roles []*Role `json:"roles" gorm:"many2many:role_menu"` // 菜单关联的角色，表示与角色的多对多关系
}

type UserAuthRole struct {
	UserAuthId int `gorm:"primaryKey;uniqueIndex:idx_user_auth_role"`
	RoleId     int `gorm:"primaryKey;uniqueIndex:idx_user_auth_role"`
}

type RoleResource struct {
	RoleId     int `json:"-" gorm:"primaryKey;uniqueIndex:idx_role_resource"`
	ResourceId int `json:"-" gorm:"primaryKey;uniqueIndex:idx_role_resource"`
}

type RoleMenu struct {
	RoleId int `json:"-" gorm:"primaryKey;uniqueIndex:idx_role_menu"`
	MenuId int `json:"-" gorm:"primaryKey;uniqueIndex:idx_role_menu"`
}

// GetRoleIdsByUserId 根据用户的 UserAuthId 查询该用户拥有的角色 ID 列表
// 参数:
//
//	db - GORM 数据库连接对象，用于执行查询操作
//	userAuthId - 用户的 UserAuthId，作为查询条件
//
// 返回:
//   - ids：用户所拥有的角色 ID 列表
//   - err：查询过程中发生的错误。如果没有错误，则返回 nil
func GetRoleIdsByUserId(db *gorm.DB, userAuthId int) (ids []int, err error) {
	// 使用 GORM 查询方式，获取用户角色表（UserAuthRole）中与 userAuthId 对应的所有 role_id
	// Model(&UserAuthRole{UserAuthId: userAuthId})：指定查询的目标表为 UserAuthRole 表，并且条件是 UserAuthId 等于传入的 userAuthId
	// Pluck("role_id", &ids)：查询所有的 role_id，并将结果存入 ids 切片
	result := db.Model(&UserAuthRole{UserAuthId: userAuthId}).Pluck("role_id", &ids)

	// 返回查询结果：ids 包含所有角色 ID，result.Error 包含可能发生的错误
	return ids, result.Error
}

// GetUserAuthInfoById 通过 id 获取对应的 UserAuth
func GetUserAuthInfoById(db *gorm.DB, id int) (*UserAuth, error) {
	var userAuth = UserAuth{
		Model: Model{ID: id},
	}
	// 查询 userAuth 表中的第一条记录，并且通过 Preload 预加载关联的 Roles 和 UserInfo 数据
	result := db.Model(&userAuth).
		// Preload("Roles") 会将 userAuth 关联的 Roles 数据一并查询出来
		Preload("Roles").
		// Preload("UserInfo") 会将 userAuth 关联的 UserInfo 数据一并查询出来
		Preload("UserInfo").
		// First(&userAuth) 查询 userAuth 表中第一条符合条件的记录，并将结果存入 userAuth 变量
		First(&userAuth)

	return &userAuth, result.Error
}

// CreateNewUser 传入用户名和密码注册新用户
func CreateNewUser(db *gorm.DB, username, password string) (*UserAuth, *UserInfo, *UserAuthRole, error) {
	// 创建 userinfo
	num, err := Count(db, &UserInfo{})
	if err != nil {
		slog.Info(err.Error())
	}

	number := strconv.Itoa(num)
	userinfo := &UserInfo{
		Email:    username,
		Nickname: "游客" + number,
		Avatar:   "https://www.bing.com/rp/ar_9isCNU2Q-VG1yEDDHnx8HAFQ.png",
		Intro:    "我是这个程序的第" + number + "个用户",
	}
	// 在 user 表中创建对应的记录
	result := db.Create(&userinfo)
	if result.Error != nil {
		return nil, nil, nil, result.Error
	}

	// 创建 userAuth
	pass, _ := utils.BcryptHash(password)
	userAuth := &UserAuth{
		Username:   username,
		Password:   pass,
		UserInfoId: userinfo.ID,
	}
	result = db.Create(&userAuth)
	if result.Error != nil {
		return nil, nil, nil, result.Error
	}

	// 创建 user - auth 关联表
	userRole := &UserAuthRole{
		UserAuthId: userAuth.ID,
		RoleId:     2, // 默认身份为游客
	}

	result = db.Create(&userRole)
	if result.Error != nil {
		return nil, nil, nil, result.Error
	}
	return userAuth, userinfo, userRole, result.Error
}

// GetAllMenuList 获取所有菜单列表（超级管理员用）
func GetAllMenuList(db *gorm.DB) (menu []Menu, err error) {
	result := db.Find(&menu)
	return menu, result.Error
}

// GetMenuListByUserId 根据 user id 获取菜单列表
func GetMenuListByUserId(db *gorm.DB, id int) (menus []Menu, err error) {
	var userAuth UserAuth
	result := db.Where(&UserAuth{Model: Model{ID: id}}).
		Preload("Roles").Preload("Roles.Menus").First(&userAuth)

	if result.Error != nil {
		return nil, result.Error
	}

	set := make(map[int]Menu)
	for _, role := range userAuth.Roles {
		for _, menu := range role.Menus {
			set[menu.ID] = menu
		}
	}

	for _, menu := range set {
		menus = append(menus, menu)
	}

	return menus, nil
}

// GetMenuList 根据 keyword 从数据库中获取 menu 菜单
func GetMenuList(db *gorm.DB, keyword string) (List []Menu, total int64, err error) {
	db = db.Model(&Menu{})
	if keyword != "" {
		db = db.Where("name like ?", "%"+keyword+"%")
	}
	result := db.Count(&total).Find(&List)
	return List, total, result.Error
}

// SaveOrUpdateMenu 新增和编辑菜单
func SaveOrUpdateMenu(db *gorm.DB, menu *Menu) error {
	var result *gorm.DB

	if menu.ID > 0 { // 编辑菜单
		result = db.Model(menu).
			Select("name", "path", "component", "icon", "redirect", "parent_id", "order_num", "catalogue", "hidden", "keep_alive", "external").
			Updates(menu)
	} else { // 新建菜单
		result = db.Create(menu)
	}

	return result.Error
}

// GetRoleOption 获取角色配置
func GetRoleOption(db *gorm.DB) (list []OptionVO, err error) {
	result := db.Model(&Role{}).Select("id", "name").Find(&list)
	if result.Error != nil {
		return nil, result.Error
	}
	return list, nil
}

// CheckMenuInUse 判断当前菜单是否正在使用中, 传入 menuId
func CheckMenuInUse(db *gorm.DB, id int) (bool, error) {
	var count int64
	result := db.Model(&RoleMenu{}).Where("menu_id = ?", id).Count(&count)
	return count > 0, result.Error
}

// GetMenuById 根据 menuId 获取对应的菜单记录
func GetMenuById(db *gorm.DB, id int) (menu *Menu, err error) {
	result := db.First(&menu, id)
	return menu, result.Error
}

// CheckMenuHasChild 根据 menuId 判断获取的菜单是否有子菜单
func CheckMenuHasChild(db *gorm.DB, id int) (bool, error) {
	var count int64
	result := db.Model(&Menu{}).Where("parent_id = ?", id).Count(&count)
	return count > 0, result.Error
}

// DeleteMenu 根据 menuId 删除对应的菜单
func DeleteMenu(db *gorm.DB, id int) error {
	result := db.Delete(&Menu{}, id)
	return result.Error
}
