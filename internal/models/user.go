package models

import (
	"github.com/google/uuid"
)

type UserData struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	JWTToken string    `json:"jwt_token,omitempty"`
}

type UserCredentials struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}
