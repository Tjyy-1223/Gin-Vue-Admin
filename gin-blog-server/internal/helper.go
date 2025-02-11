package ginblog

import (
	"context"
	"gin-blog-server/internal/global"
	"gin-blog-server/internal/model"
	"github.com/glebarez/sqlite"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"log/slog"
	"os"
	"time"
)

// InitLogger 根据配置文件初始化 slog 日志
func InitLogger(conf *global.Config) *slog.Logger {
	// 根据配置文件中的日志级别来设置日志输出的级别
	var level slog.Level
	switch conf.Log.Level {
	case "debug":
		level = slog.LevelDebug // 设置为调试级别
	case "info":
		level = slog.LevelInfo // 设置为信息级别
	case "warn":
		level = slog.LevelWarn // 设置为警告级别
	case "error":
		level = slog.LevelError // 设置为错误级别
	default:
		level = slog.LevelInfo // 默认使用信息级别
	}

	// 配置日志处理器的选项
	option := &slog.HandlerOptions{
		AddSource: false, // 不添加日志的来源信息（例如，调用栈信息）。可以根据需求修改
		Level:     level, // 设置日志的最小级别，低于该级别的日志将被忽略
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// 针对日志中的时间戳进行定制化处理
			if a.Key == slog.TimeKey {
				// 如果是时间字段，将其格式化为字符串，使用 Go 内建的时间格式
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format(time.DateTime)) // 格式化为标准时间格式
				}
			}
			return a // 返回修改后的属性
		},
	}

	// 根据配置中的日志格式选择输出格式
	var handler slog.Handler
	switch conf.Log.Format {
	case "json":
		// 如果配置为 JSON 格式，使用 JSON 格式的日志处理器
		handler = slog.NewJSONHandler(os.Stdout, option)
	case "text":
		// 默认使用文本格式
		fallthrough
	default:
		handler = slog.NewTextHandler(os.Stdout, option)
	}

	// 创建新的日志记录器并将其作为默认日志记录器
	logger := slog.New(handler)
	slog.SetDefault(logger) // 设置默认的日志记录器
	return logger           // 返回日志记录器
}

// InitDatabase 初始化数据库连接并返回一个 GORM DB 实例
// 根据配置文件中的设置来选择数据库类型、连接配置和日志级别。
// 并且在需要时执行数据库自动迁移操作。
func InitDatabase(conf *global.Config) *gorm.DB {
	// 获取配置文件中数据库类型和数据源名称（DSN）
	dbtype := conf.DbType() // 数据库类型（如 mysql 或 sqlite）
	dsn := conf.DbDSN()     // 数据源名称（数据库连接字符串）

	var db *gorm.DB
	var err error

	// 设置日志级别（根据配置文件中的日志模式选择）
	var level logger.LogLevel
	switch conf.Server.DbLogMode {
	case "silent":
		level = logger.Silent // 静默模式，不输出日志
	case "info":
		level = logger.Info // 信息级别日志，输出常规信息
	case "warn":
		level = logger.Warn // 警告级别日志，输出警告信息
	case "error":
		fallthrough
	default:
		level = logger.Error // 错误级别日志，输出错误信息
	}

	// 配置 GORM 的日志和其他参数
	config := &gorm.Config{
		Logger:                                   logger.Default.LogMode(level), // 设置日志级别
		DisableForeignKeyConstraintWhenMigrating: true,                          // 禁用外键约束迁移
		SkipDefaultTransaction:                   true,                          // 禁用默认事务
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 禁用复数表名（例如：User 表名会映射到 `user`）
		},
	}

	// 根据数据库类型选择不同的数据库驱动
	switch dbtype {
	case "mysql":
		// 使用 MySQL 驱动进行数据库连接
		db, err = gorm.Open(mysql.Open(dsn), config)
	case "sqlite":
		// 使用 SQLite 驱动进行数据库连接
		db, err = gorm.Open(sqlite.Open(dsn), config)
	default:
		// 如果数据库类型不支持，打印错误并终止程序
		log.Fatal("不支持的数据库类型: ", dbtype)
	}

	// 检查数据库连接是否成功
	if err != nil {
		log.Fatal("数据库连接失败", err) // 如果连接失败，打印错误并终止程序
	}
	// 如果连接成功，打印数据库类型和连接字符串
	log.Println("数据库连接成功", dbtype, dsn)

	// 如果配置文件中开启了数据库自动迁移功能
	if conf.Server.DbAutoMigrate {
		// 调用 model.MakeMigrate 执行数据库迁移
		if err := model.MakeMigrate(db); err != nil {
			log.Fatal("数据库迁移失败： ", err) // 如果迁移失败，打印错误并终止程序
		}
		// 迁移成功后打印日志
		log.Println("数据库自动迁移成功")
	}

	// 返回 GORM 的数据库实例
	return db
}

// InitRedis 初始化 Redis 客户端并测试连接
func InitRedis(conf *global.Config) *redis.Client {
	// 创建一个 Redis 客户端实例，配置连接参数
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,     // Redis 服务器地址（包括端口），从配置文件中获取
		Password: conf.Redis.Password, // Redis 密码，如果没有设置密码则传空字符串
		DB:       conf.Redis.DB,       // Redis 数据库索引，通常是 0、1、2 等，从配置文件中获取
	})

	// 使用 Ping 命令测试 Redis 连接是否正常
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		// 如果连接失败，打印错误信息并终止程序
		log.Fatal("Redis 连接失败：", err)
	}

	// 如果连接成功，打印成功信息
	log.Println("Redis 连接成功", conf.Redis.Addr, conf.Redis.DB, conf.Redis.Password)

	// 返回 Redis 客户端实例
	return rdb
}
