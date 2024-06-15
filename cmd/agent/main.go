package main

import (
	"context"
	"os"
	"os/signal"

	agent "github.com/SinnerUfa/practicum-metric/internal/agent"
	config "github.com/SinnerUfa/practicum-metric/internal/config"
	mlog "github.com/SinnerUfa/practicum-metric/internal/mlog"
)

var cfg agent.Config = agent.DefaultConfig

func main() {
	log := mlog.New(true)

	if err := config.Load(&cfg, os.Args[1:]); err != nil {
		log.Warning(err)
		os.Exit(1)
	}
	log.Info("\nCurrent configuration:", cfg)

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)

	if err := run(ctx, log, cfg); err != nil {
		log.Warning(err)
		cancel()
		os.Exit(1)
	}
}

func run(ctx context.Context, log mlog.Logger, cfg agent.Config) error {
	if err := agent.Run(ctx, log, cfg); err != nil {
		return err
	}
	return nil
}
