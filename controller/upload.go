package controller

import (
	"github.com/forbearing/golib/minio"
	. "github.com/forbearing/golib/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type upload struct{}

var Upload = new(upload)

func (*upload) Put(c *gin.Context) {
	slog := zap.S()
	// NOTE:字段为 file 必须和前端协商好.
	file, err := c.FormFile("file")
	if err != nil {
		slog.Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	// check file size.
	if file.Size > MAX_UPLOAD_SIZE {
		slog.Error(CodeTooLargeFile)
		ResponseJSON(c, CodeTooLargeFile)
		return
	}
	fd, err := file.Open()
	if err != nil {
		slog.Error(err)
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
	data, err := minio.Get(c.Param(PARAM_FILE))
	if err != nil {
		zap.S().Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	ResponseDATA(c, data)
}
