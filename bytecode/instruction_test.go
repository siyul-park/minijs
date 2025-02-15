package bytecode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstruction_String(t *testing.T) {
	tests := []struct {
		instruction Instruction
		expect      string
	}{
		{instruction: New(NOP), expect: "nop"},
		{instruction: New(POP), expect: "pop"},

		{instruction: New(F64LD, 0x01), expect: "f64ld 0x1"},
		{instruction: New(F64ADD), expect: "f64add"},
		{instruction: New(F64SUB), expect: "f64sub"},
		{instruction: New(F64MUL), expect: "f64mul"},
		{instruction: New(F64DIV), expect: "f64div"},
		{instruction: New(F64MOD), expect: "f64mod"},

		{instruction: New(CLD, 0x01, 0x01), expect: "cld 0x1 0x1"},
		{instruction: New(CADD), expect: "cadd"},
		{instruction: New(C2F64), expect: "c2f64"},
	}

	for _, test := range tests {
		t.Run(test.expect, func(t *testing.T) {
			assert.Equal(t, test.expect, test.instruction.String())
		})
	}
}
