package logic

import (
	m "github.com/rinosukmandityo/hexagonal-login/models"
	repo "github.com/rinosukmandityo/hexagonal-login/repositories"
	svc "github.com/rinosukmandityo/hexagonal-login/services"
)

type loginService struct {
	loginRepo repo.UserRepository
}

func NewLoginService(loginRepo repo.UserRepository) svc.LoginService {
	return &loginService{
		loginRepo,
	}
}

func (u *loginService) Authenticate(username, password string) (bool, *m.User, error) {
	return u.loginRepo.Authenticate(username, password)
}
