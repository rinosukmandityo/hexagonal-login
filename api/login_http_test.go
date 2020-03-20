// +build login_http

package api_test

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	. "github.com/rinosukmandityo/hexagonal-login/api"
	m "github.com/rinosukmandityo/hexagonal-login/models"
	repo "github.com/rinosukmandityo/hexagonal-login/repositories"
	rh "github.com/rinosukmandityo/hexagonal-login/repositories/helper"
	"github.com/rinosukmandityo/hexagonal-login/services/logic"

	"github.com/go-chi/chi"
)

/*
	==================
	RUN FROM TERMINAL
	==================
	go test -v -tags=login_http

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
	t.Run("Authenticate User", AuthenticateUser)
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
			userService.Delete(_data.ID)
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
		t.Run("Case 1.1: Authenticate with user name", func(t *testing.T) {
			_data := testdata[0]
			_data.Email = ""
			if e := PostData(t, ts, "/auth", _data); e != nil {
				t.Errorf("[ERROR] - Failed to authenticate user %s ", e.Error())
			}
		})
		t.Run("Case 1.2: Authenticate with email address", func(t *testing.T) {
			_data := testdata[0]
			_data.Username = ""
			if e := PostData(t, ts, "/auth", _data); e != nil {
				t.Errorf("[ERROR] - Failed to authenticate user %s ", e.Error())
			}
		})
	})
	t.Run("Case 2: Negative Test", func(t *testing.T) {
		t.Run("Case 2.1: Username does not exists", func(t *testing.T) {
			_data := testdata[0]
			_data.Username = "USERNAME DOES NOT EXISTS"
			_data.Email = "EMAIL DOES NOT EXISTS"
			if e := PostData(t, ts, "/auth", _data); e == nil {
				t.Error("[ERROR] - It should be error 'User Not Found'")
			}
		})
		t.Run("Case 2.2: Password does not match", func(t *testing.T) {
			_data := testdata[0]
			_data.Password = "PASSWORD DOES NOT MATCH"
			if e := PostData(t, ts, "/auth", _data); e == nil {
				t.Error("[ERROR] - It should be error 'Password does not match'")
			}
		})
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
