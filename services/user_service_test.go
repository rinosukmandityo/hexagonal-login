package services

import (
	"sync"
	"testing"

	"github.com/rinosukmandityo/hexagonal-login/helper"
	"github.com/rinosukmandityo/hexagonal-login/logic"
	m "github.com/rinosukmandityo/hexagonal-login/models"
	repo "github.com/rinosukmandityo/hexagonal-login/repositories"
)

/*
	==================
	RUN FROM TERMINAL
	==================
	set mongo_url=mongodb://localhost:27017/local
	set mongo_timeout=10
	set mongo_db=local
	set url_db=mongo
*/

var (
	userRepo repo.UserRepository
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
	userRepo = helper.ChooseRepo()
}

func TestInsertUser(t *testing.T) {
	userService := logic.NewUserService(userRepo)

	testdata := UserTestData()
	wg := sync.WaitGroup{}

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

	t.Run("Case 2: Negative Test", func(t *testing.T) {
		t.Run("Case 2.1: Duplicate username", func(t *testing.T) {
			_data := testdata[0]
			_data.ID = "userid04"
			if e := userService.Store(&_data); e == nil {
				t.Error("[ERROR] - duplicate validation username is not working")
			}
		})

		t.Run("Case 2.2: Duplicate ID", func(t *testing.T) {
			_data := testdata[0]
			_data.Username = "username04"
			if e := userService.Store(&_data); e == nil {
				t.Error("[ERROR] - duplicate validation ID is not working")
			}
		})
	})
}

func TestUpdateUser(t *testing.T) {
	testdata := UserTestData()
	userService := logic.NewUserService(userRepo)
	t.Run("Case 1: Update data", func(t *testing.T) {
		_data := testdata[0]
		_data.Username = _data.Username + "UPDATED"
		if e := userService.Update(&_data); e != nil {
			t.Errorf("[ERROR] - Failed to update data %s ", e.Error())
		}
	})
	t.Run("Case 2: Negative Test", func(t *testing.T) {
		_data := m.User{ID: "ID DID NOT EXISTS"}
		if e := userService.Update(&_data); e == nil {
			t.Error("[ERROR] - It should be error 'User Not Found'")
		}
	})
}

func TestDeleteUser(t *testing.T) {
	testdata := UserTestData()
	userService := logic.NewUserService(userRepo)
	t.Run("Case 1: Delete data", func(t *testing.T) {
		_data := testdata[1]
		if e := userService.Delete(&_data); e != nil {
			t.Errorf("[ERROR] - Failed to delete data %s ", e.Error())
		}
	})
	t.Run("Case 2: Negative Test", func(t *testing.T) {
		_data := testdata[1]
		if e := userService.Delete(&_data); e == nil {
			t.Error("[ERROR] - It should be error 'User Not Found'")
		}
	})
}

func TestGetUser(t *testing.T) {
	testdata := UserTestData()
	userService := logic.NewUserService(userRepo)
	t.Run("Case 1: Get data", func(t *testing.T) {
		_data := testdata[0]
		if _, e := userService.GetById(_data.ID); e != nil {
			t.Errorf("[ERROR] - Failed to get data %s ", e.Error())
		}
	})
	t.Run("Case 2: Negative Test", func(t *testing.T) {
		if _, e := userService.GetById("ID DID NOT EXISTS"); e == nil {
			t.Error("[ERROR] - It should be error 'User Not Found'")
		}
	})
}
