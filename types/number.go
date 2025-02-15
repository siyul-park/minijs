package types

import (
	"strconv"
)

type Float64 float64

const KindFloat64 Kind = "float64"

func NewFloat64(value float64) Float64 {
	return Float64(value)
}

func (v Float64) Kind() Kind {
	return KindFloat64
}

func (v Float64) Interface() any {
	return float64(v)
}

func (v Float64) String() string {
	return strconv.FormatFloat(float64(v), 'f', -1, 64)
}
