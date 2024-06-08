package flags

import (
	"errors"
	"flag"
	"reflect"
)
var (
	ExNoAcsess              = errors.New("No access to object")
	ExNotStructure          = errors.New("Object is not a structure")
	ExFieldNotSet           = errors.New("Field cannot be set")
	ExFieldTypeNotSupported = errors.New("Field type not supported")
)

func Load(v any, args []string) error {
	structValuePtr := reflect.ValueOf(v)
	structTypePtr := reflect.TypeOf(v)

	if structValuePtr.Kind() != reflect.Ptr {
		return ExNoAcsess
	}

	if structValuePtr.Elem().Kind() != reflect.Struct {
		return ExNotStructure
	}
	structValue := structValuePtr.Elem()
	structType := structTypePtr.Elem()
	flagMap := readFlags(structValue, structType)

	flags := flag.NewFlagSet("flags", flag.ExitOnError)
	if err := createFlags(flags, flagMap); err != nil {
		return err
	}
	if err := flags.Parse(args); err != nil {
		return err
	}
	if err := setFlags(structValue, flagMap); err != nil {
		return err
	}
	return nil
}

type Flag struct {
	name  string
	usage string
	ptr   any
	value reflect.Value
}

func readFlags(v reflect.Value, t reflect.Type) map[int]Flag {
	num := t.NumField()
	m := make(map[int]Flag, num)
	for i := 0; i < num; i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)
		tmp, ok := fieldType.Tag.Lookup("flag")
		if !ok {
			continue
		}
		if tmp == "-" || tmp == "" {
			continue
		}
		m[i] = Flag{tmp, fieldType.Name, nil, fieldValue}
	}
	return m
}

func createFlags(f *flag.FlagSet, list map[int]Flag) error {
	for k := range list {
		var tmp Flag
		tmp = list[k]
		switch tmp.value.Kind() {
		case reflect.Uint:
			tmp.ptr = f.Uint(tmp.name, uint(tmp.value.Uint()), tmp.usage)
		case reflect.String:
			tmp.ptr = f.String(tmp.name, tmp.value.String(), tmp.usage)
		case reflect.Int:
			tmp.ptr = f.Int(tmp.name, int(tmp.value.Int()), tmp.usage)
		default:
			return ExFieldTypeNotSupported
		}
		list[k] = tmp
	}
	return nil
}
func setFlags(v reflect.Value, list map[int]Flag) error {
	num := v.NumField()

	for i := 0; i < num; i++ {
		if _, ok := list[i]; !ok {
			continue
		}
		fieldValue := v.Field(i)
		if !fieldValue.CanSet() {
			return ExFieldNotSet
		}
		value := reflect.Indirect(reflect.ValueOf(list[i].ptr))

		switch fieldValue.Kind() {
		case reflect.Uint:
			fieldValue.SetUint(value.Uint())
		case reflect.String:
			fieldValue.SetString(value.String())
		case reflect.Int:
			fieldValue.SetInt(value.Int())
		default:
			return ExFieldTypeNotSupported
		}
	}
	return nil
}
