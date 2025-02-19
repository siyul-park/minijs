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

		{instruction: New(GLBLOAD), expect: "glb.load"},

		{instruction: New(OBJSET), expect: "obj.set"},
		{instruction: New(OBJGET), expect: "obj.get"},

		{instruction: New(BOOLLOAD, 0x01), expect: "bool.load 0x01"},
		{instruction: New(BOOLTOI32), expect: "bool.to_i32"},
		{instruction: New(BOOLTOSTR), expect: "bool.to_str"},

		{instruction: New(I32LOAD, 0x01), expect: "i32.load 0x00000001"},
		{instruction: New(I32MUL), expect: "i32.mul"},
		{instruction: New(I32ADD), expect: "i32.add"},
		{instruction: New(I32SUB), expect: "i32.sub"},
		{instruction: New(I32DIV), expect: "i32.div"},
		{instruction: New(I32MOD), expect: "i32.mod"},
		{instruction: New(I32TOBOOL), expect: "i32.to_bool"},
		{instruction: New(I32TOF64), expect: "i32.to_f64"},
		{instruction: New(I32TOSTR), expect: "i32.to_str"},

		{instruction: New(F64LOAD, 0x01), expect: "f64.load 0x0000000000000001"},
		{instruction: New(F64ADD), expect: "f64.add"},
		{instruction: New(F64SUB), expect: "f64.sub"},
		{instruction: New(F64MUL), expect: "f64.mul"},
		{instruction: New(F64DIV), expect: "f64.div"},
		{instruction: New(F64MOD), expect: "f64.mod"},
		{instruction: New(F64TOI32), expect: "f64.to_i32"},
		{instruction: New(F64TOSTR), expect: "f64.to_str"},

		{instruction: New(STRLOAD, 0x01, 0x01), expect: "str.load 0x00000001 0x00000001"},
		{instruction: New(STRADD), expect: "str.add"},
		{instruction: New(STRTOI32), expect: "str.to_i32"},
		{instruction: New(STRTOF64), expect: "str.to_f64"},
	}

	for _, test := range tests {
		t.Run(test.expect, func(t *testing.T) {
			assert.Equal(t, test.expect, test.instruction.String())
		})
	}
}
