package token

type Type string

type Token struct {
	Type    Type
	Literal string
}

const (
	EOF     Type = "EOF"
	ILLEGAL Type = "ILLEGAL"

	INT    Type = "INT"
	FLOAT  Type = "FLOAT"
	STRING Type = "STRING"
	FALSE  Type = "FALSE"
	TRUE   Type = "TRUE"
	IDENT  Type = "IDENT"

	PLUS     Type = "+"
	MINUS    Type = "-"
	MULTIPLY Type = "*"
	DIVIDE   Type = "/"
	PERIOD   Type = "."
	LPAREN   Type = "("
	RPAREN   Type = ")"
)

var types = map[string]Type{
	"true":  TRUE,
	"false": FALSE,
}

func NewToken(typ Type, literal string) Token {
	return Token{Type: typ, Literal: literal}
}

func TypeOf(literal string) Type {
	typ, ok := types[literal]
	if !ok {
		typ = IDENT
	}
	return typ
}
