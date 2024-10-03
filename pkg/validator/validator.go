// Validator package contains functions for validating user data.
// TODO: Add more validation functions.
// TODO: Add real validation.
package validator

import "errors"

const (
	UsernameMinLength = 3
	UsernameMaxLength = 20
	EmailMinLenght    = 6
	EmailMaxLength    = 100
	PasswordMinLength = 8
)

var (
	ErrInvalidUsername  = errors.New("invalid username")
	ErrInvalidEmail     = errors.New("invalid email")
	ErrPasswordTooShort = errors.New("password is too short")
)

// ValidateUsername validates username.
func ValidateUsername(username string) error {
	if len(username) < UsernameMinLength || len(username) > UsernameMaxLength {
		return ErrInvalidUsername
	}
	return nil
}

// ValidateEmail validates email.
func ValidateEmail(email string) error {
	if len(email) < EmailMinLenght || len(email) > EmailMaxLength {
		return ErrInvalidEmail
	}
	return nil
}

// ValidatePassword validates password.
func ValidatePassword(password string) error {
	if len(password) < PasswordMinLength {
		return ErrPasswordTooShort
	}
	return nil
}
