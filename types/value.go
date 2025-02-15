package types

type Kind string

type Value interface {
	Kind() Kind
	Interface() any
	String() string
}

const KindUnknown Kind = "<unknown>"
const KindVoid Kind = "void"
