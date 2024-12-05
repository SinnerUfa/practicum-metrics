package server

import (
	"context"
	"net/http"
	"sync"

	mlog "github.com/SinnerUfa/practicum-metric/internal/mlog"
	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
)

func Run(ctx context.Context, log mlog.Logger, cfg Config) error {
	rep := repository.New()
	hundler := NewMainHundler(log, cfg, rep)
	httpServer := &http.Server{
		Addr:    cfg.Adress,
		Handler: hundler,
	}
	go func(log mlog.Logger) {
		log.Info("start server on adress:", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("error listening and serving: ", err)
		}
	}(log)

	var wg sync.WaitGroup
	wg.Add(1)
	go func(log mlog.Logger) {
		defer wg.Done()
		<-ctx.Done()
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Error("error shutdown server: ", err)
		}
	}(log)
	wg.Wait()
	return nil
}
