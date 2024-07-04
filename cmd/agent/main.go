package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	agent "github.com/SinnerUfa/practicum-metric/internal/agent"
	config "github.com/SinnerUfa/practicum-metric/internal/config"

	mlog "github.com/SinnerUfa/practicum-metric/internal/mlog"
)

var cfg agent.Config = agent.DefaultConfig

func main() {

	log := mlog.New(mlog.SlogType)
	if err := config.Load(&cfg, os.Args[1:]); err != nil {
		log.Error("", "err", err)
	}
	log.Info("", "cfg", cfg)

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)

	if err := run(ctx, log, cfg); err != nil {
		cancel()
		log.Error("", "err", err)
	}
}

func run(ctx context.Context, log *slog.Logger, cfg agent.Config) error {
	if err := agent.Run(ctx, log, cfg); err != nil {
		return err
	}
	return nil
}
