package logic

import (
	"errors"

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
	user := new(m.User)
	param := repo.GetParam{
		Tablename: user.TableName(),
		Filter: map[string]interface{}{"$or": []map[string]interface{}{
			{"Username": username},
			{"Email": username},
		}},
		Result: user,
	}
	if e := u.loginRepo.GetBy(param); e != nil {
		return false, user, e
	}
	if !repo.IsPasswordMatch(password, user.Password) {
		return false, user, errors.New("Password does not match")
	}

	return false, user, nil
}
