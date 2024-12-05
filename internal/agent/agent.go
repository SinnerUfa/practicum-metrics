package agent

import (
	"context"
	"log/slog"

	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
	ticker "github.com/SinnerUfa/practicum-metric/internal/ticker"
)

func Run(ctx context.Context, cfg Config) error {
	rep, err := repository.New(ctx, repository.Config{})
	if err != nil {
		slog.Warn("repository start with error", "err", err)
		return err
	}
	slog.Info("repository open")
	slog.Info("repository storge type", "type", rep.Type())

	loader := ticker.NewAndRun(ctx, cfg.PollInterval, NewLoader(rep.Storage()))
	poster := ticker.NewAndRun(ctx, cfg.ReportInterval, NewPoster(rep.Storage(), cfg.Adress))

	<-ctx.Done()
	loader.Close()
	poster.Close()
	return nil
}
