package model

import "github.com/forbearing/golib/model"

type Feishu struct {
	AppID       string
	AppSecretID string

	model.Base
}
