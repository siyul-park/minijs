package interpreter

type Value interface {
	Type() Type
	Interface() any
}

type Type byte

const (
	UNKNOWN Type = iota
	VOID
	UNDEFINED
	NULL
	BOOL
	INT32
	FLOAT64
	STRING
	OBJECT
)

func (t Type) String() string {
	switch t {
	case VOID:
		return "void"
	case UNDEFINED:
		return "undefined"
	case NULL:
		return "null"
	case BOOL:
		return "bool"
	case INT32:
		return "int32"
	case FLOAT64:
		return "float64"
	case STRING:
		return "string"
	case OBJECT:
		return "object"
	default:
		return "<invalid>"
	}
}
