package database

import (
	"context"
	"database/sql"

	"github.com/SinnerUfa/practicum-metric/internal/codes"
	"github.com/SinnerUfa/practicum-metric/internal/metrics"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database struct {
	db *sql.DB
}

func New(ctx context.Context, dsn string) (*Database, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	// создание таблиц
	return &Database{db: db}, nil
}

func (d *Database) SetContext(ctx context.Context, m metrics.Metric) error {
	if m.Name == "" {
		return codes.ErrRepMetricNotSupported
	}
	switch m.Type {
	case metrics.MetricTypeCounter:
	case metrics.MetricTypeGauge:
	default:
		return codes.ErrRepMetricNotSupported
	}
	return nil
}

func (d *Database) GetContext(ctx context.Context, m *metrics.Metric) error {
	if m.Name == "ping" && m.Type == "" && m.Value.IsString() && m.Value.String() == "ping" {
		return d.db.PingContext(ctx)
	}
	switch m.Type {
	case metrics.MetricTypeCounter:
	case metrics.MetricTypeGauge:
	default:
		return codes.ErrRepMetricNotSupported
	}
	return nil
}

func (d *Database) GetListContext(ctx context.Context) (out []metrics.Metric) {
	return
}

func (d *Database) SetListContext(ctx context.Context, in []metrics.Metric) error {
	return nil
}

func (d *Database) Close() error {
	return d.db.Close()
}
