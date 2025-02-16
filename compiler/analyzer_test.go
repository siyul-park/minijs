package compiler

import (
	"testing"

	"github.com/siyul-park/minijs/ast"
	"github.com/siyul-park/minijs/interpreter"
	"github.com/siyul-park/minijs/token"
	"github.com/stretchr/testify/assert"
)

func TestAnalyzer_Analyze(t *testing.T) {
	tests := []struct {
		node ast.Node
		meta *Meta
	}{
		{
			node: ast.NewProgram(),
			meta: &Meta{Kind: interpreter.KindVoid},
		},
		{
			node: ast.NewEmptyStatement(),
			meta: &Meta{Kind: interpreter.KindVoid},
		},
		{
			node: ast.NewBlockStatement(),
			meta: &Meta{Kind: interpreter.KindVoid},
		},
		{
			node: ast.NewExpressionStatement(
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
			),
			meta: &Meta{Kind: interpreter.KindVoid},
		},
		{
			node: ast.NewPrefixExpression(
				token.PLUS,
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
			),
			meta: &Meta{Kind: interpreter.KindInt32},
		},
		{
			node: ast.NewPrefixExpression(
				token.PLUS,
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1.5"}, 1.5),
			),
			meta: &Meta{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewPrefixExpression(
				token.PLUS,
				ast.NewStringLiteral(token.Token{Type: token.STRING, Literal: "foo"}, "foo"),
			),
			meta: &Meta{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewPrefixExpression(
				token.MINUS,
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
			),
			meta: &Meta{Kind: interpreter.KindInt32},
		},
		{
			node: ast.NewPrefixExpression(
				token.MINUS,
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2.0"}, 2),
			),
			meta: &Meta{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewInfixExpression(
				token.PLUS,
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2"}, 2),
			),
			meta: &Meta{Kind: interpreter.KindInt32},
		},
		{
			node: ast.NewInfixExpression(
				token.PLUS,
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2.0"}, 2),
			),
			meta: &Meta{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewInfixExpression(
				token.PLUS,
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewStringLiteral(token.Token{Type: token.STRING, Literal: "2"}, "2"),
			),
			meta: &Meta{Kind: interpreter.KindString},
		},
		{
			node: ast.NewInfixExpression(
				token.MINUS,
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2"}, 2),
			),
			meta: &Meta{Kind: interpreter.KindInt32},
		},
		{
			node: ast.NewInfixExpression(
				token.MULTIPLE,
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2"}, 2),
			),
			meta: &Meta{Kind: interpreter.KindInt32},
		},
		{
			node: ast.NewInfixExpression(
				token.MULTIPLE,
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2.0"}, 2),
			),
			meta: &Meta{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewInfixExpression(
				token.MULTIPLE,
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewStringLiteral(token.Token{Type: token.STRING, Literal: "2"}, "2"),
			),
			meta: &Meta{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewInfixExpression(
				token.DIVIDE,
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2"}, 2),
			),
			meta: &Meta{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewInfixExpression(
				token.DIVIDE,
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 2),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2.0"}, 2),
			),
			meta: &Meta{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewInfixExpression(
				token.DIVIDE,
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewStringLiteral(token.Token{Type: token.STRING, Literal: "2"}, "2"),
			),
			meta: &Meta{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewInfixExpression(
				token.MODULO,
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2"}, 2),
			),
			meta: &Meta{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewInfixExpression(
				token.MODULO,
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 2),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2.0"}, 2),
			),
			meta: &Meta{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewInfixExpression(
				token.MODULO,
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewStringLiteral(token.Token{Type: token.STRING, Literal: "2"}, "2"),
			),
			meta: &Meta{Kind: interpreter.KindFloat64},
		},
	}

	analyzer := NewAnalyzer()

	for _, tt := range tests {
		t.Run(tt.node.String(), func(t *testing.T) {
			meta := analyzer.Analyze(tt.node)
			assert.Equal(t, tt.meta, meta)
		})
	}
}
