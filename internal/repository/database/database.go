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
	return &Database{db: db}, nil
}

func (d *Database) Set(m metrics.Metric) error {
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

func (d *Database) Get(m *metrics.Metric) error {
	if m.Name == "ping" && m.Type == "" && m.Value.IsString() && m.Value.String() == "ping" {
		return d.db.Ping()
	}
	switch m.Type {
	case metrics.MetricTypeCounter:
	case metrics.MetricTypeGauge:
	default:
		return codes.ErrRepMetricNotSupported
	}
	return nil
}

func (d *Database) GetList() (out []metrics.Metric) {
	return
}

func (d *Database) SetList(in []metrics.Metric) error {
	return nil
}

func (d *Database) Close() error {
	return d.db.Close()
}
