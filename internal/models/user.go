package models

import (
	"github.com/FlutterDizaster/EncryNest/pkg/keychain"
	"github.com/google/uuid"
)

type UserData struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	JWTToken string    `json:"jwt_token,omitempty"`
	XKeys    keychain.KeyPair
}
