package jwt

import (
	"github.com/forbearing/golib/database"
	"github.com/forbearing/golib/model"
)

// func setToken(accessToken, refreshToken string, s *model.Session) {
// 	if s == nil {
// 		return
// 	}
// 	accessTokenCache.Add(accessToken, s)
// 	refreshTokenCache.Add(refreshToken, s)
// 	database.Database[*model.Session]().Update(s)
// }
//
// func removeToken(accessToken, refreshToken string) {
// 	sessions := make([]*model.Session, 0)
// 	if len(accessToken) > 0 {
// 		accessTokenCache.Remove(accessToken)
// 		if err := database.Database[*model.Session]().WithLimit(-1).WithQuery(&model.Session{AccessToken: accessToken}).WithSelect("id").List(&sessions); err == nil {
// 			database.Database[*model.Session]().WithPurge().Delete(sessions...)
// 		}
// 	}
//
// 	if len(refreshToken) > 0 {
// 		refreshTokenCache.Remove(refreshToken)
// 		if err := database.Database[*model.Session]().WithLimit(-1).WithQuery(&model.Session{RefreshToken: refreshToken}).WithSelect("id").List(&sessions); err == nil {
// 			database.Database[*model.Session]().WithPurge().Delete(sessions...)
// 		}
// 	}
// }
//
// func GetAccessToken(accessToken string) (*model.Session, bool) {
// 	return accessTokenCache.Get(accessToken)
// }
//
// func GetRefreshToken(refreshToken string) (*model.Session, bool) {
// 	return refreshTokenCache.Get(refreshToken)
// }

func setSession(userId string, s *model.Session) {
	if len(userId) == 0 || s == nil {
		return
	}
	database.Database[*model.Session]().Update(s)
	// sessionCache.Add 必须在 database.Update 之后, 因为它的ID会在 database.Database 之后生成
	sessionCache.Add(userId, s)
}

func GetSession(userId string) (*model.Session, bool) {
	// TODO: database
	return sessionCache.Get(userId)
}

func removeSession(userId string) {
	sessionCache.Remove(userId)
	sessions := make([]*model.Session, 0)
	if err := database.Database[*model.Session]().WithLimit(-1).WithSelect("id").WithQuery(&model.Session{UserId: userId}).List(&sessions); err == nil {
		database.Database[*model.Session]().WithPurge().Delete(sessions...)
	}
}
