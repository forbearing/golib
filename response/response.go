package response

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/forbearing/golib/util"
	"github.com/gin-gonic/gin"
)

const (
	// 成功处理和失败处理状态码
	CodeSuccess Code = 0
	CodeFailure Code = -1

	// 通用状态码
	CodeInvalidParam Code = 1000 + iota
	CodeBadRequest
	CodeInvalidToken
	CodeNeedLogin
	CodeUnauthorized
	CodeNetworkTimeout
	CodeContextTimeout
	CodeTooManyRequest
	CodeNotFound
	CodeNoPermission
	CodeAlreadyExist
)

const (
	// 业务状态码
	CodeServerCallbackTimeout Code = 2000 + iota
	CodeInvalidLogin
	CodeInvalidSignup
	CodeOldPasswordNotMatch
	CodeNewPasswordNotMatch

	CodeRequireRoleId

	CodeNotFoundQueryID
	CodeNotFoundRouteID
	CodeNotFoundRoleId
	CodeNotFoundUser
	CodeNotFoundUserId
	CodeNotFoundRoomType
	CodeNotFoundCategory

	CodeAlreadyExistsUser
	CodeAlreadyExistsRole
	CodeAlreadyExistsRoomType
	CodeAlreadyExistsCategory

	CodeTooLargeFile
	CodeNotPNGJPG
)

type codeValue struct {
	Status int
	Msg    string
}

var codeValueMap = map[Code]codeValue{
	// 成功处理或失败处理的值.
	CodeSuccess: {http.StatusOK, "success"},
	CodeFailure: {http.StatusBadRequest, "failure"},

	// 通用状态码值
	CodeInvalidParam:   {http.StatusBadRequest, "invalid query parameter"},
	CodeBadRequest:     {http.StatusBadRequest, "invalid request"},
	CodeInvalidToken:   {http.StatusUnauthorized, "invalid token"},
	CodeNeedLogin:      {http.StatusUnauthorized, "need login"},
	CodeUnauthorized:   {http.StatusUnauthorized, "unauthorized"},
	CodeNetworkTimeout: {http.StatusBadRequest, "network timeout"},
	CodeContextTimeout: {http.StatusBadRequest, "context timeout"},
	CodeTooManyRequest: {http.StatusTooManyRequests, "too many requests"},
	CodeNotFound:       {http.StatusNotFound, "not found"},
	CodeNoPermission:   {http.StatusBadRequest, "no permission"},
	CodeAlreadyExist:   {http.StatusBadRequest, "already exist"},

	// 业务状态码值
	CodeServerCallbackTimeout: {http.StatusOK, "server callback timeout"},
	CodeInvalidLogin:          {http.StatusBadRequest, "invalid username or password"},
	CodeInvalidSignup:         {http.StatusBadRequest, "invalid username or password"},
	CodeOldPasswordNotMatch:   {http.StatusBadRequest, "old password not match"},
	CodeNewPasswordNotMatch:   {http.StatusBadRequest, "new password not match"},
	CodeRequireRoleId:         {http.StatusBadRequest, "role id is required"},

	CodeNotFoundQueryID:  {http.StatusBadRequest, "not found query parameter 'id'"},
	CodeNotFoundRouteID:  {http.StatusBadRequest, "not found router 'id'"},
	CodeNotFoundRoleId:   {http.StatusBadRequest, "not found role id"},
	CodeNotFoundUser:     {http.StatusBadRequest, "not found user"},
	CodeNotFoundUserId:   {http.StatusBadRequest, "not found user id"},
	CodeNotFoundRoomType: {http.StatusBadRequest, "room type not found"},
	CodeNotFoundCategory: {http.StatusBadRequest, "category not found"},

	CodeAlreadyExistsUser:     {http.StatusBadRequest, "user already exists"},
	CodeAlreadyExistsRole:     {http.StatusBadRequest, "role already exists"},
	CodeAlreadyExistsRoomType: {http.StatusBadRequest, "room type already exists"},
	CodeAlreadyExistsCategory: {http.StatusBadRequest, "category already exists"},

	CodeTooLargeFile: {http.StatusBadRequest, "too large file"},
	CodeNotPNGJPG:    {http.StatusBadGateway, "image must be png or jpg"},
}

type Code int32

func (r Code) Msg() string {
	val, ok := codeValueMap[r]
	if !ok {
		val.Msg = codeValueMap[CodeFailure].Msg
	}
	return val.Msg
}
func (r Code) String() string {
	return r.Msg()
}
func (r Code) Status() int {
	val, ok := codeValueMap[r]
	if !ok {
		val.Status = http.StatusBadRequest
	}
	return val.Status
}

func ResponseJSON(c *gin.Context, code Code, data ...any) {
	if len(data) > 0 {
		c.JSON(code.Status(), gin.H{
			"code": code,
			"msg":  code.Msg(),
			"data": data[0],
		})
	} else {
		c.JSON(code.Status(), gin.H{
			"code": code,
			"msg":  code.Msg(),
			"data": "",
		})
	}
}
func ResponseBytes(c *gin.Context, code Code, data ...[]byte) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("X-cached", "true")
	var dataStr string
	if len(data) > 0 {
		dataStr = fmt.Sprintf(`{"code":%d,"msg":"%s",,"data":%s}`, code, code.Msg(), util.BytesToString(data[0]))
	} else {
		dataStr = fmt.Sprintf(`{"code":"%d","msg":"%s"}`, code, code.Msg())
	}
	c.Writer.WriteHeader(code.Status())
	c.Writer.Write(util.StringToBytes(dataStr))
}

func ResponseBytesList(c *gin.Context, code Code, total uint64, data ...[]byte) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	var dataStr string
	if len(data) > 0 {
		dataStr = fmt.Sprintf(`{"code":%d,"msg":"%s","data":{"total":%d,"items":%s}}`, code, code.Msg(), total, util.BytesToString(data[0]))
	} else {
		dataStr = fmt.Sprintf(`{"code":"%d","msg":"%s","data":{"total":0,"items":[]}}`, code, code.Msg())
	}
	c.Writer.WriteHeader(code.Status())
	c.Writer.Write(util.StringToBytes(dataStr))
}

func ResponseTEXT(c *gin.Context, code Code, data ...any) {
	if len(data) > 0 {
		c.String(code.Status(), stringAny(data))
	} else {
		c.String(code.Status(), "")
	}
}
func ResponseDATA(c *gin.Context, data []byte, headers ...map[string]string) {
	header := make(map[string]string)
	if len(headers) > 0 {
		if headers[0] != nil {
			header = headers[0]
		}
	}
	for k, v := range header {
		c.Header(k, v)
	}
	c.Data(http.StatusOK, "application/octet-stream", data)
}
func ResponesFILE(c *gin.Context, filename string) {
	c.File(filename)
}

func stringAny(v any) string {
	if v == nil {
		return ""
	}
	val, ok := v.(fmt.Stringer)
	if ok {
		return val.String()
	}

	switch val := v.(type) {
	case string:
		return val
	case []byte:
		return string(val)
	case []string:
		return strings.Join(val, ",")
	case [][]byte:
		return string(bytes.Join(val, []byte(",")))
	default:
		data, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return string(data)
	}
}
