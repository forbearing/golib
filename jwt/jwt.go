package jwt

import (
	"errors"
	"time"

	"github.com/forbearing/golib/config"
	"github.com/golang-jwt/jwt/v4"
)

var (
	secret = []byte("hybfkuf")
	issuer = "hybfkuf"
)

type Claims struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenToken
func GenToken(userId string, username string) (string, error) {
	if username == config.App.AuthConfig.NoneExpireUser {
		return config.App.AuthConfig.NoneExpireToken, nil
	}

	// 创建一个我们自己声明的 claims
	c := Claims{
		userId, username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(config.App.TokenExpireDuration).Unix(), // 过期时间
			Issuer:    issuer,                                                // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的 secret 签名并获得完整的编码后的字符串 token
	return token.SignedString(secret)
}

// ParseToken
func ParseToken(tokenStr string) (*Claims, error) {
	if tokenStr == config.App.AuthConfig.NoneExpireToken {
		return &Claims{
			UserId: "root",
			// 这里必须写成 root 或者 admin, 但是 admin 需要作为普通管理使用,所以这里使用 root
			// 配合 casbin 使用.
			Username:       "root",
			StandardClaims: jwt.StandardClaims{},
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
