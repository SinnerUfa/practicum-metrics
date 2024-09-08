package agent

import (
	"context"
	metrics "github.com/SinnerUfa/practicum-metric/internal/metrics"
	cpu "github.com/shirou/gopsutil/cpu"
	mem "github.com/shirou/gopsutil/v4/mem"
	slog "log/slog"
)

type GMetricLoad struct {
	setter  metrics.ListSetter
	counter uint
}

func GNewLoader(setter metrics.ListSetter) *GMetricLoad {
	return &GMetricLoad{setter: setter}
}

func (m *GMetricLoad) Load() {
	l := make([]metrics.Metric, 0)
	vs, err := mem.VirtualMemory()
	if err != nil {
		l = append(l, metrics.Metric{Name: "Total", Type: metrics.MetricTypeGauge, Value: metrics.Uint(vs.Total)})
		l = append(l, metrics.Metric{Name: "Free", Type: metrics.MetricTypeGauge, Value: metrics.Uint(vs.Free)})
	}
	p, err := cpu.Percent(0, false)
	if err != nil && len(p) != 0 {
		l = append(l, metrics.Metric{Name: "CPUutilization1", Type: metrics.MetricTypeGauge, Value: metrics.Float(p[0])})
	}
	if len(l) != 0 {
		m.setter.SetList(context.Background(), l)
		m.counter++
		slog.Debug("load gmetrics", "increment", m.counter)
		return
	}
	slog.Debug("load gmetrics", "len gmetrics = ", len(l))
}

func (m *GMetricLoad) Tick() {
	m.Load()
}
