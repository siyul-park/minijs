package token

type Type string

type Token struct {
	Type    Type
	Literal string
}

const (
	ILLEGAL    Type = "ILLEGAL"
	IDENTIFIER Type = "IDENTIFIER"
	NUMBER     Type = "NUMBER"
	BOOLEAN    Type = "BOOLEAN"
	STRING     Type = "STRING"
	KEYWORD    Type = "KEYWORD"
)

var (
	EOF = New("EOF", "")

	LEFT_PAREN    = New(KEYWORD, "(")
	RIGHT_PAREN   = New(KEYWORD, ")")
	LEFT_BRACKET  = New(KEYWORD, "[")
	RIGHT_BRACKET = New(KEYWORD, "]")
	LEFT_BRACE    = New(KEYWORD, "{")
	RIGHT_BRACE   = New(KEYWORD, "}")
	COMMA         = New(KEYWORD, ",")
	PERIOD        = New(KEYWORD, ".")
	COLON         = New(KEYWORD, ":")
	SEMICOLON     = New(KEYWORD, ";")

	PLUS                        = New(KEYWORD, "+")
	MINUS                       = New(KEYWORD, "−")
	MULTIPLE                    = New(KEYWORD, "*")
	DIVIDE                      = New(KEYWORD, "/")
	MODULO                      = New(KEYWORD, "％")
	ASSIGN                      = New(KEYWORD, "＝")
	LESS_THAN                   = New(KEYWORD, "＜")
	GREATER_THAN                = New(KEYWORD, "＞")
	NOT                         = New(KEYWORD, "!")
	BIT_AND                     = New(KEYWORD, "&")
	BIT_OR                      = New(KEYWORD, "|")
	BIT_XOR                     = New(KEYWORD, "^")
	BIT_NOT                     = New(KEYWORD, "~")
	EQUAL                       = New(KEYWORD, "==")
	NOT_EQUAL                   = New(KEYWORD, "!=")
	LESS_THAN_EQUAL             = New(KEYWORD, "<=")
	GREATER_THAN_EQUAL          = New(KEYWORD, ">=")
	AND                         = New(KEYWORD, "&&")
	OR                          = New(KEYWORD, "||")
	LEFT_SHIFT                  = New(KEYWORD, "<<")
	RIGHT_SHIFT                 = New(KEYWORD, ">>")
	UNSIGNED_RIGHT_SHIFT        = New(KEYWORD, ">>>")
	PLUS_ASSIGN                 = New(KEYWORD, "+=")
	MINUS_ASSIGN                = New(KEYWORD, "-=")
	MULTIPLE_ASSIGN             = New(KEYWORD, "*=")
	DIVIDE_ASSIGN               = New(KEYWORD, "/=")
	MODULO_ASSIGN               = New(KEYWORD, "%=")
	BIT_AND_ASSIGN              = New(KEYWORD, "&=")
	BIT_OR_ASSIGN               = New(KEYWORD, "|=")
	BIT_XOR_ASSIGN              = New(KEYWORD, "^=")
	LEFT_SHIFT_ASSIGN           = New(KEYWORD, "<<=")
	RIGHT_SHIFT_ASSIGN          = New(KEYWORD, ">>=")
	UNSIGNED_RIGHT_SHIFT_ASSIGN = New(KEYWORD, ">>>=")

	TRUE      = New(BOOLEAN, "true")
	FALSE     = New(BOOLEAN, "false")
	NAN       = New(NUMBER, "NaN")
	INFINITY  = New(NUMBER, "Infinity")
	NULL      = New(KEYWORD, "null")
	UNDEFINED = New(KEYWORD, "undefined")

	IF       = New(KEYWORD, "if")
	ELSE     = New(KEYWORD, "else")
	WHILE    = New(KEYWORD, "while")
	FOR      = New(KEYWORD, "for")
	RETURN   = New(KEYWORD, "return")
	BREAK    = New(KEYWORD, "break")
	CONTINUE = New(KEYWORD, "continue")
	SWITCH   = New(KEYWORD, "switch")
	CASE     = New(KEYWORD, "case")
	DEFAULT  = New(KEYWORD, "default")
	TRY      = New(KEYWORD, "try")
	CATCH    = New(KEYWORD, "catch")
	FINALLY  = New(KEYWORD, "finally")
	THROW    = New(KEYWORD, "throw")
	NEW      = New(KEYWORD, "new")
	DELETE   = New(KEYWORD, "delete")
	IMPORT   = New(KEYWORD, "import")
	EXPORT   = New(KEYWORD, "export")
	EXTENDS  = New(KEYWORD, "extends")
	SUPER    = New(KEYWORD, "super")
	THIS     = New(KEYWORD, "this")
	STATIC   = New(KEYWORD, "static")
	YIELD    = New(KEYWORD, "yield")
	AWAIT    = New(KEYWORD, "await")

	VAR   = New(KEYWORD, "var")
	LET   = New(KEYWORD, "let")
	CONST = New(KEYWORD, "const")

	FUNCTION  = New(KEYWORD, "function")
	ASYNC     = New(KEYWORD, "async")
	GENERATOR = New(KEYWORD, "generator")
	ARGUMENTS = New(KEYWORD, "arguments")
	SET       = New(KEYWORD, "set")
	GET       = New(KEYWORD, "get")
	ARROW     = New(KEYWORD, "=>")

	TYPEOF     = New(KEYWORD, "typeof")
	INSTANCEOF = New(KEYWORD, "instanceof")
	IN         = New(KEYWORD, "in")
	VOID       = New(KEYWORD, "void")

	CLASS     = New(KEYWORD, "class")
	INTERFACE = New(KEYWORD, "interface")
	ENUM      = New(KEYWORD, "enum")
	OPERATOR  = New(KEYWORD, "operator")
)

var keywords = []Token{
	TRUE, FALSE, NAN, INFINITY, NULL, UNDEFINED,
	IF, ELSE, WHILE, FOR, RETURN, BREAK, CONTINUE,
	SWITCH, CASE, DEFAULT, TRY, CATCH, FINALLY,
	THROW, NEW, DELETE, IMPORT, EXPORT, EXTENDS, SUPER,
	THIS, STATIC, YIELD, AWAIT, VAR, LET, CONST,
	FUNCTION, ASYNC, GENERATOR, ARGUMENTS, SET, GET,
	ARROW, TYPEOF, INSTANCEOF, IN, VOID, CLASS,
	INTERFACE, ENUM, OPERATOR,
	LEFT_PAREN, RIGHT_PAREN, LEFT_BRACKET, RIGHT_BRACKET, LEFT_BRACE, RIGHT_BRACE,
	COMMA, PERIOD, COLON, SEMICOLON,
	PLUS, MINUS, MULTIPLE, DIVIDE, MODULO, ASSIGN,
	LESS_THAN, GREATER_THAN, NOT, BIT_AND, BIT_OR, BIT_XOR,
	BIT_NOT, EQUAL, NOT_EQUAL, LESS_THAN_EQUAL, GREATER_THAN_EQUAL,
	AND, OR, LEFT_SHIFT, RIGHT_SHIFT, UNSIGNED_RIGHT_SHIFT,
	PLUS_ASSIGN, MINUS_ASSIGN, MULTIPLE_ASSIGN, DIVIDE_ASSIGN,
	MODULO_ASSIGN, BIT_AND_ASSIGN, BIT_OR_ASSIGN, BIT_XOR_ASSIGN,
	LEFT_SHIFT_ASSIGN, RIGHT_SHIFT_ASSIGN, UNSIGNED_RIGHT_SHIFT_ASSIGN,
}

var types = map[string]Type{}

func init() {
	for _, keyword := range keywords {
		types[keyword.Literal] = keyword.Kind()
	}
}

func TypeOf(literal string) Type {
	typ, ok := types[literal]
	if !ok {
		typ = IDENTIFIER
	}
	return typ
}

func New(typ Type, literal string) Token {
	return Token{Type: typ, Literal: literal}
}

func (t Token) Kind() Type {
	if t.Type == KEYWORD {
		return Type(t.Literal)
	}
	return t.Type
}

func (t Token) String() string {
	return t.Literal
}
