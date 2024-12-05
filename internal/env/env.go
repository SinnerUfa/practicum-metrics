package env

import (
	"errors"
	"os"
	"reflect"
	"strconv"
)

var (
	ErrNoAcsess              = errors.New("no access to object")
	ErrNotStructure          = errors.New("object is not a structure")
	ErrFieldNotSet           = errors.New("field cannot be set")
	ErrFieldTypeNotSupported = errors.New("field type not supported")
)

func Load(v any) error {
	structValuePtr := reflect.ValueOf(v)
	structTypePtr := reflect.TypeOf(v)

	if structValuePtr.Kind() != reflect.Ptr {
		return ErrNoAcsess
	}

	if structValuePtr.Elem().Kind() != reflect.Struct {
		return ErrNotStructure
	}
	structValue := structValuePtr.Elem()
	structType := structTypePtr.Elem()

	for i := 0; i < structType.NumField(); i++ {
		fieldValue := structValue.Field(i)
		fieldType := structType.Field(i)
		tmp, ok := fieldType.Tag.Lookup("env")
		if !ok {
			continue
		}
		tmp, ok = GetEnv(tmp)
		if !ok {
			continue
		}
		err := SetEnv(fieldValue, tmp)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetEnv(key string) (string, bool) {
	if key == "-" || key == "" {
		return "", false
	}
	out := os.Getenv(key)
	if out == "" {
		return "", false
	}
	return out, true
}

func SetEnv(field reflect.Value, value string) error {
	if !field.CanSet() {
		return ErrFieldNotSet
	}
	switch field.Kind() {
	case reflect.Uint:
		newValaue, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return err
		}
		field.SetUint(newValaue)

	case reflect.String:
		field.SetString(value)
	case reflect.Int:
		newValaue, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return err
		}
		field.SetInt(newValaue)
	default:
		return ErrFieldTypeNotSupported
	}
	return nil
}
