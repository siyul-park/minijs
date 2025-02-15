package parser

import (
	"testing"

	"github.com/siyul-park/minijs/lexer"
	"github.com/stretchr/testify/assert"

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
						token.PLUS,
						ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "a"), "a"),
						ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "b"), "b"),
					),
				),
				ast.NewStatement(
					ast.NewInfixExpression(
						token.PLUS,
						ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "c"), "c"),
						ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "d"), "d"),
					),
				),
			),
		},
		{
			"1234567890",
			ast.NewProgram(
				ast.NewStatement(
					ast.NewNumberLiteral(token.New(token.NUMBER, "1234567890"), 1234567890),
				),
			),
		},
		{
			`"hello"`,
			ast.NewProgram(
				ast.NewStatement(
					ast.NewStringLiteral(token.New(token.STRING, "hello"), "hello"),
				),
			),
		},
		{
			"true",
			ast.NewProgram(
				ast.NewStatement(
					ast.NewBoolLiteral(token.New(token.BOOLEAN, "true"), true),
				),
			),
		},
		{
			"a + b",
			ast.NewProgram(
				ast.NewStatement(
					ast.NewInfixExpression(
						token.PLUS,
						ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "a"), "a"),
						ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "b"), "b"),
					),
				),
			),
		},
		{
			"a + b + c",
			ast.NewProgram(
				ast.NewStatement(
					ast.NewInfixExpression(
						token.PLUS,
						ast.NewInfixExpression(
							token.PLUS,
							ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "a"), "a"),
							ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "b"), "b"),
						),
						ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "c"), "c"),
					),
				),
			),
		},
		{
			"a * b + c",
			ast.NewProgram(
				ast.NewStatement(
					ast.NewInfixExpression(
						token.PLUS,
						ast.NewInfixExpression(
							token.MULTIPLE,
							ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "a"), "a"),
							ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "b"), "b"),
						),
						ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "c"), "c"),
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
