package store

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/tikimcrzx723/alejandrinasweb/internal/dtos"
	"github.com/tikimcrzx723/alejandrinasweb/internal/store/queries"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

const (
	ClientRole = 1
	AdminRole  = 2
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserConflict      = errors.New("user edit conflict")
	ErrDuplicateEmail    = errors.New("a user with that email already exists")
	ErrDuplicateUsername = errors.New("a user with that username already exists")
)

type UserStore struct {
	db *pgxpool.Pool
}

func GeneratePassword(text string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func ComparePassword(hash []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(password))
}

func (s *UserStore) CreateUser(ctx context.Context, user *dtos.InsertUser) error {
	ctx, cancel := context.WithTimeout(ctx, QueryTimeDuration)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return errBeginTransaction("begin transaction", err)
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(
		ctx,
		queries.UserInsert,
		user.FirstName,
		user.LastName,
		user.Username,
		user.Email,
		user.PasswordHash.Hash,
		false,
		true,
	).Scan(&user.ID)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), `"users_email_key" (SQLSTATE 23505`):
			return ErrDuplicateEmail
		case strings.Contains(err.Error(), `"users_username_key" (SQLSTATE 23505)`):
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	_, err = tx.Exec(
		ctx,
		queries.UserRoleInsert,
		user.ID,
		ClientRole,
	)
	if err != nil {
		return err
	}

	insertToken := dtos.InsertToken{
		UserID: user.ID,
		Scope:  ScopeActivation,
		Expiry: time.Now().Add(ActivationExpiryTime).Unix(),
	}

	token, err := generetaToken(insertToken)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		ctx,
		queries.TokenInsert,
		token.Hash,
		token.UserID,
		token.Expiry,
	)
	if err != nil {
		return err
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return errCommitTransaction("committing transaction", err)
	}

	return nil
}
