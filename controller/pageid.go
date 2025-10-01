package controller

import (
	"time"

	. "github.com/forbearing/gst/response"

	"github.com/gin-gonic/gin"
)

var PAGE_ID = time.Now().Unix()

func PageID(c *gin.Context) {
	ResponseJSON(c, CodeSuccess, gin.H{"page_id": PAGE_ID})
}
