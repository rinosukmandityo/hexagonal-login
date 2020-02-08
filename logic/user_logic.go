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
	userRepo repo.LoginRepository
}

func NewUserService(userRepo repo.LoginRepository) svc.UserService {
	return &userService{
		userRepo,
	}
}

func (u *userService) GetAll() ([]m.User, error) {
	res := []m.User{}
	param := repo.GetAllParam{
		Tablename: new(m.User).TableName(),
		Result:    &res,
	}
	if e := u.userRepo.GetAll(param); e != nil {
		return res, e
	}
	return res, nil
}

func (u *userService) GetById(id string) (*m.User, error) {
	res := new(m.User)
	param := repo.GetParam{
		Tablename: res.TableName(),
		Filter:    map[string]interface{}{"_id": id},
		Result:    res,
	}
	if e := u.userRepo.GetBy(param); e != nil {
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
	param := repo.StoreParam{
		Tablename: user.TableName(),
		Data:      user,
	}
	return u.userRepo.Store(param)

}
func (u *userService) Update(user *m.User) error {
	if e := validate.Validate(user); e != nil {
		return errs.Wrap(helper.ErrUserInvalid, "service.User.Update")
	}
	if user.ID == "" {
		user.ID = shortid.MustGenerate()
	}
	param := repo.UpdateParam{
		Tablename: user.TableName(),
		Filter:    map[string]interface{}{"_id": user.ID},
		Data:      user,
	}
	return u.userRepo.Update(param)

}
func (u *userService) Delete(user *m.User) error {
	if user.ID == "" {
		return errs.Wrap(helper.ErrUserNotFound, "service.User.Delete")
	}
	param := repo.DeleteParam{
		Tablename: user.TableName(),
		Filter:    map[string]interface{}{"_id": user.ID},
	}
	if e := u.userRepo.Delete(param); e != nil {
		return e
	}
	return nil

}

func (u *userService) GetByUsername(username string) (bool, *m.User, error) {
	res := new(m.User)
	param := repo.GetParam{
		Tablename: res.TableName(),
		Filter:    map[string]interface{}{"Username": username},
		Result:    res,
	}
	if e := u.userRepo.GetBy(param); e != nil {
		return false, res, e
	}

	return true, res, nil
}
