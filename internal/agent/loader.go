package agent

import (
	"log/slog"
	"math/rand"
	"reflect"
	"runtime"

	metrics "github.com/SinnerUfa/practicum-metric/internal/metrics"
	repository "github.com/SinnerUfa/practicum-metric/internal/repository"
)

type MetricLoad struct {
	rep     repository.Repository
	counter uint
}

func NewLoader(rep repository.Repository) *MetricLoad {
	return &MetricLoad{rep: rep}
}

func (m *MetricLoad) Load() {
	l := GetRuntimeMetrics()

	l = append(l, metrics.Metric{Name: "PollCount", Type: metrics.MetricTypeCounter, Value: metrics.Uint(m.counter)})
	m.counter++
	l = append(l, metrics.Metric{Name: "RandomValue", Type: metrics.MetricTypeGauge, Value: metrics.Int(rand.Int())})

	m.rep.SetList(l)
	slog.Debug("load metrics", "increment", m.counter)
}

func (m *MetricLoad) Tick() {
	m.Load()
}

func GetRuntimeMetrics() []metrics.Metric {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	return Parse(m)
}

func Parse(v any) (m []metrics.Metric) {
	structValue := reflect.Indirect(reflect.ValueOf(v))
	structType := reflect.TypeOf(v)
	if structValue.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < structType.NumField(); i++ {
		fieldValue := structValue.Field(i)
		fieldType := structType.Field(i)
		var v *metrics.Value
		switch fieldType.Type.Kind() {
		case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
			v = metrics.Uint(fieldValue.Uint())
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			v = metrics.Int(fieldValue.Int())
		case reflect.Float32, reflect.Float64:
			v = metrics.Float(fieldValue.Float())
		default:
			v = nil
		}
		if v == nil {
			continue
		}
		m = append(m, metrics.Metric{
			Name:  fieldType.Name,
			Type:  metrics.MetricTypeGauge,
			Value: v,
		})
	}
	return
}
