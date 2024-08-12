package server

import (
	"context"
	"log/slog"
	"net/http"

	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
)

func Run(ctx context.Context, cfg Config) error {
	rep, err := repository.New(ctx,
		repository.Config{
			StoreInterval:   cfg.StoreInterval,
			FileStoragePath: cfg.FileStoragePath,
			Restore:         cfg.Restore,
			DatabaseDSN:     cfg.DatabaseDSN,
		})
	if err != nil {
		slog.Warn("repository start with error", "err", err)
		return err
	}
	slog.Info("repository open")

	httpServer := &http.Server{
		Addr:    cfg.Adress,
		Handler: Routes(rep.Storage()),
	}
	errChan := make(chan error)
	go func(ch chan error) {
		slog.Info("start HTTP server on adress", "Addr:", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Warn("HTTP server stop with error", "server err", err)
			ch <- err
		}
	}(errChan)

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		slog.Info("Shutdowning...")
		if err := httpServer.Shutdown(ctx); err != nil {
			slog.Warn("HTTP server shutdown with error", "server err:", err)
			return err
		}
		slog.Info("HTTP server shutdowned")
		rep.Close()
		slog.Info("repository closed")
	}
	return nil
}
