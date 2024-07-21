package server

import (
	"context"
	"log/slog"
	"net/http"

	codes "github.com/SinnerUfa/practicum-metric/internal/codes"
	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
)

func Run(ctx context.Context, log *slog.Logger, cfg Config) error {
	rep, err := repository.New(ctx,
		repository.Config{
			StoreInterval:   cfg.StoreInterval,
			FileStoragePath: cfg.FileStoragePath,
			Restore:         true,
			DatabaseDSN:     cfg.DatabaseDSN,
			Log:             log,
		})
	if err != nil {
		log.Warn("repo error", "err", err)
		return err
	}
	httpServer := &http.Server{
		Addr:    cfg.Adress,
		Handler: Routes(log, rep),
	}
	errChan := make(chan error)
	go func(log *slog.Logger, ch chan error) {
		log.Info("start server on adress", "Addr:", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Warn("", "err", codes.ErrSrvListen, "server err", err)
			ch <- codes.ErrSrvListen
		}
	}(log, errChan)

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		log.Info("Shutdowning")
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Warn("", "err", codes.ErrSrvShutdown, "server err:", err)
			return codes.ErrSrvShutdown
		}
	}
	return nil
}
