package bytecode

import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

type Instruction []byte

type Opcode byte

type Type struct {
	Mnemonic string
	Widths   []int
}

const (
	NOP Opcode = iota
	POP
	F64LOAD
	F64ADD
	F64SUB
	F64MUL
	F64DIV
	F64MOD
)

var types = map[Opcode]*Type{
	NOP:     {Mnemonic: "nop", Widths: []int{}},
	POP:     {Mnemonic: "pop", Widths: []int{}},
	F64LOAD: {Mnemonic: "f64load", Widths: []int{int(unsafe.Sizeof(float64(0)))}},
	F64ADD:  {Mnemonic: "f64add", Widths: []int{}},
	F64SUB:  {Mnemonic: "f64sub", Widths: []int{}},
	F64MUL:  {Mnemonic: "f64mul", Widths: []int{}},
	F64DIV:  {Mnemonic: "f64div", Widths: []int{}},
	F64MOD:  {Mnemonic: "f64mod", Widths: []int{}},
}

func TypeOf(op Opcode) *Type {
	typ, ok := types[op]
	if !ok {
		return types[NOP]
	}
	return typ
}

func (t *Type) Width() int {
	width := 1
	for _, w := range t.Widths {
		width += w
	}
	return width
}

func New(op Opcode, operands ...uint64) Instruction {
	typ, ok := types[op]
	if !ok {
		return nil
	}

	width := 1
	for _, w := range typ.Widths {
		width += w
	}

	bytecode := make(Instruction, width)
	bytecode[0] = byte(op)

	offset := 1
	for i, o := range operands {
		width := typ.Widths[i]
		switch width {
		case 1:
			bytecode[offset] = byte(o)
		case 2:
			binary.BigEndian.PutUint16(bytecode[offset:], uint16(o))
		case 4:
			binary.BigEndian.PutUint32(bytecode[offset:], uint32(o))
		case 8:
			binary.BigEndian.PutUint64(bytecode[offset:], o)
		default:
			return nil
		}
		offset += width
	}
	return bytecode
}

func (i Instruction) Type() *Type {
	return TypeOf(i.Opcode())
}

func (i Instruction) Opcode() Opcode {
	return Opcode(i[0])
}

func (i Instruction) Operands() []uint64 {
	typ := i.Type()
	operands := make([]uint64, len(typ.Widths))
	offset := 0
	for j, width := range typ.Widths {
		switch width {
		case 1:
			operands[j] = uint64(i[1])
		case 2:
			operands[j] = uint64(binary.BigEndian.Uint16(i[1:]))
		case 4:
			operands[j] = uint64(binary.BigEndian.Uint32(i[1:]))
		case 8:
			operands[j] = binary.BigEndian.Uint64(i[1:])
		default:
			continue
		}
		offset += width
	}
	return operands
}

func (i Instruction) String() string {
	typ := i.Type()
	operands := i.Operands()
	switch len(operands) {
	case 0:
		return typ.Mnemonic
	case 1:
		return fmt.Sprintf("%s 0x%X", typ.Mnemonic, operands[0])
	case 2:
		return fmt.Sprintf("%s 0x%X 0x%X", typ.Mnemonic, operands[0], operands[1])
	case 4:
		return fmt.Sprintf("%s 0x%X 0x%X 0x%X 0x%X", typ.Mnemonic, operands[0], operands[1], operands[2], operands[3])
	case 8:
		return fmt.Sprintf("%s 0x%X 0x%X 0x%X 0x%X 0x%X 0x%X 0x%X 0x%X", typ.Mnemonic, operands[0], operands[1], operands[2], operands[3], operands[4], operands[5], operands[6], operands[7])
	}
	return fmt.Sprintf("ERROR: unhandled operand width for %s: %d", typ.Mnemonic, typ.Width())
}
