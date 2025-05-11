package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"anti-brutforce/internal/app"
	"anti-brutforce/internal/bucket"
	"anti-brutforce/internal/config"
	"anti-brutforce/internal/transport/grpc"
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

	loginBuckets := bucket.NewMemoryBucketStorage(&cfg.Limiter.Login, logger)
	passwordBuckets := bucket.NewMemoryBucketStorage(&cfg.Limiter.Password, logger)
	ipBuckets := bucket.NewMemoryBucketStorage(&cfg.Limiter.IP, logger)

	defer func() {
		logger.Info("Closing resources...")
		loginBuckets.Close()
		passwordBuckets.Close()
		ipBuckets.Close()
		logger.Info("Resources closed properly")
	}()

	service := app.NewService(loginBuckets, passwordBuckets, ipBuckets)
	srv := grpc.NewServer(service, cfg.Grpc.Port)

	errCh := make(chan error, 1)
	go func() {
		logger.Info("Starting server", zap.String("port", cfg.Grpc.Port))
		if err := srv.Start(); err != nil {
			errCh <- fmt.Errorf("server error: %w", err)
		}
	}()

	// Ожидаем либо ошибки при запуске, либо сигнала завершения
	select {
	case err := <-errCh:
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
