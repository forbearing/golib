package model

func init() {
	Register[*LoginLog]()
}

type LoginStatus string

const (
	LoginStatusSuccess = "success"
	LoginStatusFailure = "failure"
)

type LoginLog struct {
	UserID   string      `json:"user_id" schema:"user_id"`
	Username string      `json:"username" schema:"username"`
	ClientIP string      `json:"client_ip" schema:"client_ip"`
	Token    string      `json:"token" schema:"token"`
	Status   LoginStatus `json:"status" schema:"status"`

	UserAgent
	Base
}
