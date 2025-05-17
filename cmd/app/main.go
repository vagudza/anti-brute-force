package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/vagudza/anti-brute-force/internal/app"
	"github.com/vagudza/anti-brute-force/internal/bucket"
	"github.com/vagudza/anti-brute-force/internal/config"
	"github.com/vagudza/anti-brute-force/internal/iplist"
	"github.com/vagudza/anti-brute-force/internal/storage"
	"github.com/vagudza/anti-brute-force/internal/transport/grpc"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGHUP,
	)
	defer cancel()

	logger := initLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.Error("app panic", zap.Any("panic", r), zap.Stack("stack"))
			os.Exit(1)
		}
	}()

	cfg, err := config.New()
	if err != nil {
		logger.Fatal("Failed to create config", zap.Error(err))
	}

	pgStorage, err := storage.NewStorage(ctx, &cfg.Postgres)
	if err != nil {
		logger.Fatal("Failed to create storage", zap.Error(err))
	}

	loginBuckets := bucket.NewMemoryBucketStorage(&cfg.Limiters.Login, logger)
	passwordBuckets := bucket.NewMemoryBucketStorage(&cfg.Limiters.Password, logger)
	ipBuckets := bucket.NewMemoryBucketStorage(&cfg.Limiters.IP, logger)

	defer func() {
		logger.Info("Closing resources...")
		err = loginBuckets.Close(ctx)
		if err != nil {
			logger.Error("Failed to close login buckets", zap.Error(err))
		}

		err = passwordBuckets.Close(ctx)
		if err != nil {
			logger.Error("Failed to close password buckets", zap.Error(err))
		}

		err = ipBuckets.Close(ctx)
		if err != nil {
			logger.Error("Failed to close IP buckets", zap.Error(err))
		}

		logger.Info("Resources closed properly")
	}()

	ipListService := iplist.NewService(pgStorage)
	service := app.NewService(
		logger,
		loginBuckets,
		passwordBuckets,
		ipBuckets,
		ipListService,
	)
	srv := grpc.NewServer(service, &cfg.Grpc)

	errCh := make(chan error, 1)
	go func() {
		logger.Info("Starting server", zap.String("port", cfg.Grpc.Port))
		if err := srv.Start(); err != nil {
			errCh <- fmt.Errorf("server error: %w", err)
		}
	}()

	select {
	case err = <-errCh:
		logger.Error("Server failed", zap.Error(err))
	case <-ctx.Done():
		logger.Info("Shutdown signal received")
	}

	// Создаем контекст с таймаутом для graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	logger.Info("Shutting down server...")
	if err := srv.Stop(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	} else {
		logger.Info("Server exited properly")
	}
}

func initLogger() *zap.Logger {
	return zap.Must(zap.NewDevelopment())
}
