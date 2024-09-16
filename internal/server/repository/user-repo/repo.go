package userrepo

import (
	"context"
	"errors"

	"github.com/FlutterDizaster/EncryNest/internal/models"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailTaken   = errors.New("user with given email is already registered")
)

type UserRepo interface {
	RegisterUser(ctx context.Context, user *models.UserData) error
	GetUser(ctx context.Context, id uuid.UUID) (*models.UserData, error)
	GetUserByEmail(ctx context.Context, email string) (*models.UserData, error)
	GetUserByUsername(ctx context.Context, username string) (*models.UserData, error)
}
