package client_test

import (
	"errors"
	"fmt"
	"net"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/forbearing/golib/bootstrap"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/model"
	"github.com/forbearing/golib/pkg/client"
	"github.com/forbearing/golib/router"
	"github.com/forbearing/golib/types/consts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	token = "-"
	port  = 8000
	addr2 = fmt.Sprintf("http://localhost:%d/api/user/", port)

	id1     = "user1"
	id2     = "user2"
	id3     = "user3"
	id4     = "user4"
	id5     = "user5"
	name1   = id1
	name2   = id2
	name3   = id3
	name4   = id4
	name5   = id5
	email1  = "user1@gmail.com"
	email2  = "user2@gmail.com"
	email3  = "user3@gmail.com"
	email4  = "user4@gmail.com"
	email5  = "user5@gmail.com"
	avatar1 = "avatar1"
	avatar2 = "avatar2"
	avatar3 = "avatar3"
	avatar4 = "avatar4"
	avatar5 = "avatar5"

	name1Modified   = id1 + "_modified"
	email1Modified  = email1 + "_modified"
	avatar1Modified = avatar1 + "_modified"

	name2Modified   = id2 + "_modified"
	email2Modified  = email2 + "_modified"
	avatar2Modified = avatar2 + "_modified"

	user1 = User{Name: name1, Email: email1, Avatar: avatar1, Base: model.Base{ID: id1}}
	user2 = User{Name: name2, Email: email2, Avatar: avatar2, Base: model.Base{ID: id2}}
	user3 = User{Name: name3, Email: email3, Avatar: avatar3, Base: model.Base{ID: id3}}
	user4 = User{Name: name4, Email: email4, Avatar: avatar4, Base: model.Base{ID: id4}}
	user5 = User{Name: name5, Email: email5, Avatar: avatar5, Base: model.Base{ID: id5}}
)

func startServer() {
	model.Register[*User]()

	os.Setenv(config.DATABASE_TYPE, string(config.DBSqlite))
	os.Setenv(config.SQLITE_IS_MEMORY, "true")
	os.Setenv(config.SERVER_PORT, fmt.Sprintf("%d", port))
	os.Setenv(config.LOGGER_DIR, "/tmp/test_client")
	os.Setenv(config.AUTH_NONE_EXPIRE_TOKEN, token)

	os.Setenv(config.DATABASE_TYPE, string(config.DBMySQL))
	os.Setenv(config.MYSQL_DATABASE, "test")
	os.Setenv(config.MYSQL_USERNAME, "test")
	os.Setenv(config.MYSQL_PASSWORD, "test")

	if err := bootstrap.Bootstrap(); err != nil {
		panic(err)
	}

	go func() {
		router.Register[*User](router.API(), "user", consts.Most, consts.MostBatch)
		if err := bootstrap.Run(); err != nil {
			panic(err)
		}
		os.Exit(0)
	}()
	for {
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			l.Close()
			time.Sleep(1 * time.Second)
			continue
		}
		if errors.Is(err, syscall.EADDRINUSE) {
			break
		}
		panic(err)

	}
}

func Test_Client(t *testing.T) {
	startServer()

	cli, err := client.New(addr2, client.WithToken(token), client.WithQueryPagination(1, 2))
	assert.NoError(t, err)
	fmt.Println(cli.QueryString())
	fmt.Println(cli.RequestURL())

	_, err = cli.Create(user1)
	require.NoError(t, err)
	_, err = cli.Create(user2)
	require.NoError(t, err)
	_, err = cli.Create(user3)
	require.NoError(t, err)
	_, err = cli.Create(user4)
	require.NoError(t, err)
	_, err = cli.Create(user5)
	require.NoError(t, err)

	users := make([]User, 0)
	total := new(int64)
	user := new(User)

	// test List
	t.Run("list", func(t *testing.T) {
		assert.NoError(t, cli.List(&users, total))
		assert.Equal(t, 2, len(users))
		assert.Equal(t, int64(5), *total)
	})
	// test Get
	t.Run("get", func(t *testing.T) {
		assert.NoError(t, cli.Get(id1, user))
		assert.Equal(t, id1, user.ID)
		assert.Equal(t, name1, user.Name)
		assert.Equal(t, email1, user.Email)
		assert.Equal(t, avatar1, user.Avatar)
		assert.NoError(t, cli.Get(id2, user))
		assert.Equal(t, id2, user.ID)
		assert.Equal(t, name2, user.Name)
		assert.Equal(t, email2, user.Email)
		assert.Equal(t, avatar2, user.Avatar)
		assert.NoError(t, cli.Get(id3, user))
		assert.Equal(t, id3, user.ID)
		assert.Equal(t, name3, user.Name)
		assert.Equal(t, email3, user.Email)
		assert.Equal(t, avatar3, user.Avatar)
		assert.NoError(t, cli.Get(id4, user))
		assert.Equal(t, id4, user.ID)
		assert.Equal(t, name4, user.Name)
		assert.Equal(t, email4, user.Email)
		assert.Equal(t, avatar4, user.Avatar)
		assert.NoError(t, cli.Get(id5, user))
		assert.Equal(t, id5, user.ID)
		assert.Equal(t, name5, user.Name)
		assert.Equal(t, email5, user.Email)
		assert.Equal(t, avatar5, user.Avatar)
	})

	// Test Update
	t.Run("update", func(t *testing.T) {
		_, err = cli.Update(&User{Name: name1Modified, Email: email1Modified, Base: model.Base{ID: id1}})
		assert.NoError(t, err)
		assert.NoError(t, cli.Get(id1, user))
		assert.Equal(t, id1, user.ID)
		assert.Equal(t, name1Modified, user.Name)
		assert.Equal(t, email1Modified, user.Email)
		assert.Empty(t, user.Avatar)
	})

	// Test UpdatePartial
	t.Run("update_partial", func(t *testing.T) {
		_, err = cli.UpdatePartial(&User{Avatar: avatar1Modified, Base: model.Base{ID: id1}})
		assert.NoError(t, err)
		assert.NoError(t, cli.Get(id1, user))
		assert.Equal(t, id1, user.ID)
		assert.Equal(t, name1Modified, user.Name)
		assert.Equal(t, email1Modified, user.Email)
		assert.Equal(t, avatar1Modified, user.Avatar)

		_, err = cli.UpdatePartial(&User{Name: name1, Base: model.Base{ID: id1}})
		assert.NoError(t, err)
		assert.NoError(t, cli.Get(id1, user))
		assert.Equal(t, id1, user.ID)
		assert.Equal(t, name1, user.Name)
		assert.Equal(t, email1Modified, user.Email)
		assert.Equal(t, avatar1Modified, user.Avatar)

		_, err = cli.UpdatePartial(&User{Email: email1, Avatar: avatar1, Base: model.Base{ID: id1}})
		assert.NoError(t, err)
		assert.NoError(t, cli.Get(id1, user))
		assert.Equal(t, id1, user.ID)
		assert.Equal(t, name1, user.Name)
		assert.Equal(t, email1, user.Email)
		assert.Equal(t, avatar1, user.Avatar)
	})

	// Test BatchCreate
	t.Run("batch_create", func(t *testing.T) {
		cli, err := client.New(addr2, client.WithToken(token))
		require.NoError(t, err)
		items := make([]User, 0)
		total := *new(int64)

		// 1. delete all resources.
		_, err = cli.BatchDelete([]string{id1, id2, id3, id4, id5})
		require.NoError(t, err)
		_, err = cli.BatchCreate(user1)
		require.ErrorIs(t, err, client.ErrNotStructSlice)

		// 2.check the number of resources after create.
		err = cli.List(&items, &total)
		require.NoError(t, err)
		require.Equal(t, 0, len(items))
		require.Equal(t, int64(0), total)

		// 3.create resources.
		_, err = cli.BatchCreate([]User{user1, user2, user3, user4, user5})
		require.NoError(t, err)

		// 4.check the number of resources after create.
		err = cli.List(&items, &total)
		require.NoError(t, err)
		require.Equal(t, 5, len(items))
		require.Equal(t, int64(5), total)
	})

	// Test BatchDelete
	t.Run("batch_delete", func(t *testing.T) {
		cli, err := client.New(addr2, client.WithToken(token))
		require.NoError(t, err)
		items := make([]User, 0)
		total := *new(int64)

		// 1.create resources.
		_, err = cli.BatchUpdate([]User{user1, user2, user3, user4, user5})
		require.NoError(t, err)

		// 2.check the number of resources after create.
		err = cli.List(&items, &total)
		require.NoError(t, err)
		require.Equal(t, 5, len(items))
		require.Equal(t, int64(5), total)

		// 3.delete resources
		_, err = cli.BatchDelete([]string{id1, id2, id3, id4, id5})
		require.NoError(t, err)
		_, err = cli.BatchDelete([]int{1})
		require.ErrorIs(t, err, client.ErrNotStringSlice)

		// 4.check the number of resources after delete
		err = cli.List(&items, &total)
		require.NoError(t, err)
		require.Equal(t, 0, len(items))
		require.Equal(t, int64(0), total)
	})

	// Test BatchUpdate
	t.Run("batch_update", func(t *testing.T) {
		cli, err := client.New(addr2, client.WithToken(token))
		require.NoError(t, err)

		// 1.delete all resources
		_, err = cli.BatchDelete([]string{id1, id2, id3, id4, id5})
		require.NoError(t, err)

		// 2.creat all resources
		_, err = cli.BatchCreate([]User{user1, user2, user3, user4, user5})
		require.NoError(t, err)

		// u1 only modified email
		u1 := user1
		u1.Email = email1Modified
		// u2 only modified avator
		u2 := user2
		u2.Avatar = avatar2Modified
		_, err = cli.BatchUpdate([]User{u1, u2})
		require.NoError(t, err)

		u := new(User)
		err = cli.Get(id1, u)
		require.NoError(t, err)
		require.Equal(t, u.Name, user1.Name)
		require.Equal(t, u.Email, email1Modified)
		require.Equal(t, u.Avatar, user1.Avatar)

		err = cli.Get(id2, u)
		require.NoError(t, err)
		require.Equal(t, u.Name, user2.Name)
		require.Equal(t, u.Email, user2.Email)
		require.Equal(t, u.Avatar, avatar2Modified)

		err = cli.Get(id3, u)
		require.NoError(t, err)
		require.Equal(t, u.Name, user3.Name)
		require.Equal(t, u.Email, user3.Email)
		require.Equal(t, u.Avatar, user3.Avatar)

		err = cli.Get(id4, u)
		require.NoError(t, err)
		require.Equal(t, u.Name, user4.Name)
		require.Equal(t, u.Email, user4.Email)
		require.Equal(t, u.Avatar, user4.Avatar)
	})
	t.Run("batch_update_partial", func(t *testing.T) {
		cli, err := client.New(addr2, client.WithToken(token))
		require.NoError(t, err)

		// 1.delete all resources
		_, err = cli.BatchDelete([]string{id1, id2, id3, id4, id5})
		require.NoError(t, err)

		// 2.creat all resources
		_, err = cli.BatchCreate([]User{user1, user2, user3, user4, user5})
		require.NoError(t, err)

		// u1 only modified email
		u1 := &User{Email: email1Modified}
		u1.ID = id1
		// u2 only modified avator
		u2 := &User{Avatar: avatar2Modified}
		u2.ID = id2
		_, err = cli.BatchUpdatePartial([]*User{u1, u2})
		require.NoError(t, err)

		u := new(User)
		err = cli.Get(id1, u)
		require.NoError(t, err)
		require.Equal(t, u.Name, user1.Name)
		require.Equal(t, u.Email, email1Modified)
		require.Equal(t, u.Avatar, user1.Avatar)

		err = cli.Get(id2, u)
		require.NoError(t, err)
		require.Equal(t, u.Name, user2.Name)
		require.Equal(t, u.Email, user2.Email)
		require.Equal(t, u.Avatar, avatar2Modified)

		err = cli.Get(id3, u)
		require.NoError(t, err)
		require.Equal(t, u.Name, user3.Name)
		require.Equal(t, u.Email, user3.Email)
		require.Equal(t, u.Avatar, user3.Avatar)

		err = cli.Get(id4, u)
		require.NoError(t, err)
		require.Equal(t, u.Name, user4.Name)
		require.Equal(t, u.Email, user4.Email)
		require.Equal(t, u.Avatar, user4.Avatar)
	})
}

type User struct {
	Name   string `json:"name,omitempty"`
	Email  string `json:"email,omitempty"`
	Avatar string `json:"avatar,omitempty"`

	model.Base
}

func (u *User) GetTableName() string {
	return "test_users"
}
