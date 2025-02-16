package interpreter

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"

	"github.com/siyul-park/minijs/bytecode"
)

type Interpreter struct {
	stack  []reference
	heap   []value
	free   []uint64
	frames []*Frame
	sp     int
	fp     int
}

func New() *Interpreter {
	return &Interpreter{
		stack:  make([]reference, 64),
		heap:   make([]value, 0, 64),
		free:   make([]uint64, 0, 64),
		frames: make([]*Frame, 64),
	}
}

type reference struct {
	kind    Kind
	pointer uint64
}

type value struct {
	kind  Kind
	value Value
}

func (i *Interpreter) Pop() Value {
	return i.pop()
}

func (i *Interpreter) Execute(code bytecode.Bytecode) error {
	frame := NewFrame(code, 0)
	insns := frame.Instructions()
	consts := frame.Constants()

	i.exec(frame)
	defer i.exit()

	for frame.ip < len(insns)-1 {
		frame.ip++

		ip := frame.ip
		opcode := bytecode.Opcode(insns[ip])

		switch opcode {
		case bytecode.NOP:
		case bytecode.POP:
			i.pop()
		case bytecode.BLOAD:
			val := binary.BigEndian.Uint32(insns[frame.ip+1:])
			i.push(Bool(val))
			frame.ip += 4
		case bytecode.BTOI32:
			val, _ := i.pop().(Bool)
			i.push(Int32(val))
		case bytecode.BTOC:
			val, _ := i.pop().(Bool)
			i.push(String(val.String()))
		case bytecode.I32LOAD:
			val := Int32(binary.BigEndian.Uint32(insns[frame.ip+1:]))
			i.push(val)
			frame.ip += 4
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
		case bytecode.I32TOB:
			val, _ := i.pop().(Int32)
			if val > 0 {
				val = 1
			}
			i.push(Bool(val))
		case bytecode.I32TOF64:
			val, _ := i.pop().(Int32)
			i.push(Float64(val))
		case bytecode.I32TOC:
			val, _ := i.pop().(Int32)
			i.push(String(val.String()))
		case bytecode.F64LOAD:
			val := Float64(math.Float64frombits(binary.BigEndian.Uint64(insns[frame.ip+1:])))
			i.push(val)
			frame.ip += 8
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
		case bytecode.F64TOC:
			val, _ := i.pop().(Float64)
			i.push(String(val.String()))
		case bytecode.CLOAD:
			offset := int(binary.BigEndian.Uint32(insns[frame.ip+1:]))
			size := int(binary.BigEndian.Uint32(insns[frame.ip+5:]))
			i.push(String(consts[offset : offset+size]))
			frame.ip += 8
		case bytecode.CADD:
			val2, _ := i.pop().(String)
			val1, _ := i.pop().(String)
			i.push(val1 + val2)
		case bytecode.CTOF64:
			val, _ := i.pop().(String)
			f, err := strconv.ParseFloat(string(val), 64)
			if err != nil {
				f = math.NaN()
			}
			i.push(Float64(f))
		case bytecode.CTOI32:
			val, _ := i.pop().(String)
			n, err := strconv.Atoi(string(val))
			if err != nil {
				n = 0
			}
			i.push(Int32(n))
		default:
			typ := bytecode.TypeOf(opcode)
			if typ == nil {
				return fmt.Errorf("unknown opcode: %v", opcode)
			}
			return fmt.Errorf("unknown opcode: %v", typ.Mnemonic)
		}

		frame = i.frame()
		insns = frame.Instructions()
		consts = frame.Constants()
	}
	return nil
}

func (i *Interpreter) frame() *Frame {
	return i.frames[i.fp-1]
}

func (i *Interpreter) exec(f *Frame) {
	if len(i.frames) <= i.fp {
		frames := make([]*Frame, i.fp*2)
		copy(frames, i.frames)
		i.frames = frames
	}
	i.frames[i.fp] = f
	i.fp++
}

func (i *Interpreter) exit() {
	if i.fp == 0 {
		return
	}
	i.fp--
}

func (i *Interpreter) push(val Value) {
	if len(i.stack) <= i.sp {
		stack := make([]reference, i.sp*2)
		copy(stack, i.stack)
		i.stack = stack
	}

	kind := val.Kind()
	ref := reference{kind: kind}

	switch val := val.(type) {
	case Bool:
		ref.pointer = uint64(val)
	case Int32:
		ref.pointer = uint64(val)
	case Float64:
		ref.pointer = math.Float64bits(float64(val))
	default:
		v := value{kind: kind, value: val}

		index := -1
		if len(i.free) > 0 {
			index = int(i.free[len(i.free)-1])
			i.free = i.free[:len(i.free)-1]
		}

		if index >= 0 {
			i.heap[index] = v
		} else {
			i.heap = append(i.heap, v)
			index = len(i.heap) - 1
		}

		ref.pointer = uint64(index)
	}

	i.stack[i.sp] = ref
	i.sp++
}

func (i *Interpreter) pop() Value {
	if i.sp == 0 {
		return nil
	}

	i.sp--
	ref := i.stack[i.sp]

	switch ref.kind {
	case KindBool:
		return Bool(ref.pointer)
	case KindInt32:
		return Int32(ref.pointer)
	case KindFloat64:
		return Float64(math.Float64frombits(ref.pointer))
	default:
		v := i.heap[ref.pointer]
		i.free = append(i.free, ref.pointer)
		return v.value
	}
}
