package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"os"

	"github.com/YL-Tan/GoHomeAi/internal/config"
	"github.com/YL-Tan/GoHomeAi/internal/logger"
	"github.com/YL-Tan/GoHomeAi/internal/server"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/YL-Tan/GoHomeAi/internal/db"
)

//go:embed db/migrations/*.sql
var embedMigrations embed.FS

func main() {
	if err := godotenv.Load(); err != nil {
		logger.Log.Fatal("No .env file found", zap.Error(err))
	}

	config.LoadConfig()
	logger.InitLogger()
	defer logger.Log.Sync()

	ctx := context.Background()
    poolConfig, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
    if err != nil {
        logger.Log.Fatal("Failed to initialize database pool", zap.Error(err))
    }
    pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
    if err != nil {
        logger.Log.Fatal("Failed to create pgx pool", zap.Error(err))
    }
    defer pool.Close()

    if err := pool.Ping(ctx); err != nil {
        logger.Log.Fatal("Failed to ping pgx pool", zap.Error(err))
    }

    if err := applyMigrations(pool); err != nil {
        logger.Log.Fatal("Migration failed", zap.Error(err))
    }

	dbInstance := stdlib.OpenDBFromPool(pool)
    queries := db.New(dbInstance)

	server.InitRouter(queries)

	port := viper.GetString("server.port")
	if port == "" {
		port = "8080"
	}

	logger.Log.Info("GoHomeAi server is running", zap.String("port", port))
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Log.Fatal("Server failed", zap.Error(err))
	}
}

func applyMigrations(pool *pgxpool.Pool) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	db := stdlib.OpenDBFromPool(pool)
	if err := goose.Up(db, "db/migrations"); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	return nil
}
