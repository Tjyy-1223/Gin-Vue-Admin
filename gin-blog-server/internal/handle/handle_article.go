package handle

import (
	"gin-blog-server/internal/global"
	"gin-blog-server/internal/model"
	"github.com/gin-gonic/gin"
	"io"
	"log/slog"
	"mime/multipart"
	"strconv"
)

type Article struct{}

// ArticleQuery 文章查询输入请求
// TODO: 添加对标签数组的查询
type ArticleQuery struct {
	PageQuery
	Title      string `form:"title"`
	CategoryId int    `form:"category_id"`
	TagId      int    `form:"tag_id"`
	Type       int    `form:"type"`
	Status     int    `form:"status"`
	IsDelete   *bool  `form:"is_delete"`
}

// AddOrEditArticleReq 新增或者编辑文章的请求
type AddOrEditArticleReq struct {
	ID          int    `json:"id"`
	Title       string `json:"title" binding:"required"`
	Desc        string `json:"desc"`
	Content     string `json:"content" binding:"required"`
	Img         string `json:"img"`
	Type        int    `json:"type" binding:"required,min=1,max=3"`   // 类型: 1-原创 2-转载 3-翻译
	Status      int    `json:"status" binding:"required,min=1,max=3"` // 类型: 1-公开 2-私密 3-评论可见
	IsTop       bool   `json:"is_top"`
	OriginalUrl string `json:"original_url"`

	TagNames     []string `json:"tag_names"`
	CategoryName string   `json:"category_name"`
}

type ArticleVO struct {
	model.Article
	// gorm:"-" 是一个标签，它的主要作用是告诉 GORM 忽略结构体中的某个字段，使其在数据库操作（如创建表、插入数据、查询数据等）中不被考虑。
	// GORM 默认会将结构体名转换为蛇形命名法（snake_case）并以复数形式作为数据库表名。
	LikeCount    int `json:"like_count" gorm:"-"`
	ViewCount    int `json:"view_count" gorm:"-"`
	CommentCount int `json:"comment_count" gorm:"-"`
}

type UpdateArticleTopReq struct {
	ID    int  `json:"id"`
	IsTop bool `json:"is_top"`
}

type SoftDeleteReq struct {
	Ids      []int `json:"ids"`
	IsDelete bool  `json:"is_delete"`
}

// GetList 获取文章列表
func (*Article) GetList(c *gin.Context) {
	var query ArticleQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, global.ErrRequest, err)
		return
	}

	db := GetDB(c)
	rdb := GetRDB(c)

	list, total, err := model.GetArticleList(db, query.Page, query.Size, query.Title, query.IsDelete, query.Status, query.Type, query.CategoryId, query.TagId)
	if err != nil || list == nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}

	// 获取所有文章的点赞数
	likeCountMap := rdb.HGetAll(rctx, global.ARTICLE_LIKE_COUNT).Val()
	// 获取所有文章的观看量，并排序
	viewCountZ := rdb.ZRangeWithScores(rctx, global.ARTICLE_VIEW_COUNT, 0, -1).Val()

	viewCountMap := make(map[int]int)
	for _, article := range viewCountZ {
		id, _ := strconv.Atoi(article.Member.(string))
		viewCountMap[id] = int(article.Score)
	}

	data := make([]ArticleVO, 0)
	for _, article := range list {
		likeCount, _ := strconv.Atoi(likeCountMap[strconv.Itoa(article.ID)])
		data = append(data, ArticleVO{
			Article:   article,
			LikeCount: likeCount,
			ViewCount: viewCountMap[article.ID],
		})
	}

	ReturnSuccess(c, PageResult[ArticleVO]{
		Size:  query.Size,
		Page:  query.Page,
		Total: total,
		List:  data,
	})
}

// SaveOrUpdate 新增或者编辑文章
func (*Article) SaveOrUpdate(c *gin.Context) {
	var req AddOrEditArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, global.ErrRequest, err)
		return
	}

	db := GetDB(c)
	auth, _ := CurrentUserAuth(c)

	if req.Img == "" {
		req.Img = model.GetConfig(db, global.CONFIG_ARTICLE_COVER) // 默认图片
	}

	article := model.Article{
		Model:       model.Model{ID: req.ID},
		Title:       req.Title,
		Desc:        req.Desc,
		Content:     req.Content,
		Img:         req.Img,
		Type:        req.Type,
		Status:      req.Status,
		OriginalUrl: req.OriginalUrl,
		IsTop:       req.IsTop,
		UserId:      auth.UserInfoId,
	}

	err := model.SaveOrUpdateArticle(db, &article, req.CategoryName, req.TagNames)
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, article)
}

// UpdateTop 修改置顶信息
func (*Article) UpdateTop(c *gin.Context) {
	var req UpdateArticleTopReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, global.ErrRequest, err)
		return
	}

	err := model.UpdateArticleTop(GetDB(c), req.ID, req.IsTop)
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, nil)
}

// GetDetail 获取文章详细信息
func (*Article) GetDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ReturnError(c, global.ErrRequest, err)
		return
	}

	article, err := model.GetArticle(GetDB(c), id)
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, article)
}

// UpdateSoftDelete 软删除文章
func (*Article) UpdateSoftDelete(c *gin.Context) {
	var req SoftDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, global.ErrRequest, err)
		return
	}

	rows, err := model.UpdateArticleSoftDelete(GetDB(c), req.Ids, req.IsDelete)
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, rows)
}

// Delete 物理删除文章
func (*Article) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		ReturnError(c, global.ErrRequest, err)
		return
	}

	rows, err := model.DeleteArticle(GetDB(c), ids)
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, rows)
}

// Export 导出文章: 获取导出后的资源链接列表
// TODO: 目前是前端导出
func (*Article) Export(c *gin.Context) {
	ReturnSuccess(c, nil)
}

// Import 导入文章：题目 + 内容
func (*Article) Import(c *gin.Context) {
	db := GetDB(c)
	auth, _ := CurrentUserAuth(c)

	_, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		ReturnError(c, global.ErrFileReceive, err)
		return
	}

	fileName := fileHeader.Filename
	// 获取文章题目
	title := fileName[:len(fileName)-3]
	// 获取文章内容
	content, err := readFromFileHeader(fileHeader)
	if err != nil {
		ReturnError(c, global.ErrFileReceive, err)
		return
	}

	// 获取默认文章封面
	defaultImg := model.GetConfig(db, global.CONFIG_ARTICLE_COVER)
	err = model.ImportArticle(db, auth.ID, title, content, defaultImg, "学习", "后端开发")
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, nil)
}

// 获取文章内容
func readFromFileHeader(file *multipart.FileHeader) (string, error) {
	open, err := file.Open()
	if err != nil {
		slog.Error("文件读取，目标地址错误：", err)
		return "", err
	}
	defer open.Close()

	all, err := io.ReadAll(open)
	if err != nil {
		slog.Error("文件读取失败：", err)
		return "", err
	}

	return string(all), nil
}
