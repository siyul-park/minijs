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
		constants    []string
	}{
		{
			node: ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
			},
		},
		{
			node: ast.NewPrefixExpression(
				token.NewToken(token.PLUS, "+"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
			},
		},
		{
			node: ast.NewPrefixExpression(
				token.NewToken(token.MINUS, "-"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.I32LOAD, uint64(0xFFFFFFFFFFFFFFFF)),
				bytecode.New(bytecode.I32MUL),
			},
		},
		{
			node: ast.NewInfixExpression(
				token.NewToken(token.PLUS, "+"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2"}, 2),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.I32LOAD, 2),
				bytecode.New(bytecode.I32ADD),
			},
		},
		{
			node: ast.NewInfixExpression(
				token.NewToken(token.PLUS, "+"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2.0"}, 2),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.I32TOF64),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64ADD),
			},
		},
		{
			node: ast.NewInfixExpression(
				token.NewToken(token.PLUS, "+"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewStringLiteral(token.Token{Type: token.KindString, Literal: "2"}, "2"),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.I3TO2C),
				bytecode.New(bytecode.CLOAD, 0, 1),
				bytecode.New(bytecode.CADD),
			},
			constants: []string{"2"},
		},
		{
			node: ast.NewInfixExpression(
				token.NewToken(token.MINUS, "-"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2"}, 2),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.I32LOAD, 2),
				bytecode.New(bytecode.I32SUB),
			},
		},
		{
			node: ast.NewInfixExpression(
				token.NewToken(token.ASTERISK, "*"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2"}, 2),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.I32LOAD, 2),
				bytecode.New(bytecode.I32MUL),
			},
		},
		{
			node: ast.NewInfixExpression(
				token.NewToken(token.SLASH, "/"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2"}, 2),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.I32TOF64),
				bytecode.New(bytecode.I32LOAD, 2),
				bytecode.New(bytecode.I32TOF64),
				bytecode.New(bytecode.F64DIV),
			},
		},

		{
			node: ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1.0"}, 1),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
			},
		},
		{
			node: ast.NewPrefixExpression(
				token.NewToken(token.PLUS, "+"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1.0"}, 1),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
			},
		},
		{
			node: ast.NewPrefixExpression(
				token.NewToken(token.MINUS, "-"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1.0"}, 1),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(-1)),
				bytecode.New(bytecode.F64MUL),
			},
		},
		{
			node: ast.NewInfixExpression(
				token.NewToken(token.PLUS, "+"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1.0"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2.0"}, 2),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64ADD),
			},
		},
		{
			node: ast.NewInfixExpression(
				token.NewToken(token.PLUS, "+"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1.0"}, 1),
				ast.NewStringLiteral(token.Token{Type: token.KindString, Literal: "2"}, "2"),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64TOC),
				bytecode.New(bytecode.CLOAD, 0, 1),
				bytecode.New(bytecode.CADD),
			},
			constants: []string{"2"},
		},
		{
			node: ast.NewInfixExpression(
				token.NewToken(token.MINUS, "-"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1.0"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2.0"}, 2),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64SUB),
			},
		},
		{
			node: ast.NewInfixExpression(
				token.NewToken(token.ASTERISK, "*"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1.0"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2.0"}, 2),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64MUL),
			},
		},
		{
			node: ast.NewInfixExpression(
				token.NewToken(token.SLASH, "/"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1.0"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2.0"}, 2),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64DIV),
			},
		},
		{
			node: ast.NewInfixExpression(
				token.NewToken(token.PERCENT, "%"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1.0"}, 1),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2.0"}, 2),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64MOD),
			},
		},

		{
			node: ast.NewStringLiteral(token.Token{Type: token.KindString, Literal: "abc"}, "abc"),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.CLOAD, 0, 3),
			},
			constants: []string{"abc"},
		},
		{
			node: ast.NewInfixExpression(
				token.NewToken(token.PLUS, "+"),
				ast.NewStringLiteral(token.Token{Type: token.KindString, Literal: "foo"}, "foo"),
				ast.NewStringLiteral(token.Token{Type: token.KindString, Literal: "bar"}, "bar"),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.CLOAD, 0, 3),
				bytecode.New(bytecode.CLOAD, 4, 3),
				bytecode.New(bytecode.CADD),
			},
			constants: []string{"foo", "bar"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.node.String(), func(t *testing.T) {
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
