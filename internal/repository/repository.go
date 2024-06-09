package repository

import (
	metrics "github.com/SinnerUfa/practicum-metric/internal/metrics"
	memory "github.com/SinnerUfa/practicum-metric/internal/repository/memory"
)

type Repository interface {
	Set(m metrics.Metric) error
	Get(m *metrics.Metric) error
	List() (out []metrics.Metric)
	SetList([]metrics.Metric) error
}

func New() Repository {
	return memory.New()
}
