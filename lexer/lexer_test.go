package lexer

import (
	"testing"

	"github.com/siyul-park/minijs/token"
	"github.com/stretchr/testify/assert"
)

func TestLexer_Next(t *testing.T) {
	tests := []struct {
		source string
		tokens []token.Token
	}{
		{
			source: `"string"`,
			tokens: []token.Token{
				token.NewToken(token.STRING, `string`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: "`",
			tokens: []token.Token{
				token.NewToken(token.TEMPLATE, "`"),
				token.NewToken(token.EOF, ""),
			},
		},

		{
			source: `(`,
			tokens: []token.Token{
				token.NewToken(token.PAREN_OPEN, `(`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `)`,
			tokens: []token.Token{
				token.NewToken(token.PAREN_CLOSE, `)`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `[`,
			tokens: []token.Token{
				token.NewToken(token.BRACKET_OPEN, `[`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `]`,
			tokens: []token.Token{
				token.NewToken(token.BRACKET_CLOSE, `]`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `{`,
			tokens: []token.Token{
				token.NewToken(token.CURLY_OPEN, `{`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `}`,
			tokens: []token.Token{
				token.NewToken(token.CURLY_CLOSE, `}`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `,`,
			tokens: []token.Token{
				token.NewToken(token.COMMA, `,`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `.`,
			tokens: []token.Token{
				token.NewToken(token.PERIOD, `.`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `:`,
			tokens: []token.Token{
				token.NewToken(token.COLON, `:`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `;`,
			tokens: []token.Token{
				token.NewToken(token.SEMICOLON, `;`),
				token.NewToken(token.EOF, ""),
			},
		},

		{
			source: `+`,
			tokens: []token.Token{
				token.NewToken(token.PLUS, `+`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `+=`,
			tokens: []token.Token{
				token.NewToken(token.PLUS_ASSIGN, `+=`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `-`,
			tokens: []token.Token{
				token.NewToken(token.MINUS, `-`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `-=`,
			tokens: []token.Token{
				token.NewToken(token.MINUS_ASSIGN, `-=`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `*`,
			tokens: []token.Token{
				token.NewToken(token.MULTIPLE, `*`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `*=`,
			tokens: []token.Token{
				token.NewToken(token.MULTIPLE_ASSIGN, `*=`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `/`,
			tokens: []token.Token{
				token.NewToken(token.DIVIDE, `/`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `/=`,
			tokens: []token.Token{
				token.NewToken(token.DIVIDE_ASSIGN, `/=`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `// line comment`,
			tokens: []token.Token{
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `/* block comment */`,
			tokens: []token.Token{
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `%`,
			tokens: []token.Token{
				token.NewToken(token.MODULAR, `%`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `%=`,
			tokens: []token.Token{
				token.NewToken(token.MODULAR_ASSIGN, `%=`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `=`,
			tokens: []token.Token{
				token.NewToken(token.ASSIGN, `=`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `==`,
			tokens: []token.Token{
				token.NewToken(token.EQUAL, `==`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `=>`,
			tokens: []token.Token{
				token.NewToken(token.ARROW, `=>`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `!`,
			tokens: []token.Token{
				token.NewToken(token.NOT, `!`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `!=`,
			tokens: []token.Token{
				token.NewToken(token.NOT_EQUAL, `!=`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `&`,
			tokens: []token.Token{
				token.NewToken(token.BIT_AND, `&`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `&`,
			tokens: []token.Token{
				token.NewToken(token.BIT_AND, `&`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `&&`,
			tokens: []token.Token{
				token.NewToken(token.AND, `&&`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `&=`,
			tokens: []token.Token{
				token.NewToken(token.BIT_AND_ASSIGN, `&=`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `|`,
			tokens: []token.Token{
				token.NewToken(token.BIT_OR, `|`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `||`,
			tokens: []token.Token{
				token.NewToken(token.OR, `||`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `|=`,
			tokens: []token.Token{
				token.NewToken(token.BIT_OR_ASSIGN, `|=`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `^`,
			tokens: []token.Token{
				token.NewToken(token.BIT_XOR, `^`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `^=`,
			tokens: []token.Token{
				token.NewToken(token.BIT_XOR_ASSIGN, `^=`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `~`,
			tokens: []token.Token{
				token.NewToken(token.BIT_NOT, `~`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `<`,
			tokens: []token.Token{
				token.NewToken(token.LESS_THAN, `<`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `<=`,
			tokens: []token.Token{
				token.NewToken(token.LESS_THAN_EQUAL, `<=`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `<<`,
			tokens: []token.Token{
				token.NewToken(token.LEFT_SHIFT, `<<`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `<<=`,
			tokens: []token.Token{
				token.NewToken(token.LEFT_SHIFT_ASSIGN, `<<=`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `>`,
			tokens: []token.Token{
				token.NewToken(token.GREATER_THAN, `>`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `>=`,
			tokens: []token.Token{
				token.NewToken(token.GREATER_THAN_EQUAL, `>=`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `>>`,
			tokens: []token.Token{
				token.NewToken(token.RIGHT_SHIFT, `>>`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `>>=`,
			tokens: []token.Token{
				token.NewToken(token.RIGHT_SHIFT_ASSIGN, `>>=`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `>>>`,
			tokens: []token.Token{
				token.NewToken(token.UNSIGNED_RIGHT_SHIFT, `>>>`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `>>>=`,
			tokens: []token.Token{
				token.NewToken(token.UNSIGNED_RIGHT_SHIFT_ASSIGN, `>>>=`),
				token.NewToken(token.EOF, ""),
			},
		},

		{
			source: `123`,
			tokens: []token.Token{
				token.NewToken(token.NUMBER, `123`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `12.3`,
			tokens: []token.Token{
				token.NewToken(token.NUMBER, `12.3`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `0b01`,
			tokens: []token.Token{
				token.NewToken(token.NUMBER, `0b01`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `0o01`,
			tokens: []token.Token{
				token.NewToken(token.NUMBER, `0o01`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `0x01`,
			tokens: []token.Token{
				token.NewToken(token.NUMBER, `0x01`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `10_000`,
			tokens: []token.Token{
				token.NewToken(token.NUMBER, `10000`),
				token.NewToken(token.EOF, ""),
			},
		},

		{
			source: `true`,
			tokens: []token.Token{
				token.NewToken(token.BOOLEAN, `true`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `false`,
			tokens: []token.Token{
				token.NewToken(token.BOOLEAN, `false`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `NaN`,
			tokens: []token.Token{
				token.NewToken(token.NUMBER, `NaN`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `null`,
			tokens: []token.Token{
				token.NewToken(token.NULL, `null`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `undefined`,
			tokens: []token.Token{
				token.NewToken(token.UNDEFINED, `undefined`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `if`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `if`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `else`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `else`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `while`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `while`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `for`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `for`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `return`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `return`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `break`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `break`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `continue`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `continue`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `switch`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `switch`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `case`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `case`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `default`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `default`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `try`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `try`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `catch`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `catch`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `finally`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `finally`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `throw`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `throw`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `new`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `new`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `delete`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `delete`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `import`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `import`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `export`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `export`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `extends`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `extends`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `super`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `super`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `this`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `this`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `static`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `static`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `yield`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `yield`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `await`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `await`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `var`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `var`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `let`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `let`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `const`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `const`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `function`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `function`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `async`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `async`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `generator`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `generator`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `arguments`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `arguments`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `set`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `set`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `get`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `get`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `typeof`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `typeof`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `instanceof`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `instanceof`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `in`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `in`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `void`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `void`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `class`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `class`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `interface`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `interface`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `enum`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `enum`),
				token.NewToken(token.EOF, ""),
			},
		},
		{
			source: `operator`,
			tokens: []token.Token{
				token.NewToken(token.KEYWORD, `operator`),
				token.NewToken(token.EOF, ""),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.source, func(t *testing.T) {
			l := New(tt.source)
			for _, expect := range tt.tokens {
				actual := l.Next()
				assert.Equal(t, expect.Type, actual.Type)
				assert.Equal(t, expect.Literal, actual.Literal)
			}
		})
	}
}
