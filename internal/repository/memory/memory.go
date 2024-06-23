package memory

import (
	"fmt"
	"strconv"
	"strings"
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
	switch m.Type {
	case "counter":
		v, err := strconv.ParseInt(m.Value, 10, 64)
		if err != nil {
			return codes.ErrRepParseInt
		}
		if _, ok := mem.Counters[m.Name]; !ok {
			mem.Counters[m.Name] = &cnt.Counter{}
		}
		mem.Counters[m.Name].Set(v)

	case "gauge":
		v, err := strconv.ParseFloat(m.Value, 64)
		if err != nil {
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
	case "counter":
		c, ok := mem.Counters[m.Name]
		if !ok {
			return codes.ErrRepNotFound
		}
		m.Value = fmt.Sprint(c.Value())
	case "gauge":
		g, ok := mem.Gauges[m.Name]
		if !ok {
			return codes.ErrRepNotFound
		}
		m.Value = TrimFloat(g.Value())
	default:
		return codes.ErrRepMetricNotSupported
	}
	return nil
}

func (mem *Memory) List() (out []metrics.Metric) {
	mem.RLock()
	defer mem.RUnlock()
	for k, v := range mem.Counters {
		out = append(out, metrics.Metric{Name: k, Value: fmt.Sprint(v.Value()), Type: "counter"})
	}
	for k, v := range mem.Gauges {
		out = append(out, metrics.Metric{Name: k, Value: TrimFloat(v.Value()), Type: "gauge"})
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

func TrimFloat(v float64) string {
	s := strings.Split(fmt.Sprintf("%.5f", v), ".")
	s[1] = strings.TrimRight(s[1], ". 0")
	if s[1] == "" {
		return s[0]
	}
	return strings.Join(s, ".")
}
