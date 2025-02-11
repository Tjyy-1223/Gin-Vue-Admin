package global

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"strings"
)

// Config 结构体包含应用程序的配置信息
type Config struct {
	Server struct {
		Mode          string // 服务器模式 debug | release
		Port          string // 服务器端口
		DbType        string // 数据库类型 mysql | sqlite
		DbAutoMigrate bool   // 是否自动迁移数据库表结构
		DbLogMode     string // 数据库日志模式 silent | error | warn | info
	}
	Log struct {
		Level     string // 日志级别 debug | info | warn | error
		Prefix    string // 日志前缀
		Format    string // 日志格式 text | json
		Directory string // 日志存放目录
	}
	JWT struct {
		Secret string // JWT 密钥
		Expire int64  // JWT 过期时间（小时）
		Issuer string // JWT 签发者
	}
	Mysql struct {
		Host     string // MySQL 服务器地址
		Port     string // MySQL 端口
		Config   string // MySQL 高级配置
		Dbname   string // MySQL 数据库名称
		Username string // MySQL 用户名
		Password string // MySQL 密码
	}
	SQLite struct {
		Dsn string // SQLite 数据源名称（DSN）
	}
	Redis struct {
		DB       int    // Redis 数据库索引
		Addr     string // Redis 服务器地址：端口
		Password string // Redis 密码
	}
	Session struct {
		Name   string // Session 名称
		Salt   string // Session 盐值
		MaxAge int    // Session 最大过期时间（秒）
	}
	Email struct {
		Form     string // 发件人邮箱
		Host     string // SMTP 服务器地址（例如 smtp.qq.com）
		Port     int    // SMTP 端口（一般为 465）
		SmtpPass string // SMTP 密钥（开启SMTP时获取的密钥，而非邮箱密码）
		SmtpUser string // SMTP 用户名（邮箱账号）
	}
	Captcha struct {
		SendEmail  bool // 是否通过邮箱发送验证码
		ExpireTime int  // 验证码过期时间（秒）
	}
	Upload struct {
		Size      int    // 文件上传最大大小（单位：字节）
		OssType   string // OSS 存储类型 local | giniu
		Path      string // 本地文件访问路径
		StorePath string // 本地文件存储路径
	}
	Qiniu struct {
		ImgPath       string // 外链图片地址
		Zone          string // 存储区域
		Bucket        string // 空间名称
		AccessKey     string // 七牛云访问密钥
		SecretKey     string // 七牛云私密密钥
		UseHTTPS      bool   // 是否使用 https 协议
		UseCdnDomains bool   // 是否使用 CDN 上传加速
	}
}

// Conf 存储应用配置的全局变量
var Conf *Config

// GetConfig 返回全局配置，如果未初始化，则触发 panic
func GetConfig() *Config {
	if Conf == nil {
		log.Panic("配置文件未初始化")
		return nil
	}
	return Conf
}

// ReadConfig 读取并解析配置文件，并返回解析后的配置对象
// path: 配置文件的路径
func ReadConfig(path string) *Config {
	v := viper.New()
	v.SetConfigFile(path)                              // 设置配置文件路径
	v.AutomaticEnv()                                   // 允许从环境变量读取配置
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // 转换环境变量格式：SERVER_APP_MODE => SERVER.APPMODE

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		panic("配置文件读取失败： " + err.Error())
	}

	// 反序列化配置文件内容到 Conf 结构体
	if err := v.Unmarshal(&Conf); err != nil {
		panic("配置文件反序列化失败： " + err.Error())
	}

	log.Println("配置文件内容加载成功： ", path)
	return Conf
}

// DbType 返回数据库类型，如果未设置，则默认为 sqlite
func (*Config) DbType() string {
	if Conf.Server.DbType == "" {
		Conf.Server.DbType = "sqlite"
	}
	return Conf.Server.DbType
}

// DbDSN 返回数据库连接字符串（DSN），根据配置中的数据库类型生成不同的连接字符串
func (*Config) DbDSN() string {
	switch Conf.Server.DbType {
	case "mysql":
		// 构造 MySQL 连接字符串
		conf := Conf.Mysql
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?%s",
			conf.Username, conf.Password, conf.Host, conf.Port, conf.Dbname, conf.Config,
		)
	case "sqlite":
		// 返回 SQLite 的 DSN
		return Conf.SQLite.Dsn
	// 默认使用 sqlite，并且使用内存数据库
	default:
		Conf.Server.DbType = "sqlite"
		if Conf.SQLite.Dsn == "" {
			Conf.SQLite.Dsn = "file::memory" // 使用内存数据库
		}
		return Conf.SQLite.Dsn
	}
}
