package usersrepo

import (
	"context"
	"errors"

	"github.com/FlutterDizaster/EncryNest/internal/models"
	"github.com/FlutterDizaster/EncryNest/internal/server/repository/postgres"
	sharederrors "github.com/FlutterDizaster/EncryNest/internal/shared-errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// TODO: add retry logic

type PostgresUserRepository struct {
	poolManager *postgres.PoolManager
}

// var _ usercontroller.UserRepository = &PostgresUserRepository{}

func NewPostgresRepository(poolManager *postgres.PoolManager) *PostgresUserRepository {
	return &PostgresUserRepository{
		poolManager: poolManager,
	}
}

func (r *PostgresUserRepository) AddUser(
	ctx context.Context,
	user *models.UserCredentials,
) (uuid.UUID, error) {
	pool := r.poolManager.Pool()

	row := pool.QueryRow(ctx, CreateUserQuery, user.Username, user.PasswordHash, user.Email)

	var id uuid.UUID
	err := row.Scan(&id)

	// FIXME: wrong error check
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, sharederrors.ErrUserAlredyExists
		}
		return uuid.Nil, err
	}

	return id, nil
}

func (r *PostgresUserRepository) GetUserID(
	ctx context.Context,
	user *models.UserCredentials,
) (uuid.UUID, error) {
	pool := r.poolManager.Pool()

	row := pool.QueryRow(ctx, GetUserIDQuery, user.Username, user.PasswordHash)

	var id uuid.UUID
	err := row.Scan(&id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, sharederrors.ErrUserNotFound
		}
		return uuid.Nil, err
	}

	return id, nil
}
