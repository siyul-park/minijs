package bytecode

import (
	"fmt"
	"strings"
	"unicode"
)

type Bytecode struct {
	Instructions []byte
	Constants    []byte
}

func (b *Bytecode) Emit(instructions ...Instruction) int {
	offset := len(b.Instructions)
	for _, instruction := range instructions {
		b.Instructions = append(b.Instructions, instruction...)
	}
	return offset
}

func (b *Bytecode) Fetch(offset int) (Instruction, int) {
	if offset >= len(b.Instructions) {
		return nil, 0
	}
	typ := TypeOf(Opcode(b.Instructions[offset]))
	if typ == nil {
		return nil, 0
	}
	width := typ.Width()
	return b.Instructions[offset : offset+width], width
}

func (b *Bytecode) Store(constants []byte) int {
	offset := len(b.Constants)
	b.Constants = append(b.Constants, constants...)
	return offset
}

func (b *Bytecode) String() string {
	var out strings.Builder

	out.WriteString("section .text:\n")
	out.WriteString("\tglobal _main\n\n")
	out.WriteString("_main:\n")

	offset := 0
	for offset < len(b.Instructions) {
		bytecode, read := b.Fetch(offset)
		if read == 0 {
			break
		}
		fmt.Fprintf(&out, "\t%s\n", bytecode.String())
		offset += read
	}

	out.WriteString("\n.section .data:\n")
	for i := 0; i < len(b.Constants); i++ {
		fmt.Fprint(&out, " \t")
		for ; b.Constants[i] != 0 && i < len(b.Constants); i++ {
			if unicode.IsPrint(rune(b.Constants[i])) {
				fmt.Fprintf(&out, "%c", rune(b.Constants[i]))
			} else {
				fmt.Fprintf(&out, "0x%X", b.Constants[i])
			}
		}
		out.WriteString("\n")
	}

	return out.String()
}
