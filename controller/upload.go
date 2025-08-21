package controller

import (
	"github.com/forbearing/golib/logger"
	"github.com/forbearing/golib/provider/minio"
	. "github.com/forbearing/golib/response"
	"github.com/forbearing/golib/types"
	"github.com/forbearing/golib/types/consts"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type upload struct{}

var Upload = new(upload)

func (*upload) Put(c *gin.Context) {
	log := logger.Controller.WithControllerContext(types.NewControllerContext(c), consts.Phase("Put"))

	// NOTE:字段为 file 必须和前端协商好.
	file, err := c.FormFile("file")
	if err != nil {
		log.Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	// check file size.
	if file.Size > MAX_UPLOAD_SIZE {
		log.Error(CodeTooLargeFile)
		ResponseJSON(c, CodeTooLargeFile)
		return
	}
	fd, err := file.Open()
	if err != nil {
		log.Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	defer fd.Close()

	filename, err := minio.Put(fd, file.Size)
	if err != nil {
		zap.S().Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	ResponseJSON(c, CodeSuccess, gin.H{
		"filename": filename,
	})
}

func (*upload) Preview(c *gin.Context) {
	log := logger.Controller.WithControllerContext(types.NewControllerContext(c), consts.Phase("Preview"))
	data, err := minio.Get(c.Param(consts.PARAM_FILE))
	if err != nil {
		log.Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	ResponseDATA(c, data)
}
