package models

import "github.com/google/uuid"

type UpdateAction int

const (
	UpdateActionCreate UpdateAction = iota
	UpdateActionUpdate
	UpdateActionDelete
)

type Update struct {
	UserID   uuid.UUID
	ClientID uuid.UUID
	Action   UpdateAction
	Secret   *Secret
}
