package secretsrepo

import (
	"errors"
	"slices"
	"sync"
	"time"

	"github.com/FlutterDizaster/EncryNest/internal/models"
	"github.com/google/uuid"
)

type SecretsList struct {
	list sync.Map
}

func (l *SecretsList) Add(secret models.Secret) (uuid.UUID, string, error) {
	id := uuid.New()
	version := time.Now().Format(time.RFC3339Nano)

	secret.ID = id
	secret.Version = version

	l.list.Store(id, secret)

	return id, version, nil
}

func (l *SecretsList) Get(id uuid.UUID) (models.Secret, error) {
	secret, ok := l.list.Load(id)
	if !ok || secret == nil {
		return models.Secret{}, errors.New("secret not found")
	}
	return secret.(models.Secret), nil
}

func (l *SecretsList) Delete(id uuid.UUID) error {
	l.list.Delete(id)
	return nil
}

func (l *SecretsList) Update(id uuid.UUID, secret models.Secret) (string, error) {
	secret.Version = time.Now().Format(time.RFC3339Nano)
	l.list.Store(id, secret)
	return secret.Version, nil
}

func (l *SecretsList) GetSecretsAboveVersion(
	knownVersion string,
) ([]models.Secret, error) {
	ver, err := time.Parse(time.RFC3339Nano, knownVersion)
	if err != nil {
		return nil, errors.New("error wrong version")
	}

	newSecrets := make([]models.Secret, 0)

	l.list.Range(func(_, value any) bool {
		secret, ok := value.(models.Secret)
		if !ok {
			return true
		}

		sver, verErr := time.Parse(time.RFC3339Nano, secret.Version)
		if verErr != nil {
			return true
		}

		if sver.After(ver) {
			newSecrets = append(newSecrets, secret)
		}

		return true
	})

	return newSecrets, nil
}

func (l *SecretsList) DeleteUnknownSecretsBeforeVersion(
	version string,
	knownIDs []uuid.UUID,
) ([]uuid.UUID, error) {
	deletedIDs := make([]uuid.UUID, 0)

	ver, err := time.Parse(time.RFC3339Nano, version)
	if err != nil {
		return nil, errors.New("error wrong version")
	}

	l.list.Range(func(key, value any) bool {
		secret, ok := value.(models.Secret)
		if !ok {
			return true
		}

		sver, verErr := time.Parse(time.RFC3339Nano, secret.Version)
		if verErr != nil {
			return true
		}

		if sver.After(ver) {
			return true
		}

		if slices.Contains(knownIDs, key.(uuid.UUID)) && value != nil {
			return true
		}

		deletedIDs = append(deletedIDs, key.(uuid.UUID))

		return true
	})

	for _, id := range deletedIDs {
		l.list.Delete(id)
	}

	return deletedIDs, nil
}
