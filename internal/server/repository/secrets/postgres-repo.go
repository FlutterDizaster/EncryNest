package secretsrepo

import (
	"context"
	"log/slog"
	"time"

	"github.com/FlutterDizaster/EncryNest/internal/models"
	"github.com/FlutterDizaster/EncryNest/internal/server/repository/postgres"
	"github.com/google/uuid"
)

type PostgresSecretsRepository struct {
	poolManager *postgres.PoolManager
}

// var _ secretscontroller.SecretsRepository = (*PostgresSecretsRepository)(nil)

func NewPostgresRepository(poolManager *postgres.PoolManager) *PostgresSecretsRepository {
	return &PostgresSecretsRepository{
		poolManager: poolManager,
	}
}

func (r *PostgresSecretsRepository) AddSecret(
	ctx context.Context,
	userID uuid.UUID,
	secret *models.Secret,
) (uuid.UUID, string, error) {
	pool := r.poolManager.Pool()

	row := pool.QueryRow(ctx, AddSecretQuery, userID, secret.Kind, secret.Data)

	var version time.Time
	var id uuid.UUID

	err := row.Scan(&id, &version)
	if err != nil {
		return uuid.Nil, "", err
	}

	return id, version.Format(time.RFC3339Nano), nil
}

func (r *PostgresSecretsRepository) UpdateSecret(
	ctx context.Context,
	userID uuid.UUID,
	secret *models.Secret,
) (string, error) {
	pool := r.poolManager.Pool()

	row := pool.QueryRow(ctx, UpdateSecretQuery, secret.Kind, secret.Data, secret.ID, userID)

	var version time.Time
	err := row.Scan(&version)
	if err != nil {
		return "", err
	}

	return version.Format(time.RFC3339Nano), nil
}

func (r *PostgresSecretsRepository) RemoveSecret(
	ctx context.Context,
	userID uuid.UUID,
	id uuid.UUID,
) (string, error) {
	pool := r.poolManager.Pool()

	_, err := pool.Exec(ctx, RemoveSecretQuery, id, userID)
	if err != nil {
		return "", err
	}

	return time.Now().Format(time.RFC3339Nano), nil
}

func (r *PostgresSecretsRepository) GetSecretsAboveVersion(
	ctx context.Context,
	userID uuid.UUID,
	knownVersion string,
) ([]models.Secret, error) {
	pool := r.poolManager.Pool()

	knownVersionTime, err := time.Parse(time.RFC3339Nano, knownVersion)
	if err != nil {
		return nil, err
	}
	rows, err := pool.Query(ctx, GetSecretsAboveVersionQuery, userID, knownVersionTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	secrets := make([]models.Secret, 0)

	for rows.Next() {
		var secret models.Secret
		var version time.Time
		err = rows.Scan(&secret.ID, &version, &secret.Kind, &secret.Data)
		if err != nil {
			slog.Error("Error while scanning secret", slog.Any("err", err))
			continue
		}

		secret.Version = version.Format(time.RFC3339Nano)

		secrets = append(secrets, secret)
	}

	return secrets, nil
}

func (r *PostgresSecretsRepository) DeleteUnknownSecretsBeforeVersion(
	ctx context.Context,
	userID uuid.UUID,
	version string,
	knownIDs []uuid.UUID,
) ([]uuid.UUID, error) {
	pool := r.poolManager.Pool()

	knownVersionTime, err := time.Parse(time.RFC3339Nano, version)
	if err != nil {
		return nil, err
	}

	rows, err := pool.Query(
		ctx,
		DeleteUnknownSecretsBeforeVersionQuery,
		userID,
		knownIDs,
		knownVersionTime,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deletedIDs := make([]uuid.UUID, 0)

	for rows.Next() {
		var id uuid.UUID
		err = rows.Scan(&id)
		if err != nil {
			slog.Error("Error while scanning deleted secret ID", slog.Any("err", err))
			continue
		}

		deletedIDs = append(deletedIDs, id)
	}

	return deletedIDs, nil
}
