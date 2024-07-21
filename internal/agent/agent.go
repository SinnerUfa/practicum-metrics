package agent

import (
	"context"
	"log/slog"

	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
	ticker "github.com/SinnerUfa/practicum-metric/internal/ticker"
)

func Run(ctx context.Context, log *slog.Logger, cfg Config) error {
	rep, err := repository.New(ctx, repository.Config{Log: log})
	if err != nil {
		return err
	}
	loader := NewLoader(log, rep)
	ticker.NewAndRun(ctx, cfg.PollInterval, loader)
	poster := NewPoster(ctx, log, rep, cfg.Adress)
	ticker.NewAndRun(ctx, cfg.ReportInterval, poster)

	<-ctx.Done()
	return nil
}
