package usercontroller

import (
	"context"
	"time"

	"github.com/FlutterDizaster/EncryNest/internal/models"
	jwtresolver "github.com/FlutterDizaster/EncryNest/internal/server/jwt-resolver"
	"github.com/google/uuid"
)

const (
	TokenTTL = time.Hour * 24 * 31
	Issuer   = "EncryNest"
)

type UserRepository interface {
	AddUser(ctx context.Context, user *models.UserCredentials) (uuid.UUID, error)
	// GetUser(ctx context.Context, id uuid.UUID) (*models.UserCredentials, error)
	GetUserID(ctx context.Context, user *models.UserCredentials) (uuid.UUID, error)
}

type UserController struct {
	userRepo    UserRepository
	jwtResolver *jwtresolver.JWTResolver
}

// var _ userservice.UserController = (*UserController)(nil)

// RegisterUser registers new user in the system and returns JWT token string.
func (c *UserController) RegisterUser(
	ctx context.Context,
	user *models.UserCredentials,
) (string, error) {
	// Adding user to the repository
	id, err := c.userRepo.AddUser(ctx, user)
	if err != nil {
		return "", err
	}

	// Returning JWT token
	return c.jwtResolver.CreateToken(Issuer, user.Username, id)
}

// AuthUser authenticates user and returns JWT token string.
func (c *UserController) AuthUser(
	ctx context.Context,
	user *models.UserCredentials,
) (string, error) {
	// Getting user from the repository
	id, err := c.userRepo.GetUserID(ctx, user)
	if err != nil {
		return "", err
	}

	// Returning JWT token
	return c.jwtResolver.CreateToken(Issuer, user.Username, id)
}
