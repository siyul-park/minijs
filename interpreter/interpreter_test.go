package interpreter

import (
	"math"
	"testing"

	"github.com/siyul-park/minijs/bytecode"
	"github.com/stretchr/testify/assert"
)

func TestInterpreter_Execute(t *testing.T) {
	tests := []struct {
		instructions []bytecode.Instruction
		constants    []string
		stack        []any
	}{
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
			},
			stack: []any{int32(1)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.I32LOAD, 2),
				bytecode.New(bytecode.I32MUL),
			},
			stack: []any{int32(2)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 6),
				bytecode.New(bytecode.I32LOAD, 2),
				bytecode.New(bytecode.I32DIV),
			},
			stack: []any{int32(3)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 7),
				bytecode.New(bytecode.I32LOAD, 3),
				bytecode.New(bytecode.I32MOD),
			},
			stack: []any{int32(1)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 5),
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.I32DIV),
			},
			stack: []any{int32(5)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 5),
				bytecode.New(bytecode.I322F64),
			},
			stack: []any{float64(5)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 42),
				bytecode.New(bytecode.I322C),
			},
			stack: []any{"42"},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
			},
			stack: []any{float64(1)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64ADD),
			},
			stack: []any{float64(3)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64SUB),
			},
			stack: []any{float64(-1)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64MUL),
			},
			stack: []any{float64(2)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64DIV),
			},
			stack: []any{0.5},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(3.7)),
				bytecode.New(bytecode.F64I32),
			},
			stack: []any{int32(3)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F642C),
			},
			stack: []any{"1"},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.CLOAD, 0, 3),
			},
			constants: []string{"abc"},
			stack:     []any{"abc"},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.CLOAD, 0, 3),
				bytecode.New(bytecode.CLOAD, 0, 3),
				bytecode.New(bytecode.CADD),
			},
			constants: []string{"abc"},
			stack:     []any{"abcabc"},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.CLOAD, 0, 3),
				bytecode.New(bytecode.C2I32),
			},
			constants: []string{"123"},
			stack:     []any{int32(123)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.CLOAD, 0, 1),
				bytecode.New(bytecode.C2F64),
			},
			constants: []string{"1"},
			stack:     []any{float64(1)},
		},
	}

	for _, tt := range tests {
		var code bytecode.Bytecode
		code.Add(tt.instructions...)
		for _, c := range tt.constants {
			code.Store([]byte(c + "\x00"))
		}

		t.Run(code.String(), func(t *testing.T) {
			interpreter := New()

			err := interpreter.Execute(code)
			assert.NoError(t, err)

			for _, val := range tt.stack {
				assert.Equal(t, val, interpreter.Top())
			}
		})
	}
}

func BenchmarkInterpreter_Execute(b *testing.B) {
	tests := []struct {
		instructions []bytecode.Instruction
		constants    []string
	}{
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.I32LOAD, 2),
				bytecode.New(bytecode.I32MUL),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 6),
				bytecode.New(bytecode.I32LOAD, 2),
				bytecode.New(bytecode.I32DIV),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 7),
				bytecode.New(bytecode.I32LOAD, 3),
				bytecode.New(bytecode.I32MOD),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 5),
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.I32DIV),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 5),
				bytecode.New(bytecode.I322F64),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 42),
				bytecode.New(bytecode.I322C),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64ADD),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64SUB),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64MUL),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64DIV),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(3.7)),
				bytecode.New(bytecode.F64I32),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F642C),
			},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.CLOAD, 0, 3),
			},
			constants: []string{"abc"},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.CLOAD, 0, 3),
				bytecode.New(bytecode.CLOAD, 0, 3),
				bytecode.New(bytecode.CADD),
			},
			constants: []string{"abc"},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.CLOAD, 0, 3),
				bytecode.New(bytecode.C2I32),
			},
			constants: []string{"123"},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.CLOAD, 0, 1),
				bytecode.New(bytecode.C2F64),
			},
			constants: []string{"1"},
		},
	}

	for _, tt := range tests {
		var code bytecode.Bytecode
		code.Add(tt.instructions...)
		for _, c := range tt.constants {
			code.Store([]byte(c + "\x00"))
		}

		b.Run(code.String(), func(b *testing.B) {
			interpreter := New()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				err := interpreter.Execute(code)
				assert.NoError(b, err)
			}
		})
	}
}
