package model

import "gorm.io/gorm"

const (
	TYPE_ARTICLE = iota + 1 // 文章
	TYPE_LINK               // 友链
	TYPE_TALK               // 说说
)

/*
如果评论类型是文章，那么 topic_id 就是文章的 id
如果评论类型是友链，不需要 topic_id
*/

type Comment struct {
	Model
	UserId      int    `json:"user_id"`       // 评论者
	ReplyUserId int    `json:"reply_user_id"` // 被回复者
	TopicId     int    `json:"topic_id"`      // 评论的文章
	ParentId    int    `json:"parent_id"`     // 父评论 被回复的评论
	Content     string `gorm:"type:varchar(500);not null" json:"content"`
	Type        int    `gorm:"type:tinyint(1);not null;comment:评论类型(1.文章 2.友链 3.说说)" json:"type"` // 评论类型 1.文章 2.友链 3.说说
	IsReview    bool   `json:"is_review"`

	// Belongs To
	User      *UserAuth `gorm:"foreignKey:UserId" json:"user"`
	ReplyUser *UserAuth `gorm:"foreignKey:ReplyUserId" json:"reply_user"`
	Article   *Article  `gorm:"foreignKey:TopicId" json:"article"`
}

type CommentVO struct {
	Comment
	LikeCount  int         `json:"like_count" gorm:"-"`
	ReplyCount int         `json:"reply_count" gorm:"-"`
	ReplyList  []CommentVO `json:"reply_list" gorm:"-"`
}

// GetArticleCommentCount 获取某篇文章的评论数
func GetArticleCommentCount(db *gorm.DB, articleId int) (count int64, err error) {
	result := db.Model(&Comment{}).
		Where("topic_id = ? AND type = 1 AND is_review = 1", articleId).
		Count(&count)
	return count, result.Error
}
