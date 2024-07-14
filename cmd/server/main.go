package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	config "github.com/SinnerUfa/practicum-metric/internal/config"
	mlog "github.com/SinnerUfa/practicum-metric/internal/mlog"
	server "github.com/SinnerUfa/practicum-metric/internal/server"
)

var cfg server.Config = server.DefaultConfig

func main() {
	log := mlog.New(mlog.ZapType)

	if err := config.Load(&cfg, os.Args[1:]); err != nil {
		log.Error("", "err", err)
		return
	}
	log.Info("", "cfg", cfg)

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)

	if err := run(ctx, log, cfg); err != nil {
		cancel()
		log.Error("", "err", err)
	}
}

func run(ctx context.Context, log *slog.Logger, cfg server.Config) error {
	if err := server.Run(ctx, log, cfg); err != nil {
		return err
	}
	return nil
}
