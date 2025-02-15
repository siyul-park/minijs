package interpreter

type Mark byte
type Kind byte

const (
	PRIMITIVE = 0b10000000
	KIND      = 0b01110000
	SIZE      = 0b00001111
)

const (
	VOID    Kind = 0
	FLOAT64 Kind = 0b000010000
	STRING  Kind = 0b001000000
	UNKNOWN Kind = 0b011100000
)

var mnemonics = map[Kind]string{
	VOID:    "void",
	FLOAT64: "float64",
	STRING:  "string",
}

func (k Kind) String() string {
	v, ok := mnemonics[k]
	if !ok {
		return "unknown"
	}
	return v
}
