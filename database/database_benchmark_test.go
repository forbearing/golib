package database_test

import (
	"testing"
	"time"

	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/model"
)

func BenchmarkDatabaseQuery(b *testing.B) {
	users := make([]*model.User, 0)
	b.ResetTimer()
	for b.Loop() {
		_ = database.Database[*model.User]().WithQuery(&model.User{
			Name:         "user",
			EnName:       "user",
			Password:     "mypass",
			Email:        "user@gmail.com",
			Avatar:       "avatar",
			AvatarUrl:    "avatar_url",
			AvatarThumb:  "avatar_thumb",
			AvatarMiddle: "avatar_middle",
			AvatarBig:    "avatar_big",
			Mobile:       "mobile",
			Nickname:     "nickname",
			Introduction: "introduction",
			Status:       1,
			RoleId:       "role_id",
			DepartmentId: "department_id",
			LastLogin:    model.GormTime(time.Now()),
			LastLoginIP:  "last_login_ip",
			LockExpire:   0,
		}).WithTryRun().List(&users)
	}
}
