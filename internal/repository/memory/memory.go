package memory

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	metrics "github.com/SinnerUfa/practicum-metric/internal/metrics"
	cnt "github.com/SinnerUfa/practicum-metric/internal/metrics/counter"
	gau "github.com/SinnerUfa/practicum-metric/internal/metrics/gauge"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrNotSupported = errors.New("this type of metrics is not supported")
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
	mem.RLock()
	defer mem.RUnlock()
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
		m.Value = fmt.Sprintf("%.4f", g.Value())
	default:
		return ErrNotSupported
	}
	// fmt.Println(mem)
	return nil
}

func (mem *Memory) List() (out []metrics.Metric) {
	mem.RLock()
	defer mem.RUnlock()
	for k, v := range mem.Counters {
		out = append(out, metrics.Metric{Name: k, Value: fmt.Sprint(v.Value()), Type: "counter"})
	}
	for k, v := range mem.Gauges {
		out = append(out, metrics.Metric{Name: k, Value: fmt.Sprintf("%.4f", v.Value()), Type: "gauge"})
	}
	// fmt.Println(mem)
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
