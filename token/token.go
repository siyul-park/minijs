package token

type Type string

type Token struct {
	Type    Type   // 토큰 타입
	Literal string // 토큰의 실제 값
}

const (
	EOF        Type = "EOF"
	ILLEGAL    Type = "ILLEGAL"
	IDENTIFIER Type = "IDENTIFIER"
	NUMBER     Type = "NUMBER"
	BOOLEAN    Type = "BOOLEAN"
	STRING     Type = "STRING"
	NULL       Type = "null"
	UNDEFINED  Type = "undefined"

	PLUS     Type = "＋"
	MINUS    Type = "−"
	ASTERISK Type = "*"
	SLASH    Type = "/"
	PERCENT  Type = "％"
	ASSIGN   Type = "＝"
	LESS     Type = "＜"
	GREATER  Type = "＞"
	BANG     Type = "!"
	AMP      Type = "&"
	VBAR     Type = "|"

	EQUAL           Type = "=="
	NOT_EQUAL       Type = "!="
	LESS_EQUAL      Type = "<="
	GREATER_EQUAL   Type = ">="
	PLUS_ASSIGN     Type = "+="
	MINUS_ASSIGN    Type = "-="
	ASTERISK_ASSIGN Type = "*="
	SLASH_ASSIGN    Type = "/="
	PERCENT_ASSIGN  Type = "%="
	AMP_ASSIGN      Type = "&="
	VBAR_ASSIGN     Type = "|="
	CARAT_ASSIGN    Type = "^="
	LSHIFT_ASSIGN   Type = "<<="
	RSHIFT_ASSIGN   Type = ">>="
	URSHIFT_ASSIGN  Type = ">>>="
	ARROW           Type = "=>"

	LPAREN    Type = "("
	RPAREN    Type = ")"
	LBRACKET  Type = "]"
	RBRACKET  Type = "]"
	LBRACE    Type = "{"
	RBRACE    Type = "}"
	COMMA     Type = ","
	SEMICOLON Type = ";"
	PERIOD    Type = "."

	IF       Type = "if"
	ELSE     Type = "else"
	WHILE    Type = "while"
	FOR      Type = "for"
	RETURN   Type = "return"
	BREAK    Type = "break"
	CONTINUE Type = "continue"
	SWITCH   Type = "switch"
	CASE     Type = "case"
	DEFAULT  Type = "default"
	TRY      Type = "try"
	CATCH    Type = "catch"
	FINALLY  Type = "finally"
	THROW    Type = "throw"
	NEW      Type = "new"
	DELETE   Type = "delete"
	IMPORT   Type = "import"
	EXPORT   Type = "export"
	EXTENDS  Type = "extends"
	SUPER    Type = "super"
	THIS     Type = "this"
	STATIC   Type = "static"
	YIELD    Type = "yield"
	AWAIT    Type = "await"

	VAR   Type = "var"
	LET   Type = "let"
	CONST Type = "const"

	FUNCTION  Type = "function"
	ASYNC     Type = "async"
	GENERATOR Type = "generator"
	ARGUMENTS Type = "arguments"
	SET       Type = "set"
	GET       Type = "get"

	TYPEOF     Type = "typeof"
	INSTANCEOF Type = "instanceof"
	IN         Type = "in"
	VOID       Type = "void"

	// 추가된 타입들
	CLASS     Type = "class"
	INTERFACE Type = "interface"
	ENUM      Type = "enum"
	OPERATOR  Type = "operator"
)

var types = map[string]Type{
	"true":      BOOLEAN,
	"false":     BOOLEAN,
	"NaN":       NUMBER,
	"Infinity":  NUMBER,
	"null":      NULL,
	"undefined": UNDEFINED,

	"if":       IF,
	"else":     ELSE,
	"while":    WHILE,
	"for":      FOR,
	"return":   RETURN,
	"break":    BREAK,
	"continue": CONTINUE,
	"switch":   SWITCH,
	"case":     CASE,
	"default":  DEFAULT,
	"try":      TRY,
	"catch":    CATCH,
	"finally":  FINALLY,
	"throw":    THROW,
	"new":      NEW,
	"delete":   DELETE,
	"import":   IMPORT,
	"export":   EXPORT,
	"extends":  EXTENDS,
	"super":    SUPER,
	"this":     THIS,
	"static":   STATIC,
	"yield":    YIELD,
	"await":    AWAIT,

	"var":   VAR,
	"let":   LET,
	"const": CONST,

	"function":  FUNCTION,
	"async":     ASYNC,
	"generator": GENERATOR,
	"arguments": ARGUMENTS,
	"set":       SET,
	"get":       GET,

	"typeof":     TYPEOF,
	"instanceof": INSTANCEOF,
	"in":         IN,
	"void":       VOID,

	"class":     CLASS,
	"interface": INTERFACE,
	"enum":      ENUM,
	"operator":  OPERATOR,
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
