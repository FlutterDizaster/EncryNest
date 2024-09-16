package userrepo

import (
	"context"

	"github.com/FlutterDizaster/EncryNest/internal/models"
	"github.com/google/uuid"
)

type InMemoryUserRepo struct {
	users map[uuid.UUID]*models.UserData
}

var _ UserRepo = (*InMemoryUserRepo)(nil)

func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{
		users: make(map[uuid.UUID]*models.UserData),
	}
}

func (r *InMemoryUserRepo) RegisterUser(_ context.Context, user *models.UserData) error {
	for _, v := range r.users {
		if v.Email == user.Email {
			return ErrEmailTaken
		}
	}

	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepo) GetUser(_ context.Context, id uuid.UUID) (*models.UserData, error) {
	return r.users[id], nil
}

func (r *InMemoryUserRepo) GetUserByEmail(
	_ context.Context,
	email string,
) (*models.UserData, error) {
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

func (r *InMemoryUserRepo) GetUserByUsername(
	_ context.Context,
	username string,
) (*models.UserData, error) {
	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}
