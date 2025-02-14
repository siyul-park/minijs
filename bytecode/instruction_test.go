package bytecode

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInstruction_String(t *testing.T) {
	tests := []struct {
		instruction Instruction
		expect      string
	}{
		{instruction: New(NOP), expect: "nop"},
		{instruction: New(F64LOAD, 0x01), expect: "f64load 0x1"},
		{instruction: New(F64ADD), expect: "f64add"},
		{instruction: New(F64SUB), expect: "f64sub"},
		{instruction: New(F64MUL), expect: "f64mul"},
		{instruction: New(F64DIV), expect: "f64div"},
		{instruction: New(F64MOD), expect: "f64mod"},
	}

	for _, test := range tests {
		t.Run(test.expect, func(t *testing.T) {
			assert.Equal(t, test.expect, test.instruction.String())
		})
	}
}
