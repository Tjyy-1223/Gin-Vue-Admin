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
