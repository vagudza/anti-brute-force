package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"anti-brutforce/internal/config"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGHUP,
	)
	defer cancel()

	cfg, err := config.New()
	if err != nil {
		log.Print("can't create config", err)
		return
	}

	_ = ctx
	_ = cfg
}
