package server

import (
	"context"
	"log/slog"
	"net/http"

	unloader "github.com/SinnerUfa/practicum-metric/internal/unloader"

	codes "github.com/SinnerUfa/practicum-metric/internal/codes"
	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
)

func Run(ctx context.Context, log *slog.Logger, cfg Config) error {
	rep := repository.New()
	// создание rep с сохранением в файл сделать внутри new
	if cfg.Restore {
		unloader.Load(cfg.FileStoragePath, log, rep)
	}
	rep = unloader.Save(ctx, cfg.FileStoragePath, cfg.StoreInterval, log, rep)
	// создание rep с сохранением в файл сделать внутри new

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
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Warn("", "err", codes.ErrSrvShutdown, "server err:", err)
			return codes.ErrSrvShutdown
		}
	}
	return nil
}
