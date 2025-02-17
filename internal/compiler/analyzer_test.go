package compiler

import (
	"testing"

	"github.com/siyul-park/minijs/internal/ast"
	"github.com/siyul-park/minijs/internal/interpreter"
	"github.com/siyul-park/minijs/internal/token"

	"github.com/stretchr/testify/assert"
)

func TestAnalyzer_Analyze(t *testing.T) {
	tests := []struct {
		node   ast.Node
		symbol *Symbol
	}{
		{
			node:   ast.NewProgram(),
			symbol: &Symbol{Kind: interpreter.KindVoid},
		},
		{
			node:   ast.NewEmptyStatement(),
			symbol: &Symbol{Kind: interpreter.KindVoid},
		},
		{
			node:   ast.NewBlockStatement(),
			symbol: &Symbol{Kind: interpreter.KindVoid},
		},
		{
			node: ast.NewExpressionStatement(
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
			),
			symbol: &Symbol{Kind: interpreter.KindVoid},
		},
		{
			node: ast.NewPrefixExpression(
				token.New(token.PLUS, "+"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
			),
			symbol: &Symbol{Kind: interpreter.KindInt32},
		},
		{
			node: ast.NewPrefixExpression(
				token.New(token.PLUS, "+"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1.5"}, 1.5),
			),
			symbol: &Symbol{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewPrefixExpression(
				token.New(token.PLUS, "+"),
				ast.NewStringLiteral(token.Token{Type: token.STRING, Literal: "foo"}, "foo"),
			),
			symbol: &Symbol{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewPrefixExpression(
				token.New(token.MINUS, "-"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
			),
			symbol: &Symbol{Kind: interpreter.KindInt32},
		},
		{
			node: ast.NewPrefixExpression(
				token.New(token.MINUS, "-"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2.0"}, 2),
			),
			symbol: &Symbol{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewInfixExpression(
				token.New(token.PLUS, "+"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2"}, 2),
			),
			symbol: &Symbol{Kind: interpreter.KindInt32},
		},
		{
			node: ast.NewInfixExpression(
				token.New(token.PLUS, "+"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2.0"}, 2),
			),
			symbol: &Symbol{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewInfixExpression(
				token.New(token.PLUS, "+"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewStringLiteral(token.Token{Type: token.STRING, Literal: "2"}, "2"),
			),
			symbol: &Symbol{Kind: interpreter.KindString},
		},
		{
			node: ast.NewInfixExpression(
				token.New(token.MINUS, "-"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2"}, 2),
			),
			symbol: &Symbol{Kind: interpreter.KindInt32},
		},
		{
			node: ast.NewInfixExpression(
				token.New(token.MULTIPLY, "*"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2"}, 2),
			),
			symbol: &Symbol{Kind: interpreter.KindInt32},
		},
		{
			node: ast.NewInfixExpression(
				token.New(token.MULTIPLY, "*"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2.0"}, 2),
			),
			symbol: &Symbol{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewInfixExpression(
				token.New(token.MULTIPLY, "*"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewStringLiteral(token.Token{Type: token.STRING, Literal: "2"}, "2"),
			),
			symbol: &Symbol{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewInfixExpression(
				token.New(token.DIVIDE, "/"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2"}, 2),
			),
			symbol: &Symbol{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewInfixExpression(
				token.New(token.DIVIDE, "/"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 2),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2.0"}, 2),
			),
			symbol: &Symbol{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewInfixExpression(
				token.New(token.DIVIDE, "/"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewStringLiteral(token.Token{Type: token.STRING, Literal: "2"}, "2"),
			),
			symbol: &Symbol{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewInfixExpression(
				token.New(token.MODULUS, "%"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2"}, 2),
			),
			symbol: &Symbol{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewInfixExpression(
				token.New(token.MODULUS, "%"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 2),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2.0"}, 2),
			),
			symbol: &Symbol{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewInfixExpression(
				token.New(token.MODULUS, "%"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewStringLiteral(token.Token{Type: token.STRING, Literal: "2"}, "2"),
			),
			symbol: &Symbol{Kind: interpreter.KindFloat64},
		},
	}

	analyzer := NewAnalyzer()

	for _, tt := range tests {
		t.Run(tt.node.String(), func(t *testing.T) {
			sb := analyzer.analyze(tt.node)
			assert.Equal(t, tt.symbol, sb)
		})
	}
}
