package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/YL-Tan/GoHomeAi/internal/config"
	"github.com/YL-Tan/GoHomeAi/internal/db"
	"github.com/YL-Tan/GoHomeAi/internal/logger"
	"github.com/YL-Tan/GoHomeAi/internal/server"
	"github.com/YL-Tan/GoHomeAi/internal/workers"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	initSystem()

	workerPool := workers.NewWorkerPool()
	workerPool.Start()

	database := initDatabase(ctx)
	defer database.Pool.Close()

	wsServer := server.NewWebSocketServer()
	wsServer.Start(workerPool)

	server := startServer(database.Queries, workerPool, wsServer)

	go monitorSystemHealth(ctx)
	go enqueueBackgroundJobs(workerPool)

	waitForShutdown(ctx, server, workerPool)
}

func initSystem() {
	if err := godotenv.Load(); err != nil {
		logger.Log.Fatal("No .env file found", zap.Error(err))
	}
	config.LoadConfig()
	logger.InitLogger()
	defer logger.Log.Sync()
}

func initDatabase(ctx context.Context) *db.Database {
	database, err := db.InitDB(ctx)
	if err != nil {
		logger.Log.Fatal("Failed to initialize database", zap.Error(err))
	}
	return database
}

func startServer(queries *db.Queries, workerPool *workers.WorkerPool, wsServer *server.WebSocketServer) *http.Server {
	server := &http.Server{
		Addr:    ":" + viper.GetString("server.port"),
		Handler: server.InitRouter(queries, workerPool, wsServer),
	}

	go func() {
		logger.Log.Info("GoHomeAi server is running", zap.String("port", viper.GetString("server.port")))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("Server failed", zap.Error(err))
		}
	}()

	return server
}

func monitorSystemHealth(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			logger.Log.Info("Performing system health check...")
			// Add any health checks you need here (e.g., DB, CPU, memory, etc.)
		case <-ctx.Done():
			logger.Log.Info("Stopping system health monitoring")
			return
		}
	}
}

func enqueueBackgroundJobs(workerPool *workers.WorkerPool) {
	for i := 1; i <= 100; i++ {
		workerPool.AddJob(workers.Job{ID: i, Message: "Processing AI Task"})
		time.Sleep(50 * time.Millisecond)
	}
}

func waitForShutdown(ctx context.Context, server *http.Server, workerPool *workers.WorkerPool) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.Log.Info("Shutdown signal received, cleaning up...")

	workerPool.Stop()

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Log.Fatal("Server shutdown failed", zap.Error(err))
	}

	logger.Log.Info("GoHomeAi server shutdown complete.")
}
