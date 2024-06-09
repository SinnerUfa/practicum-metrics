package hundlers

// import (
// 	"errors"
// 	"strings"

// 	metrics "github.com/SinnerUfa/practicum-metric/internal/metrics"
// )

// var (
// 	ErrBadReqStringType = errors.New("bad request string - type")
// 	ErrBadReqStringName = errors.New("bad request string - name")
// )

// func SplitURL(u string) (mm *metrics.Metric, err error) {
// 	s := strings.Split(strings.TrimLeft(u, " /"), "/")
// 	m := metrics.Metric{}
// 	switch s[0] {
// 	case "update":
// 		if len(s) != 4 {
// 			return nil, ErrBadReqStringName
// 		}
// 		m.Type = s[1]
// 		m.Name = s[2]
// 		m.Value = s[3]
// 	case "value":
// 		if len(s) != 3 {
// 			return nil, ErrBadReqStringName
// 		}
// 		m.Type = s[1]
// 		m.Name = s[2]
// 	default:
// 		return nil, ErrBadReqStringType
// 	}
// 	return &m, nil
// }
