package server

import (
	"context"
	"net/http"

	codes "github.com/SinnerUfa/practicum-metric/internal/codes"
	mlog "github.com/SinnerUfa/practicum-metric/internal/mlog"
	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
)

func Run(ctx context.Context, log mlog.Logger, cfg Config) error {
	rep := repository.New()
	httpServer := &http.Server{
		Addr:    cfg.Adress,
		Handler: Routes(log, rep),
	}
	errChan := make(chan error)
	go func(log mlog.Logger, ch chan error) {
		log.Info("start server on adress:", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Warning(codes.ErrSrvListen, " ", err)
			ch <- codes.ErrSrvListen
		}
	}(log, errChan)

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Warning(codes.ErrSrvShutdown, " ", err)
			return codes.ErrSrvShutdown
		}
	}
	return nil
}
