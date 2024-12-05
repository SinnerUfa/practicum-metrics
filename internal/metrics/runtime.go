package metrics

import (
	"fmt"
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

		switch fieldType.Type.Kind() {
		case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
			fallthrough
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			fallthrough
		case reflect.Float32, reflect.Float64:
			m = append(m, Metric{Name: fieldType.Name, Type: "gauge", Value: fmt.Sprint(fieldValue)})
		}
	}
	return
}
