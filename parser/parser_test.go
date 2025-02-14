package parser

import (
	"github.com/siyul-park/minijs/ast"
	"github.com/siyul-park/minijs/lexer"
	"github.com/siyul-park/minijs/token"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		source  string
		program *ast.Program
	}{
		{
			source:  ``,
			program: &ast.Program{},
		},
		{
			source: `1234567890`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.NumberLiteral{Token: token.NewToken(token.DECIMAL, `1234567890`), Value: 1234567890},
			}},
		},
		{
			source: `1.234567890`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.NumberLiteral{Token: token.NewToken(token.DECIMAL, `1.234567890`), Value: 1.234567890},
			}},
		},
		{
			source: `0e-5`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.NumberLiteral{Token: token.NewToken(token.EXPONENTIAL, `0e-5`), Value: 0},
			}},
		},
		{
			source: `0e+5`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.NumberLiteral{Token: token.NewToken(token.EXPONENTIAL, `0e+5`), Value: 0},
			}},
		},
		{
			source: `5e1`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.NumberLiteral{Token: token.NewToken(token.EXPONENTIAL, `5e1`), Value: 50},
			}},
		},
		{
			source: `175e-2`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.NumberLiteral{Token: token.NewToken(token.EXPONENTIAL, `175e-2`), Value: 1.75},
			}},
		},
		{
			source: `1e3`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.NumberLiteral{Token: token.NewToken(token.EXPONENTIAL, `1e3`), Value: 1000},
			}},
		},
		{
			source: `1e-3`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.NumberLiteral{Token: token.NewToken(token.EXPONENTIAL, `1e-3`), Value: 0.001},
			}},
		},
		{
			source: `1E3`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.NumberLiteral{Token: token.NewToken(token.EXPONENTIAL, `1e3`), Value: 1000},
			}},
		},
		{
			source: `1_000_000`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.NumberLiteral{Token: token.NewToken(token.DECIMAL, `1000000`), Value: 1000000},
			}},
		},
		{
			source: `1_050.95`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.NumberLiteral{Token: token.NewToken(token.DECIMAL, `1050.95`), Value: 1050.95},
			}},
		},
		{
			source: `0b10000000000000000000000000000000`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.NumberLiteral{Token: token.NewToken(token.BINARY, `0b10000000000000000000000000000000`), Value: 2147483648},
			}},
		},
		{
			source: `0B00000000011111111111111111111111`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.NumberLiteral{Token: token.NewToken(token.BINARY, `0b00000000011111111111111111111111`), Value: 8388607},
			}},
		},
		{
			source: `0b1010_0001_1000_0101`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.NumberLiteral{Token: token.NewToken(token.BINARY, `0b1010000110000101`), Value: 41349},
			}},
		},
		{
			source: `NaN`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.NumberLiteral{Token: token.NewToken(token.NAN, `NaN`)},
			}},
		},
		{
			source: `Infinity`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.NumberLiteral{Token: token.NewToken(token.INFINITY, `Infinity`)},
			}},
		},
		{
			source: `"abcdefg"`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.StringLiteral{Token: token.NewToken(token.STRING, `abcdefg`), Value: `abcdefg`},
			}},
		},
		{
			source: `true`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.BoolLiteral{Token: token.NewToken(token.TRUE, `true`), Value: true},
			}},
		},
		{
			source: `false`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.BoolLiteral{Token: token.NewToken(token.FALSE, `false`), Value: false},
			}},
		},
		{
			source: `foo`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.IdentifierLiteral{Token: token.NewToken(token.IDENTIFIER, `foo`), Value: `foo`},
			}},
		},
		{
			source: `+1234567890`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.PrefixExpression{
					Token: token.NewToken(token.PLUS, "+"),
					Right: &ast.NumberLiteral{Token: token.NewToken(token.DECIMAL, `1234567890`), Value: 1234567890},
				},
			}},
		},
		{
			source: `-1234567890`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.PrefixExpression{
					Token: token.NewToken(token.MINUS, "-"),
					Right: &ast.NumberLiteral{Token: token.NewToken(token.DECIMAL, `1234567890`), Value: 1234567890},
				},
			}},
		},
		{
			source: `12345+67890`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.InfixExpression{
					Token: token.NewToken(token.PLUS, "+"),
					Left:  &ast.NumberLiteral{Token: token.NewToken(token.DECIMAL, `12345`), Value: 12345},
					Right: &ast.NumberLiteral{Token: token.NewToken(token.DECIMAL, `67890`), Value: 67890},
				},
			}},
		},
		{
			source: `12345-67890`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.InfixExpression{
					Token: token.NewToken(token.MINUS, "-"),
					Left:  &ast.NumberLiteral{Token: token.NewToken(token.DECIMAL, `12345`), Value: 12345},
					Right: &ast.NumberLiteral{Token: token.NewToken(token.DECIMAL, `67890`), Value: 67890},
				},
			}},
		},
		{
			source: `12345*67890`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.InfixExpression{
					Token: token.NewToken(token.MULTIPLY, "*"),
					Left:  &ast.NumberLiteral{Token: token.NewToken(token.DECIMAL, `12345`), Value: 12345},
					Right: &ast.NumberLiteral{Token: token.NewToken(token.DECIMAL, `67890`), Value: 67890},
				},
			}},
		},
		{
			source: `12345/67890`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.InfixExpression{
					Token: token.NewToken(token.DIVIDE, "/"),
					Left:  &ast.NumberLiteral{Token: token.NewToken(token.DECIMAL, `12345`), Value: 12345},
					Right: &ast.NumberLiteral{Token: token.NewToken(token.DECIMAL, `67890`), Value: 67890},
				},
			}},
		},
		{
			source: `12345%67890`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.InfixExpression{
					Token: token.NewToken(token.MODULO, "%"),
					Left:  &ast.NumberLiteral{Token: token.NewToken(token.DECIMAL, `12345`), Value: 12345},
					Right: &ast.NumberLiteral{Token: token.NewToken(token.DECIMAL, `67890`), Value: 67890},
				},
			}},
		},
		{
			source: `+12345+-67890`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.InfixExpression{
					Token: token.NewToken(token.PLUS, "+"),
					Left: &ast.PrefixExpression{
						Token: token.NewToken(token.PLUS, "+"),
						Right: &ast.NumberLiteral{Token: token.NewToken(token.DECIMAL, `12345`), Value: 12345},
					},
					Right: &ast.PrefixExpression{
						Token: token.NewToken(token.MINUS, "-"),
						Right: &ast.NumberLiteral{Token: token.NewToken(token.DECIMAL, `67890`), Value: 67890},
					},
				},
			}},
		},
		{
			source: `12*34+56/78`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.InfixExpression{
					Token: token.NewToken(token.PLUS, "+"),
					Left: &ast.InfixExpression{
						Token: token.NewToken(token.MULTIPLY, "*"),
						Left:  &ast.NumberLiteral{Token: token.NewToken(token.DECIMAL, `12`), Value: 12},
						Right: &ast.NumberLiteral{Token: token.NewToken(token.DECIMAL, `34`), Value: 34},
					},
					Right: &ast.InfixExpression{
						Token: token.NewToken(token.DIVIDE, "/"),
						Left:  &ast.NumberLiteral{Token: token.NewToken(token.DECIMAL, `56`), Value: 56},
						Right: &ast.NumberLiteral{Token: token.NewToken(token.DECIMAL, `78`), Value: 78},
					},
				},
			}},
		},
		{
			source: `(12+34)*(56-78)`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.InfixExpression{
					Token: token.NewToken(token.MULTIPLY, "*"),
					Left: &ast.InfixExpression{
						Token: token.NewToken(token.PLUS, "+"),
						Left:  &ast.NumberLiteral{Token: token.NewToken(token.DECIMAL, `12`), Value: 12},
						Right: &ast.NumberLiteral{Token: token.NewToken(token.DECIMAL, `34`), Value: 34},
					},
					Right: &ast.InfixExpression{
						Token: token.NewToken(token.MINUS, "-"),
						Left:  &ast.NumberLiteral{Token: token.NewToken(token.DECIMAL, `56`), Value: 56},
						Right: &ast.NumberLiteral{Token: token.NewToken(token.DECIMAL, `78`), Value: 78},
					},
				},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.source, func(t *testing.T) {
			l := lexer.New(tt.source)
			p := New(l)

			prg, err := p.Parse()
			assert.NoError(t, err)
			assert.Equal(t, tt.program, prg)
		})
	}
}
