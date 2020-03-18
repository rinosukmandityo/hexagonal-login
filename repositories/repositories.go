package repositories

import (
	m "github.com/rinosukmandityo/hexagonal-login/models"
)

type UserRepository interface {
	GetAll() ([]m.User, error)
	GetBy(filter map[string]interface{}) (*m.User, error)
	Store(data *m.User) error
	Update(data *m.User) error
	Delete(data *m.User) error
	Authenticate(username, password string) (bool, *m.User, error)
}
