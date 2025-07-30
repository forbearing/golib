package model

import "github.com/forbearing/golib/model"

type User struct {
	AppID       string
	AppSecretID string

	model.Base
}
