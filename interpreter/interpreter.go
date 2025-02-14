package interpreter

import (
	"encoding/binary"
	"fmt"
	"github.com/siyul-park/minijs/bytecode"
	"github.com/siyul-park/minijs/types"
	"math"
)

type Interpreter struct {
	stack     []types.Value
	frames    []*Frame
	constants []byte
	sp        int
	fp        int
}

func New(code bytecode.Bytecode) *Interpreter {
	i := &Interpreter{
		stack:     make([]types.Value, 512),
		frames:    make([]*Frame, 64),
		constants: code.Constants,
	}
	i.call(NewFrame(code, 0))
	return i
}

func (i *Interpreter) Peek(offset int) types.Value {
	sp := i.sp - offset - 1
	if sp < 0 {
		return nil
	}
	return i.stack[sp]
}

func (i *Interpreter) Execute() error {
	frame := i.frame()
	insns := frame.Instructions()

	for frame.ip < len(insns)-1 {
		frame.ip++

		ip := frame.ip
		opcode := bytecode.Opcode(insns[ip])

		switch opcode {
		case bytecode.NOP:
		case bytecode.POP:
			i.pop()
		case bytecode.F64LOAD:
			val := binary.BigEndian.Uint64(insns[frame.ip+1 : frame.ip+9])
			i.push(types.NewFloat64(math.Float64frombits(val)))
			frame.ip += 8
		case bytecode.F64ADD:
			val2, _ := i.pop().(types.Float64)
			val1, _ := i.pop().(types.Float64)
			i.push(types.NewFloat64(val1.Value + val2.Value))
		case bytecode.F64SUB:
			val2, _ := i.pop().(types.Float64)
			val1, _ := i.pop().(types.Float64)
			i.push(types.NewFloat64(val1.Value - val2.Value))
		case bytecode.F64MUL:
			val2, _ := i.pop().(types.Float64)
			val1, _ := i.pop().(types.Float64)
			i.push(types.NewFloat64(val1.Value * val2.Value))
		case bytecode.F64DIV:
			val2, _ := i.pop().(types.Float64)
			val1, _ := i.pop().(types.Float64)
			i.push(types.NewFloat64(val1.Value / val2.Value))
		case bytecode.F64MOD:
			val2, _ := i.pop().(types.Float64)
			val1, _ := i.pop().(types.Float64)
			i.push(types.NewFloat64(math.Mod(val1.Value, val2.Value)))
		default:
			return fmt.Errorf("unknown opcode: %v", opcode)
		}

		frame = i.frame()
		insns = frame.Instructions()
	}

	return nil
}

func (i *Interpreter) frame() *Frame {
	return i.frames[i.fp-1]
}

func (i *Interpreter) ret() *Frame {
	i.fp--
	return i.frames[i.fp]
}

func (i *Interpreter) call(f *Frame) {
	i.frames[i.fp] = f
	i.fp++
}

func (i *Interpreter) push(val types.Value) {
	if i.sp >= len(i.stack) {
		i.stack = append(i.stack, val)
	} else {
		i.stack[i.sp] = val
	}
	i.sp++
}

func (i *Interpreter) pop() types.Value {
	if i.sp == 0 {
		return nil
	}
	i.sp--
	return i.stack[i.sp]
}
