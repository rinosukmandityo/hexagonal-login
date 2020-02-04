package logic

import (
	m "github.com/rinosukmandityo/hexagonal-login/models"
	repo "github.com/rinosukmandityo/hexagonal-login/repositories"
	svc "github.com/rinosukmandityo/hexagonal-login/services"
)

type loginService struct {
	loginRepo repo.LoginRepository
}

func NewLoginService(loginRepo repo.LoginRepository) svc.LoginService {
	return &loginService{
		loginRepo,
	}
}

func (u *loginService) Authenticate(username, password string) (bool, *m.User, error) {
	return u.loginRepo.Authenticate(username, password)
}

func (u *loginService) ChangePassword(user m.User, newpassword string) error {
	return u.loginRepo.ChangePassword(user, newpassword)

}
func (u *loginService) UserActivation(activationkey string) error {
	return u.loginRepo.UserActivation(activationkey)

}
