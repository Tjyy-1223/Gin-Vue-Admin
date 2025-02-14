package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"gin-blog-server/internal/global"
	"github.com/k3a/html2text"
	"github.com/thanhpk/randstr"
	"github.com/vanng822/go-premailer/premailer"
	"gopkg.in/gomail.v2"
	"html/template"
	"io/fs"
	"log/slog"
	"path/filepath"
	"strconv"
	"strings"
)

// EmailData 注册的核心思想：
// 1. 发送邮件的同时创建 code 存储在本地 redis 中
// 2. 当用户点击验证链接时即向 sever 发出 code 数据，如果这个数据在 redis 中存在，则验证成功，否则失败
type EmailData struct {
	URL      template.URL // 验证链接
	UserName string       // 用户名即邮箱地址
	Subject  string       // 邮箱主题
}

// Format 将邮箱地址转换成小写，并去除空格
// 格式化邮件可以防止写错大小写重复注册，同时给用户预留犯错空间，输入空格和大小写错误也能正常处理
func Format(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// GenEmailVerificationInfo 返回生成加密后的 base64 字符串
func GenEmailVerificationInfo(email string, password string) string {
	code := GetCode()
	info := Encode(email + "|" + password + "|" + code)
	return info
}

// ParseEmailVerificationInfo 返回解析base64字符串后的 邮箱地址和code
func ParseEmailVerificationInfo(info string) (string, string, error) {
	data, err := Decode(info)
	if err != nil {
		return "", "", err
	}

	str := strings.Split(data, "|")
	if len(str) != 3 {
		return "", "", errors.New("wrong verification info format")
	}
	return str[0], str[1], nil
}

// GetCode 生成随机字符串
func GetCode() string {
	code := randstr.String(24)
	return code
}

// Encode 返回生成 base64 编码
func Encode(s string) string {
	data := base64.StdEncoding.EncodeToString([]byte(s))
	return data
}

// Decode 返回解码 base64
func Decode(s string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", errors.New("email verify failed, decode error")
	}
	return string(data), nil
}

// GetEmailData 生成邮件数据
func GetEmailData(email string, info string) *EmailData {
	return &EmailData{
		URL:      template.URL(GetEmailVerifyURL(info)), // 验证链接
		UserName: email,                                 // 用户邮箱地址
		Subject:  "请完成账号注册",                             // 邮件主题
	}
}

// GetEmailVerifyURL 生成验证链接
func GetEmailVerifyURL(info string) string {
	baseurl := global.GetConfig().Server.Port
	if baseurl[0] == ':' {
		baseurl = fmt.Sprintf("localhost:%s", baseurl)
	}
	// 如果是用docker部署,则 注释上面的代码，使用下面的代码
	// baseurl := "你的域名"   切记不需要加端口

	// 点击该链接可以触发 api/email/verify -> 进一步将账号存储到对应数据库中，完成账号注册
	return fmt.Sprintf("%s/api/email/verify?info=%s", baseurl, info)
}

// SendEmail 发送邮件
// 发送邮件需要配置邮箱服务器信息， 可以在config.yaml中配置
// 以下情况会发生错误: 1. 邮箱配置错误,smtp信息错误 2. 修改模板后,解析模板失败!
func SendEmail(email string, data *EmailData) error {
	config := global.GetConfig().Email
	from := config.Form
	Pass := config.SmtpPass
	User := config.SmtpUser
	to := email
	Host := config.Host
	Port := config.Port

	slog.Info("User: " + User + "Pass " + Pass + "Host " + Host + "Port: " + strconv.Itoa(Port))

	var body bytes.Buffer
	// 解析模版
	template, err := ParseTemplateDir("../assets/templates")
	if err != nil {
		return errors.New("解析模版失败")
	}
	slog.Info("解析模版成功！")

	fmt.Println("URL:", data.URL)
	// 执行模版
	// 把html数据存储在body中， 第二个参数是模板名称， 第三个参数是模板数据（把模板中的占位符换成data数据）
	template.ExecuteTemplate(&body, "email-verify.tpl", &data)
	//为了确保html文件在各个邮件客户端都能正常显示，把html转换成内联模式
	htmlString := body.String()
	prem, _ := premailer.NewPremailerFromString(htmlString, nil)
	htmlline, _ := prem.Transform()

	// 使用 gomail 库发送邮件
	m := gomail.NewMessage()
	slog.Info("准备发送邮件\n")
	// 设定 m 头
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	// 设定html体 内容
	m.SetBody("text/html", htmlline)
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))

	// 配置 SMTP 链接
	d := gomail.NewDialer(Host, Port, User, Pass)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	slog.Info("smtp 连接已建立")
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

// ParseTemplateDir 解析模版目录
func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	// 遍历模版目录，将所有模版文件路径添加到 paths 中
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return template.ParseFiles(paths...)
}
