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

		{instruction: New(GLBLOAD), expect: "glbload"},

		{instruction: New(OBJSET), expect: "objset"},
		{instruction: New(OBJGET), expect: "objget"},

		{instruction: New(BOOLLOAD, 0x01), expect: "boolload 0x01"},
		{instruction: New(BOOLTOI32), expect: "booltoi32"},
		{instruction: New(BOOLTOSTR), expect: "booltostr"},

		{instruction: New(I32LOAD, 0x01), expect: "i32load 0x00000001"},
		{instruction: New(I32MUL), expect: "i32mul"},
		{instruction: New(I32ADD), expect: "i32add"},
		{instruction: New(I32SUB), expect: "i32sub"},
		{instruction: New(I32DIV), expect: "i32div"},
		{instruction: New(I32MOD), expect: "i32mod"},
		{instruction: New(I32TOBOOL), expect: "i32tobool"},
		{instruction: New(I32TOF64), expect: "i32tof64"},
		{instruction: New(I32TOSTR), expect: "i32tostr"},

		{instruction: New(F64LOAD, 0x01), expect: "f64load 0x0000000000000001"},
		{instruction: New(F64ADD), expect: "f64add"},
		{instruction: New(F64SUB), expect: "f64sub"},
		{instruction: New(F64MUL), expect: "f64mul"},
		{instruction: New(F64DIV), expect: "f64div"},
		{instruction: New(F64MOD), expect: "f64mod"},
		{instruction: New(F64TOI32), expect: "f64toi32"},
		{instruction: New(F64TOSTR), expect: "f64tostr"},

		{instruction: New(STRLOAD, 0x01, 0x01), expect: "strload 0x00000001 0x00000001"},
		{instruction: New(STRADD), expect: "stradd"},
		{instruction: New(STRTOI32), expect: "strtoi32"},
		{instruction: New(STRTOF64), expect: "strtof64"},
	}

	for _, test := range tests {
		t.Run(test.expect, func(t *testing.T) {
			assert.Equal(t, test.expect, test.instruction.String())
		})
	}
}
