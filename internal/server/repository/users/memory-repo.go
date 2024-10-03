package usersrepo

import (
	"context"
	"errors"
	"sync"

	"github.com/FlutterDizaster/EncryNest/internal/models"
	sharederrors "github.com/FlutterDizaster/EncryNest/internal/shared-errors"
	"github.com/google/uuid"
)

var (
	ErrConvertingCredentials = errors.New("converting user credentials error")
)

type InMemoryRepository struct {
	users sync.Map
}

// var _ usercontroller.UserRepository = &InMemoryRepository{}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{}
}

func (r *InMemoryRepository) AddUser(
	ctx context.Context,
	user *models.UserCredentials,
) (uuid.UUID, error) {
	return r.addUser(ctx, user)
}

func (r *InMemoryRepository) GetUserID(
	_ context.Context,
	user *models.UserCredentials,
) (uuid.UUID, error) {
	return r.getUserID(user)
}

func (r *InMemoryRepository) addUser(
	_ context.Context,
	user *models.UserCredentials,
) (uuid.UUID, error) {
	// check is user alresdy exists
	var err error
	r.users.Range(func(_, value any) bool {
		u, ok := value.(models.UserCredentials)
		if !ok {
			return true
		}

		if u.Email == user.Email {
			err = sharederrors.ErrUserAlredyExists
			return false
		}

		if u.Username == user.Username {
			err = sharederrors.ErrUserAlredyExists
			return false
		}
		return true
	})

	if err != nil {
		return uuid.Nil, err
	}

	// Creating new user entry
	id := uuid.New()

	for _, ok := r.users.Load(id); ok; _, ok = r.users.Load(id) {
		id = uuid.New()
	}

	r.users.Store(id, *user)

	return id, nil
}

func (r *InMemoryRepository) getUserID(user *models.UserCredentials) (uuid.UUID, error) {
	var err error
	var id uuid.UUID
	found := false

	r.users.Range(func(key, value any) bool {
		u, ok := value.(models.UserCredentials)
		if !ok {
			err = ErrConvertingCredentials
			return false
		}

		if u.Username == user.Username && u.PasswordHash == user.PasswordHash {
			id, ok = key.(uuid.UUID)
			if !ok {
				err = errors.New("converting user id error")
			}
			found = true
			return false
		}

		return true
	})

	if found {
		return id, nil
	}

	if err == nil {
		return uuid.Nil, sharederrors.ErrUserNotFound
	}

	return uuid.Nil, err
}

func (r *InMemoryRepository) GetUser(
	_ context.Context,
	id uuid.UUID,
) (models.UserCredentials, error) {
	return r.getUser(id)
}

func (r *InMemoryRepository) getUser(id uuid.UUID) (models.UserCredentials, error) {
	u, ok := r.users.Load(id)
	if ok {
		return models.UserCredentials{}, sharederrors.ErrUserNotFound
	}

	user, ok := u.(models.UserCredentials)

	if !ok {
		return models.UserCredentials{}, ErrConvertingCredentials
	}

	return user, nil
}
