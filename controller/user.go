package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/jwt"
	"github.com/forbearing/golib/model"
	. "github.com/forbearing/golib/response"
	"github.com/forbearing/golib/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type user struct{}

var User = new(user)

// Login
func (*user) Login(c *gin.Context) {
	req := new(model.User)
	if err := c.ShouldBindJSON(req); err != nil {
		zap.S().Error(err)
		ResponseJSON(c, CodeInvalidLogin)
		return
	}

	users := make([]*model.User, 0)
	if err := database.Database[*model.User]().WithLimit(1).WithQuery(&model.User{Name: req.Name, Password: encryptPasswd(req.Password)}).List(&users); err != nil {
		zap.S().Error(err)
		ResponseJSON(c, CodeInvalidLogin)
		return
	}
	if len(users) == 0 {
		zap.S().Error("not found any accounts")
		ResponseJSON(c, CodeInvalidLogin)
		return
	}
	u := users[0]
	if len(u.ID) == 0 {
		zap.S().Error("username or password not equal")
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

// // Login
// func (*user) Login(c *gin.Context) {
// 	req := new(model.User)
// 	if err := c.ShouldBindJSON(req); err != nil {
// 		zap.S().Error(err)
// 		ResponseJSON(c, CodeInvalidLogin)
// 		return
// 	}
//
// 	users := make([]*model.User, 0)
// 	if err := database.Database[*model.User]().WithLimit(1).WithQuery(&model.User{Username: req.Username}).List(&users); err != nil {
// 		zap.S().Error(err)
// 		ResponseJSON(c, CodeInvalidLogin)
// 		return
// 	}
// 	if len(users) == 0 {
// 		zap.S().Error("not found any accounts")
// 		ResponseJSON(c, CodeInvalidLogin)
// 		return
// 	}
// 	u := users[0]
// 	if len(u.ID) == 0 {
// 		zap.S().Error("account id length is zero")
// 		ResponseJSON(c, CodeInvalidLogin)
// 		return
// 	}
// 	if u.Username == req.Username && u.Password == encryptPasswd(req.Password) {
// 		token, err := jwt.GenToken(0, req.Username)
// 		if err != nil {
// 			zap.S().Error(err)
// 			ResponseJSON(c, CodeFailure)
// 			return
// 		}
// 		u.Token = token
// 		// set password to empty.
// 		u.Password = ""
// 		u.RePassword = ""
// 		u.NewPassword = ""
// 		ResponseJSON(c, CodeSuccess, u)
// 		return
// 	} else {
// 		zap.S().Error("username or password not equal")
// 		ResponseJSON(c, CodeInvalidLogin)
// 		return
// 	}
// }

// Signup
func (*user) Signup(c *gin.Context) {
	req := new(model.User)
	if err := c.ShouldBindJSON(req); err != nil {
		zap.S().Error(err)
		ResponseJSON(c, CodeInvalidSignup)
		return
	}
	// TODO: check password complexibility.
	if len(req.Password) == 0 {
		zap.S().Error("password length is 0")
		ResponseJSON(c, CodeInvalidSignup)
		return
	}

	if req.Password != req.RePassword {
		zap.S().Error("password and rePassword not the same")
		ResponseJSON(c, CodeInvalidSignup)
		return
	}

	// // roleId is required
	// if req.RoleID == nil {
	// 	zap.S().Error(CodeRequireRoleId)
	// 	ResponseJSON(c, CodeRequireRoleId)
	// 	return
	// }
	// role := new(model.Role)
	// if err := database.Role().Get(role, req.RoleID); err != nil {
	// 	zap.S().Error(err)
	// 	ResponseJSON(c, CodeFailure)
	// 	return
	// }
	// if len(role.Name) == 0 {
	// 	zap.S().Error(CodeNotFoundRoleId)
	// 	ResponseJSON(c, CodeNotFoundRoleId)
	// 	return
	// }

	users := make([]*model.User, 0)
	if err := database.Database[*model.User]().WithLimit(1).WithQuery(&model.User{Name: req.Name}).List(&users); err != nil {
		zap.S().Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	if len(users) > 0 {
		if len(users[0].ID) > 0 {
			ResponseJSON(c, CodeAlreadyExistsUser)
			return
		}
	}
	req.Password = encryptPasswd(req.Password)
	req.Status = 1
	if err := database.Database[*model.User]().WithDebug().Create(req); err != nil {
		zap.S().Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	ResponseJSON(c, CodeSuccess)
}

func (*user) ChangePasswd(c *gin.Context) {
	req := new(model.User)
	if err := c.ShouldBindJSON(req); err != nil {
		zap.S().Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	if len(req.ID) == 0 {
		zap.S().Error(CodeNotFoundUserId)
		ResponseJSON(c, CodeNotFoundUserId)
		return
	}
	u := new(model.User)
	if err := database.Database[*model.User]().Get(u, req.ID); err != nil {
		zap.S().Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	if len(u.Name) == 0 {
		zap.S().Error(CodeNotFoundUser)
		ResponseJSON(c, CodeNotFoundUser)
		return
	}
	if encryptPasswd(req.Password) != u.Password {
		zap.S().Error(CodeOldPasswordNotMatch)
		ResponseJSON(c, CodeOldPasswordNotMatch)
		return
	}
	if req.NewPassword != req.RePassword {
		zap.S().Error(CodeNewPasswordNotMatch)
		ResponseJSON(c, CodeNewPasswordNotMatch)
		return
	}
	u.Password = encryptPasswd(req.NewPassword)
	if err := database.Database[*model.User]().Update(u); err != nil {
		zap.S().Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	ResponseJSON(c, CodeSuccess)
}

func encryptPasswd(pass string) string {
	hash := sha256.New().Sum([]byte(pass))
	return hex.EncodeToString(hash)
	// hashed, _ := bcrypt.GenerateFromPassword([]byte(pass), 8)
	// return string(hashed)
}
