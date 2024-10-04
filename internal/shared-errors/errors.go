package sharederrors

import "errors"

var (
	ErrUserAlredyExists = errors.New("username with given email or username is alredy registered")
	ErrUserNotFound     = errors.New("user not found")
)
