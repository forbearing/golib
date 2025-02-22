package model

import (
	"strconv"

	"github.com/cockroachdb/errors"
)

type User struct {
	Name         string `json:"name,omitempty"`
	EnName       string `json:"en_name,omitempty"`
	Password     string `json:"password,omitempty"`
	RePassword   string `json:"re_password,omitempty" gorm:"-"`
	NewPassword  string `json:"new_password,omitempty" gorm:"-"`
	Email        string `json:"email,omitempty" gorm:"unique"`
	Avatar       string `json:"avatar,omitempty"`
	AvatarUrl    string `json:"avatar_url,omitempty"`    // 用户头像
	AvatarThumb  string `json:"avatar_thumb,omitempty"`  // 用户头像 72x72
	AvatarMiddle string `json:"avatar_middle,omitempty"` // 用户头像 240x240
	AvatarBig    string `json:"avatar_big,omitempty"`    // 用户头像 640x640
	Mobile       string `json:"mobile,omitempty"`
	Nickname     string `json:"nickname,omitempty"`
	Introduction string `json:"introduction,omitempty"`
	Status       uint   `json:"status,omitempty" gorm:"type:smallint;default:1;comment:status(0: disabled, 1: enabled)"`
	// State 员工状态
	// 1 在职
	// 2 离职
	// 3 试用期
	// 4 实习生
	RoleId       string `json:"role_id,omitempty"`
	DepartmentId string `json:"department_id,omitempty"`

	LastLogin   GormTime `json:"last_login,omitempty"`
	LastLoginIP string   `json:"last_login_ip,omitempty"`
	LockExpire  int64    `json:"lock_expire,omitempty"`
	NumWrong    int      `json:"num_wrong,omitempty" gorm:"comment:the number of input password wrong"`

	Token           string   `json:"token,omitempty" gorm:"-"`
	AccessToken     string   `json:"access_token,omitempty" gorm:"-"`
	RefreshToken    string   `json:"refresh_token,omitempty" gorm:"-"`
	SessionId       string   `json:"session_id,omitempty" gorm:"-"`
	TokenExpiration GormTime `json:"token_expiration,omitempty"`

	Base
}

// CreateBefore check whether user mobile is valid.
func (u *User) CreateBefore() error {
	if len(u.Mobile) != 0 {
		if len(u.Mobile) != 11 {
			return errors.New("mobile number length must be 11")
		}
		if _, err := strconv.Atoi(u.Mobile); err != nil {
			return err
		}
	}
	return nil
}

// GetAfter clean the fields value of Password, RePassword, NewPassword .
func (u *User) GetAfter() error { return u.ListAfter() }
