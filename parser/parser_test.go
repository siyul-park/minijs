package parser

import (
	"github.com/siyul-park/miniscript/ast"
	"github.com/siyul-park/miniscript/lexer"
	"github.com/siyul-park/miniscript/token"
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
				&ast.IntLiteral{Token: token.NewToken(token.INT, `1234567890`), Value: 1234567890},
			}},
		},
		{
			source: `1.234567890`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.FloatLiteral{Token: token.NewToken(token.FLOAT, `1.234567890`), Value: 1.234567890},
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
				&ast.IdentifierLiteral{Token: token.NewToken(token.IDENT, `foo`), Value: `foo`},
			}},
		},
		{
			source: `+1234567890`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.PrefixExpression{
					Token: token.NewToken(token.PLUS, "+"),
					Right: &ast.IntLiteral{Token: token.NewToken(token.INT, `1234567890`), Value: 1234567890},
				},
			}},
		},
		{
			source: `-1234567890`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.PrefixExpression{
					Token: token.NewToken(token.MINUS, "-"),
					Right: &ast.IntLiteral{Token: token.NewToken(token.INT, `1234567890`), Value: 1234567890},
				},
			}},
		},
		{
			source: `12345+67890`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.InfixExpression{
					Token: token.NewToken(token.PLUS, "+"),
					Left:  &ast.IntLiteral{Token: token.NewToken(token.INT, `12345`), Value: 12345},
					Right: &ast.IntLiteral{Token: token.NewToken(token.INT, `67890`), Value: 67890},
				},
			}},
		},
		{
			source: `12345-67890`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.InfixExpression{
					Token: token.NewToken(token.MINUS, "-"),
					Left:  &ast.IntLiteral{Token: token.NewToken(token.INT, `12345`), Value: 12345},
					Right: &ast.IntLiteral{Token: token.NewToken(token.INT, `67890`), Value: 67890},
				},
			}},
		},
		{
			source: `12345*67890`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.InfixExpression{
					Token: token.NewToken(token.MULTIPLY, "*"),
					Left:  &ast.IntLiteral{Token: token.NewToken(token.INT, `12345`), Value: 12345},
					Right: &ast.IntLiteral{Token: token.NewToken(token.INT, `67890`), Value: 67890},
				},
			}},
		},
		{
			source: `12345/67890`,
			program: &ast.Program{Nodes: []ast.Node{
				&ast.InfixExpression{
					Token: token.NewToken(token.DIVIDE, "/"),
					Left:  &ast.IntLiteral{Token: token.NewToken(token.INT, `12345`), Value: 12345},
					Right: &ast.IntLiteral{Token: token.NewToken(token.INT, `67890`), Value: 67890},
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
						Right: &ast.IntLiteral{Token: token.NewToken(token.INT, `12345`), Value: 12345},
					},
					Right: &ast.PrefixExpression{
						Token: token.NewToken(token.MINUS, "-"),
						Right: &ast.IntLiteral{Token: token.NewToken(token.INT, `67890`), Value: 67890},
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
						Left:  &ast.IntLiteral{Token: token.NewToken(token.INT, `12`), Value: 12},
						Right: &ast.IntLiteral{Token: token.NewToken(token.INT, `34`), Value: 34},
					},
					Right: &ast.InfixExpression{
						Token: token.NewToken(token.DIVIDE, "/"),
						Left:  &ast.IntLiteral{Token: token.NewToken(token.INT, `56`), Value: 56},
						Right: &ast.IntLiteral{Token: token.NewToken(token.INT, `78`), Value: 78},
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
						Left:  &ast.IntLiteral{Token: token.NewToken(token.INT, `12`), Value: 12},
						Right: &ast.IntLiteral{Token: token.NewToken(token.INT, `34`), Value: 34},
					},
					Right: &ast.InfixExpression{
						Token: token.NewToken(token.MINUS, "-"),
						Left:  &ast.IntLiteral{Token: token.NewToken(token.INT, `56`), Value: 56},
						Right: &ast.IntLiteral{Token: token.NewToken(token.INT, `78`), Value: 78},
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
