package client_test

import (
	"fmt"
	"testing"

	"github.com/forbearing/golib/model"
	"github.com/forbearing/golib/pkg/client"
	"github.com/stretchr/testify/assert"
)

var (
	token = "-"
	addr2 = "http://localhost:8002/api/user//"

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

	user1 = User{Name: name1, Email: email1, Avatar: avatar1, Base: model.Base{ID: id1}}
	user2 = User{Name: name2, Email: email2, Avatar: avatar2, Base: model.Base{ID: id2}}
	user3 = User{Name: name3, Email: email3, Avatar: avatar3, Base: model.Base{ID: id3}}
	user4 = User{Name: name4, Email: email4, Avatar: avatar4, Base: model.Base{ID: id4}}
	user5 = User{Name: name5, Email: email5, Avatar: avatar5, Base: model.Base{ID: id5}}
)

func TestClient(t *testing.T) {
	client, err := client.New(addr2, client.WithToken(token), client.WithQueryPagination(1, 2))
	assert.NoError(t, err)
	fmt.Println(client.QueryString())
	fmt.Println(client.RequestURL())

	client.Create(user1)
	client.Create(user2)
	client.Create(user3)
	client.Create(user4)
	client.Create(user5)

	users := make([]User, 0)
	total := new(int64)
	user := new(User)

	// test List
	{
		assert.NoError(t, client.List(&users, total))
		assert.Equal(t, 2, len(users))
		assert.Equal(t, int64(5), *total)
	}

	// test Get
	{
		assert.NoError(t, client.Get(id1, user))
		assert.Equal(t, id1, user.ID)
		assert.Equal(t, name1, user.Name)
		assert.Equal(t, email1, user.Email)
		assert.Equal(t, avatar1, user.Avatar)
		assert.NoError(t, client.Get(id2, user))
		assert.Equal(t, id2, user.ID)
		assert.Equal(t, name2, user.Name)
		assert.Equal(t, email2, user.Email)
		assert.Equal(t, avatar2, user.Avatar)
		assert.NoError(t, client.Get(id3, user))
		assert.Equal(t, id3, user.ID)
		assert.Equal(t, name3, user.Name)
		assert.Equal(t, email3, user.Email)
		assert.Equal(t, avatar3, user.Avatar)
		assert.NoError(t, client.Get(id4, user))
		assert.Equal(t, id4, user.ID)
		assert.Equal(t, name4, user.Name)
		assert.Equal(t, email4, user.Email)
		assert.Equal(t, avatar4, user.Avatar)
		assert.NoError(t, client.Get(id5, user))
		assert.Equal(t, id5, user.ID)
		assert.Equal(t, name5, user.Name)
		assert.Equal(t, email5, user.Email)
		assert.Equal(t, avatar5, user.Avatar)
	}

	// Test Update
	{
		_, err = client.Update(&User{Name: name1Modified, Email: email1Modified, Base: model.Base{ID: id1}})
		assert.NoError(t, err)
		assert.NoError(t, client.Get(id1, user))
		assert.Equal(t, id1, user.ID)
		assert.Equal(t, name1Modified, user.Name)
		assert.Equal(t, email1Modified, user.Email)
		assert.Empty(t, user.Avatar)
	}

	// Test UpdatePartial
	{
		_, err = client.UpdatePartial(&User{Avatar: avatar1Modified, Base: model.Base{ID: id1}})
		assert.NoError(t, err)
		assert.NoError(t, client.Get(id1, user))
		assert.Equal(t, id1, user.ID)
		assert.Equal(t, name1Modified, user.Name)
		assert.Equal(t, email1Modified, user.Email)
		assert.Equal(t, avatar1Modified, user.Avatar)

		_, err = client.UpdatePartial(&User{Name: name1, Base: model.Base{ID: id1}})
		assert.NoError(t, err)
		assert.NoError(t, client.Get(id1, user))
		assert.Equal(t, id1, user.ID)
		assert.Equal(t, name1, user.Name)
		assert.Equal(t, email1Modified, user.Email)
		assert.Equal(t, avatar1Modified, user.Avatar)

		_, err = client.UpdatePartial(&User{Email: email1, Avatar: avatar1, Base: model.Base{ID: id1}})
		assert.NoError(t, err)
		assert.NoError(t, client.Get(id1, user))
		assert.Equal(t, id1, user.ID)
		assert.Equal(t, name1, user.Name)
		assert.Equal(t, email1, user.Email)
		assert.Equal(t, avatar1, user.Avatar)

	}
}

type User struct {
	Name   string `json:"name,omitempty"`
	Email  string `json:"email,omitempty"`
	Avatar string `json:"avatar,omitempty"`

	model.Base
}
