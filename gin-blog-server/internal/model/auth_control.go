package model

import "gorm.io/gorm"

func GetResource(db *gorm.DB, uri, method string) (resource Resource, err error) {
	result := db.Where(&Resource{Url: uri, Method: method}).First(&resource)
	return resource, result.Error
}

func AddResource(db *gorm.DB, name, uri, method string, anonymous bool) (*Resource, error) {
	resource := Resource{
		Name:      name,
		Method:    method,
		Url:       uri,
		Anonymous: anonymous,
	}
	result := db.Save(&resource)
	return &resource, result.Error
}

// GetResourceByRole 根据 role id 获取拥有权限的 resources
func GetResourceByRole(db *gorm.DB, rid int) (resources []Resource, err error) {
	var role Role
	// 从数据库中查询一条 Role 记录，同时预加载该 Role 关联的 Resources 数据，并将查询结果存储到 role 变量中
	result := db.Model(&Role{}).Preload("Resources").Take(&role, rid)
	return role.Resources, result.Error
}

// CheckRoleAuth 根据 role id 判断是否有运行 uri+method 的权限
func CheckRoleAuth(db *gorm.DB, rid int, uri, method string) (bool, error) {
	resources, err := GetResourceByRole(db, rid)
	if err != nil {
		return false, err
	}

	for _, r := range resources {
		if r.Url == uri && r.Method == method {
			return true, nil
		}
	}
	return false, nil
}
