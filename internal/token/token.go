package token

type Type string

type Token struct {
	Type    Type
	Literal string
}

const (
	ILLEGAL Type = "ILLEGAL"
	EOF     Type = "EOF"

	NUMBER     Type = "NUMBER"
	STRING     Type = "STRING"
	IDENTIFIER Type = "IDENTIFIER"

	NULL      Type = "null"
	UNDEFINED Type = "undefined"
	TRUE      Type = "true"
	FALSE     Type = "false"

	BREAK      Type = "break"
	DO         Type = "do"
	INSTANCEOF Type = "instanceof"
	TYPEOF     Type = "typeof"
	CASE       Type = "case"
	ELSE       Type = "else"
	NEW        Type = "new"
	VAR        Type = "var"
	CATCH      Type = "catch"
	FINALLY    Type = "finally"
	RETURN     Type = "return"
	VOID       Type = "void"
	CONTINUE   Type = "continue"
	FOR        Type = "for"
	SWITCH     Type = "switch"
	WHILE      Type = "while"
	DEBUGGER   Type = "debugger"
	FUNCTION   Type = "function"
	THIS       Type = "this"
	WITH       Type = "with"
	DEFAULT    Type = "default"
	IF         Type = "if"
	THROW      Type = "throw"
	DELETE     Type = "delete"
	IN         Type = "in"
	TRY        Type = "try"

	OPEN_BRACKET                  Type = "["
	CLOSE_BRACKET                 Type = "]"
	OPEN_PAREN                    Type = "("
	CLOSE_PAREN                   Type = ")"
	OPEN_BRACE                    Type = "{"
	CLOSE_BRACE                   Type = "}"
	SEMICOLON                     Type = ";"
	COMMA                         Type = ","
	ASSIGN                        Type = "="
	QUESTION                      Type = "?"
	COLON                         Type = ":"
	DOT                           Type = "."
	PLUS                          Type = "+"
	MINUS                         Type = "-"
	PLUS_PLUS                     Type = "++"
	MINUS_MINUS                   Type = "--"
	BIT_NOT                       Type = "~"
	NOT                           Type = "!"
	MULTIPLY                      Type = "*"
	DIVIDE                        Type = "/"
	MODULUS                       Type = "%"
	RIGHT_SHIFT_ARITHMETIC        Type = ">>"
	LEFT_SHIFT_ARITHMETIC         Type = "<<"
	RIGHT_SHIFT_LOGICAL           Type = ">>>"
	LESS_THAN                     Type = "<"
	GREATER_THAN                  Type = ">"
	LESS_THAN_OR_EQUAL            Type = "<="
	GREATER_THAN_OR_EQUAL         Type = ">="
	EQUAL                         Type = "=="
	NOT_EQUAL                     Type = "!="
	IDENTITY_EQUAL                Type = "==="
	IDENTITY_NOT_EQUAL            Type = "!=="
	BIT_AND                       Type = "&"
	BIT_OR                        Type = "|"
	AND                           Type = "&&"
	OR                            Type = "||"
	MULTIPLY_ASSIGN               Type = "*="
	DIVIDE_ASSIGN                 Type = "%="
	MODULUS_ASSIGN                Type = "%="
	PLUS_ASSIGN                   Type = "+="
	MINUS_ASSIGN                  Type = "-="
	LEFT_SHIFT_ARITHMETIC_ASSIGN  Type = "<<="
	RIGHT_SHIFT_ARITHMETIC_ASSIGN Type = ">>="
	RIGHT_SHIFT_LOGICAL_ASSIGN    Type = ">>>="
	BIT_AND_ASSIGN                Type = "&="
	BIT_OR_ASSIGN                 Type = "|="
	BIT_XOR_ASSIGN                Type = "^="
)

var reserved = []Type{
	NULL, UNDEFINED, TRUE, FALSE,
	BREAK, DO, INSTANCEOF, TYPEOF, CASE, ELSE, NEW, VAR, CATCH,
	FINALLY, RETURN, VOID, CONTINUE, FOR, SWITCH, WHILE, DEBUGGER,
	FUNCTION, THIS, WITH, DEFAULT, IF, THROW, DELETE, IN, TRY,
	OPEN_BRACKET, CLOSE_BRACKET, OPEN_PAREN, CLOSE_PAREN,
	OPEN_BRACE, CLOSE_BRACE, SEMICOLON, COMMA, ASSIGN, QUESTION,
	COLON, DOT, PLUS, MINUS, PLUS_PLUS, MINUS_MINUS, BIT_NOT, NOT,
	MULTIPLY, DIVIDE, MODULUS, RIGHT_SHIFT_ARITHMETIC,
	LEFT_SHIFT_ARITHMETIC, RIGHT_SHIFT_LOGICAL, LESS_THAN,
	GREATER_THAN, LESS_THAN_OR_EQUAL, GREATER_THAN_OR_EQUAL,
	EQUAL, NOT_EQUAL, IDENTITY_EQUAL, IDENTITY_NOT_EQUAL,
	BIT_AND, BIT_OR, AND, OR, MULTIPLY_ASSIGN, DIVIDE_ASSIGN,
	MODULUS_ASSIGN, PLUS_ASSIGN, MINUS_ASSIGN,
	LEFT_SHIFT_ARITHMETIC_ASSIGN, RIGHT_SHIFT_ARITHMETIC_ASSIGN,
	RIGHT_SHIFT_LOGICAL_ASSIGN, BIT_AND_ASSIGN, BIT_OR_ASSIGN,
	BIT_XOR_ASSIGN,
}

var types = map[string]Type{}

func init() {
	for _, t := range reserved {
		types[string(t)] = t
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

func (t Token) String() string {
	return t.Literal
}
