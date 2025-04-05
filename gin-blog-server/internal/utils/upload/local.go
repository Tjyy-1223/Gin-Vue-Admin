package upload

import (
	"errors"
	"gin-blog-server/internal/global"
	"gin-blog-server/internal/utils"
	"io"
	"log/slog"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"
)

// Local 本地文件上传
type Local struct{}

// UploadFile 文件上传到本地
func (*Local) UploadFile(file *multipart.FileHeader) (filePath, fileName string, err error) {
	ext := path.Ext(file.Filename)
	name := strings.TrimSuffix(file.Filename, ext) // 读取文件名
	name = utils.MD5(name)

	filename := name + "_" + time.Now().Format("20060102150405") + ext // 拼接新文件名

	conf := global.Conf.Upload
	mkdirErr := os.MkdirAll(conf.StorePath, os.ModePerm) // 尝试创建存储路径
	if mkdirErr != nil {
		slog.Error("function os.MkdirAll() Filed", slog.Any("err", mkdirErr.Error()))
		return "", "", errors.New("function os.MkdirAll() Filed, err:" + mkdirErr.Error())
	}

	storePath := conf.StorePath + "/" + filename // 文件存储路径
	filepath := conf.Path + "/" + filename       // 文件展示路径

	f, openError := file.Open() // 读取文件
	if openError != nil {
		slog.Error("function file.Open() Filed", slog.String("err", openError.Error()))
		return "", "", errors.New("function file.Open() Filed, err:" + openError.Error())
	}
	defer f.Close() // 创建文件 defer 关闭

	// 创建文件的保存位置，即文件写入位置
	out, createErr := os.Create(storePath)
	if createErr != nil {
		slog.Error("function os.Create() Filed", slog.String("err", createErr.Error()))
		return "", "", errors.New("function os.Create() Filed, err:" + createErr.Error())
	}
	defer out.Close()

	_, copyErr := io.Copy(out, f) // 拷贝文件
	if copyErr != nil {
		slog.Error("function io.Copy() Filed", slog.String("err", copyErr.Error()))
		return "", "", errors.New("function io.Copy() Filed, err:" + copyErr.Error())
	}
	return filepath, filename, nil
}

// DeleteFile 从本地删除文件
func (*Local) DeleteFile(key string) error {
	p := global.GetConfig().Upload.StorePath + "/" + key
	if strings.Contains(p, global.GetConfig().Upload.StorePath) {
		if err := os.Remove(p); err != nil {
			return errors.New("本地文件删除失败, err:" + err.Error())
		}
	}
	return nil
}
