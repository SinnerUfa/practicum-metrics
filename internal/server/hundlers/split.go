package hundlers

import (
	"errors"
	"strings"

	metrics "github.com/SinnerUfa/practicum-metric/internal/metrics"
)

var (
	ExBadReqStringType = errors.New("Bad request string - type")
	ExBadReqStringName = errors.New("Bad request string - name")
)

func SplitURL(u string) (mm *metrics.Metric, err error) {
	// fmt.Println("|", u, "|")
	s := strings.Split(strings.TrimLeft(u, " /"), "/")
	// for i := range s {
	// 	fmt.Println("|", s[i], "|")
	// }

	m := metrics.Metric{}
	switch s[0] {
	case "update":
		if len(s) != 4 {
			return nil, ExBadReqStringName
		}
		m.Type = s[1]
		m.Name = s[2]
		m.Value = s[3]
	case "value":
		if len(s) != 3 {
			return nil, ExBadReqStringName
		}
		m.Type = s[1]
		m.Name = s[2]
	default:
		return nil, ExBadReqStringType
	}
	return &m, nil
}
