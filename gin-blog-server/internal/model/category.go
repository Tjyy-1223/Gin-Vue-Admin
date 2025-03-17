package model

type Category struct {
	Model
	Name     string    `gorm:"unique;type:varchar(20);not null" json:"name"`
	Articles []Article `gorm:"foreignKey:CategoryId"`
}
