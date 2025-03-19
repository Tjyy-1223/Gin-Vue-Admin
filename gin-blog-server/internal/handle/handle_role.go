package handle

import (
	"gin-blog-server/internal/global"
	"gin-blog-server/internal/model"
	"github.com/gin-gonic/gin"
)

type Role struct{}

// GetOption 获取角色选项
// @Summary 获取角色选项
// @Description 获取角色选项
// @Tags role
// @Produce json
// @Success 0 {object} Response[model.OptionVO]
// @Router /role/option [get]
func (*Role) GetOption(c *gin.Context) {
	list, err := model.GetRoleOption(GetDB(c))
	if err != nil {
		ReturnError(c, global.ErrDbOp, err)
		return
	}
	ReturnSuccess(c, list)
}
