package bytecode

import (
	"fmt"
	"strings"
)

type Bytecode struct {
	Instructions []byte
	Constants    []byte
}

func (b *Bytecode) Append(instructions ...Instruction) {
	for _, instruction := range instructions {
		b.Instructions = append(b.Instructions, instruction...)
	}
}

func (b *Bytecode) Instruction(offset int) (Instruction, int) {
	if offset >= len(b.Instructions) {
		return nil, 0
	}
	typ := TypeOf(Opcode(b.Instructions[offset]))
	width := typ.Width()
	return b.Instructions[offset : offset+width], width
}

func (b *Bytecode) String() string {
	var out strings.Builder
	offset := 0
	for offset < len(b.Instructions) {
		bytecode, read := b.Instruction(offset)
		_, _ = fmt.Fprintf(&out, "%04d %s\n", offset, bytecode.String())
		offset += read
	}
	return out.String()
}
