package memory

import (
	"errors"
	"fmt"
	"strconv"

	metrics "github.com/SinnerUfa/practicum-metric/internal/metrics"
	cnt "github.com/SinnerUfa/practicum-metric/internal/metrics/counter"
	gau "github.com/SinnerUfa/practicum-metric/internal/metrics/gauge"
)

var (
	ErrNotFound     = errors.New("Not found")
	ErrNotSupported = errors.New("This type of metrics is not supported")
)

type Memory struct {
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
	switch m.Type {
	case "counter":
		v, err := strconv.ParseInt(m.Value, 10, 64)
		if err != nil {
			return err
		}
		if _, ok := mem.Counters[m.Name]; !ok {
			mem.Counters[m.Name] = &cnt.Counter{}
		}
		mem.Counters[m.Name].Set(v)

	case "gauge":
		v, err := strconv.ParseFloat(m.Value, 64)
		if err != nil {
			return err
		}
		if _, ok := mem.Gauges[m.Name]; !ok {
			mem.Gauges[m.Name] = &gau.Gauge{}
		}
		mem.Gauges[m.Name].Set(v)
	default:
		return ErrNotSupported
	}
	// fmt.Println(mem)
	return nil
}

func (mem *Memory) Get(m *metrics.Metric) error {
	switch m.Type {
	case "counter":
		c, ok := mem.Counters[m.Name]
		if !ok {
			return ErrNotFound
		}
		m.Value = fmt.Sprint(c.Value())
	case "gauge":
		g, ok := mem.Gauges[m.Name]
		if !ok {
			return ErrNotFound
		}
		m.Value = fmt.Sprint(g.Value())
	default:
		return ErrNotSupported
	}
	// fmt.Println(mem)
	return nil
}

func (mem *Memory) List() (out []metrics.Metric) {
	for k, v := range mem.Counters {
		out = append(out, metrics.Metric{k, fmt.Sprint(v.Value()), "counter"})
	}
	for k, v := range mem.Gauges {
		out = append(out, metrics.Metric{k, fmt.Sprint(v.Value()), "gauge"})
	}
	// fmt.Println(mem)
	return
}
