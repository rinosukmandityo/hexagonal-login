package helper

import (
	"errors"
)

var (
	ErrUserNotFound      = errors.New("User Not Found")
	ErrUserInvalid       = errors.New("User Invalid")
	ErrUserNameDuplicate = errors.New("User Name Already Exists")
)
