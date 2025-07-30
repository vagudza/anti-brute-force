package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/vagudza/anti-brute-force/internal/app"
	"github.com/vagudza/anti-brute-force/internal/bucket"
	"github.com/vagudza/anti-brute-force/internal/config"
	"github.com/vagudza/anti-brute-force/internal/iplist"
	"github.com/vagudza/anti-brute-force/internal/storage"
	"github.com/vagudza/anti-brute-force/internal/transport/grpc"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGHUP,
	)

	logger := initLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.Fatal("app panic", zap.Any("panic", r), zap.Stack("stack"))
		}
	}()
	defer cancel()

	cfg, err := config.New()
	if err != nil {
		logger.Fatal("failed to create config", zap.Error(err))
	}

	pgStorage, err := storage.NewStorage(ctx, &cfg.Postgres)
	if err != nil {
		logger.Fatal("failed to create storage", zap.Error(err))
	}

	loginBuckets := bucket.NewMemoryBucketStorage(&cfg.Limiters.Login, logger)
	passwordBuckets := bucket.NewMemoryBucketStorage(&cfg.Limiters.Password, logger)
	ipBuckets := bucket.NewMemoryBucketStorage(&cfg.Limiters.IP, logger)

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
		logger.Info("starting server", zap.String("port", cfg.Grpc.Port))
		if err = srv.Start(); err != nil {
			errCh <- fmt.Errorf("server error: %w", err)
		}
	}()

	select {
	case err = <-errCh:
		logger.Error("server failed", zap.Error(err))
		return
	case <-ctx.Done():
		logger.Info("shutdown signal received")
	}

	{
		logger.Info("graceful shutdown: shutting down grpc server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Stop(ctx); err != nil {
			logger.Error("server forced to shutdown", zap.Error(err))
			return
		}

		logger.Info("close buckets")
		if err = loginBuckets.Close(ctx); err != nil {
			logger.Error("failed to close login buckets", zap.Error(err))
		}
		if err = passwordBuckets.Close(ctx); err != nil {
			logger.Error("failed to close password buckets", zap.Error(err))
		}
		if err = ipBuckets.Close(ctx); err != nil {
			logger.Error("failed to close IP buckets", zap.Error(err))
		}

		logger.Info("close database connection")
		if err = pgStorage.Close(ctx); err != nil {
			logger.Error("failed to close database connection", zap.Error(err))
		}

		logger.Info("resources closed properly")
	}
}

func initLogger() *zap.Logger {
	return zap.Must(zap.NewDevelopment())
}
