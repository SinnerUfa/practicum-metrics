package agent

import (
	"context"
	"sync"

	mlog "github.com/SinnerUfa/practicum-metric/internal/mlog"
	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
)

func Run(ctx context.Context, log mlog.Logger, cfg Config) error {
	rep := repository.New()
	go func(log mlog.Logger, cfg Config, rep repository.Repository) {
		log.Info("start periodic functions")
		log.Info(cfg)
		log.Info(rep)
	}(log, cfg, rep)

	var wg sync.WaitGroup
	wg.Add(1)
	go func(log mlog.Logger) {
		defer wg.Done()
		<-ctx.Done()
		log.Info("stop periodic functions")
	}(log)
	wg.Wait()
	return nil
}
