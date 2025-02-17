package interpreter

type Value interface {
	Kind() Kind
	Interface() any
	String() string
}

type Kind byte

const (
	KindUnknown Kind = iota
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
