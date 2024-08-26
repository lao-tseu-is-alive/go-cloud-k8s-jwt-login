package login

import (
	"errors"
	"fmt"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common/pkg/database"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common/pkg/gohttp"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common/pkg/golog"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-jwt-login/pkg/crypto"
	"net/http"
)

type Service struct {
	Log        golog.MyLogger
	DbConn     database.DB
	Store      Storage
	JwtChecker gohttp.JwtChecker
}

type F5Authenticator interface {
	gohttp.Authentication
	GetJwtInfoHandler() http.HandlerFunc
}

// NewLoginService Function to create an instance of Login Service Authenticator
func NewLoginService(db database.DB, store Storage, jwtCheck gohttp.JwtChecker) F5Authenticator {
	l := jwtCheck.GetLogger()
	return &Service{
		Log:        l,
		DbConn:     db,
		Store:      store,
		JwtChecker: jwtCheck,
	}
}

func (s *Service) AuthenticateUser(user, passwordHash string) bool {
	if !s.Store.Exist(user) {
		s.Log.Warn("User %s does not exist", user)
		return false
	}
	userInfo, err := s.Store.Get(user)
	if err != nil {
		s.Log.Error("Error getting user %s from DB: %v", user, err)
		return false
	}
	if !userInfo.IsActive {
		s.Log.Warn("User %s is not active", user)
		return false
	}
	if userInfo.IsLocked {
		s.Log.Warn("User %s is locked", user)
		return false
	}
	if !crypto.ComparePasswords(passwordHash, userInfo.PasswordHash) {
		s.Log.Warn("Password for user %s is not correct", user)
		return false
	}
	return true
}

func (s *Service) GetUserInfoFromLogin(login string) (*gohttp.UserInfo, error) {
	if !s.Store.Exist(login) {
		msg := fmt.Sprintf(UserDoesNotExist, login)
		s.Log.Warn(msg)
		return nil, errors.New(msg)
	}
	UserInfo := &gohttp.UserInfo{}
	user, err := s.Store.Get(login)
	if err != nil {
		msg := fmt.Sprintf("Error getting user %s from DB: %v", login, err)
		s.Log.Error(msg)
		return nil, errors.New(msg)
	}
	UserInfo.UserId = int(*user.ExternalId)
	UserInfo.UserName = user.Name
	UserInfo.UserEmail = user.Email
	UserInfo.UserLogin = user.Username
	UserInfo.IsAdmin = user.IsAdmin
	return UserInfo, nil
}
