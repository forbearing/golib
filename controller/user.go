package controller

import (
	"fmt"
	"time"

	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/jwt"
	"github.com/forbearing/golib/logger"
	"github.com/forbearing/golib/model"
	. "github.com/forbearing/golib/response"
	"github.com/forbearing/golib/types/consts"
	"github.com/forbearing/golib/types/helper"
	"github.com/forbearing/golib/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type user struct{}

var User = new(user)

// Login
func (*user) Login(c *gin.Context) {
	log := logger.Controller.WithControllerContext(helper.NewControllerContext(c), consts.Phase("Login"))

	req := new(model.User)
	if err := c.ShouldBindJSON(req); err != nil {
		log.Error(err)
		ResponseJSON(c, CodeInvalidLogin)
		return
	}

	passwd, err := encryptPasswd(req.Password)
	if err != nil {
		log.Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	users := make([]*model.User, 0)
	if err = database.Database[*model.User](helper.NewDatabaseContext(c)).WithLimit(1).WithQuery(&model.User{Name: req.Name, Password: passwd}).List(&users); err != nil {
		log.Error(err)
		ResponseJSON(c, CodeInvalidLogin)
		return
	}
	if len(users) == 0 {
		log.Error("not found any accounts")
		ResponseJSON(c, CodeInvalidLogin)
		return
	}
	u := users[0]
	if len(u.ID) == 0 {
		log.Error("username or password not match")
		ResponseJSON(c, CodeInvalidLogin)
		return
	}
	aToken, rToken, err := jwt.GenTokens(u.ID, req.Name)
	if err != nil {
		ResponseJSON(c, CodeFailure)
		return
	}
	u.Token = aToken
	u.AccessToken = aToken
	u.RefreshToken = rToken
	u.SessionId = util.UUID()
	fmt.Println("SessionId: ", u.SessionId)
	u.TokenExpiration = model.GormTime(time.Now().Add(config.App.TokenExpireDuration))
	writeLocalSessionAndCookie(c, aToken, rToken, u)
	ResponseJSON(c, CodeSuccess, u)
}

// Signup
func (*user) Signup(c *gin.Context) {
	log := logger.Controller.WithControllerContext(helper.NewControllerContext(c), consts.Phase("Signup"))

	req := new(model.User)
	if err := c.ShouldBindJSON(req); err != nil {
		log.Error(err)
		ResponseJSON(c, CodeInvalidSignup)
		return
	}
	// TODO: check password complexibility.
	if len(req.Password) == 0 {
		log.Error("password length is 0")
		ResponseJSON(c, CodeInvalidSignup)
		return
	}

	if req.Password != req.RePassword {
		log.Error("password and rePassword not the same")
		ResponseJSON(c, CodeInvalidSignup)
		return
	}

	users := make([]*model.User, 0)
	if err := database.Database[*model.User](helper.NewDatabaseContext(c)).WithLimit(1).WithQuery(&model.User{Name: req.Name}).List(&users); err != nil {
		log.Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	if len(users) > 0 {
		if len(users[0].ID) > 0 {
			ResponseJSON(c, CodeAlreadyExistsUser)
			return
		}
	}
	passwd, err := encryptPasswd(req.Password)
	if err != nil {
		log.Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	req.Password = passwd
	req.Status = 1
	req.ID = util.UUID()
	fmt.Println("create user:", req.ID, req.Name, req.Password, req.RePassword)
	if err := database.Database[*model.User](helper.NewDatabaseContext(c)).Create(req); err != nil {
		log.Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	ResponseJSON(c, CodeSuccess)
}

func (*user) ChangePasswd(c *gin.Context) {
	log := logger.Controller.WithControllerContext(helper.NewControllerContext(c), consts.Phase("ChangePasswd"))

	req := new(model.User)
	if err := c.ShouldBindJSON(req); err != nil {
		log.Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	if len(req.ID) == 0 {
		log.Error(CodeNotFoundUserId)
		ResponseJSON(c, CodeNotFoundUserId)
		return
	}
	u := new(model.User)
	if err := database.Database[*model.User](helper.NewDatabaseContext(c)).Get(u, req.ID); err != nil {
		log.Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	if len(u.ID) == 0 {
		log.Error(CodeNotFoundUser)
		ResponseJSON(c, CodeNotFoundUser)
		return
	}
	passwd, err := encryptPasswd(req.Password)
	if err != nil {
		log.Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	if passwd != u.Password {
		log.Error(CodeOldPasswordNotMatch)
		ResponseJSON(c, CodeOldPasswordNotMatch)
		return
	}
	if req.NewPassword != req.RePassword {
		log.Error(CodeNewPasswordNotMatch)
		ResponseJSON(c, CodeNewPasswordNotMatch)
		return
	}
	passwd, err = encryptPasswd(req.NewPassword)
	if err != nil {
		log.Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	u.Password = passwd
	if err := database.Database[*model.User](helper.NewDatabaseContext(c)).Update(u); err != nil {
		log.Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	ResponseJSON(c, CodeSuccess)
}

func encryptPasswd(pass string) (string, error) {
	// hash := sha256.New().Sum([]byte(pass))
	// return hex.EncodeToString(hash)
	hashed, err := bcrypt.GenerateFromPassword([]byte(pass), 8)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}
