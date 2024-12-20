package model

import (
	"errors"
	"fmt"

	pkgmodel "github.com/forbearing/golib/model"
	"github.com/forbearing/golib/util"
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
	{Name: ValueOf("user01"), Email: ValueOf("user01@example.com")},
	{Name: ValueOf("user02"), Email: ValueOf("user02@example.com")},
	{Name: ValueOf("user03"), Email: ValueOf("user03@example.com")},
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

	Name     *string `json:"name,omitempty"`
	Email    *string `json:"email,omitempty"`
	Avatar   *string `json:"avatar,omitempty"`
	Sunname  *string `json:"sunname,omitempty"`
	Nickname *string `json:"nickname,omitempty"`
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
		u.ID = *u.Name
	}

	return nil
}

// CreateAfter is the hook that executes after create the current User record.
// current hook only invoke by user calling database.Database().Create.
//
// NOTE: database.Database().UpdateByID/Count operations do NOT triggle any hooks.
func (*User) CreateAfter() error {
	fmt.Println("User CreateAfter")
	return nil
}

// UpdateBefore is the hook that executes before update the current User record.
// current hook only invoke by user calling database.Database().Update

// NOTE: database.Database().UpdateByID/Count operations do NOT triggle any hooks.
func (*User) UpdateBefore() error {
	fmt.Println("User UpdateBefore")
	return nil
}

// UpdateBefore is the hook that executes after update the current User record.
// current hook only invoke by user calling database.Database().Update
//
// NOTE: database.Database().UpdateByID/Count operations do NOT triggle any hooks.
func (*User) UpdateAfter() error {
	fmt.Println("User UpdateAfter")
	return nil
}

// DeleteBefore is the hook that executes before delete the current User record.
// current hook only invoke by user calling database.Database().Delete
//
// NOTE: database.Database().UpdateByID/Count operations do NOT triggle any hooks.
func (*User) DeleteBefore() error {
	fmt.Println("User DeleteBefore")
	return nil
}

// DeleteAfter is the hook that executes after delete the current User record.
// current hook only invoke by user calling database.Database().Delete
//
// NOTE: database.Database().UpdateByID/Count operations do NOT triggle any hooks.
func (*User) DeleteAfter() error {
	fmt.Println("User DeleteAfter")
	return nil
}
