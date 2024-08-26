package login

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common/pkg/gohttp"
	"net/http"
)

func (s *Service) GetJwtInfoHandler() http.HandlerFunc {
	handlerName := "GetJwtInfoHandler"
	return func(w http.ResponseWriter, r *http.Request) {
		l := s.Log
		gohttp.TraceRequest(handlerName, r, l)
		myJwtInfo := &JwtInfo{
			Authenticated: false,
			Method:        "",
			Message:       "",
			Token:         "",
			Login:         "",
			Email:         "",
			UserId:        0,
		}
		// get the user from the F5 Header UserId
		login := r.Header.Get("UserId")
		if login == "" {
			l.Warn("UserId header missing")
			myJwtInfo.Message = "F5 UserId header missing"
		} else {
			myJwtInfo.Login = login
			myJwtInfo.Method = "F5"
			userInfo, err := s.GetUserInfoFromLogin(login)
			if err != nil {
				const msgErrGetUserInfo = "Error getting user info from login: %v"
				l.Error(msgErrGetUserInfo, err)
				myJwtInfo.Message = fmt.Sprintf(msgErrGetUserInfo, err)
			} else {
				l.Info(fmt.Sprintf("LoginUser(%s) succesfull login for User id (%d)", userInfo.UserLogin, userInfo.UserId))
				token, err := s.JwtChecker.GetTokenFromUserInfo(userInfo)
				if err != nil {
					myJwtInfo.Message = fmt.Sprintf("Error getting token from user info: %v", err)
				}
				myJwtInfo.Authenticated = true
				myJwtInfo.Token = token.String()
				myJwtInfo.Email = userInfo.UserEmail
				myJwtInfo.UserId = userInfo.UserId
			}
		}
		res, err := json.Marshal(myJwtInfo)
		if err != nil {
			l.Error("Error marshalling myJwtInfo: %v", err)
			http.Error(w, "Internal Server Error marshalling myJwtInfo", http.StatusInternalServerError)
			return
		}
		myJwtInfoBytes := [][]byte{[]byte("const go_info="), res, []byte(";")}
		w.Header().Set("Content-Type", "application/javascript")
		w.WriteHeader(http.StatusOK)
		response := bytes.Join(myJwtInfoBytes, []byte(""))
		w.Write(response)
	}

}
