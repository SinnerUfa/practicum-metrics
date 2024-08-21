package main

import (
	"context"
	slog "log/slog"
	"os"
	"os/signal"

	config "github.com/SinnerUfa/practicum-metric/internal/config"
	mlog "github.com/SinnerUfa/practicum-metric/internal/mlog"
	server "github.com/SinnerUfa/practicum-metric/internal/server"
)

func main() {
	cfg := server.Config{
		Adress:          server.DefaultAdress,
		StoreInterval:   server.DefaultStoreInterval,
		FileStoragePath: server.DefaultFileStoragePath,
		Restore:         server.DefaultRestore,
		DatabaseDSN:     server.DefaultDatabaseDSN,
	}
	slog.SetDefault(mlog.New(mlog.ZapType, slog.LevelDebug))

	if err := config.Load(&cfg, os.Args[1:]); err != nil {
		slog.Error("configuration load fail", "err", err)
		return
	}
	slog.Info("configuration loaded", "config", cfg)

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)

	if err := run(ctx, cfg); err != nil {
		cancel()
		slog.Error("exit with error", "err", err)
	}
}

func run(ctx context.Context, cfg server.Config) error {
	if err := server.Run(ctx, cfg); err != nil {
		return err
	}
	return nil
}
