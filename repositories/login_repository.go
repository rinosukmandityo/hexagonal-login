package repositories

import (
	m "github.com/rinosukmandityo/hexagonal-login/models"
)

type LoginRepository interface {
	Authenticate(username, password string) (bool, *m.User, error)
	ChangePassword(user m.User, newpassword string) error
	UserActivation(activationkey string) error
}
