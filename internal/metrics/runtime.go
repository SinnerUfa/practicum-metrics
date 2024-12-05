package metrics

import (
	"reflect"
	"runtime"
)

func GetRuntimeMetrics() []Metric {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	return Parse(m)
}

func Parse(v any) (m []Metric) {
	structValue := reflect.Indirect(reflect.ValueOf(v))
	structType := reflect.TypeOf(v)
	if structValue.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < structType.NumField(); i++ {
		fieldValue := structValue.Field(i)
		fieldType := structType.Field(i)
		var v *Value
		switch fieldType.Type.Kind() {
		case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
			v = Uint(fieldValue.Uint())
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			v = Int(fieldValue.Int())
		case reflect.Float32, reflect.Float64:
			v = Float(fieldValue.Float())
		}
		if v == nil {
			continue
		}
		m = append(m, Metric{
			Name:  fieldType.Name,
			Type:  MetricTypeGauge,
			Value: v,
		})
	}
	return
}
