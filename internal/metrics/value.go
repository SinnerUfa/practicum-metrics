package metrics

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	cstr "golang.org/x/exp/constraints"
)

type valueKind int

const (
	valueKindNil valueKind = iota
	valueKindString
	valueKindInt64
	valueKindUint64
	valueKindFloat64
)

type Value struct {
	kind valueKind
	data any
}

func String(v string) *Value {
	return &Value{kind: valueKindString, data: v}
}

func Float[T cstr.Float](v T) *Value {
	return &Value{kind: valueKindFloat64, data: float64(v)}
}

func Int[T cstr.Signed](v T) *Value {
	return &Value{kind: valueKindInt64, data: int64(v)}
}

func Uint[T cstr.Unsigned](v T) *Value {
	return &Value{kind: valueKindUint64, data: uint64(v)}
}

func (v *Value) IsString() bool {
	return v.kind == valueKindString
}
func (v *Value) IsFloat() bool {
	return v.kind == valueKindFloat64
}

func (v *Value) IsInt() bool {
	return v.kind == valueKindInt64
}

func (v *Value) IsUint() bool {
	return v.kind == valueKindUint64
}

func (v *Value) ToString() (out *Value, ok bool) {
	switch v.kind {
	case valueKindNil:
		out, ok = nil, false
	case valueKindString:
		out, ok = v, true
	case valueKindFloat64:
		ss := strings.Split(fmt.Sprintf("%.5f", v.data.(float64)), ".")
		ss[1] = strings.TrimRight(ss[1], ". 0")
		s := ""
		if ss[1] == "" {
			s = ss[0]
		} else {
			s = strings.Join(ss, ".")
		}
		out, ok = &Value{kind: valueKindString, data: s}, true
	case valueKindInt64:
		out, ok = &Value{kind: valueKindString, data: fmt.Sprint(v.data.(int64))}, true
	case valueKindUint64:
		out, ok = &Value{kind: valueKindString, data: fmt.Sprint(v.data.(uint64))}, true
	}
	return
}

func (v *Value) ToFloat() (out *Value, ok bool) {
	switch v.kind {
	case valueKindNil:
		out, ok = nil, false
	case valueKindString:
		f, err := strconv.ParseFloat(v.data.(string), 64)
		if err != nil {
			out, ok = nil, false
			break
		}
		out, ok = &Value{kind: valueKindFloat64, data: f}, true
	case valueKindFloat64:
		out, ok = v, true
	case valueKindInt64:
		out, ok = &Value{kind: valueKindFloat64, data: float64(v.data.(int64))}, true
	case valueKindUint64:
		out, ok = &Value{kind: valueKindFloat64, data: float64(v.data.(uint64))}, true
	}
	return
}

func (v *Value) ToInt() (out *Value, ok bool) {
	switch v.kind {
	case valueKindNil:
		out, ok = nil, false
	case valueKindString:
		f, err := strconv.ParseInt(v.data.(string), 10, 64)
		if err != nil {
			out, ok = nil, false
			break
		}
		out, ok = &Value{kind: valueKindInt64, data: f}, true
	case valueKindFloat64:
		out, ok = &Value{kind: valueKindInt64, data: int64(v.data.(float64))}, true
	case valueKindInt64:
		out, ok = v, true
	case valueKindUint64:
		if v.data.(uint64) > math.MaxInt64 {
			out, ok = nil, false
			break
		}
		out, ok = &Value{kind: valueKindInt64, data: int64(v.data.(uint64))}, true
	}
	return
}

func (v *Value) ToUint() (out *Value, ok bool) {
	switch v.kind {
	case valueKindNil:
		out, ok = nil, false
	case valueKindString:
		f, err := strconv.ParseUint(v.data.(string), 10, 64)
		if err != nil {
			out, ok = nil, false
			break
		}
		out, ok = &Value{kind: valueKindUint64, data: f}, true
	case valueKindFloat64:
		out, ok = &Value{kind: valueKindUint64, data: uint64(v.data.(float64))}, true
	case valueKindInt64:
		if v.data.(int64) < 0 {
			out, ok = nil, false
			break
		}
		out, ok = &Value{kind: valueKindUint64, data: uint64(v.data.(int64))}, true
	case valueKindUint64:
		out, ok = v, true
	}
	return
}

func (v *Value) String() string {
	switch v.kind {
	case valueKindNil:
		return ""
	case valueKindString:
		return v.data.(string)
	}
	t, ok := v.ToString()
	if !ok {
		return ""
	}
	return t.String()
}

func (v *Value) Float64() (out float64, ok bool) {
	switch v.kind {
	case valueKindNil:
		return 0, false
	case valueKindFloat64:
		return v.data.(float64), true
	}
	t, ok := v.ToFloat()
	if !ok {
		return 0, false
	}
	return t.Float64()
}

func (v *Value) Uint64() (out uint64, ok bool) {
	switch v.kind {
	case valueKindNil:
		return 0, false
	case valueKindUint64:
		return v.data.(uint64), true
	}
	t, ok := v.ToUint()
	if !ok {
		return 0, false
	}
	return t.Uint64()
}

func (v *Value) Int64() (out int64, ok bool) {
	switch v.kind {
	case valueKindNil:
		return 0, false
	case valueKindInt64:
		return v.data.(int64), true
	}
	t, ok := v.ToInt()
	if !ok {
		return 0, false
	}
	return t.Int64()
}
