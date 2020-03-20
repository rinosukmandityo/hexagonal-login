package logic

import (
	"github.com/rinosukmandityo/hexagonal-login/helper"
	m "github.com/rinosukmandityo/hexagonal-login/models"
	repo "github.com/rinosukmandityo/hexagonal-login/repositories"
	svc "github.com/rinosukmandityo/hexagonal-login/services"

	errs "github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
)

type userService struct {
	userRepo repo.UserRepository
}

func NewUserService(userRepo repo.UserRepository) svc.UserService {
	return &userService{
		userRepo,
	}
}

func (u *userService) GetAll() ([]m.User, error) {
	res := []m.User{}
	if res, e := u.userRepo.GetAll(); e != nil {
		return res, e
	}
	return res, nil
}

func (u *userService) GetById(id string) (*m.User, error) {
	res, e := u.userRepo.GetBy(map[string]interface{}{"ID": id})
	if e != nil {
		return res, e
	}

	return res, nil

}
func (u *userService) Store(data *m.User) error {
	if e := validate.Validate(data); e != nil {
		return errs.Wrap(helper.ErrUserInvalid, "service.User.Store")
	}
	if data.ID == "" {
		data.ID = shortid.MustGenerate()
	}
	if isFound, _, _ := u.GetByUsername(data.Username); isFound {
		return errs.Wrap(helper.ErrUserNameDuplicate, "service.User.Store")
	}
	data.Password = repo.EncryptPassword(data.Password)
	return u.userRepo.Store(data)

}
func (u *userService) Update(data map[string]interface{}, id string) (*m.User, error) {
	user := new(m.User)
	var e error
	if data["ID"].(string) == "" {
		return user, errs.Wrap(helper.ErrUserInvalid, "service.User.Update")
	}
	if data["Password"].(string) != "" {
		data["Password"] = repo.EncryptPassword(data["Password"].(string))
	}
	user, e = u.userRepo.Update(data, id)
	if e != nil {
		return user, errs.Wrap(e, "service.User.Update")
	}
	return user, nil

}
func (u *userService) Delete(id string) error {
	if id == "" {
		return errs.Wrap(helper.ErrUserInvalid, "service.User.Delete")
	}
	if e := u.userRepo.Delete(id); e != nil {
		return e
	}
	return nil

}

func (u *userService) GetByUsername(username string) (bool, *m.User, error) {
	res, e := u.userRepo.GetBy(map[string]interface{}{"Username": username})
	if e != nil {
		return false, res, e
	}

	return true, res, nil
}
