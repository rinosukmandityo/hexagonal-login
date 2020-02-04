package logic

import (
	"errors"

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

var (
	ErrUserNotFound      = errors.New("User Not Found")
	ErrUserInvalid       = errors.New("User Invalid")
	ErrUserNameDuplicate = errors.New("User Name Already Exists")
)

func NewUserService(userRepo repo.UserRepository) svc.UserService {
	return &userService{
		userRepo,
	}
}

func (u *userService) GetAll() ([]m.User, error) {
	return u.userRepo.GetAll()
}

func (u *userService) GetById(id string) (*m.User, error) {
	return u.userRepo.GetById(id)

}
func (u *userService) Store(user *m.User) error {
	if e := validate.Validate(user); e != nil {
		return errs.Wrap(ErrUserInvalid, "service.User.Store")
	}
	if user.ID == "" {
		user.ID = shortid.MustGenerate()
	}
	if isFound, _, _ := u.userRepo.GetByUsername(user.Username); isFound {
		return errs.Wrap(ErrUserNameDuplicate, "service.User.Store")
	}
	return u.userRepo.Store(user)

}
func (u *userService) Update(user *m.User) error {
	if e := validate.Validate(user); e != nil {
		return errs.Wrap(ErrUserInvalid, "service.User.Update")
	}
	if user.ID == "" {
		user.ID = shortid.MustGenerate()
	}
	return u.userRepo.Update(user)

}
func (u *userService) Delete(user *m.User) error {
	if user.ID == "" {
		return errs.Wrap(ErrUserNotFound, "service.User.Delete")
	}
	if e := u.userRepo.Delete(user); e != nil {
		return e
	}
	return nil

}
