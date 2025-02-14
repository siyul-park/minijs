package lexer

import (
	"github.com/siyul-park/minijs/token"
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
			tokens: []token.Token{token.NewToken(token.DECIMAL, `1234567890`), token.NewToken(token.EOF, "")},
		},
		{
			source: `1.234567890`,
			tokens: []token.Token{token.NewToken(token.DECIMAL, `1.234567890`), token.NewToken(token.EOF, "")},
		},
		{
			source: `0e-5`,
			tokens: []token.Token{token.NewToken(token.EXPONENTIAL, `0e-5`), token.NewToken(token.EOF, "")},
		},
		{
			source: `0e+5`,
			tokens: []token.Token{token.NewToken(token.EXPONENTIAL, `0e+5`), token.NewToken(token.EOF, "")},
		},
		{
			source: `5e1`,
			tokens: []token.Token{token.NewToken(token.EXPONENTIAL, `5e1`), token.NewToken(token.EOF, "")},
		},
		{
			source: `175e-2`,
			tokens: []token.Token{token.NewToken(token.EXPONENTIAL, `175e-2`), token.NewToken(token.EOF, "")},
		},
		{
			source: `1e3`,
			tokens: []token.Token{token.NewToken(token.EXPONENTIAL, `1e3`), token.NewToken(token.EOF, "")},
		},
		{
			source: `1e-3`,
			tokens: []token.Token{token.NewToken(token.EXPONENTIAL, `1e-3`), token.NewToken(token.EOF, "")},
		},
		{
			source: `1E3`,
			tokens: []token.Token{token.NewToken(token.EXPONENTIAL, `1e3`), token.NewToken(token.EOF, "")},
		},
		{
			source: `1_000_000`,
			tokens: []token.Token{token.NewToken(token.DECIMAL, `1000000`), token.NewToken(token.EOF, "")},
		},
		{
			source: `1_050.95`,
			tokens: []token.Token{token.NewToken(token.DECIMAL, `1050.95`), token.NewToken(token.EOF, "")},
		},
		{
			source: `0b10000000000000000000000000000000`,
			tokens: []token.Token{token.NewToken(token.BINARY, `0b10000000000000000000000000000000`), token.NewToken(token.EOF, "")},
		},
		{
			source: `0B00000000011111111111111111111111`,
			tokens: []token.Token{token.NewToken(token.BINARY, `0b00000000011111111111111111111111`), token.NewToken(token.EOF, "")},
		},
		{
			source: `0b1010_0001_1000_0101`,
			tokens: []token.Token{token.NewToken(token.BINARY, `0b1010000110000101`), token.NewToken(token.EOF, "")},
		},
		//{
		//	source: `0o2_2_5_6`,
		//	tokens: []token.Token{token.NewToken(token.DECIMAL, `0o2256`), token.NewToken(token.EOF, "")},
		//},
		//{
		//	source: `0xA0_B0_C0`,
		//	tokens: []token.Token{token.NewToken(token.DECIMAL, `0xA0B0C0`), token.NewToken(token.EOF, "")},
		//},
		//{
		//	source: `1_000_000_000_000_000_000_000n`,
		//	tokens: []token.Token{token.NewToken(token.DECIMAL, `1000000000000000000000000n`), token.NewToken(token.EOF, "")},
		//},
		{
			source: `NaN`,
			tokens: []token.Token{token.NewToken(token.NAN, `NaN`), token.NewToken(token.EOF, "")},
		},
		{
			source: `Infinity`,
			tokens: []token.Token{token.NewToken(token.INFINITY, `Infinity`), token.NewToken(token.EOF, "")},
		},
		{
			source: `-Infinity`,
			tokens: []token.Token{token.NewToken(token.MINUS, `-`), token.NewToken(token.INFINITY, `Infinity`), token.NewToken(token.EOF, "")},
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
			tokens: []token.Token{token.NewToken(token.IDENTIFIER, `foo`), token.NewToken(token.EOF, "")},
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
			source: `%`,
			tokens: []token.Token{token.NewToken(token.MODULO, `%`), token.NewToken(token.EOF, "")},
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
			tokens: []token.Token{token.NewToken(token.PLUS, `+`), token.NewToken(token.DECIMAL, `1234567890`), token.NewToken(token.EOF, "")},
		},
		{
			source: `-1234567890`,
			tokens: []token.Token{token.NewToken(token.MINUS, `-`), token.NewToken(token.DECIMAL, `1234567890`), token.NewToken(token.EOF, "")},
		},
		{
			source: `12345+67890`,
			tokens: []token.Token{token.NewToken(token.DECIMAL, `12345`), token.NewToken(token.PLUS, `+`), token.NewToken(token.DECIMAL, `67890`), token.NewToken(token.EOF, "")},
		},
		{
			source: `12345-67890`,
			tokens: []token.Token{token.NewToken(token.DECIMAL, `12345`), token.NewToken(token.MINUS, `-`), token.NewToken(token.DECIMAL, `67890`), token.NewToken(token.EOF, "")},
		},
		{
			source: `foo.bar`,
			tokens: []token.Token{token.NewToken(token.IDENTIFIER, `foo`), token.NewToken(token.PERIOD, `.`), token.NewToken(token.IDENTIFIER, `bar`), token.NewToken(token.EOF, "")},
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
