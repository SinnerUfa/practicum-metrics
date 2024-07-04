package agent

import (
	"log/slog"
	"math/rand"

	metrics "github.com/SinnerUfa/practicum-metric/internal/metrics"
	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
)

type MetricLoad struct {
	log     *slog.Logger
	rep     repository.Repository
	counter uint
}

func NewLoader(log *slog.Logger, rep repository.Repository) *MetricLoad {
	return &MetricLoad{log: log, rep: rep}
}

func (m *MetricLoad) Load() error {
	l := metrics.GetRuntimeMetrics()

	l = append(l, metrics.Metric{Name: "PollCount", Type: metrics.MetricTypeCounter, Value: metrics.Uint(m.counter)})
	m.counter++
	l = append(l, metrics.Metric{Name: "RandomValue", Type: metrics.MetricTypeGauge, Value: metrics.Int(rand.Int())})

	m.rep.SetList(l)
	m.log.Info("load metrics", "increment", m.counter)
	return nil
}

func (m *MetricLoad) Tick() error {
	return m.Load()
}
