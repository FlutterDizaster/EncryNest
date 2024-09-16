package secretsrepo

import (
	"context"

	"github.com/FlutterDizaster/EncryNest/internal/models/secrets"
	"github.com/google/uuid"
)

// GetSecretOptions contains options for GetSecrets method.
type GetSecretOptions struct {
	// Required
	// The ID of the secrets owner.
	Owner uuid.UUID
	// Optional
	// If nil, returns all secrets.
	// If not nil, returns only secret with given ID.
	ID *uuid.UUID
	// Optional
	// If nil, returns all secrets.
	// If not nil, returns only secret with given kind.
	Kind *secrets.SecretKind
}

// SecretsRepo represents secrets repository.
type SecretsRepo interface {
	// StoreSecret stores new secret.
	// Return stored secret with specified ID.
	// Return error if secret is not stored.
	StoreSecret(ctx context.Context, secret secrets.Secret, owner uuid.UUID) (secrets.Secret, error)
	// GetSecrets returns secrets.
	// Return secrets by given options.
	// Return error if secrets is not found.
	GetSecrets(ctx context.Context, options GetSecretOptions) (secrets.Secret, error)
	// DeleteSecret deletes secret with specified ID.
	DeleteSecret(ctx context.Context, id uuid.UUID) error
}
