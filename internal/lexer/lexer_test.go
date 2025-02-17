package lexer

import (
	"strings"
	"testing"

	"github.com/siyul-park/minijs/internal/token"

	"github.com/stretchr/testify/assert"
)

func TestLexer_Next(t *testing.T) {
	tests := []struct {
		source string
		tokens []token.Token
	}{
		{source: `// comment`, tokens: []token.Token{token.New(token.EOF, "")}},
		{source: `/* comment */`, tokens: []token.Token{token.New(token.EOF, "")}},

		{source: `123`, tokens: []token.Token{token.New(token.NUMBER, "123")}},
		{source: `12.3`, tokens: []token.Token{token.New(token.NUMBER, "12.3")}},
		{source: `0x01`, tokens: []token.Token{token.New(token.NUMBER, "0x01")}},
		{source: `0o01`, tokens: []token.Token{token.New(token.NUMBER, "0o01")}},
		{source: `01`, tokens: []token.Token{token.New(token.NUMBER, "01")}},
		{source: `0b01`, tokens: []token.Token{token.New(token.NUMBER, "0b01")}},

		{source: `"foo"`, tokens: []token.Token{token.New(token.STRING, "foo")}},
		{source: `'foo''`, tokens: []token.Token{token.New(token.STRING, "foo")}},

		{source: `null`, tokens: []token.Token{token.New(token.NULL, "null")}},
		{source: `undefined`, tokens: []token.Token{token.New(token.UNDEFINED, "undefined")}},
		{source: `true`, tokens: []token.Token{token.New(token.TRUE, "true")}},
		{source: `false`, tokens: []token.Token{token.New(token.FALSE, "false")}},

		{source: `break`, tokens: []token.Token{token.New(token.BREAK, "break")}},
		{source: `do`, tokens: []token.Token{token.New(token.DO, "do")}},
		{source: `instanceof`, tokens: []token.Token{token.New(token.INSTANCEOF, "instanceof")}},
		{source: `typeof`, tokens: []token.Token{token.New(token.TYPEOF, "typeof")}},
		{source: `case`, tokens: []token.Token{token.New(token.CASE, "case")}},
		{source: `else`, tokens: []token.Token{token.New(token.ELSE, "else")}},
		{source: `new`, tokens: []token.Token{token.New(token.NEW, "new")}},
		{source: `var`, tokens: []token.Token{token.New(token.VAR, "var")}},
		{source: `catch`, tokens: []token.Token{token.New(token.CATCH, "catch")}},
		{source: `finally`, tokens: []token.Token{token.New(token.FINALLY, "finally")}},
		{source: `return`, tokens: []token.Token{token.New(token.RETURN, "return")}},
		{source: `void`, tokens: []token.Token{token.New(token.VOID, "void")}},
		{source: `continue`, tokens: []token.Token{token.New(token.CONTINUE, "continue")}},
		{source: `for`, tokens: []token.Token{token.New(token.FOR, "for")}},
		{source: `switch`, tokens: []token.Token{token.New(token.SWITCH, "switch")}},
		{source: `while`, tokens: []token.Token{token.New(token.WHILE, "while")}},
		{source: `debugger`, tokens: []token.Token{token.New(token.DEBUGGER, "debugger")}},
		{source: `function`, tokens: []token.Token{token.New(token.FUNCTION, "function")}},
		{source: `this`, tokens: []token.Token{token.New(token.THIS, "this")}},
		{source: `with`, tokens: []token.Token{token.New(token.WITH, "with")}},
		{source: `default`, tokens: []token.Token{token.New(token.DEFAULT, "default")}},
		{source: `if`, tokens: []token.Token{token.New(token.IF, "if")}},
		{source: `throw`, tokens: []token.Token{token.New(token.THROW, "throw")}},
		{source: `delete`, tokens: []token.Token{token.New(token.DELETE, "delete")}},
		{source: `in`, tokens: []token.Token{token.New(token.IN, "in")}},
		{source: `try`, tokens: []token.Token{token.New(token.TRY, "try")}},

		{source: `[`, tokens: []token.Token{token.New(token.OPEN_BRACKET, "[")}},
		{source: `]`, tokens: []token.Token{token.New(token.CLOSE_BRACKET, "]")}},
		{source: `(`, tokens: []token.Token{token.New(token.OPEN_PAREN, "(")}},
		{source: `)`, tokens: []token.Token{token.New(token.CLOSE_PAREN, ")")}},
		{source: `{`, tokens: []token.Token{token.New(token.OPEN_BRACE, "{")}},
		{source: `}`, tokens: []token.Token{token.New(token.CLOSE_BRACE, "}")}},
		{source: `;`, tokens: []token.Token{token.New(token.SEMICOLON, ";")}},
		{source: `,`, tokens: []token.Token{token.New(token.COMMA, ",")}},
		{source: `=`, tokens: []token.Token{token.New(token.ASSIGN, "=")}},
		{source: `?`, tokens: []token.Token{token.New(token.QUESTION, "?")}},
		{source: `:`, tokens: []token.Token{token.New(token.COLON, ":")}},
		{source: `.`, tokens: []token.Token{token.New(token.DOT, ".")}},
		{source: `+`, tokens: []token.Token{token.New(token.PLUS, "+")}},
		{source: `-`, tokens: []token.Token{token.New(token.MINUS, "-")}},
		{source: `++`, tokens: []token.Token{token.New(token.PLUS_PLUS, "++")}},
		{source: `--`, tokens: []token.Token{token.New(token.MINUS_MINUS, "--")}},
		{source: `~`, tokens: []token.Token{token.New(token.BIT_NOT, "~")}},
		{source: `!`, tokens: []token.Token{token.New(token.NOT, "!")}},
		{source: `*`, tokens: []token.Token{token.New(token.MULTIPLY, "*")}},
		{source: `/`, tokens: []token.Token{token.New(token.DIVIDE, "/")}},
		{source: `%`, tokens: []token.Token{token.New(token.MODULUS, "%")}},
		{source: `>>`, tokens: []token.Token{token.New(token.RIGHT_SHIFT_ARITHMETIC, ">>")}},
		{source: `<<`, tokens: []token.Token{token.New(token.LEFT_SHIFT_ARITHMETIC, "<<")}},
		{source: `>>>`, tokens: []token.Token{token.New(token.RIGHT_SHIFT_LOGICAL, ">>>")}},
		{source: `<`, tokens: []token.Token{token.New(token.LESS_THAN, "<")}},
		{source: `>`, tokens: []token.Token{token.New(token.GREATER_THAN, ">")}},
		{source: `<=`, tokens: []token.Token{token.New(token.LESS_THAN_OR_EQUAL, "<=")}},
		{source: `>=`, tokens: []token.Token{token.New(token.GREATER_THAN_OR_EQUAL, ">=")}},
		{source: `==`, tokens: []token.Token{token.New(token.EQUAL, "==")}},
		{source: `!=`, tokens: []token.Token{token.New(token.NOT_EQUAL, "!=")}},
		{source: `===`, tokens: []token.Token{token.New(token.IDENTITY_EQUAL, "===")}},
		{source: `!==`, tokens: []token.Token{token.New(token.IDENTITY_NOT_EQUAL, "!==")}},
		{source: `&`, tokens: []token.Token{token.New(token.BIT_AND, "&")}},
		{source: `|`, tokens: []token.Token{token.New(token.BIT_OR, "|")}},
		{source: `&&`, tokens: []token.Token{token.New(token.AND, "&&")}},
		{source: `||`, tokens: []token.Token{token.New(token.OR, "||")}},
		{source: `*=`, tokens: []token.Token{token.New(token.MULTIPLY_ASSIGN, "*=")}},
		{source: `/=`, tokens: []token.Token{token.New(token.DIVIDE_ASSIGN, "/=")}},
		{source: `%=`, tokens: []token.Token{token.New(token.MODULUS_ASSIGN, "%=")}},
		{source: `+=`, tokens: []token.Token{token.New(token.PLUS_ASSIGN, "+=")}},
		{source: `-=`, tokens: []token.Token{token.New(token.MINUS_ASSIGN, "-=")}},
		{source: `<<=`, tokens: []token.Token{token.New(token.LEFT_SHIFT_ARITHMETIC_ASSIGN, "<<=")}},
		{source: `>>=`, tokens: []token.Token{token.New(token.RIGHT_SHIFT_ARITHMETIC_ASSIGN, ">>=")}},
		{source: `>>>=`, tokens: []token.Token{token.New(token.RIGHT_SHIFT_LOGICAL_ASSIGN, ">>>=")}},
		{source: `&=`, tokens: []token.Token{token.New(token.BIT_AND_ASSIGN, "&=")}},
		{source: `|=`, tokens: []token.Token{token.New(token.BIT_OR_ASSIGN, "|=")}},
		{source: `^=`, tokens: []token.Token{token.New(token.BIT_XOR_ASSIGN, "^=")}},
	}

	for _, tt := range tests {
		t.Run(tt.source, func(t *testing.T) {
			l := New(strings.NewReader(tt.source))
			for _, expect := range tt.tokens {
				actual := l.Next()
				assert.Equal(t, expect, actual)
			}
		})
	}
}
