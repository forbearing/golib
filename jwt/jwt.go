package jwt

import (
	"time"

	"github.com/cockroachdb/errors"
	"github.com/forbearing/golib/config"
	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenDuration  = 2 * time.Hour
	RefreshTokenDuration = 7 * 24 * time.Hour
	MinUserIDLength      = 1
	MinUsernameLength    = 3
)

var (
	ErrInvalidToken        = errors.New("invalid token")
	ErrInvalidAccessToken  = errors.New("invalid access token")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrTokenExpired        = errors.New("token expired")
	ErrTokenMalformed      = errors.New("token malformed")
	ErrTokenNotValidYet    = errors.New("token not valid yet")
)

type TokenType int

const (
	AccessToken TokenType = iota
	RefreshToken
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
func GenToken(userId string, username string) (token string, err error) {
	if len(userId) < MinUserIDLength || len(username) < MinUsernameLength {
		return "", errors.New("invalid user id or username")
	}
	if username == config.App.AuthConfig.NoneExpireUser {
		return config.App.AuthConfig.NoneExpireToken, nil
	}
	now := time.Now()
	claims := Claims{
		userId, username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(config.App.TokenExpireDuration)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(now),                                     // 签发时间
			NotBefore: jwt.NewNumericDate(now),                                     // 生效时间
			Issuer:    issuer,                                                      // 签发人
			Subject:   userId,
		},
	}
	// NewWithClaims 使用指定的签名方法创建签名对象
	// SignedString 使用指定的 secret 签名并获得完整的编码后的字符串 token
	if token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret); err != nil {
		return "", errors.Wrap(err, "failed to generate access token")
	}
	return token, nil
}

// GenTokens 生成 access token 和 refresh token
func GenTokens(userId string, username string) (aToken, rToken string, err error) {
	if aToken, err = GenToken(userId, username); err != nil {
		return "", "", err
	}
	now := time.Now()
	// refresh token 不需要任何自定义数据
	// 使用指定的 secret 签名并获得完整的编码后的字符串 token
	if rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(RefreshTokenDuration)), // 过期时间
		IssuedAt:  jwt.NewNumericDate(now),                           // 签发时间
		NotBefore: jwt.NewNumericDate(now),                           // 生效时间
		Issuer:    issuer,                                            // 签发人
		Subject:   userId,
	}).SignedString(secret); err != nil {
		return "", "", errors.Wrap(err, "failed to generate refresh token")
	}
	return aToken, rToken, nil
}

// RefreshTokens 通过 refresh token 刷新一个新的 AccessToken
func RefreshTokens(accessToken, refreshToken string) (newAccessToken, newRefreshToken string, err error) {
	// verify refresh token
	refreshClaims := new(Claims)
	token := new(jwt.Token)
	if token, err = jwt.ParseWithClaims(refreshToken, refreshClaims, keyFunc); err != nil {
		return "", "", errors.Wrap(err, ErrInvalidRefreshToken.Error())
	}
	if !token.Valid {
		return "", "", ErrInvalidRefreshToken
	}
	if time.Now().After(refreshClaims.ExpiresAt.Time) {
		return "", "", ErrTokenExpired
	}

	// verify access token
	accessClaims := new(Claims)
	if token, err = jwt.ParseWithClaims(accessToken, accessClaims, keyFunc); err != nil {
		if !errors.Is(err, jwt.ErrTokenExpired) {
			return "", "", errors.Wrap(err, ErrInvalidAccessToken.Error())
		}
	} else if !token.Valid {
		return "", "", ErrInvalidAccessToken
	}
	// verify whether subject is the same
	if refreshClaims.Subject != accessClaims.Subject {
		return "", "", ErrTokenMalformed
	}

	return GenTokens(accessClaims.UserId, accessClaims.Username)
}

// ParseToken
func ParseToken(tokenStr string) (*Claims, error) {
	if len(tokenStr) == 0 {
		return nil, ErrTokenMalformed
	}
	if tokenStr == config.App.AuthConfig.NoneExpireToken {
		return &Claims{
			UserId: "root",
			// 这里必须写成 root 或者 admin, 但是 admin 需要作为普通管理使用,所以这里使用 root
			// 配合 casbin 使用.
			Username:         "root",
			RegisteredClaims: jwt.RegisteredClaims{Issuer: issuer, Subject: "root"},
		}, nil
	}

	claims := new(Claims)
	token, err := jwt.ParseWithClaims(tokenStr, claims, keyFunc)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, ErrTokenExpired
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, ErrTokenNotValidYet
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, ErrTokenMalformed
		default:
			return nil, errors.Wrap(err, "failed to parse token")
		}
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}
	if claims.Issuer != issuer {
		return nil, errors.New("invalid token issuer")
	}
	return claims, nil
}
func keyFunc(token *jwt.Token) (any, error) { return secret, nil }
