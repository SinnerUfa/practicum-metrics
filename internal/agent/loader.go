package agent

import (
	"fmt"
	"math/rand"

	metrics "github.com/SinnerUfa/practicum-metric/internal/metrics"
	mlog "github.com/SinnerUfa/practicum-metric/internal/mlog"
	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
)

type MetricLoad struct {
	log     mlog.Logger
	rep     repository.Repository
	counter uint
}

func NewLoader(log mlog.Logger, rep repository.Repository) *MetricLoad {
	return &MetricLoad{log: log, rep: rep}
}

func (m *MetricLoad) Load() error {
	l := metrics.GetRuntimeMetrics()

	l = append(l, metrics.Metric{Name: "PollCount", Type: "counter", Value: fmt.Sprint(m.counter)})
	m.counter++
	l = append(l, metrics.Metric{Name: "RandomValue", Type: "gauge", Value: fmt.Sprint(rand.Int())})

	m.rep.SetList(l)
	m.log.Info("Load increment:", m.counter)
	return nil
}

func (m *MetricLoad) Tick() error {
	return m.Load()
}
