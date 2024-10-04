package jwtresolver

import (
	"crypto/rand"
	"encoding/hex"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTResolver(t *testing.T) {
	b := make([]byte, 512)
	_, _ = rand.Read(b)
	secret := hex.EncodeToString(b)
	tokenTTL := time.Hour

	type test struct {
		name    string
		issuer  string
		subject string
		userID  uuid.UUID
	}
	tests := []test{
		{
			name:    "success",
			issuer:  "test_issuer",
			subject: "test_subject",
			userID:  uuid.New(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := &JWTResolver{
				secret:   secret,
				tokenTTL: tokenTTL,
			}

			token, err := res.CreateToken(tt.issuer, tt.subject, tt.userID)
			require.NoError(t, err)

			claims, err := res.DecryptToken(token)
			require.NoError(t, err)

			assert.Equal(t, tt.issuer, claims.Issuer)
			assert.Equal(t, tt.subject, claims.Subject)
			assert.Equal(t, tt.userID, claims.UserID)
		})
	}
}
