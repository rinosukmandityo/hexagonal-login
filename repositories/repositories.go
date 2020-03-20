package repositories

import (
	m "github.com/rinosukmandityo/hexagonal-login/models"
)

type UserRepository interface {
	GetAll() ([]m.User, error)
	GetBy(filter map[string]interface{}) (*m.User, error)
	Store(data *m.User) error
	Update(data map[string]interface{}, id string) (*m.User, error)
	Delete(id string) error
	Authenticate(username, password string) (bool, *m.User, error)
}
