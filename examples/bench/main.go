package main

import (
	"fmt"
	"os"

	"github.com/forbearing/golib/bootstrap"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/model"
	"github.com/forbearing/golib/router"
	. "github.com/forbearing/golib/util"
)

var (
	done      struct{}
	fakeToken = "-"
)

func init() {
	model.Register[*User]()
	model.Register[*Group]()
}

/*
List
hey -n 100000 -c 50 -m GET -H "Content-Type: application/json" -H "Authorization: Bearer -" 'http://localhost:8002/api/group'


Batch Create
hey -n 10000 -c 100 -m POST \
-H "Content-Type: application/json" \
-H "Authorization: Bearer -" \
-d '{
    "items": [
        {
            "name": "group01",
            "member_count": 1
        },
        {
            "name": "group02",
            "member_count": 2
        }
    ],
    "options":{
        "atomic": false
    }
}' \
'http://localhost:8002/api/group/batch'

*/

func main() {
	os.Setenv(config.SERVER_PORT, "8002")
	os.Setenv(config.LOGGER_DIR, "/tmp/logs")
	os.Setenv(config.AUTH_NONE_EXPIRE_TOKEN, fakeToken)
	os.Setenv(config.DATABASE_TYPE, string(config.DBMySQL))
	os.Setenv(config.MYSQL_DATABASE, "bench")
	os.Setenv(config.MYSQL_USERNAME, "bench")
	os.Setenv(config.MYSQL_PASSWORD, "bench")

	RunOrDie(bootstrap.Bootstrap)

	prepareData(100)
	router.Register[*Group](router.API, "group")

	RunOrDie(bootstrap.Run)
}

func prepareData(size int) {
	users := make([]*User, size)
	groups := make([]*Group, size)
	for i := range size {
		users[i] = &User{
			Name:     fmt.Sprintf("user-%d", i),
			EnName:   fmt.Sprintf("user-%d", i),
			Password: fmt.Sprintf("user-%d-passwd", i),
			Email:    fmt.Sprintf("user-%d@gmail.com", i),
			Avatar:   fmt.Sprintf("user-%d-avatar", i),
		}
		groups[i] = &Group{Name: fmt.Sprintf("group-%d", i)}
	}

	userResult := make([]*User, 0)
	if err := database.Database[*User]().WithLimit(-1).List(&userResult); err != nil {
		panic(err)
	}
	if err := database.Database[*User]().WithLimit(-1).WithPurge().Delete(userResult...); err != nil {
		panic(err)
	}
	if err := database.Database[*User]().Update(users...); err != nil {
		panic(err)
	}

	groupResult := make([]*Group, 0)
	if err := database.Database[*Group]().WithLimit(-1).List(&groupResult); err != nil {
		panic(err)
	}
	if err := database.Database[*Group]().WithLimit(-1).WithPurge().Delete(groupResult...); err != nil {
		panic(err)
	}
	if err := database.Database[*Group]().Update(groups...); err != nil {
		panic(err)
	}
}

type User struct {
	Name         string `json:"name,omitempty"`
	EnName       string `json:"en_name,omitempty"`
	Password     string `json:"password,omitempty"`
	Email        string `json:"email,omitempty"`
	Avatar       string `json:"avatar,omitempty"`
	AvatarUrl    string `json:"avatar_url,omitempty"`
	AvatarThumb  string `json:"avatar_thumb,omitempty"`
	AvatarMiddle string `json:"avatar_middle,omitempty"`
	AvatarBig    string `json:"avatar_big,omitempty"`
	Mobile       string `json:"mobile,omitempty"`
	Nickname     string `json:"nickname,omitempty"`
	Introduction string `json:"introduction,omitempty"`
	Status       uint   `json:"status,omitempty"`
	RoleId       string `json:"role_id,omitempty"`
	DepartmentId string `json:"department_id,omitempty"`

	model.Base
}

type Group struct {
	Name        string `json:"name,omitempty" schema:"name" binding:"required"`
	Desc        string `json:"desc,omitempty" schema:"desc"`
	MemberCount int    `json:"member_count" gorm:"default:0"`

	model.Base
}
