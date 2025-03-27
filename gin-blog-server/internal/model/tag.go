package model

import (
	"gorm.io/gorm"
	"time"
)

type Tag struct {
	Model
	Name     string    `gorm:"unique;type:varchar(20);not null" json:"name"`
	Articles []Article `gorm:"many2many:article_tag;" json:"articles,omitempty"`
}

type TagVO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name         string `json:"name"`
	ArticleCount int    `json:"article_count"`
}

// GetTagList 获取标签列表
func GetTagList(db *gorm.DB, page, size int, keyword string) (list []TagVO, total int64, err error) {
	db = db.Table("tag t").
		Joins("LEFT JOIN article_tag at ON t.id = at.tag_id").
		Select("t.id", "t.name", "COUNT(at.article_id) AS article_count", "t.created_at", "t.updated_at")

	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}

	result := db.Group("t.id").Order("t.updated_at DESC").
		Count(&total).
		Scopes(Paginate(page, size)).
		Find(&list)
	return list, total, result.Error
}

// SaveOrUpdateTag 添加或者修改标签
func SaveOrUpdateTag(db *gorm.DB, id int, name string) (*Tag, error) {
	tag := Tag{
		Model: Model{ID: id},
		Name:  name,
	}

	var result *gorm.DB
	if id > 0 {
		result = db.Updates(&tag)
	} else {
		result = db.Create(&tag)
	}

	return &tag, result.Error
}

// GetTagOption 获取标签选项列表
func GetTagOption(db *gorm.DB) ([]OptionVO, error) {
	list := make([]OptionVO, 0)
	result := db.Model(&Tag{}).Select("id", "name").Find(&list)
	return list, result.Error
}
