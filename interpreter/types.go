package interpreter

import (
	"math"
	"strconv"
)

type Value interface {
	Kind() Kind
	Interface() any
	String() string
}

type Kind byte

const (
	KindInvalid Kind = iota
	KindVoid
	KindBool
	KindInt32
	KindFloat64
	KindString
)

func (k Kind) String() string {
	switch k {
	case KindVoid:
		return "void"
	case KindBool:
		return "bool"
	case KindInt32:
		return "int32"
	case KindFloat64:
		return "float64"
	case KindString:
		return "string"
	default:
		return "<invalid>"
	}
}

type Bool int32

func (b Bool) Kind() Kind {
	return KindBool
}

func (b Bool) Interface() any {
	return b > 0
}

func (b Bool) String() string {
	return strconv.FormatBool(b > 0)
}

type Int32 int32

func (i Int32) Kind() Kind {
	return KindInt32
}

func (i Int32) Interface() any {
	return int32(i)
}

func (i Int32) String() string {
	return strconv.Itoa(int(i))
}

type Float64 float64

func (f Float64) Kind() Kind {
	return KindFloat64
}

func (f Float64) Interface() any {
	return float64(f)
}

func (f Float64) String() string {
	if math.IsNaN(float64(f)) {
		return "NaN"
	}
	if math.IsInf(float64(f), 1) {
		return "Infinity"
	}
	if math.IsInf(float64(f), -1) {
		return "-Infinity"
	}
	return strconv.FormatFloat(float64(f), 'f', -1, 64)
}

type String string

func (s String) Kind() Kind {
	return KindString
}

func (s String) Interface() any {
	return string(s)
}

func (s String) String() string {
	return "\"" + string(s) + "\""
}
