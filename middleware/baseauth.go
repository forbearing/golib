package middleware

import (
	"github.com/forbearing/gst/config"
	"github.com/gin-gonic/gin"
)

func BaseAuth() gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{
		config.App.Auth.BaseAuthUsername: config.App.Auth.BaseAuthPassword,
	})
}
