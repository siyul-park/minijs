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
		{";", ast.NewProgram(ast.NewEmptyStatement())},
		{
			"{ 1; 2; }",
			ast.NewProgram(
				ast.NewBlockStatement(
					ast.NewExpressionStatement(
						ast.NewNumberLiteral(token.New(token.NUMBER, "1"), 1),
					),
					ast.NewExpressionStatement(
						ast.NewNumberLiteral(token.New(token.NUMBER, "2"), 2),
					),
				),
			),
		},
		{
			"a + b; c + d",
			ast.NewProgram(
				ast.NewExpressionStatement(
					ast.NewInfixExpression(
						token.PLUS,
						ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "a"), "a"),
						ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "b"), "b"),
					),
				),
				ast.NewExpressionStatement(
					ast.NewInfixExpression(
						token.PLUS,
						ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "c"), "c"),
						ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "d"), "d"),
					),
				),
			),
		},
		{
			"123",
			ast.NewProgram(
				ast.NewExpressionStatement(
					ast.NewNumberLiteral(token.New(token.NUMBER, "123"), 123),
				),
			),
		},
		{
			"1.23",
			ast.NewProgram(
				ast.NewExpressionStatement(
					ast.NewNumberLiteral(token.New(token.NUMBER, "1.23"), 1.23),
				),
			),
		},
		{
			"0b01",
			ast.NewProgram(
				ast.NewExpressionStatement(
					ast.NewNumberLiteral(token.New(token.NUMBER, "0b01"), 0b01),
				),
			),
		},
		{
			"0o01",
			ast.NewProgram(
				ast.NewExpressionStatement(
					ast.NewNumberLiteral(token.New(token.NUMBER, "0o01"), 0o01),
				),
			),
		},
		{
			"0x01",
			ast.NewProgram(
				ast.NewExpressionStatement(
					ast.NewNumberLiteral(token.New(token.NUMBER, "0x01"), 0x01),
				),
			),
		},
		{
			"true",
			ast.NewProgram(
				ast.NewExpressionStatement(
					ast.NewBoolLiteral(token.New(token.BOOLEAN, "true"), true),
				),
			),
		},
		{
			"foo",
			ast.NewProgram(
				ast.NewExpressionStatement(
					ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "foo"), "foo"),
				),
			),
		},
		{
			`"hello"`,
			ast.NewProgram(
				ast.NewExpressionStatement(
					ast.NewStringLiteral(token.New(token.STRING, "hello"), "hello"),
				),
			),
		},
		{
			"-1",
			ast.NewProgram(
				ast.NewExpressionStatement(
					ast.NewPrefixExpression(
						token.MINUS,
						ast.NewNumberLiteral(token.New(token.NUMBER, "1"), 1),
					),
				),
			),
		},
		{
			"a + b",
			ast.NewProgram(
				ast.NewExpressionStatement(
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
				ast.NewExpressionStatement(
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
				ast.NewExpressionStatement(
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
