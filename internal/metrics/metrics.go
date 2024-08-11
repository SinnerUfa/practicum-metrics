package metrics

import (
	"encoding/json"
	"fmt"
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

type MetricJSON struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func (m *Metric) UnmarshalJSON(data []byte) (err error) {
	alias := &MetricJSON{}
	if err = json.Unmarshal(data, alias); err != nil {
		return
	}
	m.Name = alias.ID
	m.Type = MetricType(alias.MType)
	switch m.Type {
	case MetricTypeGauge:
		m.Value = Float(*alias.Value)
	case MetricTypeCounter:
		m.Value = Int(*alias.Delta)
	default:
		m.Value = Int(0)
	}
	return
}

func (m Metric) MarshalJSON() (data []byte, err error) {
	alias := &MetricJSON{
		ID:    m.Name,
		MType: string(m.Type),
		Delta: new(int64),
		Value: new(float64),
	}
	switch m.Type {
	case MetricTypeGauge:
		*alias.Value, _ = m.Value.Float64()
		alias.Delta = nil
	case MetricTypeCounter:
		*alias.Delta, _ = m.Value.Int64()
		alias.Value = nil
	default:
		alias.Value = nil
		alias.Delta = nil
	}
	return json.Marshal(&alias)
}
