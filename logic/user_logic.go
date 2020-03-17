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
	res, e := u.userRepo.GetBy(map[string]interface{}{"_id": id})
	if e != nil {
		return res, e
	}

	return res, nil

}
func (u *userService) Store(user *m.User) error {
	if e := validate.Validate(user); e != nil {
		return errs.Wrap(helper.ErrUserInvalid, "service.User.Store")
	}
	if user.ID == "" {
		user.ID = shortid.MustGenerate()
	}
	if isFound, _, _ := u.GetByUsername(user.Username); isFound {
		return errs.Wrap(helper.ErrUserNameDuplicate, "service.User.Store")
	}
	user.Password = repo.EncryptPassword(user.Password)
	return u.userRepo.Store(user)

}
func (u *userService) Update(user *m.User) error {
	if e := validate.Validate(user); e != nil {
		return errs.Wrap(helper.ErrUserInvalid, "service.User.Update")
	}
	if user.ID == "" {
		user.ID = shortid.MustGenerate()
	}
	if user.Password != "" {
		user.Password = repo.EncryptPassword(user.Password)
	}
	return u.userRepo.Update(user, map[string]interface{}{"_id": user.ID})

}
func (u *userService) Delete(user *m.User) error {
	if user.ID == "" {
		return errs.Wrap(helper.ErrUserNotFound, "service.User.Delete")
	}
	if e := u.userRepo.Delete(map[string]interface{}{"_id": user.ID}); e != nil {
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
