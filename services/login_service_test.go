// +build login_service

package services_test

import (
	"sync"
	"testing"

	"github.com/rinosukmandityo/hexagonal-login/logic"
	m "github.com/rinosukmandityo/hexagonal-login/models"
	rh "github.com/rinosukmandityo/hexagonal-login/repositories/helper"
	. "github.com/rinosukmandityo/hexagonal-login/services"
)

/*
	==================
	RUN FROM TERMINAL
	==================
	go test -v -tags=login_service

	===================================
	TO SET DATABASE INFO FROM TERMINAL
	===================================
	set mongo_url=mongodb://localhost:27017/local
	set mongo_timeout=10
	set mongo_db=local
	set url_db=mongo
*/

/*
	==================
	RUN FROM TERMINAL
	==================
	go test -v -tags=login_service

	===================================
	TO SET DATABASE INFO FROM TERMINAL
	===================================
	set mongo_url=mongodb://localhost:27017/local
	set mongo_timeout=10
	set mongo_db=local
	set url_db=mongo
*/

var (
	userService  UserService
	loginService LoginService
)

func UserTestData() []m.User {
	return []m.User{{
		Name:     "User 01",
		Username: "username01",
		Password: "Password.1",
		ID:       "userid01",
		Email:    "usermail01@gmail.com",
		Address:  "User Address 01",
		IsActive: false,
	}, {
		Name:     "User 02",
		Username: "username02",
		ID:       "userid02",
		Password: "Password.1",
		Email:    "usermail02@gmail.com",
		Address:  "User Address 02",
		IsActive: false,
	}, {
		Name:     "User 03",
		ID:       "userid03",
		Password: "Password.1",
		Username: "username03",
		Email:    "usermail03@gmail.com",
		Address:  "User Address 03",
		IsActive: false,
	}}
}

func init() {
	loginRepo := rh.ChooseRepo()
	userService = logic.NewUserService(loginRepo)
	loginService = logic.NewLoginService(loginRepo)
}

func TestUserService(t *testing.T) {
	t.Run("Insert User", InsertUser)
	t.Run("Authenticate User", AuthenticateUser)
}

func InsertUser(t *testing.T) {
	testdata := UserTestData()
	wg := sync.WaitGroup{}

	// Clean test data if any
	for _, data := range testdata {
		wg.Add(1)
		go func(_data m.User) {
			userService.Delete(&_data)
			wg.Done()
		}(data)
	}
	wg.Wait()

	t.Run("Case 1: Save data", func(t *testing.T) {
		for _, data := range testdata {
			wg.Add(1)
			go func(_data m.User) {
				if e := userService.Store(&_data); e != nil {
					t.Errorf("[ERROR] - Failed to save data %s ", e.Error())
				}
				wg.Done()
			}(data)
		}
		wg.Wait()

		for _, data := range testdata {
			res, e := userService.GetById(data.ID)
			if e != nil || res.ID == "" {
				t.Errorf("[ERROR] - Failed to get data")
			}
		}
	})
}

func AuthenticateUser(t *testing.T) {
	testdata := UserTestData()
	t.Run("Case 1: Authenticate user", func(t *testing.T) {
		_data := testdata[0]
		if _, _, e := loginService.Authenticate(_data.Username, _data.Password); e != nil {
			t.Errorf("[ERROR] - Failed to authenticate user %s ", e.Error())
		}
	})
	t.Run("Case 2: Negative Test", func(t *testing.T) {
		t.Run("Case 2.1: Username does not exists", func(t *testing.T) {
			if _, _, e := loginService.Authenticate("USERNAME DOES NOT EXISTS", ""); e == nil {
				t.Error("[ERROR] - It should be error 'User Not Found'")
			}
		})
		t.Run("Case 2.2: Password does not match", func(t *testing.T) {
			_data := testdata[0]
			if _, _, e := loginService.Authenticate(_data.Username, "PASSWORD DOES NOT MATCH"); e == nil {
				t.Error("[ERROR] - It should be error 'Password does not match'")
			}
		})
	})
}
