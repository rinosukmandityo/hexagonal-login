package services

import (
	m "github.com/rinosukmandityo/hexagonal-login/models"
)

type UserService interface {
	GetAll() ([]m.User, error)
	GetById(id string) (*m.User, error)
	Store(user *m.User) error
	Update(user *m.User) error
	Delete(user *m.User) error
}
