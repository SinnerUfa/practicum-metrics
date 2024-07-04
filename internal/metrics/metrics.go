package metrics

import (
	"encoding/json"
	"fmt"

	"github.com/SinnerUfa/practicum-metric/internal/api/json_metric"
)

type MetricType string

const (
	MetricTypeGauge   MetricType = "gauge"
	MetricTypeCounter MetricType = "counter"
)

type Metric struct {
	Type  MetricType
	Name  string
	Value *Value
}

func (m Metric) ReguestString(head string) string {
	switch head {
	case "update":
		return fmt.Sprint(head, "/", m.Type, "/", m.Name, "/", m.Value)
	case "value":
		return fmt.Sprint(head, "/", m.Type, "/", m.Name)
	default:
		return ""
	}
}

func (m *Metric) UnmarshalJSON(data []byte) (err error) {
	type VisitorAlias json_metric.Metrics
	jm := json_metric.New()
	alias := &struct {
		*VisitorAlias
		Delta int64   `json:"delta,omitempty"`
		Value float64 `json:"value,omitempty"`
	}{
		VisitorAlias: (*VisitorAlias)(jm),
	}
	if err = json.Unmarshal(data, alias); err != nil {
		return
	}
	m.Name = alias.ID
	m.Type = MetricType(alias.MType)
	switch m.Type {
	case MetricTypeGauge:
		m.Value = Float(alias.Value)
	case MetricTypeCounter:
		m.Value = Int(alias.Delta)
	default:
		m.Value = Int(0)
	}
	return
}

func (m Metric) MarshalJSON() (data []byte, err error) {
	jm := json_metric.New()
	jm.ID = m.Name
	jm.MType = string(m.Type)
	switch m.Type {
	case MetricTypeGauge:
		*(jm.Value), _ = m.Value.Float64()
	case MetricTypeCounter:
		*(jm.Delta), _ = m.Value.Int64()
	default:
		*(jm.Value) = 0
	}
	return json.Marshal(*jm)
}
