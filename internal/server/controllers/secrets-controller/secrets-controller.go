package secretscontroller

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"sync"

	"github.com/FlutterDizaster/EncryNest/internal/models/secrets"
	updatemanager "github.com/FlutterDizaster/EncryNest/internal/server/controllers/secrets-controller/update-manager"
	"github.com/google/uuid"
)

type SecretsRepository interface {
	// AddSecret adds new secret for specified user.
	// Returns new secret ID and secret version.
	AddSecret(
		ctx context.Context,
		userID uuid.UUID,
		secret *secrets.Secret,
	) (uuid.UUID, string, error)

	// UpdateSecret updates secret with specified ID.
	// Returns new secret version.
	UpdateSecret(ctx context.Context, userID uuid.UUID, secret *secrets.Secret) (string, error)

	// RemoveSecret removes secret with specified ID.
	// Returns secret version.
	RemoveSecret(ctx context.Context, userID uuid.UUID, id uuid.UUID) (string, error)

	// GetSecrets returns secrets above specified version for specified user.
	GetSecretsAboveVersion(
		ctx context.Context,
		userID uuid.UUID,
		knownVersion string,
	) ([]secrets.Secret, error)

	// DeleteUnknownSecretsBeforeVersion deletes secrets before specified version
	// and returns IDs of deleted secrets.
	DeleteUnknownSecretsBeforeVersion(
		ctx context.Context,
		userID uuid.UUID,
		version string,
		knownIDs []uuid.UUID,
	) ([]uuid.UUID, error)
}

// SecretsController is a controller for secrets.
// It is used to create/update/delete secrets in the system and sync secrets across clients.
// Must be created with NewSecretsController function.
type SecretsController struct {
	repo          SecretsRepository
	updateManager *updatemanager.UpdateManager
	wg            sync.WaitGroup
}

// NewSecretsController creates new SecretsController.
func NewSecretsController(repo SecretsRepository) *SecretsController {
	return &SecretsController{
		repo:          repo,
		updateManager: &updatemanager.UpdateManager{},
		wg:            sync.WaitGroup{},
	}
}

// MakeUpdate used to Create/Update/Delete secret in the system
// and send update to clients.
// Returns new secret ID and secret version.
func (c *SecretsController) MakeUpdate(
	ctx context.Context,
	update *secrets.Update,
) (string, uuid.UUID, error) {
	var version string
	newID := uuid.Nil
	var err error

	// Make right action with secrets repo
	switch update.Action {
	case secrets.UpdateActionCreate:
		newID, version, err = c.repo.AddSecret(ctx, update.UserID, update.Secret)
		if err != nil {
			return "", uuid.Nil, err
		}

	case secrets.UpdateActionUpdate:
		version, err = c.repo.UpdateSecret(ctx, update.UserID, update.Secret)
		newID = update.Secret.ID
		if err != nil {
			return "", uuid.Nil, err
		}

	case secrets.UpdateActionDelete:
		version, err = c.repo.RemoveSecret(ctx, update.UserID, update.Secret.ID)
		if err != nil {
			return "", uuid.Nil, err
		}

	default:
		return "", uuid.Nil, errors.New("unknown action")
	}

	// Send update to clients
	err = c.updateManager.SendUpdateFrom(update.UserID, update.ClientID, *update)
	if err != nil {
		return "", uuid.Nil, err
	}

	return version, newID, nil
}

// syncSecrets synchronizes secrets with secrets repo for specified user.
// Returns error if something went wrong.
// It also sends deleted updates to clients.
func (c *SecretsController) syncSecrets(
	ctx context.Context,
	userID, clientID uuid.UUID,
	knownVersion string,
	knownIDs []string,
) error {
	// Sort knownIDs
	slices.Sort(knownIDs)

	// Delete unknown secrets where version <= knownVersion
	err := c.deleteUnknownSecrets(ctx, userID, clientID, knownVersion, knownIDs)
	if err != nil {
		slog.Error(
			"Error while deleting secrets",
			slog.String("userID", userID.String()),
			slog.String("clientID", clientID.String()),
			slog.Any("err", err),
		)
	}

	// Get secrets from secrets repo where version > knownVersion
	newSecrets, err := c.repo.GetSecretsAboveVersion(ctx, userID, knownVersion)
	if err != nil {
		slog.Error(
			"Error while getting secrets",
			slog.String("userID", userID.String()),
			slog.String("clientID", clientID.String()),
			slog.Any("err", err),
		)
		return err
	}

	// Send updates to client
	for i := range newSecrets {
		upd := secrets.Update{
			Secret: &newSecrets[i],
		}
		// If knownIDs contains newSecrets[i].ID, action update, else create
		if slices.Contains(knownIDs, newSecrets[i].ID.String()) {
			upd.Action = secrets.UpdateActionUpdate
		} else {
			upd.Action = secrets.UpdateActionCreate
		}
		// Send update to update manager
		err = c.updateManager.SendUpdateTo(userID, clientID, upd)

		if err != nil {
			slog.Error(
				"Error while sending update",
				slog.String("userID", userID.String()),
				slog.String("clientID", clientID.String()),
				slog.Any("err", err),
			)
			return err
		}
	}
	return nil
}

// Delete secrets where ID not in knownIDs and version <= knownVersion
// and send deleted secrets to update manager.
func (c *SecretsController) deleteUnknownSecrets(
	ctx context.Context,
	userID, clientID uuid.UUID,
	knownVersion string,
	knownIDs []string,
) error {
	// Convert knownIDs to uuid.UUID
	knownIDsUUID := make([]uuid.UUID, 0, len(knownIDs))
	for i := range knownIDs {
		id, err := uuid.Parse(knownIDs[i])

		// If error - skip
		if err != nil {
			slog.Error(
				"Error while parsing secret ID",
				slog.String("knownID", knownIDs[i]),
			)
			continue
		}

		knownIDsUUID = append(knownIDsUUID, id)
	}

	// Delete secret from secrets repo where ID not in knownIDs and version <= knownVersion
	deletedSecrets, err := c.repo.DeleteUnknownSecretsBeforeVersion(
		ctx,
		userID,
		knownVersion,
		knownIDsUUID,
	)
	if err != nil {
		return err
	}

	// Send deleted secrets to update manager
	for i := range deletedSecrets {
		// Make update instance
		upd := secrets.Update{
			Action: secrets.UpdateActionDelete,
			Secret: &secrets.Secret{
				ID: deletedSecrets[i],
			},
		}
		// Send update
		err = c.updateManager.SendUpdateFrom(
			userID,
			clientID,
			upd,
		)
		// Log error
		if err != nil {
			slog.Error(
				"Error while sending update",
				slog.String("userID", userID.String()),
				slog.String("clientID", clientID.String()),
				slog.Any("err", err),
			)
		}
	}

	return nil
}

func (c *SecretsController) SubscribeUpdates(
	ctx context.Context,
	userID uuid.UUID,
	clientID uuid.UUID,
	knownVersion string,
	knownIDs []string,
) <-chan secrets.Update {
	// Subscribe for updates
	updateCh, err := c.updateManager.SubscribeClient(userID, clientID)
	if err != nil {
		return nil
	}

	// Start subscription watcher in background
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		c.unsubscribeWhenContextClosed(ctx, userID, clientID)
	}()

	// Start getting updates in background
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		// Get updates
		updErr := c.syncSecrets(ctx, userID, clientID, knownVersion, knownIDs)
		if updErr != nil {
			slog.Error(
				"Error while getting updates",
				slog.String("userID", userID.String()),
				slog.String("clientID", clientID.String()),
				slog.Any("err", updErr),
			)
			return
		}
	}()

	return updateCh
}

func (c *SecretsController) unsubscribeWhenContextClosed(
	ctx context.Context,
	userID, clientID uuid.UUID,
) {
	<-ctx.Done()
	err := c.updateManager.UnsubscribeClient(userID, clientID)

	if err != nil {
		slog.Error(
			"Error while unsubscribing from updates",
			slog.String("userID", userID.String()),
			slog.String("clientID", clientID.String()),
		)
	}
}
