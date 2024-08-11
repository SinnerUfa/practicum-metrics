package server

import (
	"context"
	"log/slog"
	"net/http"

	codes "github.com/SinnerUfa/practicum-metric/internal/codes"
	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
)

func Run(ctx context.Context, cfg Config) error {
	rep, err := repository.New(ctx,
		repository.Config{
			StoreInterval:   cfg.StoreInterval,
			FileStoragePath: cfg.FileStoragePath,
			Restore:         true,
			DatabaseDSN:     cfg.DatabaseDSN,
		})
	if err != nil {
		slog.Warn("repo error", "err", err)
		return err
	}
	httpServer := &http.Server{
		Addr:    cfg.Adress,
		Handler: Routes(rep.Storage()),
	}
	errChan := make(chan error)
	go func(ch chan error) {
		slog.Info("start server on adress", "Addr:", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Warn("", "err", codes.ErrSrvListen, "server err", err)
			ch <- codes.ErrSrvListen
		}
	}(errChan)

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		slog.Info("Shutdowning")
		if err := httpServer.Shutdown(ctx); err != nil {
			slog.Warn("", "err", codes.ErrSrvShutdown, "server err:", err)
			return codes.ErrSrvShutdown
		}
	}
	return nil
}
