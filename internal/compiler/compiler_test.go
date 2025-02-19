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
			node: ast.NewNullLiteral(token.Token{Type: token.NULL, Literal: "null"}),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.NULLLOAD),
			},
		},
		{
			node: ast.NewUndefinedLiteral(token.Token{Type: token.UNDEFINED, Literal: "undefined"}),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.UNDEFLOAD),
			},
		},
		{
			node: ast.NewBoolLiteral(token.Token{Type: token.TRUE, Literal: "true"}, true),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.BOOLLOAD, 1),
			},
		},
		{
			node: ast.NewPrefixExpression(
				token.New(token.PLUS, "+"),
				ast.NewBoolLiteral(token.Token{Type: token.TRUE, Literal: "true"}, true),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.BOOLLOAD, 1),
				bytecode.New(bytecode.BOOLTOI32),
			},
		},
		{
			node: ast.NewPrefixExpression(
				token.New(token.MINUS, "-"),
				ast.NewBoolLiteral(token.Token{Type: token.TRUE, Literal: "true"}, true),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.BOOLLOAD, 1),
				bytecode.New(bytecode.BOOLTOI32),
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
				bytecode.New(bytecode.BOOLLOAD, 1),
				bytecode.New(bytecode.BOOLTOI32),
				bytecode.New(bytecode.I32TOF64),
				bytecode.New(bytecode.F64SUB),
			},
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
				bytecode.New(bytecode.I32TOSTR),
				bytecode.New(bytecode.STRLOAD, 0, 1),
				bytecode.New(bytecode.STRADD),
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
				bytecode.New(bytecode.F64TOSTR),
				bytecode.New(bytecode.STRLOAD, 0, 1),
				bytecode.New(bytecode.STRADD),
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
				bytecode.New(bytecode.STRLOAD, 0, 3),
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
				bytecode.New(bytecode.STRLOAD, 0, 3),
				bytecode.New(bytecode.STRLOAD, 4, 3),
				bytecode.New(bytecode.STRADD),
			},
			literals: []string{"foo", "bar"},
		},
		{
			node: ast.NewExpressionStatement(
				ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "foo"), "foo"),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.SLTLOAD, 0),
				bytecode.New(bytecode.POP),
			},
		},
		{
			node: ast.NewVariableStatement(
				token.New(token.VAR, "var"),
				ast.NewAssignmentExpression(
					token.New(token.ASSIGN, "="),
					ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "foo"), "foo"),
					ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.SLTSTORE, 0),
				bytecode.New(bytecode.SLTLOAD, 0),
				bytecode.New(bytecode.POP),
			},
		},
		{
			node: ast.NewExpressionStatement(
				ast.NewAssignmentExpression(
					token.New(token.ASSIGN, "="),
					ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "foo"), "foo"),
					ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
				),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.SLTSTORE, 0),
				bytecode.New(bytecode.SLTLOAD, 0),
				bytecode.New(bytecode.POP),
			},
		},
		{
			node: ast.NewBlockStatement(
				ast.NewVariableStatement(
					token.New(token.VAR, "var"),
					ast.NewAssignmentExpression(
						token.New(token.ASSIGN, "="),
						ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "foo"), "foo"),
						ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
					),
				),
				ast.NewExpressionStatement(
					ast.NewInfixExpression(
						token.New(token.PLUS, "+"),
						ast.NewIdentifierLiteral(token.New(token.IDENTIFIER, "foo"), "foo"),
						ast.NewNumberLiteral(token.Token{Type: token.NUMBER, Literal: "1"}, 1),
					),
				),
			),
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.SLTSTORE, 0),
				bytecode.New(bytecode.SLTLOAD, 0),
				bytecode.New(bytecode.POP),
				bytecode.New(bytecode.SLTLOAD, 0),
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.I32ADD),
				bytecode.New(bytecode.POP),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.node.String(), func(t *testing.T) {
			compiler := New()

			expected := bytecode.Bytecode{}
			expected.Emit(tt.instructions...)
			for _, c := range tt.literals {
				expected.Store([]byte(c + "\x00"))
			}

			actual, err := compiler.Compile(tt.node)
			assert.NoError(t, err)
			assert.Equal(t, expected.String(), actual.String())
		})
	}
}
