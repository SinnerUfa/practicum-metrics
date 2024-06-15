package server

import (
	"context"
	"net/http"

	mlog "github.com/SinnerUfa/practicum-metric/internal/mlog"
	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
)

func Run(ctx context.Context, log mlog.Logger, cfg Config) error {
	rep := repository.New()
	httpServer := &http.Server{
		Addr:    cfg.Adress,
		Handler: Routes(log, cfg, rep),
	}
	go func(log mlog.Logger) {
		log.Info("start server on adress:", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("error listening and serving: ", err)
		}
	}(log)

	<-ctx.Done()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error("error shutdown server: ", err)
		return err
	}
	return nil
}
