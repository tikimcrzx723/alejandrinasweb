package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRecordNotFound     = errors.New("record not found")
	ErrRecordConflict     = errors.New("record edit confilct")
	ErrReordInvalidUUID   = errors.New("invalid uuid")
	ErrRecordRelationShip = errors.New("record relation ship")
)

const (
	QueryTimeDuration = time.Second * 5
)

func errBeginTransaction(errStr string, err error) error {
	return fmt.Errorf("%s: %w", errStr, err)
}

func errCommitTransaction(errStr string, err error) error {
	return fmt.Errorf("%s: %w", errStr, err)
}

func withTx(db *pgxpool.Pool, ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
