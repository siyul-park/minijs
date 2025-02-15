package token

type Type string

type Token struct {
	Type    Type
	Literal string
}

const (
	EOF        Type = "EOF"
	ILLEGAL    Type = "ILLEGAL"
	IDENTIFIER Type = "IDENTIFIER"
	NUMBER     Type = "NUMBER"
	BOOLEAN    Type = "BOOLEAN"
	STRING     Type = "STRING"
	NULL       Type = "NULL"
	UNDEFINED  Type = "UNDEFINED"
	KEYWORD    Type = "KEYWORD"

	PLUS                        Type = "＋"
	MINUS                       Type = "−"
	MULTIPLE                    Type = "*"
	DIVIDE                      Type = "/"
	MODULAR                     Type = "％"
	ASSIGN                      Type = "＝"
	LESS_THAN                   Type = "＜"
	GREATER_THAN                Type = "＞"
	NOT                         Type = "!"
	BIT_AND                     Type = "&"
	BIT_OR                      Type = "|"
	BIT_XOR                     Type = "^"
	BIT_NOT                     Type = "~"
	EQUAL                       Type = "=="
	NOT_EQUAL                   Type = "!="
	LESS_THAN_EQUAL             Type = "<="
	GREATER_THAN_EQUAL          Type = ">="
	AND                         Type = "&&"
	OR                          Type = "||"
	LEFT_SHIFT                  Type = "<<"
	RIGHT_SHIFT                 Type = ">>"
	UNSIGNED_RIGHT_SHIFT        Type = ">>>"
	PLUS_ASSIGN                 Type = "+="
	MINUS_ASSIGN                Type = "-="
	MULTIPLE_ASSIGN             Type = "*="
	DIVIDE_ASSIGN               Type = "/="
	MODULAR_ASSIGN              Type = "%="
	BIT_AND_ASSIGN              Type = "&="
	BIT_OR_ASSIGN               Type = "|="
	BIT_XOR_ASSIGN              Type = "^="
	LEFT_SHIFT_ASSIGN           Type = "<<="
	RIGHT_SHIFT_ASSIGN          Type = ">>="
	UNSIGNED_RIGHT_SHIFT_ASSIGN Type = ">>>="
	ARROW                       Type = "=>"

	PAREN_OPEN    Type = "("
	PAREN_CLOSE   Type = ")"
	BRACKET_OPEN  Type = "["
	BRACKET_CLOSE Type = "]"
	CURLY_OPEN    Type = "{"
	CURLY_CLOSE   Type = "}"
	COMMA         Type = ","
	PERIOD        Type = "."
	COLON         Type = ":"
	SEMICOLON     Type = ";"
	TEMPLATE      Type = "`"
)

var keywords = map[string]Type{
	"true":      BOOLEAN,
	"false":     BOOLEAN,
	"NaN":       NUMBER,
	"Infinity":  NUMBER,
	"null":      NULL,
	"undefined": UNDEFINED,

	"if":       KEYWORD,
	"else":     KEYWORD,
	"while":    KEYWORD,
	"for":      KEYWORD,
	"return":   KEYWORD,
	"break":    KEYWORD,
	"continue": KEYWORD,
	"switch":   KEYWORD,
	"case":     KEYWORD,
	"default":  KEYWORD,
	"try":      KEYWORD,
	"catch":    KEYWORD,
	"finally":  KEYWORD,
	"throw":    KEYWORD,
	"new":      KEYWORD,
	"delete":   KEYWORD,
	"import":   KEYWORD,
	"export":   KEYWORD,
	"extends":  KEYWORD,
	"super":    KEYWORD,
	"this":     KEYWORD,
	"static":   KEYWORD,
	"yield":    KEYWORD,
	"await":    KEYWORD,

	"var":   KEYWORD,
	"let":   KEYWORD,
	"const": KEYWORD,

	"function":  KEYWORD,
	"async":     KEYWORD,
	"generator": KEYWORD,
	"arguments": KEYWORD,
	"set":       KEYWORD,
	"get":       KEYWORD,

	"typeof":     KEYWORD,
	"instanceof": KEYWORD,
	"in":         KEYWORD,
	"void":       KEYWORD,

	"class":     KEYWORD,
	"interface": KEYWORD,
	"enum":      KEYWORD,
	"operator":  KEYWORD,
}

func NewToken(typ Type, literal string) Token {
	return Token{Type: typ, Literal: literal}
}

func TypeOf(literal string) Type {
	typ, ok := keywords[literal]
	if !ok {
		typ = IDENTIFIER
	}
	return typ
}
