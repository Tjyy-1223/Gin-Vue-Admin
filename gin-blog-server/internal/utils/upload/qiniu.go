package upload

import (
	"context"
	"errors"
	"fmt"
	"gin-blog-server/internal/global"
	"gin-blog-server/internal/utils"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"mime/multipart"
	"path"
	"time"
)

// Qiniu 七牛云文件上传
type Qiniu struct{}

func (*Qiniu) UploadFile(file *multipart.FileHeader) (filePath, fileName string, err error) {
	putPolicy := storage.PutPolicy{Scope: global.GetConfig().Qiniu.Bucket}
	mac := qbox.NewMac(global.GetConfig().Qiniu.AccessKey, global.GetConfig().Qiniu.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	formUploader := storage.NewFormUploader(qiniuConfig())

	ret := storage.PutRet{}
	putExtra := storage.PutExtra{Params: map[string]string{"x:name": "github logo"}}

	f, openError := file.Open()
	if openError != nil {
		return "", "", errors.New("function file.Open() Filed, err:" + openError.Error())
	}
	defer f.Close()

	// 文件名格式 建议保证唯一性
	fileKey := fmt.Sprintf("%d%s%s", time.Now().Unix(), utils.MD5(file.Filename), path.Ext(file.Filename))
	putErr := formUploader.Put(context.Background(), &ret, upToken, fileKey, f, file.Size, &putExtra)
	if putErr != nil {
		return "", "", errors.New("function formUploader.Put() Filed, err:" + putErr.Error())
	}
	return global.GetConfig().Qiniu.ImgPath + "/" + ret.Key, ret.Key, nil
}

func (*Qiniu) DeleteFile(key string) error {
	mac := qbox.NewMac(global.GetConfig().Qiniu.AccessKey, global.GetConfig().Qiniu.SecretKey)
	cfg := qiniuConfig()
	bucketManager := storage.NewBucketManager(mac, cfg)

	if err := bucketManager.Delete(global.GetConfig().Qiniu.Bucket, key); err != nil {
		return errors.New("function bucketManager.Delete() Filed, err:" + err.Error())
	}
	return nil
}

// 七牛云配置信息
func qiniuConfig() *storage.Config {
	cfg := storage.Config{
		UseHTTPS:      global.GetConfig().Qiniu.UseHTTPS,
		UseCdnDomains: global.GetConfig().Qiniu.UseCdnDomains,
	}
	switch global.GetConfig().Qiniu.Zone { // 根据配置文件进行初始化空间对应的机房
	case "ZoneHuadong":
		cfg.Zone = &storage.ZoneHuadong
	case "ZoneHuabei":
		cfg.Zone = &storage.ZoneHuabei
	case "ZoneHuanan":
		cfg.Zone = &storage.ZoneHuanan
	case "ZoneBeimei":
		cfg.Zone = &storage.ZoneBeimei
	case "ZoneXinjiapo":
		cfg.Zone = &storage.ZoneXinjiapo
	}
	return &cfg
}
