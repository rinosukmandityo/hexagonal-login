// +build user_http

package api_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	. "github.com/rinosukmandityo/hexagonal-login/api"
	"github.com/rinosukmandityo/hexagonal-login/logic"
	m "github.com/rinosukmandityo/hexagonal-login/models"
	repo "github.com/rinosukmandityo/hexagonal-login/repositories"
	rh "github.com/rinosukmandityo/hexagonal-login/repositories/helper"

	"github.com/go-chi/chi"
)

/*
	==================
	RUN FROM TERMINAL
	==================
	go test -v -tags=user_http

	===================================
	TO SET DATABASE INFO FROM TERMINAL
	===================================
	=======
	MongoDB
	=======
	set mongo_url=mongodb://localhost:27017/local
	set mongo_timeout=10
	set mongo_db=local
	set url_db=mongo
	=======
	MySQL
	=======
	set mysql_url=root:Password.1@tcp(127.0.0.1:3306)/tes
	set mysql_timeout=10
	set mysql_db=tes
	set url_db=mysql
*/

var (
	userRepo repo.UserRepository
	r        *chi.Mux
	ts       *httptest.Server
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
	userRepo = rh.ChooseRepo()
	r = RegisterHandler()
}

func TestUserHTTP(t *testing.T) {
	ts = httptest.NewServer(r)
	defer ts.Close()

	t.Run("Insert User", InsertUser)
	t.Run("Update User", UpdateUser)
	t.Run("Delete User", DeleteUser)
	t.Run("Get User", GetDataById)
}

func readUserData(resp *http.Response) (*m.User, error) {
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return nil, e
	}
	defer resp.Body.Close()
	user, _ := GetSerializer(ContentTypeJson).Decode(body)
	return user, nil
}

func PostData(t *testing.T, ts *httptest.Server, url string, _data m.User) error {
	dataBytes, e := getBytes(_data)
	if e != nil {
		return e
	}
	resp, _, e := makeRequest(t, ts, "POST", url, bytes.NewReader(dataBytes))
	if e != nil {
		return e
	}

	switch url {
	case "/":
		if resp.StatusCode != http.StatusCreated {
			return errors.New("status should be 'Status Created' (201)")
		}
	default:
		if resp.StatusCode != http.StatusOK {
			return errors.New("status should be 'Status OK' (200)")
		}
	}

	return nil
}

func GetData(t *testing.T, ts *httptest.Server, url, expected string) error {
	resp, body, e := makeRequest(t, ts, "GET", url, nil)
	if e != nil {
		return e
	}
	if resp.StatusCode != http.StatusFound && strings.Contains(body, expected) {
		return errors.New("status should be 'Status Found' (302)")
	}

	return nil
}

func getBytes(_data m.User) ([]byte, error) {
	dataBytes, e := GetSerializer(ContentTypeJson).Encode(&_data)
	if e != nil {
		return dataBytes, e
	}
	return dataBytes, nil
}

func InsertUser(t *testing.T) {
	userService := logic.NewUserService(userRepo)

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
				if e := PostData(t, ts, "/", _data); e != nil {
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
			if e := PostData(t, ts, "/", _data); e == nil {
				t.Error("[ERROR] - duplicate validation username is not working")
			}
		})

		t.Run("Case 2.2: Duplicate ID", func(t *testing.T) {
			_data := testdata[0]
			_data.Username = "username04"
			if e := PostData(t, ts, "/", _data); e == nil {
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
		if e := PostData(t, ts, "/update", _data); e != nil {
			t.Errorf("[ERROR] - Failed to update data %s ", e.Error())
		}
	})
	t.Run("Case 2: Negative Test", func(t *testing.T) {
		_data := m.User{ID: "ID DID NOT EXISTS"}
		if e := PostData(t, ts, "/update", _data); e == nil {
			t.Error("[ERROR] - It should be error 'User Not Found'")
		}
	})
}

func DeleteUser(t *testing.T) {
	testdata := UserTestData()
	t.Run("Case 1: Delete data", func(t *testing.T) {
		_data := testdata[1]
		if e := PostData(t, ts, "/delete", _data); e != nil {
			t.Errorf("[ERROR] - Failed to delete data %s ", e.Error())
		}
	})
	t.Run("Case 2: Negative Test", func(t *testing.T) {
		_data := testdata[1]
		if e := PostData(t, ts, "/delete", _data); e == nil {
			t.Error("[ERROR] - It should be error 'User Not Found'")
		}
	})
}

func GetDataById(t *testing.T) {
	testdata := UserTestData()
	t.Run("Case 1: Get Data", func(t *testing.T) {
		_data := testdata[0]
		if e := GetData(t, ts, fmt.Sprintf("/%s", _data.ID), _data.ID); e != nil {
			t.Errorf("[ERROR] - Failed to get data %s", e.Error())
		}
	})
	t.Run("Case 2: Negative Test", func(t *testing.T) {
		if e := GetData(t, ts, "/ID-DID-NOT-EXISTS", ""); e == nil {
			t.Error("[ERROR] - It should be error 'Data Not Found'")
		}
	})
}

func makeRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string, error) {
	req, e := http.NewRequest(method, ts.URL+path, body)
	if e != nil {
		return nil, "", e
	}
	req.Header.Set("Content-Type", ContentTypeJson)

	var resp *http.Response
	switch method {
	case "GET":
		resp, e = http.DefaultTransport.RoundTrip(req)
	default:
		resp, e = http.DefaultClient.Do(req)
	}
	if e != nil {
		return nil, "", e
	}

	respBody, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return nil, "", e
	}
	defer resp.Body.Close()

	return resp, string(respBody), nil
}
