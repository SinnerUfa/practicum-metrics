package main

import (
	"context"
	"os"
	"os/signal"

	config "github.com/SinnerUfa/practicum-metric/internal/config"
	mlog "github.com/SinnerUfa/practicum-metric/internal/mlog"
	server "github.com/SinnerUfa/practicum-metric/internal/server"
)

var cfg server.Config = server.DefaultConfig

func main() {
	log := mlog.New(true)

	if err := config.Load(&cfg, os.Args[1:]); err != nil {
		log.Error(err)
	}
	log.Info("\nCurrent configuration:", cfg)

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)

	if err := run(ctx, log, cfg); err != nil {
		log.Info(err)
		cancel()
		os.Exit(1)
	}
}

func run(ctx context.Context, log mlog.Logger, cfg server.Config) error {
	if err := server.Run(ctx, log, cfg); err != nil {
		return err
	}
	return nil
}
