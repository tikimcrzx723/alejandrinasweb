package db

import (
	"context"
	"fmt"

	"github.com/tikimcrzx723/alejandrinasweb/internal/env"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New() (*pgxpool.Pool, error) {
	minConn := env.GetInt("DB_MIN_CONN", 3)
	maxConn := env.GetInt("DB_MAX_CONN", 100)
	user := env.GetString("DB_USER", "postgres")
	pass := env.GetString("DB_PASSWORD", "postgres")
	host := env.GetString("DB_HOST", "localhost")
	port := env.GetString("DB_PORT", "5432")
	dbName := env.GetString("DB_NAME", "tikim-ecommerce")
	sslMode := env.GetString("DB_SSL_MODE", "disable")

	dnsDB := makeDNS(user, pass, host, port, dbName, sslMode, minConn, maxConn)
	config, err := pgxpool.ParseConfig(dnsDB)
	if err != nil {
		return nil, fmt.Errorf("%s %w", "pgxpool.ParseConfig()", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("%s %w", "pgxpool.NewWithConfig", err)
	}

	return pool, nil
}

func makeDNS(user, pass, host, port, dbName, sslMode string, minConn, maxConn int) string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s pool_min_conns=%d pool_max_conns=%d",
		user,
		pass,
		host,
		port,
		dbName,
		sslMode,
		minConn,
		maxConn,
	)
}
