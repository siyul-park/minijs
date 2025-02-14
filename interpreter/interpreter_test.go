package interpreter

import (
	"github.com/siyul-park/minijs/bytecode"
	"github.com/siyul-park/minijs/types"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestInterpreter_Execute(t *testing.T) {
	tests := []struct {
		instructions []bytecode.Instruction
		stack        []types.Value
	}{
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
			},
			stack: []types.Value{types.NewFloat64(1)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64ADD),
			},
			stack: []types.Value{types.NewFloat64(3)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64SUB),
			},
			stack: []types.Value{types.NewFloat64(-1)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64MUL),
			},
			stack: []types.Value{types.NewFloat64(2)},
		},
		{
			instructions: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
				bytecode.New(bytecode.F64LOAD, math.Float64bits(2)),
				bytecode.New(bytecode.F64DIV),
			},
			stack: []types.Value{types.NewFloat64(0.5)},
		},
	}

	for _, tt := range tests {
		var code bytecode.Bytecode
		code.Append(tt.instructions...)

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
