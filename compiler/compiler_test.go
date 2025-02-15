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
		name         string
		node         ast.Node
		instructions []bytecode.Instruction
		constants    []string
	}{
		{
			name: "Number Literal",
			node: ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
			},
		},
		{
			name: "String Literal",
			node: ast.NewStringLiteral(token.Token{Type: token.STRING, Literal: "abc"}, "abc"),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.CLOAD, 0, 3),
			},
			constants: []string{"abc"},
		},
		{
			name: "String Concatenation",
			node: ast.NewInfixExpression(
				token.NewToken(token.PLUS, "+"),
				ast.NewStringLiteral(token.Token{Type: token.STRING, Literal: "foo"}, "foo"),
				ast.NewStringLiteral(token.Token{Type: token.STRING, Literal: "bar"}, "bar"),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.CLOAD, 0, 3),
				bytecode.New(bytecode.CLOAD, 4, 3),
				bytecode.New(bytecode.CADD),
			},
			constants: []string{"foo", "bar"},
		},
		{
			name: "Addition",
			node: ast.NewInfixExpression(
				token.NewToken(token.PLUS, "+"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2"}, 2),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64ADD),
			},
		},
		{
			name: "Subtraction",
			node: ast.NewInfixExpression(
				token.NewToken(token.MINUS, "-"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2"}, 2),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64SUB),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiler := New()
			code := bytecode.Bytecode{}
			code.Add(tt.instructions...)
			for _, c := range tt.constants {
				code.Store([]byte(c + "\x00"))
			}

			result, err := compiler.Compile(tt.node)
			assert.NoError(t, err)
			assert.Equal(t, code.String(), result.String())
		})
	}
}
