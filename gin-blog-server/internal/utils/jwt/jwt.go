package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

var (
	// 错误信息常量，表示不同类型的 Token 错误
	ErrTokenExpired     = errors.New("token 已经过期，请重新登录")  // Token 已过期
	ErrTokenNotValidYet = errors.New("token 无效，请重新登录")    // Token 尚未生效
	ErrTokenMalFormed   = errors.New("token 不正确， 请重新登陆")  // Token 格式错误
	ErrTokenInvalid     = errors.New("这不是一个 token，请重新登录") // Token 无效
)

// MyClaims 自定义的 Claims 结构体，用于保存自定义的 payload 数据
type MyClaims struct {
	UserId               int   `json:"user_id"`  // 用户ID
	RoleIds              []int `json:"role_ids"` // 用户角色ID列表
	jwt.RegisteredClaims       // 内嵌 jwt.RegisteredClaims，包含 JWT 的标准注册字段（如过期时间、签发者等）
}

// GenToken 生成一个新的 JWT Token
// 参数解释：
// secret：用于签名的密钥（通常是一个私钥或者密钥）
// issuer：签发者的标识
// expireHour：Token 过期的小时数
// userId：用户ID
// roleIds：用户的角色ID数组
func GenToken(secret, issuer string, expireHour, userId int, roleIds []int) (string, error) {
	// 创建 MyClaims 实例，填充 JWT 的 Claims 数据
	claims := MyClaims{
		UserId:  userId,  // 设置用户ID
		RoleIds: roleIds, // 设置用户角色ID列表
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,                                                                    // 设置 Token 的签发者
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHour) * time.Hour)), // 设置 Token 的过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                            // 设置 Token 的签发时间
		},
	}

	// 使用 HS256 签名方法创建一个新的 JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 返回签名后的 token 字符串
	return token.SignedString([]byte(secret)) // 使用 secret 对 Token 进行签名并返回
}

// ParseToken 解析 JWT Token 并验证其合法性
// 参数解释：
// secret：用于验证签名的密钥（通常是一个私钥或公钥）
// token：要解析的 JWT Token 字符串
func ParseToken(secret, token string) (*MyClaims, error) {
	// 解析 Token，并将解析出来的 Claims 存入 MyClaims 结构体中
	jwtToken, err := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 使用 secret 来验证 token 的签名
		return []byte(secret), nil
	})

	if err != nil {
		// 错误类型判断
		switch vError, ok := err.(jwt.ValidationError); ok {
		case vError.Errors&jwt.ValidationErrorMalformed != 0:
			// Token 格式错误
			return nil, ErrTokenMalFormed
		case vError.Errors&jwt.ValidationErrorExpired != 0:
			// Token 已经过期
			return nil, ErrTokenExpired
		case vError.Errors&jwt.ValidationErrorNotValidYet != 0:
			// Token 尚未生效
			return nil, ErrTokenNotValidYet
		default:
			// 其他验证错误
			return nil, ErrTokenInvalid
		}
	}

	// 判断 Token 是否有效，如果有效则返回 Claims，否则返回无效的错误
	if claims, ok := jwtToken.Claims.(*MyClaims); ok && jwtToken.Valid {
		// 返回有效的 Claims
		return claims, nil
	}

	// Token 无效，返回错误
	return nil, ErrTokenInvalid
}
