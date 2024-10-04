package secretsrepo

import (
	"context"
	"sync"
	"time"

	"github.com/FlutterDizaster/EncryNest/internal/models"
	"github.com/google/uuid"
)

type SecretsRepo struct {
	userSecrets sync.Map
}

func NewInMemoryRepository() *SecretsRepo {
	return &SecretsRepo{
		userSecrets: sync.Map{},
	}
}

// var _ secretscontroller.SecretsRepository = (*SecretsRepo)(nil)

func (r *SecretsRepo) AddSecret(
	_ context.Context,
	userID uuid.UUID,
	secret *models.Secret,
) (uuid.UUID, string, error) {
	secretsList, ok := r.userSecrets.Load(userID)
	if !ok {
		secretsList = &SecretsList{}
		r.userSecrets.Store(userID, secretsList)
	}

	return secretsList.(*SecretsList).Add(*secret)
}

func (r *SecretsRepo) UpdateSecret(
	_ context.Context,
	userID uuid.UUID,
	secret *models.Secret,
) (string, error) {
	secretsList, ok := r.userSecrets.Load(userID)
	if !ok {
		secretsList = &SecretsList{}
		r.userSecrets.Store(userID, secretsList)
	}

	return secretsList.(*SecretsList).Update(secret.ID, *secret)
}

func (r *SecretsRepo) RemoveSecret(
	_ context.Context,
	userID uuid.UUID,
	id uuid.UUID,
) (string, error) {
	secretsList, ok := r.userSecrets.Load(userID)
	if !ok {
		secretsList = &SecretsList{}
		r.userSecrets.Store(userID, secretsList)
	}

	return time.Now().Format(time.RFC3339), secretsList.(*SecretsList).Delete(id)
}

func (r *SecretsRepo) GetSecretsAboveVersion(
	_ context.Context,
	userID uuid.UUID,
	knownVersion string,
) ([]models.Secret, error) {
	secretsList, ok := r.userSecrets.Load(userID)
	if !ok {
		secretsList = &SecretsList{}
		r.userSecrets.Store(userID, secretsList)
	}

	return secretsList.(*SecretsList).GetSecretsAboveVersion(knownVersion)
}

func (r *SecretsRepo) DeleteUnknownSecretsBeforeVersion(
	_ context.Context,
	userID uuid.UUID,
	version string,
	knownIDs []uuid.UUID,
) ([]uuid.UUID, error) {
	secretsList, ok := r.userSecrets.Load(userID)
	if !ok {
		secretsList = &SecretsList{}
		r.userSecrets.Store(userID, secretsList)
	}

	return secretsList.(*SecretsList).DeleteUnknownSecretsBeforeVersion(
		version,
		knownIDs,
	)
}
