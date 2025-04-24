package controller

import (
	"fmt"
	"net/http"
	"regexp"
	"time"
	"unicode"

	"github.com/cockroachdb/errors"
	"github.com/forbearing/golib/authn/jwt"
	"github.com/forbearing/golib/config"
	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/logger"
	"github.com/forbearing/golib/model"
	"github.com/forbearing/golib/response"
	. "github.com/forbearing/golib/response"
	"github.com/forbearing/golib/types/consts"
	"github.com/forbearing/golib/types/helper"
	"github.com/forbearing/golib/util"
	"github.com/gin-gonic/gin"
	cmap "github.com/orcaman/concurrent-map/v2"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
)

var (
	loginRatelimiterMap  = cmap.New[*rate.Limiter]()
	signupRatelimiterMap = cmap.New[*rate.Limiter]()
)

type user struct{}

var User = new(user)

// Login 多次登陆之后，使用先前的 token 会报错 "access token not match"
func (*user) Login(c *gin.Context) {
	log := logger.Controller.WithControllerContext(helper.NewControllerContext(c), consts.Phase("Login"))
	limiter, found := loginRatelimiterMap.Get(c.ClientIP())
	if !found {
		limiter = rate.NewLimiter(rate.Every(1000*time.Millisecond), 10)
		loginRatelimiterMap.Set(c.ClientIP(), limiter)
	}
	if !limiter.Allow() {
		log.Error("too many login requests")
		ResponseJSON(c, response.NewCode(http.StatusTooManyRequests, "too many login requests"))
		return
	}

	req := new(model.User)
	var err error
	if err = c.ShouldBindJSON(req); err != nil {
		log.Error(err)
		ResponseJSON(c, CodeInvalidLogin)
		return
	}
	users := make([]*model.User, 0)
	if err = database.Database[*model.User](helper.NewDatabaseContext(c)).WithLimit(1).WithQuery(&model.User{Name: req.Name}).List(&users); err != nil {
		log.Error(err)
		ResponseJSON(c, CodeInvalidLogin)
		return
	}
	if len(users) == 0 {
		log.Error("not found any accounts")
		ResponseJSON(c, CodeInvalidLogin)
		return
	}
	if len(users) != 1 {
		log.Errorf("too many accounts: %d", len(users))
		ResponseJSON(c, CodeFailure)
		return
	}
	u := users[0]
	if err = comparePasswd(req.Password, u.Password); err != nil {
		log.Errorf("user password not match: %v", err)
		ResponseJSON(c, CodeInvalidLogin)
		return
	}
	// TODO: 把以前的 token 失效掉
	aToken, rToken, err := jwt.GenTokens(u.ID, req.Name, CreateSession(c))
	if err != nil {
		ResponseJSON(c, CodeFailure)
		return
	}
	u.Token = aToken
	u.AccessToken = aToken
	u.RefreshToken = rToken
	u.SessionId = util.UUID()
	fmt.Println("SessionId: ", u.SessionId)
	u.TokenExpiration = model.GormTime(time.Now().Add(config.App.AccessTokenExpireDuration))
	writeLocalSessionAndCookie(c, aToken, rToken, u)
	// WARN: you must clean password before response to user.
	u.Password = ""

	u.LastLogin = model.GormTime(time.Now())
	u.LastLoginIP = util.IPv6ToIPv4(c.ClientIP())
	if err = database.Database[*model.User](helper.NewDatabaseContext(c)).UpdateById(u.ID, "last_login", u.LastLogin); err != nil {
		log.Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	if err = database.Database[*model.User](helper.NewDatabaseContext(c)).UpdateById(u.ID, "last_login_ip", u.LastLoginIP); err != nil {
		log.Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	ResponseJSON(c, CodeSuccess, u)
}

func (*user) Logout(c *gin.Context) {
	log := logger.Controller.WithControllerContext(helper.NewControllerContext(c), consts.Phase("Logout"))
	_, claims, err := jwt.ParseTokenFromHeader(c.Request.Header)
	if err != nil {
		log.Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	jwt.RevokeTokens(claims.Subject)

	ResponseJSON(c, CodeSuccess)
}

func (*user) RefreshToken(c *gin.Context) {
}

func (*user) Signup(c *gin.Context) {
	log := logger.Controller.WithControllerContext(helper.NewControllerContext(c), consts.Phase("Signup"))
	limiter, found := signupRatelimiterMap.Get(c.ClientIP())
	if !found {
		limiter = rate.NewLimiter(rate.Every(1000*time.Millisecond), 1)
		signupRatelimiterMap.Set(c.ClientIP(), limiter)
	}
	if !limiter.Allow() {
		log.Error("too many signup requests")
		ResponseJSON(c, response.NewCode(http.StatusTooManyRequests, "too many signup requests"))
		return
	}

	req := new(model.User)
	var err error
	if err = c.ShouldBindJSON(req); err != nil {
		log.Error(err)
		ResponseJSON(c, CodeInvalidSignup)
		return
	}
	if err = validateUsername(req.Name); err != nil {
		log.Error(err)
		ResponseJSON(c, response.NewCode(http.StatusBadRequest, err.Error()))
		return
	}
	if err = validatePassword(req.Password); err != nil {
		log.Error(err)
		ResponseJSON(c, response.NewCode(http.StatusBadRequest, err.Error()))
		return
	}
	if req.Password != req.RePassword {
		log.Error("password and rePassword not the same")
		ResponseJSON(c, CodeInvalidSignup)
		return
	}

	users := make([]*model.User, 0)
	if err = database.Database[*model.User](helper.NewDatabaseContext(c)).WithLimit(1).WithQuery(&model.User{Name: req.Name}).List(&users); err != nil {
		log.Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	if len(users) > 0 {
		ResponseJSON(c, CodeAlreadyExistsUser)
		return
	}
	hashedPasswd, err := encryptPasswd(req.Password)
	if err != nil {
		log.Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	req.Password = hashedPasswd
	req.Status = 1
	req.ID = util.UUID()
	req.LastLogin = model.GormTime(time.Now())
	req.LastLoginIP = util.IPv6ToIPv4(c.ClientIP())
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
	hashedPasswd, err := encryptPasswd(req.Password)
	if err != nil {
		log.Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	if hashedPasswd != u.Password {
		log.Error(CodeOldPasswordNotMatch)
		ResponseJSON(c, CodeOldPasswordNotMatch)
		return
	}
	if req.NewPassword != req.RePassword {
		log.Error(CodeNewPasswordNotMatch)
		ResponseJSON(c, CodeNewPasswordNotMatch)
		return
	}
	hashedPasswd, err = encryptPasswd(req.NewPassword)
	if err != nil {
		log.Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	u.Password = hashedPasswd
	if err = database.Database[*model.User](helper.NewDatabaseContext(c)).Update(u); err != nil {
		log.Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	_, claims, err := jwt.ParseTokenFromHeader(c.Request.Header)
	if err != nil {
		log.Error(err)
		ResponseJSON(c, CodeFailure)
		return
	}
	jwt.RevokeTokens(claims.Subject)
	ResponseJSON(c, CodeSuccess)
}

func encryptPasswd(pass string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pass), 8)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func comparePasswd(pass string, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pass))
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password too short")
	}
	var (
		hasNumber      = false
		hasLowerCase   = false
		hasUpperCase   = false
		hasSpecialChar = false
	)
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsLower(c):
			hasLowerCase = true
		case unicode.IsUpper(c):
			hasUpperCase = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecialChar = true
		}
	}

	if !hasNumber {
		return errors.New("password must contain at least one number")
	}
	if !hasLowerCase {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasUpperCase {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasSpecialChar {
		return errors.New("password must contain at least one special character")
	}

	// if !hasNumber || !hasLowerCase || !hasUpperCase || !hasSpecialChar {
	// 	return fmt.Errorf("password too weak")
	// }
	return nil
}

func validateUsername(username string) error {
	if len(username) < 3 || len(username) > 32 {
		return fmt.Errorf("username length must be between 3 and 32")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(username) {
		return fmt.Errorf("username can only contain letters, numbers, underscores and hyphens")
	}
	return nil
}
