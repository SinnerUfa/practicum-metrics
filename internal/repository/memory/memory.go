package memory

import (
	"sync"

	codes "github.com/SinnerUfa/practicum-metric/internal/codes"
	metrics "github.com/SinnerUfa/practicum-metric/internal/metrics"
	cnt "github.com/SinnerUfa/practicum-metric/internal/metrics/counter"
	gau "github.com/SinnerUfa/practicum-metric/internal/metrics/gauge"
)

type Memory struct {
	sync.RWMutex
	Counters map[string]*cnt.Counter
	Gauges   map[string]*gau.Gauge
}

func New() *Memory {
	return &Memory{
		Counters: make(map[string]*cnt.Counter),
		Gauges:   make(map[string]*gau.Gauge),
	}
}

func (mem *Memory) Set(m metrics.Metric) error {
	mem.Lock()
	defer mem.Unlock()
	if m.Name == "" {
		return codes.ErrRepMetricNotSupported
	}
	switch m.Type {
	case metrics.MetricTypeCounter:
		v, ok := m.Value.Int64()
		if !ok {
			return codes.ErrRepParseInt
		}
		if _, ok := mem.Counters[m.Name]; !ok {
			mem.Counters[m.Name] = &cnt.Counter{}
		}
		mem.Counters[m.Name].Set(v)

	case metrics.MetricTypeGauge:
		v, ok := m.Value.Float64()
		if !ok {
			return codes.ErrRepParseFloat
		}
		if _, ok := mem.Gauges[m.Name]; !ok {
			mem.Gauges[m.Name] = &gau.Gauge{}
		}
		mem.Gauges[m.Name].Set(v)
	default:
		return codes.ErrRepMetricNotSupported
	}
	return nil
}

func (mem *Memory) Get(m *metrics.Metric) error {
	mem.RLock()
	defer mem.RUnlock()
	switch m.Type {
	case metrics.MetricTypeCounter:
		c, ok := mem.Counters[m.Name]
		if !ok {
			return codes.ErrRepNotFound
		}
		m.Value = metrics.Int(c.Value())
	case metrics.MetricTypeGauge:
		g, ok := mem.Gauges[m.Name]
		if !ok {
			return codes.ErrRepNotFound
		}
		m.Value = metrics.Float(g.Value())
	default:
		return codes.ErrRepMetricNotSupported
	}
	return nil
}

func (mem *Memory) List() (out []metrics.Metric) {
	mem.RLock()
	defer mem.RUnlock()
	for k, v := range mem.Counters {
		out = append(out, metrics.Metric{Name: k, Value: metrics.Int(v.Value()), Type: metrics.MetricTypeCounter})
	}
	for k, v := range mem.Gauges {
		out = append(out, metrics.Metric{Name: k, Value: metrics.Float(v.Value()), Type: metrics.MetricTypeGauge})
	}
	return
}

func (mem *Memory) SetList(in []metrics.Metric) error {
	for _, v := range in {
		if err := mem.Set(v); err != nil {
			return err
		}
	}
	return nil
}
