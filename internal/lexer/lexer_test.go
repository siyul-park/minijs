package lexer

import (
	"testing"

	"github.com/siyul-park/minijs/internal/token"

	"github.com/stretchr/testify/assert"
)

func TestLexer_Next(t *testing.T) {
	tests := []struct {
		source string
		tokens []token.Token
	}{
		{source: `"string"`, tokens: []token.Token{token.New(token.STRING, `string`), token.EOF}},
		{source: `'string'`, tokens: []token.Token{token.New(token.STRING, `string`), token.EOF}},

		{source: `(`, tokens: []token.Token{token.LEFT_PAREN, token.EOF}},
		{source: `)`, tokens: []token.Token{token.RIGHT_PAREN, token.EOF}},
		{source: `[`, tokens: []token.Token{token.LEFT_BRACKET, token.EOF}},
		{source: `]`, tokens: []token.Token{token.RIGHT_BRACKET, token.EOF}},
		{source: `{`, tokens: []token.Token{token.LEFT_BRACE, token.EOF}},
		{source: `}`, tokens: []token.Token{token.RIGHT_BRACE, token.EOF}},
		{source: `,`, tokens: []token.Token{token.COMMA, token.EOF}},
		{source: `.`, tokens: []token.Token{token.PERIOD, token.EOF}},
		{source: `:`, tokens: []token.Token{token.COLON, token.EOF}},
		{source: `;`, tokens: []token.Token{token.SEMICOLON, token.EOF}},

		{source: `+`, tokens: []token.Token{token.PLUS, token.EOF}},
		{source: `+=`, tokens: []token.Token{token.PLUS_ASSIGN, token.EOF}},
		{source: `-`, tokens: []token.Token{token.MINUS, token.EOF}},
		{source: `-=`, tokens: []token.Token{token.MINUS_ASSIGN, token.EOF}},
		{source: `*`, tokens: []token.Token{token.MULTIPLE, token.EOF}},
		{source: `*=`, tokens: []token.Token{token.MULTIPLE_ASSIGN, token.EOF}},
		{source: `/`, tokens: []token.Token{token.DIVIDE, token.EOF}},
		{source: `/=`, tokens: []token.Token{token.DIVIDE_ASSIGN, token.EOF}},
		{source: `// line comment`, tokens: []token.Token{token.EOF}},
		{source: `/* block comment */`, tokens: []token.Token{token.EOF}},
		{source: `%`, tokens: []token.Token{token.MODULO, token.EOF}},
		{source: `%=`, tokens: []token.Token{token.MODULO_ASSIGN, token.EOF}},
		{source: `=`, tokens: []token.Token{token.ASSIGN, token.EOF}},
		{source: `==`, tokens: []token.Token{token.EQUAL, token.EOF}},
		{source: `=>`, tokens: []token.Token{token.ARROW, token.EOF}},
		{source: `!`, tokens: []token.Token{token.NOT, token.EOF}},
		{source: `!=`, tokens: []token.Token{token.NOT_EQUAL, token.EOF}},
		{source: `&`, tokens: []token.Token{token.BIT_AND, token.EOF}},
		{source: `&&`, tokens: []token.Token{token.AND, token.EOF}},
		{source: `&=`, tokens: []token.Token{token.BIT_AND_ASSIGN, token.EOF}},
		{source: `|`, tokens: []token.Token{token.BIT_OR, token.EOF}},
		{source: `||`, tokens: []token.Token{token.OR, token.EOF}},
		{source: `|=`, tokens: []token.Token{token.BIT_OR_ASSIGN, token.EOF}},
		{source: `^`, tokens: []token.Token{token.BIT_XOR, token.EOF}},
		{source: `^=`, tokens: []token.Token{token.BIT_XOR_ASSIGN, token.EOF}},
		{source: `~`, tokens: []token.Token{token.BIT_NOT, token.EOF}},
		{source: `<`, tokens: []token.Token{token.LESS_THAN, token.EOF}},
		{source: `<=`, tokens: []token.Token{token.LESS_THAN_EQUAL, token.EOF}},
		{source: `<<`, tokens: []token.Token{token.LEFT_SHIFT, token.EOF}},
		{source: `<<=`, tokens: []token.Token{token.LEFT_SHIFT_ASSIGN, token.EOF}},
		{source: `>`, tokens: []token.Token{token.GREATER_THAN, token.EOF}},
		{source: `>=`, tokens: []token.Token{token.GREATER_THAN_EQUAL, token.EOF}},
		{source: `>>`, tokens: []token.Token{token.RIGHT_SHIFT, token.EOF}},
		{source: `>>=`, tokens: []token.Token{token.RIGHT_SHIFT_ASSIGN, token.EOF}},
		{source: `>>>`, tokens: []token.Token{token.UNSIGNED_RIGHT_SHIFT, token.EOF}},
		{source: `>>>=`, tokens: []token.Token{token.UNSIGNED_RIGHT_SHIFT_ASSIGN, token.EOF}},

		{source: `123`, tokens: []token.Token{token.New(token.NUMBER, `123`), token.EOF}},
		{source: `12.3`, tokens: []token.Token{token.New(token.NUMBER, `12.3`), token.EOF}},
		{source: `0b01`, tokens: []token.Token{token.New(token.NUMBER, `0b01`), token.EOF}},
		{source: `0o01`, tokens: []token.Token{token.New(token.NUMBER, `0o01`), token.EOF}},
		{source: `0x01`, tokens: []token.Token{token.New(token.NUMBER, `0x01`), token.EOF}},
		{source: `10_000`, tokens: []token.Token{token.New(token.NUMBER, `10000`), token.EOF}},

		{source: `true`, tokens: []token.Token{token.TRUE, token.EOF}},
		{source: `false`, tokens: []token.Token{token.FALSE, token.EOF}},
		{source: `NaN`, tokens: []token.Token{token.NAN, token.EOF}},
		{source: `Infinity`, tokens: []token.Token{token.INFINITY, token.EOF}},
		{source: `null`, tokens: []token.Token{token.NULL, token.EOF}},
		{source: `undefined`, tokens: []token.Token{token.UNDEFINED, token.EOF}},
		{source: `if`, tokens: []token.Token{token.IF, token.EOF}},
		{source: `else`, tokens: []token.Token{token.ELSE, token.EOF}},
		{source: `while`, tokens: []token.Token{token.WHILE, token.EOF}},
		{source: `for`, tokens: []token.Token{token.FOR, token.EOF}},
		{source: `return`, tokens: []token.Token{token.RETURN, token.EOF}},
		{source: `break`, tokens: []token.Token{token.BREAK, token.EOF}},
		{source: `continue`, tokens: []token.Token{token.CONTINUE, token.EOF}},
		{source: `switch`, tokens: []token.Token{token.SWITCH, token.EOF}},
		{source: `case`, tokens: []token.Token{token.CASE, token.EOF}},
		{source: `default`, tokens: []token.Token{token.DEFAULT, token.EOF}},
		{source: `try`, tokens: []token.Token{token.TRY, token.EOF}},
		{source: `catch`, tokens: []token.Token{token.CATCH, token.EOF}},
		{source: `finally`, tokens: []token.Token{token.FINALLY, token.EOF}},
		{source: `throw`, tokens: []token.Token{token.THROW, token.EOF}},
		{source: `new`, tokens: []token.Token{token.NEW, token.EOF}},
		{source: `delete`, tokens: []token.Token{token.DELETE, token.EOF}},
		{source: `import`, tokens: []token.Token{token.IMPORT, token.EOF}},
		{source: `export`, tokens: []token.Token{token.EXPORT, token.EOF}},
		{source: `extends`, tokens: []token.Token{token.EXTENDS, token.EOF}},
		{source: `super`, tokens: []token.Token{token.SUPER, token.EOF}},
		{source: `this`, tokens: []token.Token{token.THIS, token.EOF}},
		{source: `static`, tokens: []token.Token{token.STATIC, token.EOF}},
		{source: `yield`, tokens: []token.Token{token.YIELD, token.EOF}},
		{source: `await`, tokens: []token.Token{token.AWAIT, token.EOF}},
		{source: `var`, tokens: []token.Token{token.VAR, token.EOF}},
		{source: `let`, tokens: []token.Token{token.LET, token.EOF}},
		{source: `const`, tokens: []token.Token{token.CONST, token.EOF}},
		{source: `function`, tokens: []token.Token{token.FUNCTION, token.EOF}},
		{source: `async`, tokens: []token.Token{token.ASYNC, token.EOF}},
		{source: `generator`, tokens: []token.Token{token.GENERATOR, token.EOF}},
		{source: `arguments`, tokens: []token.Token{token.ARGUMENTS, token.EOF}},
		{source: `set`, tokens: []token.Token{token.SET, token.EOF}},
		{source: `get`, tokens: []token.Token{token.GET, token.EOF}},
		{source: `typeof`, tokens: []token.Token{token.TYPEOF, token.EOF}},
		{source: `instanceof`, tokens: []token.Token{token.INSTANCEOF, token.EOF}},
		{source: `in`, tokens: []token.Token{token.IN, token.EOF}},
		{source: `void`, tokens: []token.Token{token.VOID, token.EOF}},
		{source: `class`, tokens: []token.Token{token.CLASS, token.EOF}},
		{source: `interface`, tokens: []token.Token{token.INTERFACE, token.EOF}},
		{source: `enum`, tokens: []token.Token{token.ENUM, token.EOF}},
		{source: `operator`, tokens: []token.Token{token.OPERATOR, token.EOF}},
	}

	for _, tt := range tests {
		t.Run(tt.source, func(t *testing.T) {
			l := New(tt.source)
			for _, expect := range tt.tokens {
				actual := l.Next()
				assert.Equal(t, expect.Kind(), actual.Kind())
				assert.Equal(t, expect.Literal, actual.Literal)
			}
		})
	}
}
