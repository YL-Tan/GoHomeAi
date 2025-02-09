package db

import (
	"context"
	"embed"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"

	"github.com/YL-Tan/GoHomeAi/internal/logger"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type Database struct {
	Pool    *pgxpool.Pool
	Queries *Queries
}

func InitDB(ctx context.Context) (*Database, error) {
	dbURL := os.Getenv("DATABASE_URL")

	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		logger.Log.Fatal("Failed to initialize database pool", zap.Error(err))
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		logger.Log.Fatal("Failed to create pgx pool", zap.Error(err))
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		logger.Log.Fatal("Failed to ping pgx pool", zap.Error(err))
		return nil, err
	}

	if err := runMigrations(pool); err != nil {
		pool.Close()
		logger.Log.Fatal("Migration failed", zap.Error(err))
		return nil, err
	}

	// Convert pgxpool.Pool to *sql.DB for sqlc
	queries := New(stdlib.OpenDBFromPool(pool))

	return &Database{Pool: pool, Queries: queries}, nil
}

func runMigrations(pool *pgxpool.Pool) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	db := stdlib.OpenDBFromPool(pool)
	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	return nil
}
