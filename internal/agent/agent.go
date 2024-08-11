package agent

import (
	"context"

	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
	ticker "github.com/SinnerUfa/practicum-metric/internal/ticker"
)

func Run(ctx context.Context, cfg Config) error {
	rep, err := repository.New(ctx, repository.Config{})
	if err != nil {
		return err
	}
	loader := NewLoader(rep)
	ticker.NewAndRun(ctx, cfg.PollInterval, loader)
	poster := NewPoster(ctx, rep, cfg.Adress)
	ticker.NewAndRun(ctx, cfg.ReportInterval, poster)

	<-ctx.Done()
	return nil
}
