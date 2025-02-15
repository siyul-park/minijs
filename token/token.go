package token

type Type string

type Token struct {
	Type    Type   // 토큰 타입
	Literal string // 토큰의 실제 값
}

const (
	EOF     Type = "EOF"
	ILLEGAL Type = "ILLEGAL"

	IDENTIFIER Type = "IDENTIFIER"
	NUMBER     Type = "NUMBER"
	STRING     Type = "STRING"
	BOOLEAN    Type = "BOOLEAN"
	NULL       Type = "NULL"

	PLUS     Type = "＋"
	MINUS    Type = "−"
	ASTERISK Type = "*"
	SLASH    Type = "/"
	PERCENT  Type = "％"
	EQUAL    Type = "＝"
	LESS     Type = "＜"
	GREATER  Type = "＞"
	BANG     Type = "!"
	AMP      Type = "&"
	VBAR     Type = "|"

	LPAREN    Type = "("
	RPAREN    Type = ")"
	LBRACE    Type = "{"
	RBRACE    Type = "}"
	COMMA     Type = ","
	SEMICOLON Type = ";"
	DOT       Type = "."

	IF     Type = "if"
	ELSE   Type = "else"
	WHILE  Type = "while"
	FOR    Type = "for"
	RETURN Type = "return"

	VAR   Type = "var"
	LET   Type = "let"
	CONST Type = "const"
	FUNC  Type = "function"
)

var types = map[string]Type{
	"true":      BOOLEAN,
	"false":     BOOLEAN,
	"NaN":       NUMBER,
	"Infinity":  NUMBER,
	"null":      NULL,
	"undefined": IDENTIFIER,

	"if":     IF,
	"else":   ELSE,
	"while":  WHILE,
	"for":    FOR,
	"return": RETURN,

	"var":      VAR,
	"let":      LET,
	"const":    CONST,
	"function": FUNC,
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
