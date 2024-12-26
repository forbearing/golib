package model

import (
	"fmt"

	"github.com/forbearing/golib/database"
	pkgmodel "github.com/forbearing/golib/model"
	"github.com/forbearing/golib/util"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

const PREFIX = "golib"

var (
	ValueOf = util.ValueOf[string]
	Deref   = util.Deref[string]
)

// users without ID: the ID will be generated automatically.
// WARN: Each application boostrap will insert all record again
// - If boostrapped N times, total len(users) * N records will be inserted.
// - This may result in unintended duplicate data.
var users []*User = []*User{
	{Name: ValueOf("user11"), Email: ValueOf("user11@example.com")},
	{Name: ValueOf("user12"), Email: ValueOf("user12@example.com")},
	{Name: ValueOf("user13"), Email: ValueOf("user13@example.com")},
}

// users with predefined ID, the ID will not be generated automatically.
// Multiples boostrap cycles will mantains exactly len(usersWithIds) records in database.
// NOTE: Only field `updated_at` will be updated.
var usersWithIds []*User = []*User{
	{Name: ValueOf("user01_with_id"), Email: ValueOf("user01@example.com"), Base: pkgmodel.Base{ID: "user01_with_id"}},
	{Name: ValueOf("user02_with_id"), Email: ValueOf("user02@example.com"), Base: pkgmodel.Base{ID: "user02_with_id"}},
	{Name: ValueOf("user03_with_id"), Email: ValueOf("user03@example.com"), Base: pkgmodel.Base{ID: "user03_with_id"}},
}

// usersDeferredId defines users without initial ID.
// IDs will be manually assigned later.
var usersDeferredId []*User = []*User{
	{Name: ValueOf("user01_deferred_id"), Email: ValueOf("user01@example.com")},
	{Name: ValueOf("user02_deferred_id"), Email: ValueOf("user02@example.com")},
	{Name: ValueOf("user03_deferred_id"), Email: ValueOf("user03@example.com")},
}

func init() {
	// Assign IDs based on user names and register users
	// Multiples boostrap only insert len(usersDeferredId) records in database.
	for i := range usersDeferredId {
		usersDeferredId[i].ID = *usersDeferredId[i].Name
	}

	// create table `users` automatically
	// Ensure the package `model` will imported directly or indirectly in `main.go`.
	pkgmodel.Register(users...)
	pkgmodel.Register(usersWithIds...)
	pkgmodel.Register(usersDeferredId...)
}

type User struct {
	pkgmodel.Base

	Name     *string `json:"name,omitempty" schema:"name"`
	Email    *string `json:"email,omitempty" schema:"email"`
	Avatar   *string `json:"avatar,omitempty" schema:"avatar"`
	Sunname  *string `json:"sunname,omitempty" schema:"sunname"`
	Nickname *string `json:"nickname,omitempty" schema:"nickname"`
}

// GetTableName specifiecs the table name of the `User` model.
// the method determinates the table name that maps to the `User` model in the database.
func (*User) GetTableName() string {
	return PREFIX + "_" + pkgmodel.GetTableName[*User]()
}

// CreateBefore is a hook that executes before create the current User record.
// current hook only invoke by user calling database.Database().Create.
//
// NOTE: database.Database().UpdateByID/Count operations do NOT triggle any hooks.
func (u *User) CreateBefore() error {
	// ==============================
	// Add your bussiness logic here.
	// ==============================

	// example 1: initial default value
	if len(Deref(u.Name)) == 0 {
		return errors.New("user name is required")
	}
	if len(u.ID) == 0 {
		u.ID = Deref(u.Name)
	}

	return nil
}

// CreateAfter is the hook that executes after create the current User record.
// current hook only invoke by user calling database.Database().Create.
//
// NOTE: database.Database().UpdateByID/Count operations do NOT triggle any hooks.
func (*User) CreateAfter() error {
	// ==============================
	// Add your bussiness logic here.
	// ==============================

	fmt.Println("User CreateAfter")
	return nil
}

// UpdateBefore is the hook that executes before update the current User record.
// current hook only invoke by user calling database.Database().Update

// NOTE: database.Database().UpdateByID/Count operations do NOT triggle any hooks.
func (*User) UpdateBefore() error {
	// ==============================
	// Add your bussiness logic here.
	// ==============================

	fmt.Println("User UpdateBefore")
	return nil
}

// UpdateBefore is the hook that executes after update the current User record.
// current hook only invoke by user calling database.Database().Update
//
// NOTE: database.Database().UpdateByID/Count operations do NOT triggle any hooks.
func (*User) UpdateAfter() error {
	// ==============================
	// Add your bussiness logic here.
	// ==============================

	fmt.Println("User UpdateAfter")
	return nil
}

// DeleteBefore is the hook that executes before delete the current User record.
// current hook only invoke by user calling database.Database().Delete
//
// NOTE: database.Database().UpdateByID/Count operations do NOT triggle any hooks.
func (u *User) DeleteBefore() error {
	// ==============================
	// Add your bussiness logic here.
	// ==============================

	fmt.Println("User DeleteBefore")
	return nil
}

// DeleteAfter is the hook that executes after delete the current User record.
// current hook only invoke by user calling database.Database().Delete
//
// NOTE: database.Database().UpdateByID/Count operations do NOT triggle any hooks.
func (u *User) DeleteAfter() error {
	// ==============================
	// Add your bussiness logic here.
	// ==============================

	// examples: delete all debug logs associated with the current user
	logs := make([]*Log, 0)
	if err := database.Database[*Log]().WithLimit(-1).WithQuery(&Log{Level: LogLevelDebug, SourceID: u.ID}).List(&logs); err != nil {
		return errors.Wrap(err, "failed to list logs")
	}
	if err := database.Database[*Log]().WithLimit(-1).WithPurge().Delete(logs...); err != nil {
		return errors.Wrap(err, "failed to delete logs")
	}

	fmt.Println("User DeleteAfter")
	return nil
}

func (u *User) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if u == nil {
		return nil
	}
	if len(Deref(u.Name)) > 0 {
		enc.AddString("name", Deref(u.Name))
	}
	if len(Deref(u.Email)) > 0 {
		enc.AddString("email", Deref(u.Email))
	}
	if len(Deref(u.Avatar)) > 0 {
		enc.AddString("avatar", Deref(u.Avatar))
	}
	if len(Deref(u.Sunname)) > 0 {
		enc.AddString("sunname", Deref(u.Sunname))
	}
	if len(Deref(u.Nickname)) > 0 {
		enc.AddString("nickname", Deref(u.Nickname))
	}
	return nil
}
