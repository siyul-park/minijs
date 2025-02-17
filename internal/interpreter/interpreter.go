package interpreter

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"

	"github.com/siyul-park/minijs/internal/bytecode"
)

type Interpreter struct {
	stack  []mark
	heap   []Value
	free   []uint64
	frames []*Frame
	sp     int
	fp     int
}

type mark struct {
	kind    Kind
	pointer uint64
}

func New() *Interpreter {
	return &Interpreter{
		stack:  make([]mark, 64),
		heap:   make([]Value, 0, 64),
		free:   make([]uint64, 0, 64),
		frames: make([]*Frame, 64),
	}
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
			i.bpush(Bool(val))
			frame.ip += 4
		case bytecode.BTOI32:
			val := i.bpop()
			i.i32push(Int32(val))
		case bytecode.BTOC:
			val := i.bpop()
			i.cpush(String(val.String()))
		case bytecode.I32LOAD:
			val := Int32(binary.BigEndian.Uint32(insns[frame.ip+1:]))
			i.i32push(val)
			frame.ip += 4
		case bytecode.I32ADD:
			val2 := i.i32pop()
			val1 := i.i32pop()
			i.i32push(val1 + val2)
		case bytecode.I32SUB:
			val2 := i.i32pop()
			val1 := i.i32pop()
			i.i32push(val1 - val2)
		case bytecode.I32MUL:
			val2 := i.i32pop()
			val1 := i.i32pop()
			i.i32push(val1 * val2)
		case bytecode.I32DIV:
			val2 := i.i32pop()
			val1 := i.i32pop()
			i.i32push(val1 / val2)
		case bytecode.I32MOD:
			val2 := i.i32pop()
			val1 := i.i32pop()
			i.i32push(val1 % val2)
		case bytecode.I32TOB:
			val := i.i32pop()
			if val > 0 {
				val = 1
			}
			i.bpush(Bool(val))
		case bytecode.I32TOF64:
			val := i.i32pop()
			i.f64push(Float64(val))
		case bytecode.I32TOC:
			val := i.i32pop()
			i.cpush(String(val.String()))
		case bytecode.F64LOAD:
			val := Float64(math.Float64frombits(binary.BigEndian.Uint64(insns[frame.ip+1:])))
			i.f64push(val)
			frame.ip += 8
		case bytecode.F64ADD:
			val2 := i.f64pop()
			val1 := i.f64pop()
			i.f64push(val1 + val2)
		case bytecode.F64SUB:
			val2 := i.f64pop()
			val1 := i.f64pop()
			i.f64push(val1 - val2)
		case bytecode.F64MUL:
			val2 := i.f64pop()
			val1 := i.f64pop()
			i.f64push(val1 * val2)
		case bytecode.F64DIV:
			val2 := i.f64pop()
			val1 := i.f64pop()
			i.f64push(val1 / val2)
		case bytecode.F64MOD:
			val2 := i.f64pop()
			val1 := i.f64pop()
			i.f64push(Float64(math.Mod(float64(val1), float64(val2))))
		case bytecode.F64TOI32:
			val := i.f64pop()
			i.i32push(Int32(val))
		case bytecode.F64TOC:
			val := i.f64pop()
			i.cpush(String(val.String()))
		case bytecode.CLOAD:
			offset := int(binary.BigEndian.Uint32(insns[frame.ip+1:]))
			size := int(binary.BigEndian.Uint32(insns[frame.ip+5:]))
			i.cpush(String(consts[offset : offset+size]))
			frame.ip += 8
		case bytecode.CADD:
			val2 := i.cpop()
			val1 := i.cpop()
			i.cpush(val1 + val2)
		case bytecode.CTOI32:
			val := i.cpop()
			n, err := strconv.Atoi(string(val))
			if err != nil {
				n = 0
			}
			i.i32push(Int32(n))
		case bytecode.CTOF64:
			val := i.cpop()
			f, err := strconv.ParseFloat(string(val), 64)
			if err != nil {
				f = math.NaN()
			}
			i.f64push(Float64(f))
		default:
			typ := bytecode.TypeOf(opcode)
			if typ == nil {
				return fmt.Errorf("unknown opcode: %v", opcode)
			}
			return fmt.Errorf("unknown opcode: %v", typ.Mnemonic)
		}
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

func (i *Interpreter) resize() {
	if len(i.stack) <= i.sp {
		stack := make([]mark, i.sp*2)
		copy(stack, i.stack)
		i.stack = stack
	}
}

func (i *Interpreter) pop() Value {
	if i.sp == 0 {
		return nil
	}
	m := i.stack[i.sp-1]
	switch m.kind {
	case KindBool:
		return i.bpop()
	case KindInt32:
		return i.i32pop()
	case KindFloat64:
		return i.f64pop()
	default:
		return i.hpop()
	}
}

func (i *Interpreter) bpush(val Bool) {
	i.resize()
	i.stack[i.sp] = mark{
		kind:    KindBool,
		pointer: uint64(val),
	}
	i.sp++
}

func (i *Interpreter) bpop() Bool {
	if i.sp == 0 {
		return 0
	}
	i.sp--
	m := i.stack[i.sp]
	return Bool(m.pointer)
}

func (i *Interpreter) i32push(val Int32) {
	i.resize()
	i.stack[i.sp] = mark{
		kind:    KindInt32,
		pointer: uint64(val),
	}
	i.sp++
}

func (i *Interpreter) i32pop() Int32 {
	if i.sp == 0 {
		return 0
	}
	i.sp--
	m := i.stack[i.sp]
	return Int32(m.pointer)
}

func (i *Interpreter) f64push(val Float64) {
	i.resize()
	i.stack[i.sp] = mark{
		kind:    KindFloat64,
		pointer: math.Float64bits(float64(val)),
	}
	i.sp++
}

func (i *Interpreter) f64pop() Float64 {
	if i.sp == 0 {
		return 0.0
	}
	i.sp--
	m := i.stack[i.sp]
	return Float64(math.Float64frombits(m.pointer))
}

func (i *Interpreter) cpush(val String) {
	i.hpush(val)
}

func (i *Interpreter) cpop() String {
	val, _ := i.hpop().(String)
	return val
}

func (i *Interpreter) hpush(val Value) {
	i.resize()

	index := -1
	if len(i.free) > 0 {
		index = int(i.free[len(i.free)-1])
		i.free = i.free[:len(i.free)-1]
	}

	if index >= 0 {
		i.heap[index] = val
	} else {
		i.heap = append(i.heap, val)
		index = len(i.heap) - 1
	}

	i.stack[i.sp] = mark{
		kind:    val.Kind(),
		pointer: uint64(index),
	}
	i.sp++
}

func (i *Interpreter) hpop() Value {
	if i.sp == 0 {
		return nil
	}

	i.sp--
	m := i.stack[i.sp]
	v := i.heap[m.pointer]
	i.free = append(i.free, m.pointer)
	return v
}
