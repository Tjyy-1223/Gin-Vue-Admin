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
