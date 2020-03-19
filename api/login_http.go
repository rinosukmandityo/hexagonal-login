package api

import (
	"io/ioutil"
	"net/http"

	svc "github.com/rinosukmandityo/hexagonal-login/services"
)

type LoginHandler interface {
	Auth(http.ResponseWriter, *http.Request)
}

type loginhandler struct {
	loginService svc.LoginService
}

func NewLoginHandler(loginService svc.LoginService) LoginHandler {
	return &loginhandler{loginService}
}

func (u *loginhandler) Auth(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	requestBody, e := ioutil.ReadAll(r.Body)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	user, e := GetSerializer(contentType).Decode(requestBody)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	key := user.Username
	if key == "" {
		key = user.Email
	}
	if _, _, e = u.loginService.Authenticate(key, user.Password); e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	respBody, e := GetSerializer(contentType).Encode(user)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	SetupResponse(w, contentType, respBody, http.StatusOK)
}
