package agent

import (
	"context"
	"sync"

	mlog "github.com/SinnerUfa/practicum-metric/internal/mlog"
	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
	ticker "github.com/SinnerUfa/practicum-metric/internal/ticker"
)

func Run(ctx context.Context, log mlog.Logger, cfg Config) error {
	rep := repository.New()
	loader := NewLoader(log, rep)
	ticker.NewAndRun(ctx, cfg.PollInterval, loader)
	poster := NewPoster(ctx, log, rep, cfg.Adress)
	ticker.NewAndRun(ctx, cfg.ReportInterval, poster)

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
