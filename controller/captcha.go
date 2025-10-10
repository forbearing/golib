package controller

import (
	"crypto/rand"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wenlng/go-captcha-assets/resources/imagesv2"
	"github.com/wenlng/go-captcha-assets/resources/tiles"
	"github.com/wenlng/go-captcha/v2/slide"
)

const (
	captchaExpire = 3 * time.Minute // 验证码有效期
	verifyOffset  = 6.0             // 容忍误差像素
	cleanupPeriod = 1 * time.Minute // 定期清理频率
)

type sliderInfo struct {
	AnswerX   float64
	ExpiredAt time.Time
}

var captchaStore = sync.Map{} // key: id, value: sliderInfo

// 生成随机ID
func randomID(n int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var sb strings.Builder
	for range n {
		bigIdx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		sb.WriteByte(chars[bigIdx.Int64()])
	}
	return sb.String()
}

// SliderCaptchaGen 滑块验证码生成接口
func SliderCaptchaGen(c *gin.Context) {
	// 加载内置背景图
	bgImages, err := imagesv2.GetImages()
	if err != nil || len(bgImages) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "load bg images failed"})
		return
	}
	// 加载内置滑块形状
	tileGraphs, err := tiles.GetTiles()
	if err != nil || len(tileGraphs) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "load tiles failed"})
		return
	}
	slideGraphs := make([]*slide.GraphImage, 0, len(tileGraphs))
	for _, g := range tileGraphs {
		slideGraphs = append(slideGraphs, &slide.GraphImage{
			OverlayImage: g.OverlayImage,
			MaskImage:    g.MaskImage,
			ShadowImage:  g.ShadowImage,
		})
	}
	builder := slide.NewBuilder()
	builder.SetResources(
		slide.WithBackgrounds(bgImages),
		slide.WithGraphImages(slideGraphs),
	)
	captcha := builder.Make()

	captData, err := captcha.Generate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "generate captcha failed"})
		return
	}

	// base64图片
	masterBase64, err := captData.GetMasterImage().ToBase64()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "encode master failed"})
		return
	}
	tileBase64, err := captData.GetTileImage().ToBase64()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "encode tile failed"})
		return
	}
	block := captData.GetData() // *slide.Block

	// 生成唯一ID
	captchaID := randomID(18)

	// 存储答案
	captchaStore.Store(captchaID, sliderInfo{
		AnswerX:   float64(block.X),
		ExpiredAt: time.Now().Add(captchaExpire),
	})

	// 返回前端
	c.JSON(http.StatusOK, gin.H{
		"id":         captchaID,
		"master":     "data:image/jpeg;base64," + masterBase64,
		"tile":       "data:image/png;base64," + tileBase64,
		"width":      block.Width,
		"height":     block.Height,
		"tileWidth":  block.DX,
		"tileHeight": block.DY,
		// "answerX": block.X, // 不要返回给前端！
	})
}

// SliderCaptchaVerify 滑块验证码校验接口
func SliderCaptchaVerify(c *gin.Context) {
	var req struct {
		ID     string  `json:"id" binding:"required"`
		BlockX float64 `json:"blockX" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数有误"})
		return
	}
	val, ok := captchaStore.Load(req.ID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码不存在或已过期"})
		return
	}
	//nolint:errcheck
	info := val.(sliderInfo)
	if time.Now().After(info.ExpiredAt) {
		captchaStore.Delete(req.ID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码已过期"})
		return
	}
	// 校验滑块位置
	pass := abs(req.BlockX-info.AnswerX) < verifyOffset
	if pass {
		captchaStore.Delete(req.ID)
		c.JSON(http.StatusOK, gin.H{"verify": true})
	} else {
		c.JSON(http.StatusOK, gin.H{"verify": false})
	}
}

func abs(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}

// // 定时清理过期验证码
// func captchaCleanup() {
// 	for {
// 		time.Sleep(cleanupPeriod)
// 		now := time.Now()
// 		captchaStore.Range(func(key, value interface{}) bool {
// 			info := value.(sliderInfo)
// 			if now.After(info.ExpiredAt) {
// 				captchaStore.Delete(key)
// 			}
// 			return true
// 		})
// 	}
// }
