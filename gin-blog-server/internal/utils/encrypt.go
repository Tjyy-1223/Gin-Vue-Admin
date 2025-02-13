package utils

import (
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

// BcryptHash 使用 bcrypt 算法将明文字符串进行哈希处理
// 参数:
//
//	str - 需要进行哈希处理的明文字符串（通常是密码）
//
// 返回:
//   - 哈希后的字符串（加密后的密码）
//   - 如果哈希过程中发生错误，返回错误信息
func BcryptHash(str string) (string, error) {
	// 使用 bcrypt.GenerateFromPassword 方法生成哈希值
	// bcrypt.DefaultCost 是默认的加密强度，通常是10，值越大表示计算越复杂
	bytes, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)

	// 将生成的字节切片转换为字符串，并返回哈希值和可能的错误
	return string(bytes), err
}

// BcryptCheck 用于检查明文密码和 bcrypt 哈希值是否匹配
// 参数:
//
//	plain - 明文密码
//	hash - 存储的 bcrypt 哈希密码
//
// 返回:
//   - 如果明文密码与哈希值匹配，返回 true，否则返回 false
func BcryptCheck(plain, hash string) bool {
	// 使用 bcrypt.CompareHashAndPassword 比较明文密码和哈希密码
	// 如果密码匹配，err 会是 nil
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))

	// 如果没有错误，说明密码匹配，返回 true；否则返回 false
	return err == nil
}

// MD5 生成 MD5 哈希值
// 参数:
//
//	str - 需要进行哈希处理的字符串
//	b   - 可选的附加字节，用于对结果进行进一步处理（例如添加盐）
//
// 返回:
//   - 生成的 MD5 哈希值（32 个十六进制字符）
func MD5(str string, b ...byte) string {
	// 创建一个新的 MD5 哈希实例
	h := md5.New()

	// 写入字符串内容
	h.Write([]byte(str))

	// 如果有附加字节，将它们添加到哈希计算中
	return hex.EncodeToString(h.Sum(b))
}
