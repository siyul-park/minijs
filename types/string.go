package types

import "fmt"

type String struct {
	Value string
}

const KindString Kind = "string"

func NewString(value string) String {
	return String{Value: value}
}

func (v String) Kind() Kind {
	return KindString
}

func (v String) Interface() any {
	return v.Value
}

func (v String) String() string {
	return fmt.Sprintf("%q", v.Value)
}
