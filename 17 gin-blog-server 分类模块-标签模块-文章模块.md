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



manager.go

```go

```

handle/handle_tag.go

```go

```

model/tag.go

```go

```

对应的请求和响应如下：



