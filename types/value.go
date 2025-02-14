package types

import "fmt"

type Kind string

type Value interface {
	Kind() Kind
	Interface() any
	String() string
}

const (
	KindUnknown Kind = "<unknown>"
	KindVoid    Kind = "void"
	KindFloat64 Kind = "float64"
)

type Float64 struct {
	Value float64
}

func NewFloat64(value float64) Float64 {
	return Float64{Value: value}
}

func (v Float64) Kind() Kind {
	return KindFloat64
}

func (v Float64) Interface() any {
	return v.Value
}

func (v Float64) String() string {
	return fmt.Sprintf("%f", v.Value)
}
