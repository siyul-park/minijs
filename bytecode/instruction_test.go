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

		{instruction: New(BLOAD, 0x01), expect: "bload 0x00000001"},
		{instruction: New(BTOI32), expect: "btoi32"},
		{instruction: New(BTOC), expect: "btoc"},

		{instruction: New(I32LOAD, 0x01), expect: "i32load 0x00000001"},
		{instruction: New(I32MUL), expect: "i32mul"},
		{instruction: New(I32ADD), expect: "i32add"},
		{instruction: New(I32SUB), expect: "i32sub"},
		{instruction: New(I32DIV), expect: "i32div"},
		{instruction: New(I32MOD), expect: "i32mod"},
		{instruction: New(I32TOB), expect: "i32tob"},
		{instruction: New(I32TOF64), expect: "i32tof64"},
		{instruction: New(I3TO2C), expect: "i32toc"},

		{instruction: New(F64LOAD, 0x01), expect: "f64load 0x0000000000000001"},
		{instruction: New(F64ADD), expect: "f64add"},
		{instruction: New(F64SUB), expect: "f64sub"},
		{instruction: New(F64MUL), expect: "f64mul"},
		{instruction: New(F64DIV), expect: "f64div"},
		{instruction: New(F64MOD), expect: "f64mod"},
		{instruction: New(F64I32), expect: "f64i32"},
		{instruction: New(F64TOC), expect: "f64toc"},

		{instruction: New(CLOAD, 0x01, 0x01), expect: "cload 0x00000001 0x00000001"},
		{instruction: New(CADD), expect: "cadd"},
		{instruction: New(CTOI32), expect: "ctoi32"},
		{instruction: New(CTOF64), expect: "ctof64"},
	}

	for _, test := range tests {
		t.Run(test.expect, func(t *testing.T) {
			assert.Equal(t, test.expect, test.instruction.String())
		})
	}
}
