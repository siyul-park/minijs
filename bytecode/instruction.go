package bytecode

import (
	"encoding/binary"
	"fmt"
	"strings"
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

	I32LOAD
	I32MUL
	I32ADD
	I32SUB
	I32DIV
	I32MOD
	I32TOF64
	I3TO2C

	F64LOAD
	F64ADD
	F64SUB
	F64MUL
	F64DIV
	F64MOD
	F64I32
	F64TOC

	CLOAD
	CADD
	CTOI32
	CTOF64
)

var types = map[Opcode]*Type{
	NOP: {Mnemonic: "nop"},
	POP: {Mnemonic: "pop"},

	I32LOAD:  {Mnemonic: "i32load", Widths: []int{4}},
	I32MUL:   {Mnemonic: "i32mul"},
	I32ADD:   {Mnemonic: "i32add"},
	I32SUB:   {Mnemonic: "i32sub"},
	I32DIV:   {Mnemonic: "i32div"},
	I32MOD:   {Mnemonic: "i32mod"},
	I32TOF64: {Mnemonic: "i32tof64"},
	I3TO2C:   {Mnemonic: "i32toc"},

	F64LOAD: {Mnemonic: "f64load", Widths: []int{8}},
	F64ADD:  {Mnemonic: "f64add"},
	F64SUB:  {Mnemonic: "f64sub"},
	F64MUL:  {Mnemonic: "f64mul"},
	F64DIV:  {Mnemonic: "f64div"},
	F64MOD:  {Mnemonic: "f64mod"},
	F64I32:  {Mnemonic: "f64i32"},
	F64TOC:  {Mnemonic: "f64toc"},

	CLOAD:  {Mnemonic: "cload", Widths: []int{4, 4}},
	CADD:   {Mnemonic: "cadd"},
	CTOI32: {Mnemonic: "ctoi32"},
	CTOF64: {Mnemonic: "ctof64"},
}

func TypeOf(op Opcode) *Type {
	typ, ok := types[op]
	if !ok {
		return nil
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
			operands[j] = uint64(i[1+offset])
		case 2:
			operands[j] = uint64(binary.BigEndian.Uint16(i[1+offset:]))
		case 4:
			operands[j] = uint64(binary.BigEndian.Uint32(i[1+offset:]))
		case 8:
			operands[j] = binary.BigEndian.Uint64(i[1+offset:])
		default:
			continue
		}
		offset += width
	}
	return operands
}

func (i Instruction) String() string {
	typ := i.Type()
	if len(typ.Widths) == 0 {
		return typ.Mnemonic
	}

	operands := i.Operands()
	widths := typ.Widths

	var ops []string
	for idx, operand := range operands {
		ops = append(ops, fmt.Sprintf("0x%0*X", widths[idx]*2, operand))
	}
	return fmt.Sprintf("%s %s", typ.Mnemonic, strings.Join(ops, " "))
}
