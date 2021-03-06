// +build user_service

package services_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	m "github.com/rinosukmandityo/hexagonal-login/models"
	rh "github.com/rinosukmandityo/hexagonal-login/repositories/helper"
	. "github.com/rinosukmandityo/hexagonal-login/services"
	"github.com/rinosukmandityo/hexagonal-login/services/logic"
)

/*
	==================
	RUN FROM TERMINAL
	==================
	go test -v -tags=user_service

	===================================
	TO SET DATABASE INFO FROM TERMINAL
	===================================
	=======
	MongoDB
	=======
	set url=mongodb://localhost:27017/local
	set timeout=10
	set db=local
	set driver=mongo
	=======
	MySQL
	=======
	set url=root:Password.1@tcp(127.0.0.1:3306)/tes
	set timeout=10
	set db=tes
	set driver=mysql
*/

var (
	userService UserService
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
	userRepo := rh.ChooseRepo()
	userService = logic.NewUserService(userRepo)
}

func TestUserService(t *testing.T) {
	// t.Run("Delete All", DeleteAll)
	t.Run("Insert User", InsertUser)
	t.Run("Update User", UpdateUser)
	t.Run("Delete User", DeleteUser)
	t.Run("Get User", GetUser)
}

func InsertUser(t *testing.T) {
	testdata := UserTestData()
	wg := sync.WaitGroup{}

	// Clean test data if any
	for _, data := range testdata {
		wg.Add(1)
		go func(_data m.User) {
			userService.Delete(_data.ID)
			wg.Done()
		}(data)
	}
	wg.Wait()
	time.Sleep(time.Second * 1)

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

		time.Sleep(time.Second * 1)
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

func UpdateUser(t *testing.T) {
	testdata := UserTestData()
	t.Run("Case 1: Update data", func(t *testing.T) {
		_data := testdata[0]
		_data.Username = _data.Username + "UPDATED"
		_dataMap := _data.GetMapFormat()
		_, e := userService.Update(_dataMap, _data.ID)
		if e != nil {
			t.Errorf("[ERROR] - Failed to update data %s ", e.Error())
		}
	})
	t.Run("Case 2: Negative Test", func(t *testing.T) {
		_data := m.User{ID: "ID DID NOT EXISTS"}
		dataMap := _data.GetMapFormat()
		_, e := userService.Update(dataMap, _data.ID)
		if e == nil {
			t.Error("[ERROR] - It should be error 'User Not Found'")
		}
	})
}

func DeleteUser(t *testing.T) {
	testdata := UserTestData()
	t.Run("Case 1: Delete data", func(t *testing.T) {
		_data := testdata[1]
		if e := userService.Delete(_data.ID); e != nil {
			t.Errorf("[ERROR] - Failed to delete data %s ", e.Error())
		}
	})
	t.Run("Case 2: Negative Test", func(t *testing.T) {
		_data := testdata[1]
		if e := userService.Delete(_data.ID); e == nil {
			t.Error("[ERROR] - It should be error 'User Not Found'")
		}
	})
}

func GetUser(t *testing.T) {
	testdata := UserTestData()
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

func DeleteAll(t *testing.T) {
	for i := 1; i <= 4; i++ {
		userService.Delete(fmt.Sprintf("userid0%d", i))
	}
}
