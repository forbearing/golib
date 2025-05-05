package response

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/forbearing/golib/types/consts"
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
	CodeTooManyRequests
	CodeNotFound
	CodeForbidden
	CodeAlreadyExist
)

const (
	// 业务状态码
	CodeInvalidLogin Code = 2000 + iota
	CodeInvalidSignup
	CodeOldPasswordNotMatch
	CodeNewPasswordNotMatch

	CodeNotFoundQueryID
	CodeNotFoundRouteID
	CodeNotFoundUser
	CodeNotFoundUserId

	CodeAlreadyExistsUser
	CodeAlreadyExistsRole

	CodeTooLargeFile
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
	CodeNotFoundRouteID:     {http.StatusBadRequest, "not found router 'id'"},
	CodeNotFoundUser:        {http.StatusBadRequest, "not found user"},
	CodeNotFoundUserId:      {http.StatusBadRequest, "not found user id"},
	CodeAlreadyExistsUser:   {http.StatusConflict, "user already exists"},
	CodeAlreadyExistsRole:   {http.StatusConflict, "role already exists"},
	CodeTooLargeFile:        {http.StatusBadRequest, "too large file"},
}

type Code int32

func (r Code) Msg() string {
	val, ok := codeValueMap[r]
	if !ok {
		val.Msg = codeValueMap[CodeFailure].Msg
	}
	return val.Msg
}

func (r Code) WithStatus(status int) Code {
	return NewCode(r, status, r.Msg())
}

func (r Code) WithErr(err error) Code {
	return NewCode(r, r.Status(), err.Error())
}

func (r Code) Status() int {
	val, ok := codeValueMap[r]
	if !ok {
		val.Status = http.StatusBadRequest
	}
	return val.Status
}

func (r Code) Code() int {
	return int(r)
}

func NewCode(code Code, status int, msg string) Code {
	codeValueMap[Code(code)] = codeValue{
		Status: status,
		Msg:    msg,
	}
	return Code(code)
}

func ResponseJSON(c *gin.Context, code Code, data ...any) {
	if len(data) > 0 {
		c.JSON(code.Status(), gin.H{
			"code":            code,
			"msg":             code.Msg(),
			"data":            data[0],
			consts.REQUEST_ID: c.GetString(consts.REQUEST_ID),
		})
	} else {
		c.JSON(code.Status(), gin.H{
			"code":            code,
			"msg":             code.Msg(),
			"data":            nil,
			consts.REQUEST_ID: c.GetString(consts.REQUEST_ID),
		})
	}
}

func ResponseBytes(c *gin.Context, code Code, data ...[]byte) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("X-cached", "true")
	var dataStr string
	if len(data) > 0 {
		dataStr = fmt.Sprintf(`{"code":%d,"msg":"%s","data":%s,"request_id":"%s"}`, code, code.Msg(), util.BytesToString(data[0]), c.GetString(consts.REQUEST_ID))
	} else {
		dataStr = fmt.Sprintf(`{"code":%d,"msg":"%s","data":"","request_id":"%s"}`, code, code.Msg(), c.GetString(consts.REQUEST_ID))
	}
	c.Writer.WriteHeader(code.Status())
	c.Writer.Write(util.StringToBytes(dataStr))
}

func ResponseBytesList(c *gin.Context, code Code, total uint64, data ...[]byte) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	var dataStr string
	if len(data) > 0 {
		dataStr = fmt.Sprintf(`{"code":%d,"msg":"%s","data":{"total":%d,"items":%s},"request_id":"%s"}`, code, code.Msg(), total, util.BytesToString(data[0]), c.GetString(consts.REQUEST_ID))
	} else {
		dataStr = fmt.Sprintf(`{"code":%d,"msg":"%s","data":{"total":0,"items":[]},"request_id":"%s"}`, code, code.Msg(), c.GetString(consts.REQUEST_ID))
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
