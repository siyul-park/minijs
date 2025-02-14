package parser

import (
	"github.com/siyul-park/minijs/lexer"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/siyul-park/minijs/ast"
	"github.com/siyul-park/minijs/token"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		source  string
		program *ast.Program
	}{
		{"", ast.NewProgram()},
		{";", ast.NewProgram(ast.NewStatement(nil))},
		{
			"a + b; c + d",
			ast.NewProgram(
				ast.NewStatement(
					ast.NewInfixExpression(
						token.NewToken(token.PLUS, "+"),
						ast.NewIdentifierLiteral(token.NewToken(token.IDENTIFIER, "a"), "a"),
						ast.NewIdentifierLiteral(token.NewToken(token.IDENTIFIER, "b"), "b"),
					),
				),
				ast.NewStatement(
					ast.NewInfixExpression(
						token.NewToken(token.PLUS, "+"),
						ast.NewIdentifierLiteral(token.NewToken(token.IDENTIFIER, "c"), "c"),
						ast.NewIdentifierLiteral(token.NewToken(token.IDENTIFIER, "d"), "d"),
					),
				),
			),
		},
		{
			"1234567890",
			ast.NewProgram(
				ast.NewStatement(
					ast.NewNumberLiteral(token.NewToken(token.NUMBER, "1234567890"), 1234567890),
				),
			),
		},
		{
			`"hello"`,
			ast.NewProgram(
				ast.NewStatement(
					ast.NewStringLiteral(token.NewToken(token.STRING, "hello"), "hello"),
				),
			),
		},
		{
			"true",
			ast.NewProgram(
				ast.NewStatement(
					ast.NewBoolLiteral(token.NewToken(token.BOOLEAN, "true"), true),
				),
			),
		},
		{
			"a + b",
			ast.NewProgram(
				ast.NewStatement(
					ast.NewInfixExpression(
						token.NewToken(token.PLUS, "+"),
						ast.NewIdentifierLiteral(token.NewToken(token.IDENTIFIER, "a"), "a"),
						ast.NewIdentifierLiteral(token.NewToken(token.IDENTIFIER, "b"), "b"),
					),
				),
			),
		},
		{
			"a + b + c",
			ast.NewProgram(
				ast.NewStatement(
					ast.NewInfixExpression(
						token.NewToken(token.PLUS, "+"),
						ast.NewInfixExpression(
							token.NewToken(token.PLUS, "+"),
							ast.NewIdentifierLiteral(token.NewToken(token.IDENTIFIER, "a"), "a"),
							ast.NewIdentifierLiteral(token.NewToken(token.IDENTIFIER, "b"), "b"),
						),
						ast.NewIdentifierLiteral(token.NewToken(token.IDENTIFIER, "c"), "c"),
					),
				),
			),
		},
		{
			"a * b + c",
			ast.NewProgram(
				ast.NewStatement(
					ast.NewInfixExpression(
						token.NewToken(token.PLUS, "+"),
						ast.NewInfixExpression(
							token.NewToken(token.ASTERISK, "*"),
							ast.NewIdentifierLiteral(token.NewToken(token.IDENTIFIER, "a"), "a"),
							ast.NewIdentifierLiteral(token.NewToken(token.IDENTIFIER, "b"), "b"),
						),
						ast.NewIdentifierLiteral(token.NewToken(token.IDENTIFIER, "c"), "c"),
					),
				),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.source, func(t *testing.T) {
			l := lexer.New(tt.source)
			p := New(l)
			program, err := p.Parse()
			assert.NoError(t, err)
			assert.Equal(t, tt.program, program)
		})
	}
}
