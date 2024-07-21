package repository

import (
	"context"
	"log/slog"

	"github.com/SinnerUfa/practicum-metric/internal/repository/memory"

	"github.com/SinnerUfa/practicum-metric/internal/repository/unload"

	"github.com/SinnerUfa/practicum-metric/internal/metrics"
)

type Repository interface {
	Set(m metrics.Metric) error
	Get(m *metrics.Metric) error
	List() (out []metrics.Metric)
	SetList([]metrics.Metric) error
}

type Config struct {
	StoreInterval   uint
	FileStoragePath string
	Restore         bool
	DatabaseDSN     string
	Log             *slog.Logger
}

func New(ctx context.Context, cfg Config) (Repository, error) {
	if cfg.Restore {
		return unload.New(ctx, cfg.FileStoragePath, cfg.StoreInterval, cfg.Log)
	}
	return memory.New(), nil
}
