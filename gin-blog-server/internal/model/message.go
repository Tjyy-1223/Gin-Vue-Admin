package model

import "gorm.io/gorm"

type Message struct {
	Model
	Nickname  string `gorm:"type:varchar(50);comment:昵称" json:"nickname"`
	Avatar    string `gorm:"type:varchar(255);comment:头像地址" json:"avatar"`
	Content   string `gorm:"type:varchar(255);comment:留言内容" json:"content"`
	IpAddress string `gorm:"type:varchar(50);comment:IP 地址" json:"ipAddress"`
	IpSource  string `gorm:"type:varchar(255);comment:IP 来源" json:"ipSource"`
	Speed     int    `gorm:"type:tinyint(1);comment:弹幕速度" json:"speed"`
	IsReview  bool   `json:"is_review"`
}

func UpdateMessageReview(db *gorm.DB, ids []int, isReview bool) (int64, error) {
	result := db.Model(&Message{}).Where("id in ?", ids).Update("is_review", isReview)
	return result.RowsAffected, result.Error
}
