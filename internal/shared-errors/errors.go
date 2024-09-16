package sharederrors

import "errors"

var (
	ErrEmailTaken    = errors.New("user with given email is already registered")
	ErrUsernameTaken = errors.New("user with given username is already registered")
	ErrUserNotFound  = errors.New("user not found")
)
