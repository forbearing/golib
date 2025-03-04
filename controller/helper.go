package controller

import (
	"github.com/forbearing/golib/model"
	"github.com/forbearing/golib/types/consts"
	"github.com/gin-gonic/gin"
	"github.com/mssola/useragent"
)

func CreateSession(c *gin.Context) *model.Session {
	ua := useragent.New(c.Request.UserAgent())
	engineName, engineVersion := ua.Engine()
	browserName, browserVersion := ua.Browser()
	return &model.Session{
		UserId:         c.GetString(consts.CTX_USER_ID),
		Username:       c.GetString(consts.CTX_USERNAME),
		Platform:       ua.Platform(),
		OS:             ua.OS(),
		EngineName:     engineName,
		EngineVersion:  engineVersion,
		BrowserName:    browserName,
		BrowserVersion: browserVersion,
	}
}
