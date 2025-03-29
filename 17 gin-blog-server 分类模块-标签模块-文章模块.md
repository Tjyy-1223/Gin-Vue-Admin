# 17 gin-blog-server 分类模块-标签模块-文章模块

## 1 分类模块

该模块提供了对分类信息进行管理的接口，涵盖分类的查询、新增、编辑和删除等操作，同时还能获取分类选项列表，方便在某些场景下进行选择。

1. **分类列表获取**：通过 `GET /category/list` 接口，可获取所有分类的列表信息，用于展示分类列表页面等场景。
2. **分类新增 / 编辑**：使用 `POST /category` 接口，既可以新增一个分类，也可以对已有的分类信息进行编辑修改。
3. **分类删除**：调用 `DELETE /category` 接口，能够删除指定的分类信息。
4. **分类选项列表获取**：`GET /category/option` 接口用于获取分类的选项列表，可用于下拉选择框等交互场景。

```go
// 分类模块
category := auth.Group("/category")
{
   category.GET("/list", categoryAPI.GetList)     // 分类列表
   category.POST("", categoryAPI.SaveOrUpdate)    // 新增/编辑分类
   category.DELETE("", categoryAPI.Delete)        // 删除分类
   category.GET("/option", categoryAPI.GetOption) // 分类选项列表
}
```



### 1.1 分类列表获取 /category/list

manager.go

```go
// 分类模块
category := auth.Group("/category")
{
  category.GET("/list", categoryAPI.GetList)
}
```

handle/handle_catogory.go

```go
// GetList 获取分类列表
// @Summary 获取分类列表
// @Description 根据条件查询获取分类列表
// @Tags Category
// @Param page_size query int false "当前页数"
// @Param page_num query int false "每页条数"
// @Param keyword query string false "搜索关键字"
// @Accept json
// @Produce json
// @Success 0 {object} Response[PageResult[model.CategoryVO]]
// @Security ApiKeyAuth
// @Router /category/list [get]
func (*Category) GetList(c *gin.Context) {
	var query PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, global.ErrRequest, err)
		return
	}

	data, total, err := model.GetCategoryList(GetDB(c), query.Page, query.Size, query.Keyword)
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, PageResult[model.CategoryVO]{
		Total: total,
		List:  data,
		Size:  query.Size,
		Page:  query.Page,
	})
}
```

model/catogory.go

```go
// GetCategoryList 获取分类列表
func GetCategoryList(db *gorm.DB, num, size int, keyword string) ([]CategoryVO, int64, error) {
	var list = make([]CategoryVO, 0)
	var total int64

	db = db.Table("category c").
		Joins("LEFT JOIN article a ON c.id = a.category_id AND a.is_delete = 0 AND a.status = 1").
		Select("c.id", "c.name", "COUNT(a.id) as article_count", "c.created_at", "c.updated_at")

	if keyword != "" {
		db = db.Where("name LIKE ?", "%"+keyword+"%")
	}

	result := db.Group("c.id").
		Order("c.updated_at DESC").
		Scopes(Paginate(num, size)).
		Find(&list)

	return list, total, result.Error
}
```

 与之前的过程类似，对应的请求和响应如下：

![image-20250326164154928](./assets/image-20250326164154928.png)

![image-20250326164209294](./assets/image-20250326164209294.png)

我们可以在测试完 1.2 分类新增之后来重新查看分类列表，会显示出新增加的分类。



### 1.2 分类新增/编辑 /category POST

manager.go

```go
category.POST("", categoryAPI.SaveOrUpdate) // 新增/编辑分类
```

handle/handle_catogory.go

```go
// SaveOrUpdate 添加或修改分类
// @Summary 添加或修改分类
// @Description 添加或修改分类
// @Tags Category
// @Param form body AddOrEditCategoryReq true "添加或修改分类"
// @Accept json
// @Produce json
// @Success 0 {object} Response[model.Category]
// @Security ApiKeyAuth
// @Router /category [post]
func (*Category) SaveOrUpdate(c *gin.Context) {
	var req AddOrEditCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, global.ErrRequest, err)
		return
	}

	category, err := model.SaveOrUpdateCategory(GetDB(c), req.ID, req.Name)
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, category)
}

```

model/catogory.go

```go
// SaveOrUpdateCategory 添加或修改分类
func SaveOrUpdateCategory(db *gorm.DB, id int, name string) (*Category, error) {
	category := Category{
		Model: Model{ID: id},
		Name:  name,
	}

	var result *gorm.DB
	if id > 0 {
		result = db.Updates(&category)
	} else {
		result = db.Create(&category)
	}

	return &category, result.Error
}

```

对应的请求和响应如下：

![image-20250326164802661](./assets/image-20250326164802661.png)

![image-20250326164835561](./assets/image-20250326164835561.png)

![image-20250326164910873](./assets/image-20250326164910873.png)



### 1.3 分类删除 /category DELETE

manager.go

```go
category.DELETE("", categoryAPI.Delete)        // 删除分类
```

handle/handle_catogory.go

```go
// Delete 删除分类（批量）
// @Summary 删除分类（批量）
// @Description 根据 ID 数组删除分类
// @Tags Category
// @Param ids body []int true "分类 ID 数组"
// @Accept json
// @Produce json
// @Success 0 {object} Response[int]
// @Security ApiKeyAuth
// @Router /category [delete]
func (*Category) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		ReturnError(c, global.ErrRequest, err)
		return
	}

	db := GetDB(c)

	// 检查分类下面是否存在文章
	count, err := model.Count(db, &model.Article{}, "category_id in ?", ids)
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}

	if count > 0 {
		ReturnError(c, global.ErrCateHasArt, nil)
		return
	}

	rows, err := model.DeleteCategory(db, ids)
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, rows)
}
```

model/catogory.go

```go
// DeleteCategory 删除分类（批量）
func DeleteCategory(db *gorm.DB, ids []int) (int64, error) {
	result := db.Where("id IN ?", ids).Delete(Category{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
```

对应的请求和响应如下：

![image-20250326165651678](./assets/image-20250326165651678.png)

![image-20250326165740735](./assets/image-20250326165740735.png)

![image-20250326165811195](./assets/image-20250326165811195.png)



### 1.4 分类选项列表获取 / category/option

manager.go

```go
category.GET("/option", categoryAPI.GetOption) // 分类选项列表
```

handle/handle_catogory.go

```go
// GetOption 获取分类选项列表
// @Summary 获取分类选项列表
// @Description 获取标签选项列表
// @Tags Category
// @Accept json
// @Produce json
// @Success 0 {object} Response[[]model.OptionVO]
// @Security ApiKeyAuth
// @Router /category/option [get]
func (*Category) GetOption(c *gin.Context) {
	list, err := model.GetCategoryOption(GetDB(c))
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, list)
}

```

model/catogory.go

```go
// GetCategoryOption 获取分类选项列表
func GetCategoryOption(db *gorm.DB) ([]OptionVO, error) {
	var list = make([]OptionVO, 0)
	result := db.Model(&Category{}).Select("id", "name").Find(&list)
	return list, result.Error
}
```

对应的请求和响应如下：

在点击文章列表时，会自动调用分类选项列表加载对应的分类选项，所以会发送如下请求：

<img src="./assets/image-20250326170447291.png" alt="image-20250326170447291" style="zoom:67%;" />

![image-20250326170430278](./assets/image-20250326170430278.png)

得到的响应为分类的 ID-NAME 所构成的组合：

```go
type OptionVO struct {
	ID   int    `json:"value"`
	Name string `json:"name"`
}
```



### 1.5 补充

对于博客前台相关接口，我们也需要进行对应的接口补充：

首先在 manager.go 中补充如下操作：

```go
category := base.Group("/category")
{
  category.GET("/list", frontAPI.GetCategoryList) // 前台分类列表
}
```

同时在 handle_front.go 中进行补齐：

```go
// GetCategoryList 查询分类列表
func (*Front) GetCategoryList(c *gin.Context) {
  list, _, err := model.GetCategoryList(GetDB(c), 1, 1000, "")
  if err != nil {
    ReturnError(c, global.ErrDbOp, err)
    return
  }
  ReturnSuccess(c, list)
}
```



## 2 标签模块

此模块提供了标签的相关管理接口，功能与分类模块类似，主要对标签进行增删改查操作，以及获取标签选项列表。

1. **标签列表获取**：`GET /tag/list` 接口可获取所有标签的列表信息，可用于标签管理页面展示。
2. **标签新增 / 编辑**：通过 `POST /tag` 接口，可实现标签的新增或对已有标签信息的编辑。
3. **标签删除**：`DELETE /tag` 接口用于删除指定的标签。
4. **标签选项列表获取**：`GET /tag/option` 接口获取标签的选项列表，适用于需要选择标签的场景。

```go
// 标签模块
tag := auth.Group("/tag")
{
   tag.GET("/list", tagAPI.GetList)     // 标签列表
   tag.POST("", tagAPI.SaveOrUpdate)    // 新增/编辑标签
   tag.DELETE("", tagAPI.Delete)        // 删除标签
   tag.GET("/option", tagAPI.GetOption) // 标签选项列表
}
```

标签模块的思路与分类模块很像，其主要的功能如下：

### 2.1 标签列表 /tag/list

manager.go

```go
tag.GET("/list", tagAPI.GetList)     // 标签列表
```

handle/handle_tag.go

```go
type Tag struct{}

type AddOrEditTagReq struct {
	ID   int    `json:"id"`
	Name string `json:"name" binding:"required"`
}

// GetList 获取标签列表
// @Summary 获取标签列表
// @Description 根据条件查询获取标签列表
// @Tags Tag
// @Param page_size query int false "当前页数"
// @Param page_num query int false "每页条数"
// @Param keyword query string false "搜索关键字"
// @Accept json
// @Produce json
// @Success 0 {object} Response[PageResult[model.TagVO]] "成功"
// @Security ApiKeyAuth
// @Router /tag/list [get]
func (*Tag) GetList(c *gin.Context) {
	var query PageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		ReturnError(c, global.ErrRequest, err)
		return
	}

	data, total, err := model.GetTagList(GetDB(c), query.Page, query.Size, query.Keyword)
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, PageResult[model.TagVO]{
		Total: total,
		List:  data,
		Size:  query.Size,
		Page:  query.Page,
	})
}

```

model/tag.go

```go
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
```

对应的请求和响应如下：

![image-20250327175304608](./assets/image-20250327175304608.png)

<img src="./assets/image-20250327175324003.png" alt="image-20250327175324003"  />

<img src="./assets/image-20250327175344609.png" alt="image-20250327175344609" style="zoom:67%;" />



### 2.2 新增/编辑标签 /tag POST

manager.go

```go
tag.POST("", tagAPI.SaveOrUpdate)    // 新增/编辑标签
```

handle/handle_tag.go

```go

type AddOrEditTagReq struct {
	ID   int    `json:"id"`
	Name string `json:"name" binding:"required"`
}

// SaveOrUpdate 添加或者修改标签
// @Summary 添加或修改标签
// @Description 添加或修改标签
// @Tags Tag
// @Param form body AddOrEditTagReq true "添加或修改标签"
// @Accept json
// @Produce json
// @Success 0 {object} Response[model.Tag]
// @Security ApiKeyAuth
// @Router /tag [post]
func (*Tag) SaveOrUpdate(c *gin.Context) {
	var form AddOrEditTagReq
	if err := c.ShouldBindJSON(&form); err != nil {
		ReturnError(c, global.ErrRequest, err)
		return
	}

	tag, err := model.SaveOrUpdateTag(GetDB(c), form.ID, form.Name)
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, tag)
}

```

绑定形式：

+ **JSON 标签**：`ShouldBindJSON` 会依据结构体字段的 `json` 标签来匹配 JSON 数据中的字段名。比如在你给出的 `AddOrEditTagReq` 结构体里，`ID` 字段对应 JSON 数据里的 `"id"`，`Name` 字段对应 `"name"`。
+ **无标签情况**：若结构体字段没有 `json` 标签，`ShouldBindJSON` 会使用字段名的小写形式进行匹配。

model/tag.go

```go
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
```

对应的请求和响应如下：

<img src="./assets/image-20250327174016436.png" alt="image-20250327174016436" style="zoom: 67%;" />

![image-20250327174100630](/Users/tianjiangyu/MyStudy/Go-learning/B2-Gin-Vue-Admin/assets/image-20250327174100630.png)

对于新建一个新的标签，可以看到其对应的 AddOrEditTagReq 中  ID 为 0，之后进行的操作为 result = db.Updates(&tag)

同时我们尝试编辑一下 Tag 标签的名称，如下：

<img src="./assets/image-20250327174931600.png" alt="image-20250327174931600" style="zoom:67%;" />

![image-20250327175108504](./assets/image-20250327175108504.png)

可以看到其对应的 AddOrEditTagReq 中  ID 为 1，之后进行的操作为 result = db.Updates(&tag)，从而将标签进行删除。



### 2.3 删除标签 /tag DELETE

manager.go

```go
tag.DELETE("", tagAPI.Delete)        // 删除标签
```

handle/handle_tag.go

```go
// Delete 删除标签（可以批量操作）
// TODO: 删除行为, 添加强制删除: 有关联数据则将删除关联数据
// @Summary 删除标签（批量）
// @Description 根据 ID 数组删除标签
// @Tags Tag
// @Param ids body []int true "标签 ID 数组"
// @Accept json
// @Produce json
// @Success 0 {object} Response[int]
// @Security ApiKeyAuth
// @Router /tag [delete]
func (*Tag) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		ReturnError(c, global.ErrRequest, err)
		return
	}

	db := GetDB(c)
	// 检查标签下面是否有文章
	count, err := model.Count(db, &model.ArticleTag{}, "tag_id in ?", ids)
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}

	if count > 0 {
		ReturnError(c, global.ErrTagHasArt, nil)
		return
	}

	result := db.Delete(model.Tag{}, "id in ?", ids)
	if result.Error != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, result.RowsAffected)
}
```

对应的请求和响应如下：

![image-20250327175425379](./assets/image-20250327175425379.png)

![image-20250327180139789](./assets/image-20250327180139789.png)

可以看到，标签被成功删除。



### 2.4 标签选项列表 /tag/option

manager.go

```go
tag.GET("/option", tagAPI.GetOption) // 标签选项列表
```

handle/handle_tag.go

```go
// GetOption 获取标签选项列表
// @Summary 获取标签选项列表
// @Description 获取标签选项列表
// @Tags Tag
// @Accept json
// @Produce json
// @Success 0 {object} Response[model.OptionVO]
// @Security ApiKeyAuth
// @Router /tag/option [get]
func (*Tag) GetOption(c *gin.Context) {
	list, err := model.GetTagOption(GetDB(c))
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, list)
}
```

model/tag.go

```go
// GetTagOption 获取标签选项列表
func GetTagOption(db *gorm.DB) ([]OptionVO, error) {
	list := make([]OptionVO, 0)
	result := db.Model(&Tag{}).Select("id", "name").Find(&list)
	return list, result.Error
}
```

对应的请求和响应如下：

![image-20250327180219370](./assets/image-20250327180219370.png)

![image-20250327180243294](./assets/image-20250327180243294.png)

得到的响应为分类的 ID-NAME 所构成的组合：

```go
type OptionVO struct {
	ID   int    `json:"value"`
	Name string `json:"name"`
}
```



### 2.5 补充

对于博客前台相关接口，我们也需要进行对应的接口补充：

首先在 manager.go 中补充如下操作：

```go
tag := base.Group("/tag")
{
  tag.GET("/list", frontAPI.GetTagList) // 前台标签列表
}
```

同时在 handle_front.go 中进行补齐：

```go
// GetTagList 查询标签列表
func (*Front) GetTagList(c *gin.Context) {
	list, _, err := model.GetTagList(GetDB(c), 1, 1000, "")
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, list)
}
```





## 3 文章模块

文章模块提供了一系列丰富的接口，用于对文章进行全面的管理，包括文章的增删改查、置顶设置、软删除、物理删除以及导入导出等操作。

1. **文章列表获取**：`GET /article/list` 接口用于获取文章的列表信息，可用于文章列表展示页面。
2. **文章新增 / 编辑**：`POST /article` 接口可实现文章的新增或者对已有文章内容的编辑修改。
3. **文章置顶更新**：`PUT /article/top` 接口可更新文章的置顶状态，使文章在列表中置顶显示。
4. **文章详情获取**：`GET /article/:id` 接口根据文章的 ID 获取文章的详细信息，用于文章详情展示页面。
5. **文章软删除**：`PUT /article/soft - delete` 接口对文章进行软删除操作，文章不会真正从数据库中移除，只是标记为已删除状态。
6. **文章物理删除**：`DELETE /article` 接口会将文章从数据库中彻底删除。
7. **文章导出**：`POST /article/export` 接口可将文章信息导出，可能以某种文件格式（如 CSV、Excel 等）保存。
8. **文章导入**：`POST /article/import` 接口允许用户将外部文件中的文章信息导入到系统中。

```go
// 文章模块
articles := auth.Group("/article")
{
   articles.GET("/list", articleAPI.GetList)                 // 文章列表
   articles.POST("", articleAPI.SaveOrUpdate)                // 新增/编辑文章
   articles.PUT("/top", articleAPI.UpdateTop)                // 更新文章置顶
   articles.GET("/:id", articleAPI.GetDetail)                // 文章详情
   articles.PUT("/soft-delete", articleAPI.UpdateSoftDelete) // 软删除文章
   articles.DELETE("", articleAPI.Delete)                    // 物理删除文章
   articles.POST("/export", articleAPI.Export)               // 导出文章
   articles.POST("/import", articleAPI.Import)               // 导入文章
}
```



### 3.1 文章列表获取 /article/list

manager.go

```go
var (
	// 后端管理系统接口
	...
	articleAPI  handle.Article  // 文章
)


// 文章模块
articles := auth.Group("/article")
{
  articles.GET("/list", articleAPI.GetList) // 文章列表
}
```

handle/handle_article.go

```go
package handle

import (
	"gin-blog-server/internal/global"
	"gin-blog-server/internal/model"
	"github.com/gin-gonic/gin"
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

type ArticleVO struct {
	model.Article
	// gorm:"-" 是一个标签，它的主要作用是告诉 GORM 忽略结构体中的某个字段，使其在数据库操作（如创建表、插入数据、查询数据等）中不被考虑。
	// GORM 默认会将结构体名转换为蛇形命名法（snake_case）并以复数形式作为数据库表名。
	LikeCount    int `json:"like_count" gorm:"-"`
	ViewCount    int `json:"view_count" gorm:"-"`
	CommentCount int `json:"comment_count" gorm:"-"`
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
```

model/article.go

```go
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
		Joins("LEFT JOIN article_tag ON article_tag.article_id = article.id")
	// .Group("id") // 去重（这行似乎没有作用，先注释掉）

	if tagId != 0 {
		db = db.Where("tag_id = ?", tagId)
	}

	result := db.Count(&total).
		Scopes(Paginate(page, size)).
		Order("is_top DESC, article.id DESC").
		Find(&list)
	return list, total, result.Error
}
```

**我们先去完成功能  3.2 文章新增/编辑  /article POST，有了文章之后，再来测试文章列表获取的功能是否可以正常使用。**

对应的请求和响应如下：

![image-20250328122047990](./assets/image-20250328122047990.png)

![image-20250328122106345](./assets/image-20250328122106345.png)

这里需要注意的是，在获取文章列表时，我们使用的 sql 为：

```go
db = db.Preload("Category").Preload("Tags").
		Joins("LEFT JOIN article_tag ON article_tag.article_id = article.id").
		Group("id") // 去重
```

其中：

- `Preload("Category")` 和 `Preload("Tags")` 预加载 `Category` 和 `Tags` 关联数据。
- `LEFT JOIN article_tag ON article_tag.article_id = article.id`
   这是一个 `LEFT JOIN`，`article_tag` 可能有多条记录，如果 `article` 关联了多个 `Tags`，则 `article` 会被重复返回。
- `Group("id")`
   这个 `GROUP BY id` 用来去重，让结果只保留 `article.id` 唯一的记录。

**如果没有 Group("id") 的存在，则会获取到两条同样的 Article 信息。**

`LEFT JOIN` 连接 `article_tag` 表时，如果 `article` 关联了多个 `tags`，则 `article` 会被重复返回。例如：

| article.id | article.title | article.content | article_tag.tag_id |
| ---------- | ------------- | --------------- | ------------------ |
| 1          | "Gorm"        | "Go ORM"        | 101                |
| 1          | "Gorm"        | "Go ORM"        | 102                |

在 `SQL` 结果集中，`article.id=1` 被返回了两次，因为 `article_tag` 表中有两条数据，导致 `JOIN` 结果中 `article` 被**重复**。

添加 `GROUP BY id` 后，SQL 只保留每个 `article.id` 的**一条记录**，避免重复返回。

但是需要注意：

- `GROUP BY` 可能导致 `SELECT` 语句中的其他字段（如 `Tags`）变得不确定。
- `Gorm` 的 `Preload("Tags")` 会在 `GROUP BY` 后**额外查询** `Tags`，所以不会影响 `Tags` 的完整性。



### 3.2 文章新增/编辑  /article POST

manager.go

```go
articles.POST("", articleAPI.SaveOrUpdate)                // 新增/编辑文章
```

handle/handle_article.go

```go
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
```

model/article.go

```go
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
```

为了测试上述功能，我们首先要对数据库中的 config 表进行补全，通过config.sql 导入即可，补全后的 config 内容如如下，我们可以从中获取默认的头像设置。

![image-20250328120133017](./assets/image-20250328120133017.png)

对应的请求和响应如下：

1. 首先编写文章并选择文章分类、文章标签等信息

![image-20250328121243837](./assets/image-20250328121243837.png)

![image-20250328121332804](./assets/image-20250328121332804.png)

2. 点击确认之后查看发送的请求：

<img src="./assets/image-20250328121521147.png" alt="image-20250328121521147"  />

请求中携带的 body 字段如下：

![image-20250328121557301](./assets/image-20250328121557301.png)

对应的返回信息如下，证明成功插入了一篇新的文章：

![image-20250328121809120](./assets/image-20250328121809120.png)

之后，我们可以到 3.1 文章列表获取中查看文章内容。



### 3.3 文章置顶更新  PUT /article/top

manager.go

```go
articles.PUT("/top", articleAPI.UpdateTop)                // 更新文章置顶
```

handle/handle_article.go

```go
type UpdateArticleTopReq struct {
	ID    int  `json:"id"`
	IsTop bool `json:"is_top"`
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

```

model/article.go

```go
// UpdateArticleTop 修改置顶信息
func UpdateArticleTop(db *gorm.DB, id int, isTop bool) error {
	result := db.Model(&Article{Model: Model{ID: id}}).Update("is_top", isTop)
	return result.Error
}
```

> 1. `db.Model(&Article{Model: Model{ID: id}}).Update("is_top", isTop)`
>
> - **功能**：这种方式使用 `Model` 方法指定要操作的记录，通过 `Update` 方法只更新单个字段 `is_top` 的值。`Update` 方法接受两个参数，第一个参数是要更新的字段名，第二个参数是要更新的值。
> - **适用场景**：当你只需要更新数据库中某条记录的单个字段时，使用这种方式较为合适。它只会更新指定的字段，不会影响其他字段的值。
> - **示例代码解释**：
>
> ```go
> // UpdateArticleTop 修改置顶信息
> func UpdateArticleTop(db *gorm.DB, id int, isTop bool) error {
>     result := db.Model(&Article{Model: Model{ID: id}}).Update("is_top", isTop)
>     return result.Error
> }
> ```
>
> 在这个函数中，它会找到 `Article` 表中 `ID` 为 `id` 的记录，并将 `is_top` 字段更新为 `isTop` 的值，其他字段不会受到影响。
>
> 2. `db.Updates(&Article{Model: Model{ID: id}, IsTop: IsTop})`
>
> - **功能**：`Updates` 方法会更新结构体中所有非零值字段。它会根据结构体中设置的字段值，将这些值更新到数据库中对应 `ID` 的记录里。
> - **适用场景**：当你需要同时更新数据库中某条记录的多个字段时，使用这种方式比较方便。不过要注意，如果结构体中的其他字段有默认值，这些默认值也会被更新到数据库中。
> - **示例代码解释**：
>
> ```go
> db.Updates(&Article{Model: Model{ID: id}, IsTop: IsTop})
> ```
>
> 这里会找到 `Article` 表中 `ID` 为 `id` 的记录，然后将结构体中设置的**非零值字段** 更新到该记录中。如果 `Article` 结构体还有其他字段且有默认值，这些字段也会被更新到数据库中。

对应的请求和响应如下：

![image-20250328141617146](./assets/image-20250328141617146.png)

![image-20250328141652161](./assets/image-20250328141652161.png)



### 3.4 文章详情获取 GET /article/:id

manager.go

```go
articles.GET("/:id", articleAPI.GetDetail) // 文章详情
```

handle/handle_article.go

```go
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
```

model/article.go

```go
// GetArticle 文章的详细信息
func GetArticle(db *gorm.DB, id int) (data *Article, err error) {
	result := db.Preload("Category").Preload("Tags").
		Where(Article{Model: Model{ID: id}}).
		First(&data)
	return data, result.Error
}
```

对应的请求和响应如下：

点击查看某个具体文章时，会发送对应的请求：

![image-20250328143358343](./assets/image-20250328143358343.png)

返回具体的文章信息：

![image-20250328143720889](./assets/image-20250328143720889.png)



### 3.5 文件软删除 PUT /article/soft

manager.go

```go
articles.PUT("/soft-delete", articleAPI.UpdateSoftDelete) // 软删除文章
```

handle/handle_article.go

```go
type SoftDeleteReq struct {
	Ids      []int `json:"ids"`
	IsDelete bool  `json:"is_delete"`
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
```

model/article.go

```go
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
```

对应的请求和响应如下：

操作点击删除：

![image-20250328145443379](./assets/image-20250328145443379.png)

查看发送的请求：

![image-20250328145516914](./assets/image-20250328145516914.png)

之后，更新文章列表会发现，文章已经被成功删除：

<img src="./assets/image-20250328145554845.png" alt="image-20250328145554845" style="zoom:67%;" />

可以到数据库中查看，当前文章的记录在数据库中并没有被删除，只不过是将 is_delete 字段修改为 1，如果将其修改为 0 的话，我们刷新界面，会发现文章重新出现在了文章列表中。

![image-20250328145851211](./assets/image-20250328145851211.png)



### 3.6 文章物理删除 DELETE /article

文章物理删除的主要目的是，撤出清除数据库中的记录条数，软删除和物理删除的区别在于：

+ 对于文章，点击删除 -> 软删除 -> 文章列表不可见，查看回收站可见
+ 对于回收站中的文章，其 is_delete 属性已经为 1，此时点击删除，会触发物理删除 -> 彻底删除文章

前端代码中的核心逻辑如下：

```javascript
const extraParams = ref({
    is_delete: null, // 未删除 | 回收站
    status: null, // null-all, 1-公开, 2-私密, 3-草稿
})

function updateOrDeleteArticles(ids) {
    extraParams.value.is_delete
        ? api.deleteArticle(ids)
        : api.softDeleteArticle(JSON.parse(ids), true)
}

// 切换标签页: [全部, 公开, 私密, 草稿箱, 回收站]
function handleChangeTab(value) {
    switch (value) {
        case 'all':
            extraParams.value.is_delete = 0
            extraParams.value.status = null
            break
        case 'public':
            extraParams.value.is_delete = 0
            extraParams.value.status = 1
            break
        case 'secret':
            extraParams.value.is_delete = 0
            extraParams.value.status = 2
            break
        case 'draft':
            extraParams.value.is_delete = 0
            extraParams.value.status = 3
            break
        case 'delete':
            extraParams.value.is_delete = 1
            extraParams.value.status = null
            break
    }
    $table.value?.handleSearch()
}
```

manager.go

```go
articles.DELETE("", articleAPI.Delete)                    // 物理删除文章
```

handle/handle_article.go

```go
// Delete 物理删除文章
func (*Article) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		ReturnError(c, global.ErrRequest, err)
		return
	}

	rows, err := model.UpdateArticle(GetDB(c), ids)
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}

	ReturnSuccess(c, rows)
}
```

model/article.go

```go
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

```

对应的请求和响应如下：

当我们点击回收站时，会将 is_delete = 1 的文章全部展示

![image-20250328154027686](./assets/image-20250328154027686.png)

然后点击删除，会将这些文章从数据库记录中删除。



### 3.7 文章导出 POST /article/export

manager.go

```go
articles.POST("/export", articleAPI.Export)               // 导出文章
```

handle/handle_article.go

```go
// Export 导出文章: 获取导出后的资源链接列表
// TODO: 目前是前端导出
func (*Article) Export(c *gin.Context) {
	ReturnSuccess(c, nil)
}
```

前端对于导出的实现如下：

```javascript
// 导出文章
async function exportArticles(ids) {
    // 方式一: 前端根据文章内容和标题进行导出
    const list = $table.value?.tableData.filter(e => ids.includes(e.id))
    for (const item of list)
        downloadFile(item.content, `${item.title}.md`)

    // 方式二: 后端导出返回链接, 前端根据链接下载
    // const res = await api.exportArticles(ids)
    // for (const url of res.data)
    // downloadFile(url)
}

function downloadFile(content, fileName) {
    const aEle = document.createElement('a') // 创建下载链接
    aEle.download = fileName // 设置下载的名称
    aEle.style.display = 'none'// 隐藏的可下载链接
    // 字符内容转变成 blob 地址
    const blob = new Blob([content])
    aEle.href = URL.createObjectURL(blob)
    // 绑定点击时间
    document.body.appendChild(aEle)
    aEle.click()
    // 然后移除
    document.body.removeChild(aEle)
}
```

对应的操作过程如下：

![image-20250329095149340](./assets/image-20250329095149340.png)

![image-20250329095413387](./assets/image-20250329095413387.png)



### 3.8 文章导入 POST /article/import

manager.go

```go
articles.POST("/import", articleAPI.Import)               // 导入文章
```

handle/handle_article.go

```go
// Import 倒入文章：题目 + 内容
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
```

model/article.go

```go
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
```

对应的请求和响应如下：

1. 点击批量导入，随机选中一个 md 文件进行上传：

<img src="/Users/tianjiangyu/MyStudy/Go-learning/B2-Gin-Vue-Admin/assets/image-20250329110039172.png" alt="image-20250329110039172" style="zoom:67%;" />

初步怀疑

+ 这里我发现会报错： level=INFO msg="[Func-ReturnError] TOKEN 不存在，请重新登陆"

+ 报错的原因是 middleware/JWTAuth 进入到了 "没有找到的资源，不需要鉴权，跳过后续的验证过程" 分支，这种情况应该是数据库 - resource 表中缺少了导入操作的操作权限
+ 但是，我从数据库中检查，发现数据库配置没有问题

我对  middleware/JWTAuth 进行 debug 排查，发现前端发送 import 请求时没有携带Authorization 中 token，所以才会报错 token 不存在

这里去前端代码中进行查看发现，导入文章时候，并没有调用 api 接口，应该是从按钮处直接发送的导入请求：

![image-20250329112143470](./assets/image-20250329112143470.png)

 具体代码如下：这里直接 action 调用 /api/article/import，所以没有带 token 进行请求

```html
<div class="inline-block">
  <NUpload action="/api/article/import" :show-file-list="false" multiple @before-upload="beforeUpload"
           @finish="afterUpload">
    <NButton type="success">
      <template #icon>
        <p class="i-mdi:import" />
      </template>
      批量导入
    </NButton>
  </NUpload>
</div>
```

我们将前端代码进行修改：

```html
<div class="inline-block">
  <NUpload  action="/api/article/import" :show-file-list="false" :headers="uploadHeaders" multiple @before-upload="beforeUpload" @finish="afterUpload">
    <NButton type="success">
      <template #icon>
        <p class="i-mdi:import" />
          </template>
批量导入
  </NButton>
</NUpload>
</div>
```

```javascript
import { useAuthStore } from '@/store'

// 获取 token（通常在 store 中存储）
const { token } = useAuthStore()

// 定义上传请求头
const uploadHeaders = ref({
  'Authorization': `Bearer ${token}`
});
```

回到前端页面，再次点击上传：

![image-20250329125730793](./assets/image-20250329125730793.png)

我们可以查看其携带的 payload，可以看到他将文章的内容通过 payload 进行二进制传输：

![image-20250329125915470](./assets/image-20250329125915470.png)

可以看到文件被正确上传：

![image-20250329130238431](./assets/image-20250329130238431.png)