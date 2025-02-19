package interpreter

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"

	"github.com/siyul-park/minijs/internal/bytecode"
)

type Interpreter struct {
	global *Object
	stack  []Value
	sp     int
}

func New() *Interpreter {
	return &Interpreter{
		global: NewObject(nil),
		stack:  make([]Value, 64),
	}
}

func (i *Interpreter) Pop() Value {
	return i.pop()
}

func (i *Interpreter) Execute(code bytecode.Bytecode) error {
	instructions := code.Instructions
	constants := code.Constants

	ip := -1
	for ip < len(instructions)-1 {
		ip++

		opcode := bytecode.Opcode(instructions[ip])

		switch opcode {
		case bytecode.NOP:
		case bytecode.POP:
			i.pop()
		case bytecode.GLBLOAD:
			i.push(i.global)
		case bytecode.OBJSET:
			val3 := i.pop()
			val2 := i.pop()
			val1 := i.pop().(*Object)
			val1.Set(val2, val3)
			i.push(val3)
		case bytecode.OBJGET:
			val2 := i.pop()
			val1, _ := i.pop().(*Object)
			val3, ok := val1.Get(val2)
			if !ok {
				val3 = Undefined{}
			}
			i.push(val3)
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

		instructions = code.Instructions
		constants = code.Constants
	}
	return nil
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
