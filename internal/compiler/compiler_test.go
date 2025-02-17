package compiler

import (
	"math"
	"testing"

	"github.com/siyul-park/minijs/internal/ast"
	"github.com/siyul-park/minijs/internal/bytecode"
	"github.com/siyul-park/minijs/internal/token"

	"github.com/stretchr/testify/assert"
)

func TestCompiler_Compile(t *testing.T) {
	tests := []struct {
		node         ast.Node
		instructions []bytecode.Instruction
		literals     []string
	}{
		{
			node: ast.NewEmptyStatement(),
		},
		{
			node: ast.NewExpressionStatement(
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.POP),
			},
		},
		{
			node: ast.NewBlockStatement(
				ast.NewExpressionStatement(
					ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				),
				ast.NewExpressionStatement(
					ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "2"}, 2),
				),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.POP),
				bytecode.New(bytecode.I32LOAD, 2),
				bytecode.New(bytecode.POP),
			},
		},
		{
			node: ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
			},
		},
		{
			node: ast.NewPrefixExpression(
				token.New(token.PLUS, "+"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
			},
		},
		{
			node: ast.NewPrefixExpression(
				token.New(token.MINUS, "-"),
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
				token.New(token.PLUS, "+"),
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
				token.New(token.PLUS, "+"),
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
				token.New(token.PLUS, "+"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				ast.NewStringLiteral(token.Token{Type: token.STRING, Literal: "2"}, "2"),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.I32TOC),
				bytecode.New(bytecode.CLOAD, 0, 1),
				bytecode.New(bytecode.CADD),
			},
			literals: []string{"2"},
		},
		{
			node: ast.NewInfixExpression(
				token.New(token.MINUS, "-"),
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
				token.New(token.MULTIPLY, "*"),
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
				token.New(token.DIVIDE, "/"),
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
				token.New(token.PLUS, "+"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1.0"}, 1),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
			},
		},
		{
			node: ast.NewPrefixExpression(
				token.New(token.MINUS, "-"),
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
				token.New(token.PLUS, "+"),
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
				token.New(token.PLUS, "+"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1.0"}, 1),
				ast.NewStringLiteral(token.Token{Type: token.STRING, Literal: "2"}, "2"),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64TOC),
				bytecode.New(bytecode.CLOAD, 0, 1),
				bytecode.New(bytecode.CADD),
			},
			literals: []string{"2"},
		},
		{
			node: ast.NewInfixExpression(
				token.New(token.MINUS, "-"),
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
				token.New(token.MULTIPLY, "*"),
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
				token.New(token.DIVIDE, "/"),
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
				token.New(token.MODULUS, "%"),
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
			node: ast.NewStringLiteral(token.Token{Type: token.STRING, Literal: "abc"}, "abc"),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.CLOAD, 0, 3),
			},
			literals: []string{"abc"},
		},
		{
			node: ast.NewInfixExpression(
				token.New(token.PLUS, "+"),
				ast.NewStringLiteral(token.Token{Type: token.STRING, Literal: "foo"}, "foo"),
				ast.NewStringLiteral(token.Token{Type: token.STRING, Literal: "bar"}, "bar"),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.CLOAD, 0, 3),
				bytecode.New(bytecode.CLOAD, 4, 3),
				bytecode.New(bytecode.CADD),
			},
			literals: []string{"foo", "bar"},
		},
		{
			node: ast.NewBoolLiteral(token.Token{Type: token.TRUE, Literal: "true"}, true),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.BLOAD, 1),
			},
		},
		{
			node: ast.NewPrefixExpression(
				token.New(token.PLUS, "+"),
				ast.NewBoolLiteral(token.Token{Type: token.TRUE, Literal: "true"}, true),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.BLOAD, 1),
				bytecode.New(bytecode.BTOI32),
			},
		},
		{
			node: ast.NewPrefixExpression(
				token.New(token.MINUS, "-"),
				ast.NewBoolLiteral(token.Token{Type: token.TRUE, Literal: "true"}, true),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.BLOAD, 1),
				bytecode.New(bytecode.BTOI32),
				bytecode.New(bytecode.I32LOAD, uint64(0xFFFFFFFFFFFFFFFF)),
				bytecode.New(bytecode.I32MUL),
			},
		},
		{
			node: ast.NewInfixExpression(
				token.New(token.MINUS, "-"),
				ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1.0"}, 1.0),
				ast.NewBoolLiteral(token.Token{Type: token.TRUE, Literal: "true"}, true),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.BLOAD, 1),
				bytecode.New(bytecode.BTOI32),
				bytecode.New(bytecode.I32TOF64),
				bytecode.New(bytecode.F64SUB),
			},
		},
	}

	compiler := New()

	for _, tt := range tests {
		t.Run(tt.node.String(), func(t *testing.T) {
			code := bytecode.Bytecode{}
			code.Emit(tt.instructions...)
			for _, c := range tt.literals {
				code.Store([]byte(c + "\x00"))
			}

			result, err := compiler.Compile(tt.node)
			assert.NoError(t, err)
			assert.Equal(t, code.String(), result.String())
		})
	}
}
