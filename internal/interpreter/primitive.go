package interpreter

import (
	"math"
	"strconv"
)

type Undefined struct{}

func (Undefined) Type() Type {
	return UNDEFINED
}

func (Undefined) Interface() any {
	return nil
}

func (Undefined) String() string {
	return "undefined"
}

type Null struct{}

func (Null) Type() Type {
	return NULL
}

func (Null) Interface() any {
	return nil
}

func (Null) String() string {
	return "null"
}

type Bool int32

func (b Bool) Type() Type {
	return BOOL
}

func (b Bool) Interface() any {
	return b > 0
}

func (b Bool) String() string {
	return strconv.FormatBool(b > 0)
}

type Int32 int32

func (i Int32) Type() Type {
	return INT32
}

func (i Int32) Interface() any {
	return int32(i)
}

func (i Int32) String() string {
	return strconv.Itoa(int(i))
}

type Float64 float64

func (f Float64) Type() Type {
	return FLOAT64
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

func (s String) Type() Type {
	return STRING
}

func (s String) Interface() any {
	return string(s)
}

func (s String) String() string {
	return "\"" + string(s) + "\""
}
