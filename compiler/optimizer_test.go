package compiler

import (
	"math"
	"testing"

	"github.com/siyul-park/minijs/bytecode"
	"github.com/stretchr/testify/assert"
)

func TestOptimizer_Optimize(t *testing.T) {
	tests := []struct {
		commands []bytecode.Instruction
		expected []bytecode.Instruction
		literals []string
	}{
		{
			commands: []bytecode.Instruction{
				bytecode.New(bytecode.BLOAD, 1),
				bytecode.New(bytecode.BTOI32),
			},
			expected: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
			},
		},
		{
			commands: []bytecode.Instruction{
				bytecode.New(bytecode.BLOAD, 1),
				bytecode.New(bytecode.BTOC),
			},
			expected: []bytecode.Instruction{
				bytecode.New(bytecode.CLOAD, 0, 4),
			},
		},

		{
			commands: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.I32TOB),
			},
			expected: []bytecode.Instruction{
				bytecode.New(bytecode.BLOAD, 1),
			},
		},
		{
			commands: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.I32TOF64),
			},
			expected: []bytecode.Instruction{
				bytecode.New(bytecode.F64LOAD, math.Float64bits(1)),
			},
		},
		{
			commands: []bytecode.Instruction{
				bytecode.New(bytecode.I32LOAD, 1),
				bytecode.New(bytecode.I32TOC),
			},
			expected: []bytecode.Instruction{
				bytecode.New(bytecode.CLOAD, 0, 1),
			},
		},
		{
			commands: []bytecode.Instruction{
				bytecode.New(bytecode.CLOAD, 0, 3),
				bytecode.New(bytecode.CTOF64),
				bytecode.New(bytecode.F64TOC),
			},
			literals: []string{"foo"},
			expected: []bytecode.Instruction{
				bytecode.New(bytecode.CLOAD, 0, 3),
			},
		},
	}

	optimizer := NewOptimizer()

	for _, tt := range tests {
		commands := bytecode.Bytecode{}
		commands.Add(tt.commands...)
		for _, c := range tt.literals {
			commands.Store([]byte(c + "\x00"))
		}

		expected := bytecode.Bytecode{}
		expected.Add(tt.expected...)

		t.Run(commands.String(), func(t *testing.T) {
			acturl, err := optimizer.Optimize(commands)
			assert.NoError(t, err)

			expected.Constants = acturl.Constants
			assert.Equal(t, expected.String(), acturl.String())
		})
	}
}
