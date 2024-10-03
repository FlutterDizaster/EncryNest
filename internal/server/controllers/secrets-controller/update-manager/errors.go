package updatemanager

import "errors"

var (
	ErrTypeAssertion = errors.New("type assertion error")
	ErrNotFound      = errors.New("not found")
)
