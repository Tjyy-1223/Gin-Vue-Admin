package handle

import (
	"gin-blog-server/internal/global"
	"gin-blog-server/internal/utils/upload"
	"github.com/gin-gonic/gin"
)

type Upload struct{}

// UploadFile 上传文件
// @Summary 上传文件
// @Description 上传文件
// @Tags upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "文件"
// @Success 0 {object} Response[string]
// @Router /upload/file [post]
func (*Upload) UploadFile(c *gin.Context) {
	// 获取文件头信息
	_, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		ReturnError(c, global.ErrFileReceive, err)
		return
	}

	// 获取 oss 对象存储接口
	oss := upload.NewOSS()
	// 实现图片上传
	filePath, _, err := oss.UploadFile(fileHeader)
	if err != nil {
		ReturnError(c, global.ErrFileUpload, err)
		return
	}

	ReturnSuccess(c, filePath)
}
