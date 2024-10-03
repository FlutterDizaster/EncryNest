package secretsrepo

import (
	"errors"
	"slices"
	"sync"
	"time"

	"github.com/FlutterDizaster/EncryNest/internal/models/secrets"
	"github.com/google/uuid"
)

type SecretsList struct {
	list sync.Map
}

func (l *SecretsList) Add(secret secrets.Secret) (uuid.UUID, string, error) {
	id := uuid.New()
	version := time.Now().Format(time.RFC3339)

	secret.ID = id
	secret.Version = version

	l.list.Store(id, secret)

	return id, version, nil
}

func (l *SecretsList) Get(id uuid.UUID) (secrets.Secret, error) {
	secret, ok := l.list.Load(id)
	if !ok || secret == nil {
		return secrets.Secret{}, errors.New("secret not found")
	}
	return secret.(secrets.Secret), nil
}

func (l *SecretsList) Delete(id uuid.UUID) error {
	l.list.Delete(id)
	return nil
}

func (l *SecretsList) Update(id uuid.UUID, secret secrets.Secret) (string, error) {
	secret.Version = time.Now().Format(time.RFC3339)
	l.list.Store(id, secret)
	return secret.Version, nil
}

func (l *SecretsList) GetSecretsAboveVersion(
	knownVersion string,
) ([]secrets.Secret, error) {
	ver, err := time.Parse(time.RFC3339, knownVersion)
	if err != nil {
		return nil, errors.New("error wrong version")
	}

	newSecrets := make([]secrets.Secret, 0)

	l.list.Range(func(_, value any) bool {
		secret, ok := value.(secrets.Secret)
		if !ok {
			return true
		}

		sver, verErr := time.Parse(time.RFC3339, secret.Version)
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

	ver, err := time.Parse(time.RFC3339, version)
	if err != nil {
		return nil, errors.New("error wrong version")
	}

	l.list.Range(func(key, value any) bool {
		secret, ok := value.(secrets.Secret)
		if !ok {
			return true
		}

		sver, verErr := time.Parse(time.RFC3339, secret.Version)
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
