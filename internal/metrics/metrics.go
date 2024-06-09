package metrics

import (
	"fmt"
)

type Metric struct {
	Name  string
	Value string
	Type  string
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
