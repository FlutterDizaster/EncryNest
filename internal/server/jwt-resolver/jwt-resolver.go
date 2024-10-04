package jwtresolver

import (
	"errors"
	"time"

	"github.com/FlutterDizaster/EncryNest/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Settings struct {
	Secret   string
	TokenTTL time.Duration
}

type JWTResolver struct {
	secret   string
	tokenTTL time.Duration
}

func New(settings Settings) *JWTResolver {
	return &JWTResolver{
		secret:   settings.Secret,
		tokenTTL: settings.TokenTTL,
	}
}

func (res *JWTResolver) DecryptToken(tokenString string) (*models.Claims, error) {
	// Создание структуры models.Token
	claims := &models.Claims{}

	// Парсинг токена
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("error unexpected signing method")
		}
		return []byte(res.secret), nil
	})

	// Проверка токена на валидность
	if !token.Valid {
		return claims, errors.New("error invalid token")
	}

	return claims, err
}

func (res *JWTResolver) CreateToken(issuer, subject string, userID uuid.UUID) (string, error) {
	// Создание данных токена
	claims := models.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(res.tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: userID,

		// FIXME: temporary solution.
		// In future it should be best way to get client ID from DB for registered clients.
		// Or store it on client side and send it with auth request.
		ClientID: uuid.New(),
	}

	// Создание токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Создание зашифрованной строки токена
	tokenString, err := token.SignedString([]byte(res.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
