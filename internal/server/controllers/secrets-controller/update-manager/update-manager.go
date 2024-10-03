package updatemanager

import (
	"errors"
	"log/slog"
	"sync"

	"github.com/FlutterDizaster/EncryNest/internal/models/secrets"
	"github.com/google/uuid"
)

type UpdateManager struct {
	updatesList sync.Map
}

func (m *UpdateManager) SubscribeClient(userID, clientID uuid.UUID) (<-chan secrets.Update, error) {
	// Try to get userUpdater from map
	updater, ok := m.updatesList.Load(userID)
	if !ok {
		// Create new userUpdater
		updater = &userUpdater{}
		m.updatesList.Store(userID, updater)
	}

	// Type assertion
	clientUpdater, ok := updater.(*userUpdater)
	if !ok {
		slog.Error(
			"Type assertion error",
			slog.String("userID", userID.String()),
			slog.String("clientID", clientID.String()),
		)
		return nil, errors.New("type assertion error")
	}

	// Register client
	return clientUpdater.SubscribeClient(clientID)
}

func (m *UpdateManager) UnsubscribeClient(userID, clientID uuid.UUID) error {
	// Try to get userUpdater from map
	updater, ok := m.updatesList.Load(userID)
	if !ok {
		return errors.New("user not found")
	}

	// Type assertion
	uUpdater, ok := updater.(*userUpdater)
	if !ok {
		slog.Error(
			"Type assertion error",
			slog.String("userID", userID.String()),
			slog.String("clientID", clientID.String()),
		)
		return errors.New("type assertion error")
	}

	// Unregister client
	return uUpdater.UnsubscribeClient(clientID)
}

func (m *UpdateManager) SendUpdateFrom(userID, clientID uuid.UUID, update secrets.Update) error {
	// Try to get userUpdater from map
	updater, ok := m.updatesList.Load(userID)
	if !ok {
		return errors.New("user not found")
	}

	// Type assertion
	uUpdater, ok := updater.(*userUpdater)
	if !ok {
		slog.Error(
			"Type assertion error",
			slog.String("userID", userID.String()),
			slog.String("clientID", clientID.String()),
		)
		return errors.New("type assertion error")
	}

	// Send update
	uUpdater.SendUpdateFrom(clientID, update)

	return nil
}

func (m *UpdateManager) SendUpdateTo(userID, clientID uuid.UUID, update secrets.Update) error {
	// Try to get userUpdater from map
	updater, ok := m.updatesList.Load(userID)
	if !ok {
		return errors.New("user not found")
	}

	// Type assertion
	uUpdater, ok := updater.(*userUpdater)
	if !ok {
		slog.Error(
			"Type assertion error",
			slog.String("userID", userID.String()),
			slog.String("clientID", clientID.String()),
		)
		return errors.New("type assertion error")
	}

	// Send update
	uUpdater.SendUpdateTo(userID, update)

	return nil
}
