package updatemanager

import (
	"errors"
	"log/slog"
	"sync"

	"github.com/FlutterDizaster/EncryNest/internal/models/secrets"
	"github.com/google/uuid"
)

type userUpdater struct {
	clients sync.Map
}

func (u *userUpdater) SubscribeClient(clientID uuid.UUID) (<-chan secrets.Update, error) {
	// Try to get client from map
	ch, ok := u.clients.Load(clientID)
	if ok && ch != nil {
		return nil, errors.New("client already subscribed")
	}

	// Create new updates Channel with capacity 1
	updatesChan := make(chan secrets.Update, 1)

	// Register client
	u.clients.Store(clientID, updatesChan)

	return updatesChan, nil
}

func (u *userUpdater) UnsubscribeClient(clientID uuid.UUID) error {
	// Try to get client from map
	ch, ok := u.clients.Load(clientID)
	if !ok || ch == nil {
		return ErrNotFound
	}

	// Type assertion
	updatesChan, ok := ch.(chan secrets.Update)
	if !ok {
		slog.Error(
			"Type assertion error",
			slog.String("clientID", clientID.String()),
		)
		return ErrTypeAssertion
	}

	// Close channel and delete it from map
	close(updatesChan)

	u.clients.Delete(clientID)

	return nil
}

func (u *userUpdater) SendUpdateFrom(clientID uuid.UUID, update secrets.Update) {
	u.clients.Range(func(key, value any) bool {
		// Skip current client
		if key.(uuid.UUID) == clientID {
			return true
		}

		// Type assertion
		updatesChan, ok := value.(chan secrets.Update)
		if !ok {
			slog.Error(
				"Type assertion error",
				slog.String("clientID", clientID.String()),
			)
			return true
		}

		// Send update
		updatesChan <- update

		return true
	})
}

func (u *userUpdater) SendUpdateTo(clientID uuid.UUID, update secrets.Update) {
	value, ok := u.clients.Load(clientID)

	if !ok || value == nil {
		return
	}

	// Type assertion
	updatesChan, ok := value.(chan secrets.Update)
	if !ok {
		slog.Error(
			"Type assertion error",
			slog.String("clientID", clientID.String()),
		)
		return
	}

	// Send update
	updatesChan <- update
}
