package compiler

import (
	"math"
	"testing"

	"github.com/siyul-park/minijs/ast"
	"github.com/siyul-park/minijs/bytecode"
	"github.com/siyul-park/minijs/token"
	"github.com/stretchr/testify/assert"
)

func TestCompiler_Compile(t *testing.T) {
	tests := []struct {
		node         ast.Node
		instructions []bytecode.Instruction
		constants    [][]byte
	}{
		{
			node: ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LD, math.Float64bits(1)),
			},
		},
		{
			node: ast.NewStringLiteral(token.Token{Type: token.STRING, Literal: "abc"}, "abc"),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.CLD, 0, 3),
			},
			constants: [][]byte{[]byte("abc")},
		},
		{
			node: ast.NewPrefixExpression(
				token.NewToken(token.MINUS, "-"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LD, math.Float64bits(-1)),
				bytecode.New(bytecode.F64MUL),
			},
		},
		{
			node: ast.NewInfixExpression(
				token.NewToken(token.PLUS, "+"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2"}, 2),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LD, math.Float64bits(2)),
				bytecode.New(bytecode.F64ADD),
			},
		},
		{
			node: ast.NewInfixExpression(
				token.NewToken(token.MINUS, "-"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2"}, 2),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LD, math.Float64bits(2)),
				bytecode.New(bytecode.F64SUB),
			},
		},
		{
			node: ast.NewInfixExpression(
				token.NewToken(token.ASTERISK, "*"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2"}, 2),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LD, math.Float64bits(2)),
				bytecode.New(bytecode.F64MUL),
			},
		},
		{
			node: ast.NewInfixExpression(
				token.NewToken(token.SLASH, "/"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2"}, 2),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LD, math.Float64bits(2)),
				bytecode.New(bytecode.F64DIV),
			},
		},
	}

	for _, tt := range tests {
		var code bytecode.Bytecode
		code.Append(tt.instructions...)
		for _, c := range tt.constants {
			code.Store(c)
		}

		t.Run(tt.node.String(), func(t *testing.T) {
			compiler := New(tt.node)

			result, err := compiler.Compile()
			assert.NoError(t, err)
			assert.Equal(t, code.String(), result.String())
		})
	}
}
