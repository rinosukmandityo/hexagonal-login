package repositories

import (
	m "github.com/rinosukmandityo/hexagonal-login/models"
)

type UserRepository interface {
	GetAll() ([]m.User, error)
	GetById(id string) (*m.User, error)
	GetByUsername(username string) (bool, *m.User, error)
	Store(user *m.User) error
	Update(user *m.User) error
	Delete(user *m.User) error
}
