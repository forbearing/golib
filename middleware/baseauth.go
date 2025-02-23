package middleware

import (
	"github.com/forbearing/golib/config"
	"github.com/gin-gonic/gin"
)

func BaseAuth() gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{
		config.App.AuthConfig.BaseAuthUsername: config.App.AuthConfig.BaseAuthPassword,
	})
}
