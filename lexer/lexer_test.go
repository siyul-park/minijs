package lexer

import (
	"github.com/siyul-park/miniscript/token"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLexer_Next(t *testing.T) {
	tests := []struct {
		source string
		tokens []token.Token
	}{
		{
			source: ``,
			tokens: []token.Token{token.NewToken(token.EOF, "")},
		},
		{
			source: `1234567890`,
			tokens: []token.Token{token.NewToken(token.INT, `1234567890`), token.NewToken(token.EOF, "")},
		},
		{
			source: `1.234567890`,
			tokens: []token.Token{token.NewToken(token.FLOAT, `1.234567890`), token.NewToken(token.EOF, "")},
		},
		{
			source: `"abcdefg"`,
			tokens: []token.Token{token.NewToken(token.STRING, `abcdefg`), token.NewToken(token.EOF, "")},
		},
		{
			source: `true`,
			tokens: []token.Token{token.NewToken(token.TRUE, `true`), token.NewToken(token.EOF, "")},
		},
		{
			source: `false`,
			tokens: []token.Token{token.NewToken(token.FALSE, `false`), token.NewToken(token.EOF, "")},
		},
		{
			source: `foo`,
			tokens: []token.Token{token.NewToken(token.IDENT, `foo`), token.NewToken(token.EOF, "")},
		},
		{
			source: `+`,
			tokens: []token.Token{token.NewToken(token.PLUS, `+`), token.NewToken(token.EOF, "")},
		},
		{
			source: `-`,
			tokens: []token.Token{token.NewToken(token.MINUS, `-`), token.NewToken(token.EOF, "")},
		},
		{
			source: `*`,
			tokens: []token.Token{token.NewToken(token.MULTIPLY, `*`), token.NewToken(token.EOF, "")},
		},
		{
			source: `/`,
			tokens: []token.Token{token.NewToken(token.DIVIDE, `/`), token.NewToken(token.EOF, "")},
		},
		{
			source: `.`,
			tokens: []token.Token{token.NewToken(token.PERIOD, `.`), token.NewToken(token.EOF, "")},
		},
		{
			source: `(`,
			tokens: []token.Token{token.NewToken(token.LPAREN, `(`), token.NewToken(token.EOF, "")},
		},
		{
			source: `)`,
			tokens: []token.Token{token.NewToken(token.RPAREN, `)`), token.NewToken(token.EOF, "")},
		},
		{
			source: `+1234567890`,
			tokens: []token.Token{token.NewToken(token.PLUS, `+`), token.NewToken(token.INT, `1234567890`), token.NewToken(token.EOF, "")},
		},
		{
			source: `-1234567890`,
			tokens: []token.Token{token.NewToken(token.MINUS, `-`), token.NewToken(token.INT, `1234567890`), token.NewToken(token.EOF, "")},
		},
		{
			source: `12345+67890`,
			tokens: []token.Token{token.NewToken(token.INT, `12345`), token.NewToken(token.PLUS, `+`), token.NewToken(token.INT, `67890`), token.NewToken(token.EOF, "")},
		},
		{
			source: `12345-67890`,
			tokens: []token.Token{token.NewToken(token.INT, `12345`), token.NewToken(token.MINUS, `-`), token.NewToken(token.INT, `67890`), token.NewToken(token.EOF, "")},
		},
		{
			source: `foo.bar`,
			tokens: []token.Token{token.NewToken(token.IDENT, `foo`), token.NewToken(token.PERIOD, `.`), token.NewToken(token.IDENT, `bar`), token.NewToken(token.EOF, "")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.source, func(t *testing.T) {
			lexer := New(tt.source)
			for _, token := range tt.tokens {
				assert.Equal(t, token, lexer.Next())
			}
		})
	}
}
