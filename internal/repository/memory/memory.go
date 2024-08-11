package memory

import (
	"github.com/SinnerUfa/practicum-metric/internal/codes"
	"github.com/SinnerUfa/practicum-metric/internal/metrics"
)

type Counter struct {
	value int64
}

func (c *Counter) Set(i int64) {
	c.value += i
}

func (c *Counter) Value() int64 {
	return c.value
}

type Gauge struct {
	value float64
}

func (c *Gauge) Set(i float64) {
	c.value = i
}

func (c *Gauge) Value() float64 {
	return c.value
}

type Memory struct {
	// sync.RWMutex
	Counters map[string]*Counter
	Gauges   map[string]*Gauge
}

func New() *Memory {
	return &Memory{
		Counters: make(map[string]*Counter),
		Gauges:   make(map[string]*Gauge),
	}
}

func (mem *Memory) Set(m metrics.Metric) error {
	if m.Name == "" {
		return codes.ErrRepMetricNotSupported
	}
	switch m.Type {
	case metrics.MetricTypeCounter:
		v, ok := m.Value.Int64()
		if !ok {
			return codes.ErrRepParseInt
		}
		// mem.Lock()
		if _, ok := mem.Counters[m.Name]; !ok {
			mem.Counters[m.Name] = &Counter{}
		}
		mem.Counters[m.Name].Set(v)
		// mem.Unlock()
	case metrics.MetricTypeGauge:
		v, ok := m.Value.Float64()
		if !ok {
			return codes.ErrRepParseFloat
		}
		// mem.Lock()
		if _, ok := mem.Gauges[m.Name]; !ok {
			mem.Gauges[m.Name] = &Gauge{}
		}
		mem.Gauges[m.Name].Set(v)
		// mem.Unlock()
	default:
		return codes.ErrRepMetricNotSupported
	}
	return nil
}

func (mem *Memory) Get(m *metrics.Metric) error {
	switch m.Type {
	case metrics.MetricTypeCounter:
		// mem.RLock()
		c, ok := mem.Counters[m.Name]
		// mem.RUnlock()
		if !ok {
			return codes.ErrRepNotFound
		}
		m.Value = metrics.Int(c.Value())
	case metrics.MetricTypeGauge:
		// mem.RLock()
		g, ok := mem.Gauges[m.Name]
		// mem.RUnlock()
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
	// mem.RLock()
	// defer mem.RUnlock()
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
