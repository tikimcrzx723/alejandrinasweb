package store

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"errors"
	"time"

	"github.com/tikimcrzx723/alejandrinasweb/internal/dtos"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrTokenNotFound = errors.New("token not found")
	ErrTokenConflict = errors.New("token conflict")
	ErrTokenExpired  = errors.New("token expired")
	ErrTokenInvalid  = errors.New("token invalid")
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
	ScopePasswordReset  = "password-reset"

	ActivationExpiryTime = time.Hour * 12
)

func generetaToken(insertToken dtos.InsertToken) (*dtos.InsertToken, error) {
	token := &dtos.InsertToken{
		UserID: insertToken.UserID,
		Scope:  insertToken.Scope,
		Expiry: time.Now().Add(time.Duration(insertToken.Expiry) * time.Second).Unix(),
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

type TokenStore struct {
	db *pgxpool.Pool
}
