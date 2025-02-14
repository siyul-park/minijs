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
			source: `;`,
			program: &ast.Program{Statements: []*ast.Statement{
				{Node: nil},
			}},
		},
		{
			source: `1234567890`,
			program: &ast.Program{Statements: []*ast.Statement{
				{Node: &ast.NumberLiteral{
					Token: token.NewToken(token.NUMBER, `1234567890`),
					Value: 1234567890,
				}},
			}},
		},
		{
			source: `1.234567890`,
			program: &ast.Program{Statements: []*ast.Statement{
				{Node: &ast.NumberLiteral{
					Token: token.NewToken(token.NUMBER, `1.234567890`),
					Value: 1.234567890,
				}},
			}},
		},
		{
			source: `0e-5`,
			program: &ast.Program{Statements: []*ast.Statement{
				{Node: &ast.NumberLiteral{
					Token: token.NewToken(token.NUMBER, `0e-5`),
					Value: 0,
				}},
			}},
		},
		{
			source: `0e+5`,
			program: &ast.Program{Statements: []*ast.Statement{
				{Node: &ast.NumberLiteral{
					Token: token.NewToken(token.NUMBER, `0e+5`),
					Value: 0,
				}},
			}},
		},
		{
			source: `5e1`,
			program: &ast.Program{Statements: []*ast.Statement{
				{Node: &ast.NumberLiteral{
					Token: token.NewToken(token.NUMBER, `5e1`),
					Value: 50,
				}},
			}},
		},
		{
			source: `175e-2`,
			program: &ast.Program{Statements: []*ast.Statement{
				{Node: &ast.NumberLiteral{
					Token: token.NewToken(token.NUMBER, `175e-2`),
					Value: 1.75,
				}},
			}},
		},
		{
			source: `"hello"`,
			program: &ast.Program{Statements: []*ast.Statement{
				{Node: &ast.StringLiteral{
					Token: token.NewToken(token.STRING, `hello`),
					Value: `hello`,
				}},
			}},
		},
		{
			source: `true`,
			program: &ast.Program{Statements: []*ast.Statement{
				{Node: &ast.BoolLiteral{
					Token: token.NewToken(token.BOOLEAN, `true`),
					Value: true,
				}},
			}},
		},
		{
			source: `false`,
			program: &ast.Program{Statements: []*ast.Statement{
				{Node: &ast.BoolLiteral{
					Token: token.NewToken(token.BOOLEAN, `false`),
					Value: false,
				}},
			}},
		},
		//{
		//	source: `!a`,
		//	program: &ast.Program{Statements: []*ast.Statement{
		//		{Node: &ast.PrefixExpression{
		//			Token: token.NewToken(token.BANG, `!`),
		//			Right: &ast.IdentifierLiteral{
		//				Token: token.NewToken(token.IDENTIFIER, `a`),
		//				Value: `a`,
		//			},
		//		}},
		//	}},
		//},
		{
			source: `a + b`,
			program: &ast.Program{Statements: []*ast.Statement{
				{Node: &ast.InfixExpression{
					Token: token.NewToken(token.PLUS, `+`),
					Left: &ast.IdentifierLiteral{
						Token: token.NewToken(token.IDENTIFIER, `a`),
						Value: `a`,
					},
					Right: &ast.IdentifierLiteral{
						Token: token.NewToken(token.IDENTIFIER, `b`),
						Value: `b`,
					},
				}},
			}},
		},
		{
			source: `a + b + c`,
			program: &ast.Program{Statements: []*ast.Statement{
				{Node: &ast.InfixExpression{
					Token: token.NewToken(token.PLUS, `+`),
					Left: &ast.InfixExpression{
						Token: token.NewToken(token.PLUS, `+`),
						Left: &ast.IdentifierLiteral{
							Token: token.NewToken(token.IDENTIFIER, `a`),
							Value: `a`,
						},
						Right: &ast.IdentifierLiteral{
							Token: token.NewToken(token.IDENTIFIER, `b`),
							Value: `b`,
						},
					},
					Right: &ast.IdentifierLiteral{
						Token: token.NewToken(token.IDENTIFIER, `c`),
						Value: `c`,
					},
				}},
			}},
		},
		{
			source: `a * b + c`,
			program: &ast.Program{Statements: []*ast.Statement{
				{Node: &ast.InfixExpression{
					Token: token.NewToken(token.PLUS, `+`),
					Left: &ast.InfixExpression{
						Token: token.NewToken(token.ASTERISK, `*`),
						Left: &ast.IdentifierLiteral{
							Token: token.NewToken(token.IDENTIFIER, `a`),
							Value: `a`,
						},
						Right: &ast.IdentifierLiteral{
							Token: token.NewToken(token.IDENTIFIER, `b`),
							Value: `b`,
						},
					},
					Right: &ast.IdentifierLiteral{
						Token: token.NewToken(token.IDENTIFIER, `c`),
						Value: `c`,
					},
				}},
			}},
		},
		{
			source: `(a + b) * c`,
			program: &ast.Program{Statements: []*ast.Statement{
				{Node: &ast.InfixExpression{
					Token: token.NewToken(token.ASTERISK, `*`),
					Left: &ast.InfixExpression{
						Token: token.NewToken(token.PLUS, `+`),
						Left: &ast.IdentifierLiteral{
							Token: token.NewToken(token.IDENTIFIER, `a`),
							Value: `a`,
						},
						Right: &ast.IdentifierLiteral{
							Token: token.NewToken(token.IDENTIFIER, `b`),
							Value: `b`,
						},
					},
					Right: &ast.IdentifierLiteral{
						Token: token.NewToken(token.IDENTIFIER, `c`),
						Value: `c`,
					},
				}},
			}},
		},
		//{
		//	source: `a = b`,
		//	program: &ast.Program{Statements: []*ast.Statement{
		//		{Node: &ast.InfixExpression{
		//			Token: token.NewToken(token.EQ, `=`),
		//			Left: &ast.IdentifierLiteral{
		//				Token: token.NewToken(token.IDENTIFIER, `a`),
		//				Value: `a`,
		//			},
		//			Right: &ast.IdentifierLiteral{
		//				Token: token.NewToken(token.IDENTIFIER, `b`),
		//				Value: `b`,
		//			},
		//		}},
		//	}},
		//},
		//{
		//	source: `foo()`,
		//	program: &ast.Program{Statements: []*ast.Statement{
		//		{Node: &ast.PrefixExpression{
		//			Token: token.NewToken(token.LPAREN, `(`),
		//			Right: &ast.IdentifierLiteral{
		//				Token: token.NewToken(token.IDENTIFIER, `foo`),
		//				Value: `foo`,
		//			},
		//		}},
		//	}},
		//},
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
