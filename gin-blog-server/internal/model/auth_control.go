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

// CheckResourceInUse 检查资源是否在使用中
func CheckResourceInUse(db *gorm.DB, id int) (bool, error) {
	var count int64
	result := db.Model(&RoleResource{}).Where("resource_id = ?", id).Count(&count)
	return count > 0, result.Error
}

// GetResourceById 根据 id 获取资源
func GetResourceById(db *gorm.DB, id int) (resource Resource, err error) {
	result := db.First(&resource, id)
	return resource, result.Error
}

// CheckResourceHasChild 检查该资源中是否有子资源
func CheckResourceHasChild(db *gorm.DB, id int) (bool, error) {
	var count int64
	result := db.Model(&Resource{}).Where("parent_id = ?", id).Count(&count)
	return count > 0, result.Error
}

// DeleteResource 删除资源
func DeleteResource(db *gorm.DB, id int) (int, error) {
	result := db.Delete(&Resource{}, id)
	if result.Error != nil {
		return 0, result.Error
	}
	return int(result.RowsAffected), nil
}

// UpdateResourceAnonymous 更新资源匿名状态
func UpdateResourceAnonymous(db *gorm.DB, id int, anonymous bool) error {
	result := db.Model(&Resource{}).Where("id = ?", id).Update("anonymous", anonymous)
	return result.Error
}
