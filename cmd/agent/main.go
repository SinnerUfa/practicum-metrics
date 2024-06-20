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
		log.Fatal(err)
	}
	log.Info("\nCurrent configuration:", cfg)

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)

	if err := run(ctx, log, cfg); err != nil {
		cancel()
		log.Fatal(err)
	}
}

func run(ctx context.Context, log mlog.Logger, cfg agent.Config) error {
	if err := agent.Run(ctx, log, cfg); err != nil {
		return err
	}
	return nil
}
