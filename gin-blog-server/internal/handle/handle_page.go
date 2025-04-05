package handle

import (
	"gin-blog-server/internal/global"
	"gin-blog-server/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

type Page struct{}

// GetList 获取页面列表
// @Summary 获取页面列表
// @Description 根据条件查询获取页面列表
// @Tags Page
// @Accept json
// @Produce json
// @Success 0 {object} Response[[]model.Page]
// @Security ApiKeyAuth
// @Router /page/list [get]
func (*Page) GetList(c *gin.Context) {
	db := GetDB(c)
	rdb := GetRDB(c)

	// get from cache
	cache, err := getPageCache(rdb)
	if cache != nil && err == nil {
		slog.Debug("[handle-page-GetList] get page list from cache")
		ReturnSuccess(c, cache)
		return
	}

	switch err {
	case redis.Nil:
		break
	default:
		ReturnError(c, global.ErrRedisOp, err)
		return
	}

	// get from db
	data, _, err := model.GetPageList(db)
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}

	// add to cache
	if err := addPageCache(GetRDB(c), data); err != nil {
		ReturnError(c, global.ErrRedisOp, err)
		return
	}

	ReturnSuccess(c, data)
}

// SaveOrUpdate 添加或修改页面
// @Summary 添加或修改页面
// @Description 添加或修改页面
// @Tags Page
// @Param form body AddOrEditPageReq true "添加或修改页面"
// @Accept json
// @Produce json
// @Success 0 {object} Response[model.Page]
// @Security ApiKeyAuth
// @Router /page [post]
func (*Page) SaveOrUpdate(c *gin.Context) {
	var req model.Page
	if err := c.ShouldBindJSON(&req); err != nil {
		ReturnError(c, global.ErrRequest, err)
		return
	}

	db := GetDB(c)
	rdb := GetRDB(c)

	page, err := model.SaveOrUpdatePage(db, req.ID, req.Name, req.Label, req.Cover)
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}

	// delete cache
	if err := removePageCache(rdb); err != nil {
		ReturnError(c, global.ErrRedisOp, err)
		return
	}

	ReturnSuccess(c, page)
}

// Delete 删除封面
// @Summary 删除页面
// @Description 根据 ID 数组删除页面
// @Tags Page
// @Param ids body []int true "页面 ID 数组"
// @Accept json
// @Produce json
// @Success 0 {object} Response[int]
// @Security ApiKeyAuth
// @Router /page [delete]
func (*Page) Delete(c *gin.Context) {
	var ids []int
	if err := c.ShouldBindJSON(&ids); err != nil {
		ReturnError(c, global.ErrRequest, err)
		return
	}

	result := GetDB(c).Delete(model.Page{}, "id in ?", ids)
	if result.Error != nil {
		ReturnError(c, global.ErrDbOp, result.Error)
		return
	}

	// delete cache
	if err := removePageCache(GetRDB(c)); err != nil {
		ReturnError(c, global.ErrRedisOp, err)
		return
	}

	ReturnSuccess(c, result.RowsAffected)
}
