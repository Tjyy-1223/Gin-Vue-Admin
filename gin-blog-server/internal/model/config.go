package model

import "gorm.io/gorm"

type Config struct {
	Model
	Key   string `gorm:"unique;type:varchar(256)" json:"key"`
	Value string `gorm:"type:varchar(256)" json:"value"`
	Desc  string `gorm:"type:varchar(256)" json:"desc"`
}

// GetConfigMap 获取配置信息并返回一个映射（map），键为配置信息的 Key，值为配置的 Value。
// db: 传入的数据库连接对象，用于查询配置数据。
// 返回值：
//
//	map[string]string: 返回一个键值对映射，key 为配置信息的 Key，value 为配置的 Value。
//	error: 如果查询过程中出现错误，返回错误信息。
func GetConfigMap(db *gorm.DB) (map[string]string, error) {
	// 定义一个 Config 类型的切片，用于存储查询到的配置项。
	var configs []Config

	// 使用 GORM 的 Find 方法从数据库中查询所有配置信息，结果存储在 configs 切片中。
	result := db.Find(&configs)

	// 如果查询过程中发生错误，返回空的映射和错误信息。
	if result.Error != nil {
		return nil, result.Error
	}

	// 创建一个空的 map，用于存储配置信息的键值对。
	m := make(map[string]string)

	// 遍历查询到的配置项，将 Key 和 Value 存入到 map 中。
	for _, config := range configs {
		m[config.Key] = config.Value
	}

	// 返回构造好的 map 和 nil 错误。
	return m, nil
}

// CheckConfigMap 检查并更新配置信息。
// db: 传入的数据库连接对象，用于执行数据库操作。
// m: 一个 map，包含配置信息的键值对，其中 Key 是配置信息的名称，Value 是要更新的配置值。
// 返回值：
//
//	error: 如果在更新过程中发生错误，则返回该错误；否则返回 nil。
func CheckConfigMap(db *gorm.DB, m map[string]string) error {
	// 使用数据库事务（Transaction）确保操作的原子性。所有操作都在同一个事务中进行，如果任何操作失败，事务会回滚。
	return db.Transaction(func(tx *gorm.DB) error {
		// 遍历传入的配置 map，对于每一个键值对执行更新操作
		for k, v := range m {
			// 使用 GORM 的 Update 方法更新配置项的值
			// `Model(Config{})` 表示要操作的模型是 Config，`Where("key", k)` 指定查询条件（根据 key 查找配置项），
			// `Update("value", v)` 表示将该配置项的 value 更新为传入的 v 值。
			result := tx.Model(Config{}).Where("key", k).Update("value", v)

			// 如果更新操作发生错误，立即返回错误并终止事务
			if result.Error != nil {
				return result.Error
			}
		}
		// 如果所有更新操作成功，返回 nil，表示事务成功
		return nil
	})
}

// GetConfig 获取 key 对应的配置信息
func GetConfig(db *gorm.DB, key string) string {
	var config Config
	result := db.Where("key", key).First(&config)
	if result.Error != nil {
		return ""
	}
	return config.Value
}

// CheckConfig 更新 Config
func CheckConfig(db *gorm.DB, key, value string) error {
	var config Config
	result := db.Where("key", key).FirstOrCreate(&config)
	if result.Error != nil {
		return result.Error
	}

	config.Value = value
	result = db.Save(&config)

	return result.Error
}

// GetConfigBool 获取配置
func GetConfigBool(db *gorm.DB, key string) bool {
	val := GetConfig(db, key)
	if val == "" {
		return false
	}
	return val == "true"
}
