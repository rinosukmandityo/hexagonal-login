package repositories

import (
	m "github.com/rinosukmandityo/hexagonal-login/models"
)

type UserRepository interface {
	GetAll() ([]m.User, error)
	GetBy(filter map[string]interface{}) (*m.User, error)
	Store(data *m.User) error
	Update(data *m.User, filter map[string]interface{}) error
	Delete(filter map[string]interface{}) error
	Authenticate(username, password string) (bool, *m.User, error)
}
