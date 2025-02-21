package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/database/redis"
	"github.com/forbearing/golib/model"
	"github.com/forbearing/golib/util"
	"github.com/gin-gonic/gin"
	"github.com/mssola/useragent"
	"go.uber.org/zap"
)

// writeCookie 写 cookie 并重定向
func writeCookie(c *gin.Context, token, userId, name string, redirect ...bool) {
	zap.S().Info("writeCookie")
	zap.S().Info("'TokenExpireDuration:' ", config.App.TokenExpireDuration)
	http.SetCookie(c.Writer, &http.Cookie{
		Path:    "/",
		Name:    TOKEN,
		Value:   token,
		Expires: time.Now().Add(config.App.TokenExpireDuration),
	})
	http.SetCookie(c.Writer, &http.Cookie{
		Path:    "/",
		Name:    ID,
		Value:   userId,
		Expires: time.Now().Add(config.App.TokenExpireDuration),
	})
	http.SetCookie(c.Writer, &http.Cookie{
		Path:    "/",
		Name:    NAME,
		Value:   base64.StdEncoding.EncodeToString([]byte(name)), // 中文名,需要转码
		Expires: time.Now().Add(config.App.TokenExpireDuration),
	})
	if len(redirect) > 0 {
		if redirect[0] {
			c.Redirect(http.StatusTemporaryRedirect, config.App.Domain)
		}
	}
}

// writeLocalSessionAndCookie
func writeLocalSessionAndCookie(c *gin.Context, aToken, rToken string, user *model.User) {
	if user == nil {
		zap.S().Info("user is nil")
		return
	}
	name := user.Name
	userId := user.ID
	sessionId := user.SessionId
	zap.S().Info("user.SessionId: ", user.SessionId)
	sessionData, err := json.Marshal(user)
	if err != nil {
		zap.S().Error(err)
		return
	}
	if err := redis.SetSession(sessionId, sessionData); err != nil {
		zap.S().Error(err)
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Path:    "/",
		Name:    TOKEN,
		Value:   aToken,
		Expires: time.Now().Add(config.App.TokenExpireDuration),
	})
	http.SetCookie(c.Writer, &http.Cookie{
		Path:    "/",
		Name:    ACCESS_TOKEN,
		Value:   aToken,
		Expires: time.Now().Add(config.App.TokenExpireDuration),
	})
	http.SetCookie(c.Writer, &http.Cookie{
		Path:  "/",
		Name:  REFRESH_TOKEN,
		Value: rToken,
		// FIXME: refresh token expire duration should defined by config.
		Expires: time.Now().Add(7 * 24 * time.Hour),
	})
	http.SetCookie(c.Writer, &http.Cookie{
		Path:    "/",
		Name:    SESSION_ID,
		Value:   sessionId,
		Expires: time.Now().Add(config.App.TokenExpireDuration),
	})
	http.SetCookie(c.Writer, &http.Cookie{
		Path:    "/",
		Name:    ID,
		Value:   userId,
		Expires: time.Now().Add(config.App.TokenExpireDuration),
	})
	http.SetCookie(c.Writer, &http.Cookie{
		Path:    "/",
		Name:    NAME,
		Value:   base64.StdEncoding.EncodeToString([]byte(name)), // 中文名,需要转码
		Expires: time.Now().Add(config.App.TokenExpireDuration),
	})
}

// writeFeishuSessionAndCookie 写 cookie 并重定向
func writeFeishuSessionAndCookie(c *gin.Context, aToken, rToken string, userInfo *model.UserInfo) {
	if userInfo == nil {
		zap.S().Error("userInfo is nil")
		return
	}
	name := userInfo.Name
	userId := userInfo.UserId
	sessionData, err := json.Marshal(userInfo)
	if err != nil {
		zap.S().Error(err)
		return
	}
	sessionId := util.UUID()
	if err := redis.SetSession(sessionId, sessionData); err != nil {
		zap.S().Error(err)
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Path:    "/",
		Name:    TOKEN,
		Value:   aToken,
		Expires: time.Now().Add(config.App.TokenExpireDuration),
	})
	http.SetCookie(c.Writer, &http.Cookie{
		Path:    "/",
		Name:    ACCESS_TOKEN,
		Value:   aToken,
		Expires: time.Now().Add(config.App.TokenExpireDuration),
	})
	http.SetCookie(c.Writer, &http.Cookie{
		Path:  "/",
		Name:  REFRESH_TOKEN,
		Value: rToken,
		// FIXME: refresh token expire duration should defined by config.
		Expires: time.Now().Add(7 * 24 * time.Hour),
	})
	http.SetCookie(c.Writer, &http.Cookie{
		Path:    "/",
		Name:    SESSION_ID,
		Value:   sessionId,
		Expires: time.Now().Add(config.App.TokenExpireDuration),
	})
	http.SetCookie(c.Writer, &http.Cookie{
		Path:    "/",
		Name:    ID,
		Value:   userId,
		Expires: time.Now().Add(config.App.TokenExpireDuration),
	})
	http.SetCookie(c.Writer, &http.Cookie{
		Path:    "/",
		Name:    NAME,
		Value:   base64.StdEncoding.EncodeToString([]byte(name)), // 中文名,需要转码
		Expires: time.Now().Add(config.App.TokenExpireDuration),
	})
	ua := useragent.New(c.Request.UserAgent())
	engineName, engineVersion := ua.Engine()
	browserName, browserVersion := ua.Browser()
	database.Database[*model.LoginLog]().Create(&model.LoginLog{
		UserID:   userInfo.UserId,
		Username: userInfo.Name,
		Token:    aToken,
		Status:   model.LoginStatusSuccess,
		ClientIP: c.ClientIP(),
		UserAgent: model.UserAgent{
			Source:   c.Request.UserAgent(),
			Platform: fmt.Sprintf("%s %s", ua.Platform(), ua.OS()),
			Engine:   fmt.Sprintf("%s %s", engineName, engineVersion),
			Browser:  fmt.Sprintf("%s %s", browserName, browserVersion),
		},
	})
	domain := config.App.Domain
	if len(util.ParseScheme(c.Request)) > 0 && len(c.Request.Host) > 0 {
		domain = fmt.Sprintf("%s://%s", util.ParseScheme(c.Request), c.Request.Host)
	}
	c.Redirect(http.StatusTemporaryRedirect, domain)
}
