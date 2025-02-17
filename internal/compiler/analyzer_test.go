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
		node ast.Node
		meta *Semantic
	}{
		{
			node: ast.NewProgram(),
			meta: &Semantic{Kind: interpreter.KindVoid},
		},
		{
			node: ast.NewEmptyStatement(),
			meta: &Semantic{Kind: interpreter.KindVoid},
		},
		{
			node: ast.NewBlockStatement(),
			meta: &Semantic{Kind: interpreter.KindVoid},
		},
		{
			node: ast.NewExpressionStatement(
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
			),
			meta: &Semantic{Kind: interpreter.KindVoid},
		},
		{
			node: ast.NewPrefixExpression(
				token.New(token.PLUS, "+"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
			),
			meta: &Semantic{Kind: interpreter.KindInt32},
		},
		{
			node: ast.NewPrefixExpression(
				token.New(token.PLUS, "+"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1.5"}, 1.5),
			),
			meta: &Semantic{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewPrefixExpression(
				token.New(token.PLUS, "+"),
				ast.NewStringLiteral(token.Token{Type: token.STRING, Literal: "foo"}, "foo"),
			),
			meta: &Semantic{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewPrefixExpression(
				token.New(token.MINUS, "-"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
			),
			meta: &Semantic{Kind: interpreter.KindInt32},
		},
		{
			node: ast.NewPrefixExpression(
				token.New(token.MINUS, "-"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2.0"}, 2),
			),
			meta: &Semantic{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewInfixExpression(
				token.New(token.PLUS, "+"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2"}, 2),
			),
			meta: &Semantic{Kind: interpreter.KindInt32},
		},
		{
			node: ast.NewInfixExpression(
				token.New(token.PLUS, "+"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2.0"}, 2),
			),
			meta: &Semantic{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewInfixExpression(
				token.New(token.PLUS, "+"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewStringLiteral(token.Token{Type: token.STRING, Literal: "2"}, "2"),
			),
			meta: &Semantic{Kind: interpreter.KindString},
		},
		{
			node: ast.NewInfixExpression(
				token.New(token.MINUS, "-"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2"}, 2),
			),
			meta: &Semantic{Kind: interpreter.KindInt32},
		},
		{
			node: ast.NewInfixExpression(
				token.New(token.MULTIPLY, "*"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2"}, 2),
			),
			meta: &Semantic{Kind: interpreter.KindInt32},
		},
		{
			node: ast.NewInfixExpression(
				token.New(token.MULTIPLY, "*"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2.0"}, 2),
			),
			meta: &Semantic{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewInfixExpression(
				token.New(token.MULTIPLY, "*"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewStringLiteral(token.Token{Type: token.STRING, Literal: "2"}, "2"),
			),
			meta: &Semantic{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewInfixExpression(
				token.New(token.DIVIDE, "/"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2"}, 2),
			),
			meta: &Semantic{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewInfixExpression(
				token.New(token.DIVIDE, "/"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 2),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2.0"}, 2),
			),
			meta: &Semantic{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewInfixExpression(
				token.New(token.DIVIDE, "/"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewStringLiteral(token.Token{Type: token.STRING, Literal: "2"}, "2"),
			),
			meta: &Semantic{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewInfixExpression(
				token.New(token.MODULUS, "%"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2"}, 2),
			),
			meta: &Semantic{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewInfixExpression(
				token.New(token.MODULUS, "%"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 2),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2.0"}, 2),
			),
			meta: &Semantic{Kind: interpreter.KindFloat64},
		},
		{
			node: ast.NewInfixExpression(
				token.New(token.MODULUS, "%"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewStringLiteral(token.Token{Type: token.STRING, Literal: "2"}, "2"),
			),
			meta: &Semantic{Kind: interpreter.KindFloat64},
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
