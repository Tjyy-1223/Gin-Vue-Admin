package model

import "gorm.io/gorm"

type Page struct {
	Model
	Name  string `gorm:"unique;type:varchar(20)" json:"name"`
	Label string `gorm:"unique;type:varchar(20)" json:"label"`
	Cover string `gorm:"type:varchar(255)" json:"cover"`
}

// GetPageList 获取数据库中的所有 Page 记录
func GetPageList(db *gorm.DB) ([]Page, int64, error) {
	var pages = make([]Page, 0)
	var total int64

	result := db.Model(&Page{}).Count(&total).Find(&pages)
	return pages, total, result.Error
}

// SaveOrUpdatePage 保存或更新一个新的 Page 记录
func SaveOrUpdatePage(db *gorm.DB, id int, name, label, cover string) (*Page, error) {
	page := Page{
		Model: Model{ID: id},
		Name:  name,
		Label: label,
		Cover: cover,
	}

	var result *gorm.DB
	if id > 0 {
		result = db.Updates(&page)
	} else {
		result = db.Create(&page)
	}

	return &page, result.Error
}
