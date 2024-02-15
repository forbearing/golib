package controller

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/jwt"
	. "github.com/forbearing/golib/response"
	"github.com/gin-gonic/gin"
)

type portal struct{}

var Portal = new(portal)

func (p *portal) Tick(c *gin.Context) {
	feishuAuthUrl := "https://open.feishu.cn/open-apis/authen/v1/user_auth_page_beta?app_id=%s&redirect_uri="
	targetUrl := "%s/api/feishu/qrlogin"
	// redirectUrl := fmt.Sprintf(feishuAuthUrl, config.App.FeishuConfig.AppID) + url.QueryEscape(fmt.Sprintf(targetUrl, config.App.ServerConfig.Domain))
	redirectUrl := fmt.Sprintf(feishuAuthUrl, config.App.FeishuConfig.AppID) + url.QueryEscape(fmt.Sprintf(targetUrl, config.App.ServerConfig.Domain))
	if config.App.Mode == config.ModeDev {
		redirectUrl = fmt.Sprintf(feishuAuthUrl, config.App.FeishuConfig.AppID) + url.QueryEscape(fmt.Sprintf(targetUrl, "http://172.31.8.8:8001"))
	}
	fmt.Println("============= redirect: ", redirectUrl)

	header := c.Request.Header.Get("Authorization")
	if len(header) == 0 {
		// ResponseJSON(c, CodeNeedLogin)
		ResponseJSON(c, CodeSuccess, gin.H{
			"redirect": redirectUrl,
		})
		return
	}

	// 按空格分割
	items := strings.SplitN(header, " ", 2)
	if len(items) != 2 {
		// ResponseJSON(c, CodeInvalidToken)
		ResponseJSON(c, CodeSuccess, gin.H{
			"redirect": redirectUrl,
		})
		return
	}
	if items[0] != "Bearer" {
		// ResponseJSON(c, CodeInvalidToken)
		ResponseJSON(c, CodeSuccess, gin.H{
			"redirect": redirectUrl,
		})
		return
	}

	// items[1] 是获取到的 tokenString, 我们使用之前定义好的解析 jwt 的函数来解析它
	if _, err := jwt.ParseToken(items[1]); err != nil {
		c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
		ResponseJSON(c, CodeSuccess, gin.H{
			"redirect": redirectUrl,
		})
		return
	}
	ResponseJSON(c, CodeSuccess, gin.H{
		"redirect": config.App.Domain,
	})
}
