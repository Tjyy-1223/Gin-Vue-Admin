package model

import "gorm.io/gorm"

const (
	STATUS_PUBLIC = iota + 1 // 公开
	STATUS_SECRET            // 私密
	STATUS_DRAFT             // 草稿
)

const (
	TYPE_ORIGINAL  = iota + 1 // 原创
	TYPE_REPRINT              // 转载
	TYPE_TRANSLATE            // 翻译
)

// Article
// belongTo: 一个文章 属于 一个分类
// belongTo: 一个文章 属于 一个用户
// many2many: 一个文章 可以拥有 多个标签, 多个文章 可以使用 一个标签
type Article struct {
	Model

	Title       string `gorm:"type:varchar(100);not null" json:"title"`
	Desc        string `json:"desc"`
	Content     string `json:"content"`
	Img         string `json:"img"`
	Type        int    `gorm:"type:tinyint;comment:类型(1-原创 2-转载 3-翻译)" json:"type"`
	Status      int    `gorm:"type:tinyint;comment:状态(1-公开 2-私密)" json:"status"`
	IsTop       bool   `json:"is_top"`
	IsDelete    bool   `json:"is_delete"`
	OriginalUrl string `json:"original_url"`

	CategoryId int `json:"category_id"`
	UserId     int `json:"-"` // user_auth_id

	Tags     []*Tag    `gorm:"many2many:article_tag;joinForeignKey:article_id" json:"tags"`
	Category *Category `gorm:"foreignkey:CategoryId" json:"category"`
	User     *UserAuth `gorm:"foreignkey:UserId" json:"user"`
}

type ArticleTag struct {
	ArticleId int
	TagId     int
}

// GetArticleList 获取文章列表
func GetArticleList(db *gorm.DB, page, size int, title string, isDelete *bool, status, typ, categoryId, tagId int) (list []Article, total int64, err error) {
	db = db.Model(Article{})

	if title != "" {
		db = db.Where("title LIKE ?", "%"+title+"%")
	}
	if isDelete != nil {
		db = db.Where("is_delete", *isDelete)
	}
	if status != 0 {
		db = db.Where("status", status)
	}
	if categoryId != 0 {
		db = db.Where("category_id", categoryId)
	}
	if typ != 0 {
		db = db.Where("type", typ)
	}

	db = db.Preload("Category").Preload("Tags").
		Joins("LEFT JOIN article_tag ON article_tag.article_id = article.id").
		Group("id") // 去重

	if tagId != 0 {
		db = db.Where("tag_id = ?", tagId)
	}

	result := db.Count(&total).
		Scopes(Paginate(page, size)).
		Order("is_top DESC, article.id DESC").
		Find(&list)
	return list, total, result.Error
}

// SaveOrUpdateArticle 新增/编辑文章, 同时根据 分类名称, 标签名称 维护关联表
func SaveOrUpdateArticle(db *gorm.DB, article *Article, categoryName string, tagNames []string) error {
	// 由于要操作多个数据库表，所以要开启事务
	return db.Transaction(func(tx *gorm.DB) error {
		// 分类不存在则创建
		category := Category{Name: categoryName}
		result := db.Model(&Category{}).Where("name", categoryName).FirstOrCreate(&category)
		if result.Error != nil {
			return result.Error
		}

		// 设置文章的分类
		article.CategoryId = category.ID

		// 先 添加/更新 文章, 获取到其 ID
		if article.ID == 0 {
			result = db.Create(article)
		} else {
			result = db.Model(article).Where("id", article.ID).Updates(article)
		}
		if result.Error != nil {
			return result.Error
		}

		// 清空文章标签关联
		result = db.Delete(&ArticleTag{}, "article_id", article.ID)
		if result.Error != nil {
			return result.Error
		}

		// 并重新开始新建 文章-标签 关联
		var articleTags []ArticleTag
		for _, tagName := range tagNames {
			// 标签不存在则创建
			tag := Tag{Name: tagName}
			result := db.Model(&Tag{}).Where("name", tagName).FirstOrCreate(&tag)
			if result.Error != nil {
				return result.Error
			}
			articleTags = append(articleTags, ArticleTag{
				ArticleId: article.ID,
				TagId:     tag.ID,
			})
		}

		result = db.Create(&articleTags)
		return result.Error
	})
}

// UpdateArticleTop 修改置顶信息
func UpdateArticleTop(db *gorm.DB, id int, isTop bool) error {
	result := db.Model(&Article{Model: Model{ID: id}}).Update("is_top", isTop)
	return result.Error
}

// GetArticle 文章的详细信息
func GetArticle(db *gorm.DB, id int) (data *Article, err error) {
	result := db.Preload("Category").Preload("Tags").
		Where(Article{Model: Model{ID: id}}).
		First(&data)
	return data, result.Error
}

// UpdateArticleSoftDelete 软删除文章（修改）
func UpdateArticleSoftDelete(db *gorm.DB, ids []int, isDelete bool) (int64, error) {
	result := db.Model(Article{}).
		Where("id IN ?", ids).
		Update("is_delete", isDelete)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// DeleteArticle 物理删除文章
func DeleteArticle(db *gorm.DB, ids []int) (int64, error) {
	// 删除 [文章-标签] 关联
	result := db.Where("article_id IN ?", ids).Delete(&ArticleTag{})
	if result.Error != nil {
		return 0, result.Error
	}

	// 删除 [文章]
	result = db.Where("id IN ?", ids).Delete(&Article{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

// ImportArticle 导入文章：题目 + 内容
// TODO：如果原来的文件中有图片的话，直接上传图片会由于链接错误无法显示？如何解决图片的自动化上传云+正常显示
func ImportArticle(db *gorm.DB, userAuthId int, title string, content string, img string, categoryName string, tagName string) error {
	article := Article{
		Title:   title,
		Content: content,
		Img:     img,
		Status:  STATUS_DRAFT,
		Type:    TYPE_ORIGINAL,
		UserId:  userAuthId,
	}

	// 生成对应的分类
	category := Category{Name: categoryName}
	result := db.Model(&Category{}).Where("name", categoryName).FirstOrCreate(&category)
	if result.Error != nil {
		return result.Error
	}
	article.CategoryId = category.ID

	// 插入文章
	result = db.Create(&article)
	if result.Error != nil {
		return result.Error
	}

	// 生成对应的文章-标签记录
	var articleTag ArticleTag
	tag := Tag{Name: tagName}
	result = db.Model(&Tag{}).Where("name", tagName).FirstOrCreate(&tag)
	if result.Error != nil {
		return result.Error
	}

	// 插入 文章-标签
	articleTag.ArticleId = article.ID
	articleTag.TagId = tag.ID
	result = db.Create(&articleTag)

	return result.Error
}
