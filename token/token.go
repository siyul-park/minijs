package token

type Type string

type Token struct {
	Type    Type
	Literal string
}

const (
	EOF     Type = "EOF"
	ILLEGAL Type = "ILLEGAL"

	DECIMAL     Type = "DECIMAL"
	EXPONENTIAL Type = "EXPONENTIAL"
	BINARY      Type = "BINARY"
	STRING      Type = "STRING"
	FALSE       Type = "FALSE"
	TRUE        Type = "TRUE"
	IDENTIFIER  Type = "IDENTIFIER"

	PLUS     Type = "+"
	MINUS    Type = "-"
	MULTIPLY Type = "*"
	DIVIDE   Type = "/"
	MODULO   Type = "%"
	PERIOD   Type = "."
	LPAREN   Type = "("
	RPAREN   Type = ")"

	NAN      Type = "NAN"
	INFINITY Type = "INFINITY"
)

var types = map[string]Type{
	"true":     TRUE,
	"false":    FALSE,
	"NaN":      NAN,
	"Infinity": INFINITY,
}

func NewToken(typ Type, literal string) Token {
	return Token{Type: typ, Literal: literal}
}

func TypeOf(literal string) Type {
	typ, ok := types[literal]
	if !ok {
		typ = IDENTIFIER
	}
	return typ
}
