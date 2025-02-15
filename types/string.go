package types

import "fmt"

type String string

const KindString Kind = "string"

func NewString(value string) String {
	return String(value)
}

func (v String) Kind() Kind {
	return KindString
}

func (v String) Interface() any {
	return string(v)
}

func (v String) String() string {
	return fmt.Sprintf("%q", string(v))
}
