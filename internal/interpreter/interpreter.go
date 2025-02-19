package interpreter

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"

	"github.com/siyul-park/minijs/internal/bytecode"
)

type Interpreter struct {
	stack  []Value
	frames []Frame
	sp     int
	fp     int
}

func New() *Interpreter {
	i := &Interpreter{
		stack:  make([]Value, 64),
		frames: make([]Frame, 64),
	}
	i.call(Frame{ip: -1})
	return i
}

func (i *Interpreter) Pop() Value {
	return i.pop()
}

func (i *Interpreter) Execute(code bytecode.Bytecode) error {
	instructions := code.Instructions
	constants := code.Constants

	i.frames[i.fp-1].ip = -1
	for i.frames[i.fp-1].ip < len(instructions)-1 {
		i.frames[i.fp-1].ip++

		ip := i.frames[i.fp-1].ip
		opcode := bytecode.Opcode(instructions[ip])

		switch opcode {
		case bytecode.NOP:
		case bytecode.POP:
			i.pop()
		case bytecode.SLTLOAD:
			idx := binary.BigEndian.Uint16(instructions[ip+1:])
			var val Value = Undefined{}
			if v, ok := i.frames[i.fp-1].Slot(int(idx)); ok {
				val = v
			}
			i.push(val)
			ip += 2
		case bytecode.SLTSTORE:
			idx := binary.BigEndian.Uint16(instructions[ip+1:])
			val := i.pop()
			i.frames[i.fp-1].SetSlot(int(idx), val)
			ip += 2
		case bytecode.UNDEFLOAD:
			i.push(Undefined{})
		case bytecode.UNDEFTOF64:
			i.pop()
			i.push(Float64(math.NaN()))
		case bytecode.UNDEFTOSTR:
			val, _ := i.pop().(Undefined)
			i.push(String(val.String()))
		case bytecode.NULLLOAD:
			i.push(Null{})
		case bytecode.NULLTOI32:
			i.pop()
			i.push(Int32(0))
		case bytecode.NULLTOSTR:
			val, _ := i.pop().(Null)
			i.push(String(val.String()))
		case bytecode.BOOLLOAD:
			val := instructions[ip+1]
			i.push(Bool(val))
			ip += 1
		case bytecode.BOOLTOI32:
			val, _ := i.pop().(Bool)
			i.push(Int32(val))
		case bytecode.BOOLTOSTR:
			val, _ := i.pop().(Bool)
			i.push(String(val.String()))
		case bytecode.I32LOAD:
			val := Int32(binary.BigEndian.Uint32(instructions[ip+1:]))
			i.push(val)
			ip += 4
		case bytecode.I32ADD:
			val2, _ := i.pop().(Int32)
			val1, _ := i.pop().(Int32)
			i.push(val1 + val2)
		case bytecode.I32SUB:
			val2, _ := i.pop().(Int32)
			val1, _ := i.pop().(Int32)
			i.push(val1 - val2)
		case bytecode.I32MUL:
			val2, _ := i.pop().(Int32)
			val1, _ := i.pop().(Int32)
			i.push(val1 * val2)
		case bytecode.I32DIV:
			val2, _ := i.pop().(Int32)
			val1, _ := i.pop().(Int32)
			i.push(val1 / val2)
		case bytecode.I32MOD:
			val2, _ := i.pop().(Int32)
			val1, _ := i.pop().(Int32)
			i.push(val1 % val2)
		case bytecode.I32TOBOOL:
			val, _ := i.pop().(Int32)
			if val > 0 {
				val = 1
			}
			i.push(Bool(val))
		case bytecode.I32TOF64:
			val, _ := i.pop().(Int32)
			i.push(Float64(val))
		case bytecode.I32TOSTR:
			val, _ := i.pop().(Int32)
			i.push(String(val.String()))
		case bytecode.F64LOAD:
			val := Float64(math.Float64frombits(binary.BigEndian.Uint64(instructions[ip+1:])))
			i.push(val)
			ip += 8
		case bytecode.F64ADD:
			val2, _ := i.pop().(Float64)
			val1, _ := i.pop().(Float64)
			i.push(val1 + val2)
		case bytecode.F64SUB:
			val2, _ := i.pop().(Float64)
			val1, _ := i.pop().(Float64)
			i.push(val1 - val2)
		case bytecode.F64MUL:
			val2, _ := i.pop().(Float64)
			val1, _ := i.pop().(Float64)
			i.push(val1 * val2)
		case bytecode.F64DIV:
			val2, _ := i.pop().(Float64)
			val1, _ := i.pop().(Float64)
			i.push(val1 / val2)
		case bytecode.F64MOD:
			val2, _ := i.pop().(Float64)
			val1, _ := i.pop().(Float64)
			i.push(Float64(math.Mod(float64(val1), float64(val2))))
		case bytecode.F64TOI32:
			val, _ := i.pop().(Float64)
			i.push(Int32(val))
		case bytecode.F64TOSTR:
			val, _ := i.pop().(Float64)
			i.push(String(val.String()))
		case bytecode.STRLOAD:
			offset := int(binary.BigEndian.Uint32(instructions[ip+1:]))
			size := int(binary.BigEndian.Uint32(instructions[ip+5:]))
			i.push(String(constants[offset : offset+size]))
			ip += 8
		case bytecode.STRADD:
			val2, _ := i.pop().(String)
			val1, _ := i.pop().(String)
			i.push(val1 + val2)
		case bytecode.STRTOI32:
			val, _ := i.pop().(String)
			n, err := strconv.Atoi(string(val))
			if err != nil {
				n = 0
			}
			i.push(Int32(n))
		case bytecode.STRTOF64:
			val, _ := i.pop().(String)
			f, err := strconv.ParseFloat(string(val), 64)
			if err != nil {
				f = math.NaN()
			}
			i.push(Float64(f))
		default:
			typ := bytecode.TypeOf(opcode)
			if typ == nil {
				return fmt.Errorf("unknown opcode: %v", opcode)
			}
			return fmt.Errorf("unknown opcode: %v", typ.Mnemonic)
		}

		i.frames[i.fp-1].ip = ip
	}
	return nil
}

func (i *Interpreter) call(frame Frame) {
	if len(i.frames) <= i.fp {
		i.frames = append(i.frames, make([]Frame, len(i.frames)+1)...)
	}
	i.frames[i.fp] = frame
	i.fp++
}

func (i *Interpreter) exit() Frame {
	if i.fp == 0 {
		return Frame{}
	}
	i.fp--
	i.frames[i.fp] = Frame{}
	return i.frames[i.fp]
}

func (i *Interpreter) push(val Value) {
	if len(i.stack) <= i.sp {
		i.stack = append(i.stack, make([]Value, len(i.stack)+1)...)
	}
	i.stack[i.sp] = val
	i.sp++
}

func (i *Interpreter) pop() Value {
	if i.sp == 0 {
		return nil
	}
	i.sp--
	return i.stack[i.sp]
}
