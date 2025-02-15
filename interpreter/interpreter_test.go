package interpreter

import (
	"math"
	"testing"

	"github.com/siyul-park/minijs/bytecode"
	"github.com/siyul-park/minijs/types"
	"github.com/stretchr/testify/assert"
)

func TestInterpreter_Execute(t *testing.T) {
	tests := []struct {
		instructions []bytecode.Instruction
		constants    []string
		stack        []types.Value
	}{
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
			},
			stack: []types.Value{types.Float64(1)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64ADD),
			},
			stack: []types.Value{types.Float64(3)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64SUB),
			},
			stack: []types.Value{types.Float64(-1)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64MUL),
			},
			stack: []types.Value{types.Float64(2)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64DIV),
			},
			stack: []types.Value{types.Float64(0.5)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F642C),
			},
			stack: []types.Value{types.String("1")},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.CLOAD, 0, 3),
			},
			constants: []string{"abc"},
			stack:     []types.Value{types.String("abc")},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.CLOAD, 0, 3),
				bytecode.New(bytecode.CLOAD, 0, 3),
				bytecode.New(bytecode.CADD),
			},
			constants: []string{"abc"},
			stack:     []types.Value{types.String("abcabc")},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.CLOAD, 0, 1),
				bytecode.New(bytecode.C2F64),
			},
			constants: []string{"1"},
			stack:     []types.Value{types.Float64(1)},
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

			for i, val := range tt.stack {
				assert.Equal(t, val, interpreter.Peek(i))
			}
		})
	}
}
