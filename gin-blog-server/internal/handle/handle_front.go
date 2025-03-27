package handle

import (
	"gin-blog-server/internal/global"
	"gin-blog-server/internal/model"
	"github.com/gin-gonic/gin"
)

type Front struct{}

// GetHomeInfo 前台首页信息
func (*Front) GetHomeInfo(c *gin.Context) {
	db := GetDB(c)
	rdb := GetRDB(c)

	data, err := model.GetFrontStatistics(db)
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}
	data.ViewCount, _ = rdb.Get(rctx, global.VIEW_COUNT).Int64()

	ReturnSuccess(c, data)
}

// GetCategoryList 查询分类列表
func (*Front) GetCategoryList(c *gin.Context) {
	list, _, err := model.GetCategoryList(GetDB(c), 1, 1000, "")
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, list)
}

// GetTagList 查询标签列表
func (*Front) GetTagList(c *gin.Context) {
	list, _, err := model.GetTagList(GetDB(c), 1, 1000, "")
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, list)
}
