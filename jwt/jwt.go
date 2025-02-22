package jwt

import (
	"time"

	"github.com/cockroachdb/errors"
	"github.com/forbearing/golib/config"
	"github.com/golang-jwt/jwt/v5"
)

var (
	secret = []byte("defaultSecret")
	issuer = "golib"
)

type Claims struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenToken 生成 access token
func GenToken(userId string, username string) (string, error) {
	if username == config.App.AuthConfig.NoneExpireUser {
		return config.App.AuthConfig.NoneExpireToken, nil
	}

	// 创建一个我们自己声明的 claims
	claims := Claims{
		userId, username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.App.TokenExpireDuration)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                     // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                                     // 生效时间
			Issuer:    issuer,                                                             // 签发人
		},
	}
	// NewWithClaims 使用指定的签名方法创建签名对象
	// SignedString 使用指定的 secret 签名并获得完整的编码后的字符串 token
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

// GenTokens 生成 access token 和 refresh token
func GenTokens(userId string, username string) (aToken, rToken string, err error) {
	// 创建一个自己的声明
	c := Claims{
		userId,   // 自定义字段, userID
		username, // 自定义字段, username
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.App.TokenExpireDuration)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                     // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                                     // 生效时间
			Issuer:    issuer,                                                             // 签发人
		},
	}

	// 加密并获得完整的编码后的字符串 token
	if aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(secret); err != nil {
		return "", "", err
	}
	// refresh token 不需要任何自定义数据
	// 使用指定的 secret 签名并获得完整的编码后的字符串 token
	if rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 过期时间
		IssuedAt:  jwt.NewNumericDate(time.Now()),                         // 签发时间
		NotBefore: jwt.NewNumericDate(time.Now()),                         // 生效时间
		Issuer:    issuer,                                                 // 签发人
	}).SignedString(secret); err != nil {
		return "", "", err
	}
	return aToken, rToken, nil
}

// RefreshTokens 通过 refresh token 刷新一个新的 AccessToken
func RefreshTokens(aToken, rToken string) (newAToken, newRToken string, err error) {
	// refresh token 无效直接返回
	if _, err = jwt.Parse(rToken, keyFunc); err != nil {
		return
	}

	// 从旧的 access token 解析出 claims 数据
	var claims Claims
	token, err := jwt.ParseWithClaims(aToken, &claims, keyFunc)
	// 检查错误类型
	if err != nil {
		// 检查是否是过期错误
		if errors.Is(err, jwt.ErrTokenExpired) {
			return GenTokens(claims.UserId, claims.Username)
		}
		return
	}
	if !token.Valid {
		return "", "", errors.New("invalid token")
	}
	return GenTokens(claims.UserId, claims.Username)
	// // 如果没有错误，token 仍然有效，不需要刷新
	// return aToken, rToken, nil
}

// ParseToken
func ParseToken(tokenStr string) (*Claims, error) {
	if tokenStr == config.App.AuthConfig.NoneExpireToken {
		return &Claims{
			UserId: "root",
			// 这里必须写成 root 或者 admin, 但是 admin 需要作为普通管理使用,所以这里使用 root
			// 配合 casbin 使用.
			Username:         "root",
			RegisteredClaims: jwt.RegisteredClaims{},
		}, nil
	}

	claims := new(Claims)
	token, err := jwt.ParseWithClaims(tokenStr, claims, keyFunc)
	if err != nil {
		return nil, err
	}
	if token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
func keyFunc(token *jwt.Token) (any, error) { return secret, nil }
