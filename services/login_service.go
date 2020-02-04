package services

import (
	m "github.com/rinosukmandityo/hexagonal-login/models"
)

type LoginService interface {
	Authenticate(username, password string) (bool, *m.User, error)
	ChangePassword(user m.User, newpassword string) error
	UserActivation(activationkey string) error
}
