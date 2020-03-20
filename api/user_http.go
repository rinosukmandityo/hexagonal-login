package api

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/rinosukmandityo/hexagonal-login/helper"
	m "github.com/rinosukmandityo/hexagonal-login/models"
	svc "github.com/rinosukmandityo/hexagonal-login/services"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

type UserHandler interface {
	UserCtx(http.Handler) http.Handler
	Get(http.ResponseWriter, *http.Request)
	Post(http.ResponseWriter, *http.Request)
	Update(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
}

type userhandler struct {
	userService svc.UserService
}

func NewUserHandler(userService svc.UserService) UserHandler {
	return &userhandler{userService}
}

func (u *userhandler) UserCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		user, e := u.userService.GetById(id)
		if e != nil {
			if errors.Cause(e) == helper.ErrUserNotFound {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (u *userhandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data, ok := ctx.Value("user").(*m.User)
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	contentType := r.Header.Get("Content-Type")
	respBody, e := GetSerializer(contentType).Encode(data)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	SetupResponse(w, contentType, respBody, http.StatusFound)
}

func (u *userhandler) Post(w http.ResponseWriter, r *http.Request) {
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
	if e = u.userService.Store(user); e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	respBody, e := GetSerializer(contentType).Encode(user)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	SetupResponse(w, contentType, respBody, http.StatusCreated)
}

func (u *userhandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	existingData, ok := ctx.Value("user").(*m.User)
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id := existingData.ID
	contentType := r.Header.Get("Content-Type")
	requestBody, e := ioutil.ReadAll(r.Body)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	data, e := GetSerializer(contentType).DecodeMap(requestBody)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	user, e := u.userService.Update(data, id)
	if e != nil {
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

func (u *userhandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	existingData, ok := ctx.Value("user").(*m.User)
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id := existingData.ID
	contentType := r.Header.Get("Content-Type")
	if e := u.userService.Delete(id); e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	respBody, e := GetSerializer(contentType).EncodeMap(map[string]interface{}{"ID": id})
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	SetupResponse(w, contentType, respBody, http.StatusOK)
}
