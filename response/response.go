package response

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/forbearing/gst/types/consts"
	"github.com/forbearing/gst/util"
	"github.com/gin-gonic/gin"
)

// 成功处理和失败处理状态码
const (
	CodeSuccess Code = 0
	CodeFailure Code = -1
)

// 通用状态码
const (
	CodeInvalidParam Code = 1000 + iota
	CodeBadRequest
	CodeInvalidToken
	CodeNeedLogin
	CodeUnauthorized
	CodeNetworkTimeout
	CodeContextTimeout
	CodeTooManyRequests
	CodeNotFound
	CodeForbidden
	CodeAlreadyExist
)

// 业务状态码
const (
	CodeInvalidLogin Code = 2000 + iota
	CodeInvalidSignup
	CodeOldPasswordNotMatch
	CodeNewPasswordNotMatch

	CodeNotFoundQueryID
	CodeNotFoundRouteParam
	CodeNotFoundUser
	CodeNotFoundUserID

	CodeAlreadyExistsUser
	CodeAlreadyExistsRole

	CodeTooLargeFile
)

type codeValue struct {
	Status int
	Msg    string
}

// 原始默认的错误码映射
var defaultCodeValueMap = map[Code]codeValue{
	// 成功处理或失败处理的值.
	CodeSuccess: {http.StatusOK, "success"},
	CodeFailure: {http.StatusBadRequest, "failure"},

	// 通用状态码值
	CodeInvalidParam:    {http.StatusBadRequest, "Invalid parameters provided in the request."},
	CodeBadRequest:      {http.StatusBadRequest, "Malformed or illegal request."},
	CodeInvalidToken:    {http.StatusUnauthorized, "Invalid or expired authentication token."},
	CodeNeedLogin:       {http.StatusUnauthorized, "Authentication required to access the requested resource."},
	CodeUnauthorized:    {http.StatusUnauthorized, "Unauthorized access to the requested resource."},
	CodeNetworkTimeout:  {http.StatusGatewayTimeout, "Network operation timed out."},
	CodeContextTimeout:  {http.StatusGatewayTimeout, "Request context timed out."},
	CodeTooManyRequests: {http.StatusTooManyRequests, "too many requests, please try again later."},
	CodeNotFound:        {http.StatusNotFound, "Requested resource not found."},
	CodeForbidden:       {http.StatusForbidden, "Forbidden: Inadequate privileges for the requested operation."},
	CodeAlreadyExist:    {http.StatusConflict, "Resource already exists."},

	// 业务状态码值
	CodeInvalidLogin:        {http.StatusBadRequest, "invalid username or password"},
	CodeInvalidSignup:       {http.StatusBadRequest, "invalid username or password"},
	CodeOldPasswordNotMatch: {http.StatusBadRequest, "old password not match"},
	CodeNewPasswordNotMatch: {http.StatusBadRequest, "new password not match"},
	CodeNotFoundQueryID:     {http.StatusBadRequest, "not found query parameter 'id'"},
	CodeNotFoundRouteParam:  {http.StatusBadRequest, "not found router param"},
	CodeNotFoundUser:        {http.StatusBadRequest, "not found user"},
	CodeNotFoundUserID:      {http.StatusBadRequest, "not found user id"},
	CodeAlreadyExistsUser:   {http.StatusConflict, "user already exists"},
	CodeAlreadyExistsRole:   {http.StatusConflict, "role already exists"},
	CodeTooLargeFile:        {http.StatusBadRequest, "too large file"},
}

// 用于存储自定义的错误码映射
var customCodeValueMap = make(map[Code]codeValue)

type Code int32

// CodeInstance 表示一个错误码实例，包含自定义的状态和消息
type CodeInstance struct {
	code   Code
	status *int    // 自定义状态码，nil 表示使用默认值
	msg    *string // 自定义消息，nil 表示使用默认值
}

func (r Code) Msg() string {
	// 先查找自定义映射
	if val, ok := customCodeValueMap[r]; ok {
		return val.Msg
	}
	// 再查找默认映射
	if val, ok := defaultCodeValueMap[r]; ok {
		return val.Msg
	}
	return defaultCodeValueMap[CodeFailure].Msg
}

func (r Code) WithStatus(status int) CodeInstance {
	return CodeInstance{
		code:   r,
		status: &status,
		msg:    nil, // 保持原有消息
	}
}

func (r Code) WithErr(err error) CodeInstance {
	msg := err.Error()
	return CodeInstance{
		code:   r,
		status: nil, // 保持原有状态码
		msg:    &msg,
	}
}

func (r Code) WithMsg(msg string) CodeInstance {
	return CodeInstance{
		code:   r,
		status: nil, // 保持原有状态码
		msg:    &msg,
	}
}

func (r Code) Status() int {
	// 先查找自定义映射
	if val, ok := customCodeValueMap[r]; ok {
		return val.Status
	}
	// 再查找默认映射
	if val, ok := defaultCodeValueMap[r]; ok {
		return val.Status
	}
	return http.StatusBadRequest
}

func (r Code) Code() int {
	return int(r)
}

func (ci CodeInstance) Msg() string {
	if ci.msg != nil {
		return *ci.msg
	}
	return ci.code.Msg()
}

func (ci CodeInstance) Status() int {
	if ci.status != nil {
		return *ci.status
	}
	return ci.code.Status()
}

func (ci CodeInstance) Code() int {
	return ci.code.Code()
}

func (ci CodeInstance) WithStatus(status int) CodeInstance {
	return CodeInstance{
		code:   ci.code,
		status: &status,
		msg:    ci.msg,
	}
}

func (ci CodeInstance) WithErr(err error) CodeInstance {
	msg := err.Error()
	return CodeInstance{
		code:   ci.code,
		status: ci.status,
		msg:    &msg,
	}
}

func (ci CodeInstance) WithMsg(msg string) CodeInstance {
	return CodeInstance{
		code:   ci.code,
		status: ci.status,
		msg:    &msg,
	}
}

// Responder 响应接口，统一处理 Code 和 CodeInstance
type Responder interface {
	Msg() string
	Status() int
	Code() int
}

// 确保 Code 和 CodeInstance 都实现了 Responder 接口
var (
	_ Responder = Code(0)
	_ Responder = CodeInstance{}
)

func NewCode(code Code, status int, msg string) Code {
	customCodeValueMap[code] = codeValue{
		Status: status,
		Msg:    msg,
	}
	return code
}

func ResponseJSON(c *gin.Context, responder Responder, data ...any) {
	if len(data) > 0 {
		c.JSON(responder.Status(), gin.H{
			"code":            responder.Code(),
			"msg":             responder.Msg(),
			"data":            data[0],
			consts.REQUEST_ID: c.GetString(consts.REQUEST_ID),
		})
	} else {
		c.JSON(responder.Status(), gin.H{
			"code":            responder.Code(),
			"msg":             responder.Msg(),
			"data":            nil,
			consts.REQUEST_ID: c.GetString(consts.REQUEST_ID),
		})
	}
}

func ResponseBytes(c *gin.Context, responder Responder, data ...[]byte) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("X-cached", "true")
	var dataStr string
	if len(data) > 0 {
		dataStr = fmt.Sprintf(`{"code":%d,"msg":"%s","data":%s,"request_id":"%s"}`, responder.Code(), responder.Msg(), util.BytesToString(data[0]), c.GetString(consts.REQUEST_ID))
	} else {
		dataStr = fmt.Sprintf(`{"code":%d,"msg":"%s","data":"","request_id":"%s"}`, responder.Code(), responder.Msg(), c.GetString(consts.REQUEST_ID))
	}
	c.Writer.WriteHeader(responder.Status())
	_, _ = c.Writer.Write(util.StringToBytes(dataStr))
}

func ResponseBytesList(c *gin.Context, responder Responder, total int64, data ...[]byte) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	var dataStr string
	if len(data) > 0 {
		dataStr = fmt.Sprintf(`{"code":%d,"msg":"%s","data":{"total":%d,"items":%s},"request_id":"%s"}`, responder.Code(), responder.Msg(), total, util.BytesToString(data[0]), c.GetString(consts.REQUEST_ID))
	} else {
		dataStr = fmt.Sprintf(`{"code":%d,"msg":"%s","data":{"total":0,"items":[]},"request_id":"%s"}`, responder.Code(), responder.Msg(), c.GetString(consts.REQUEST_ID))
	}
	c.Writer.WriteHeader(responder.Status())
	_, _ = c.Writer.Write(util.StringToBytes(dataStr))
}

func ResponseTEXT(c *gin.Context, responder Responder, data ...any) {
	if len(data) > 0 {
		c.String(responder.Status(), stringAny(data))
	} else {
		c.String(responder.Status(), "")
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
