package interpreter

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"

	"github.com/siyul-park/minijs/bytecode"
	"github.com/siyul-park/minijs/types"
)

type Interpreter struct {
	stack  []types.Value
	frames []*Frame
	sp     int
	fp     int
}

func New() *Interpreter {
	return &Interpreter{
		stack:  make([]types.Value, 64),
		frames: make([]*Frame, 64),
	}
}

func (i *Interpreter) Peek(offset int) types.Value {
	sp := i.sp - offset - 1
	if sp < 0 {
		return nil
	}
	return i.stack[sp]
}

func (i *Interpreter) Execute(code bytecode.Bytecode) error {
	frame := NewFrame(code, 0)
	insns := frame.Instructions()
	consts := frame.Constants()

	i.call(frame)

	for frame.ip < len(insns)-1 {
		frame.ip++

		ip := frame.ip
		opcode := bytecode.Opcode(insns[ip])

		switch opcode {
		case bytecode.NOP:
		case bytecode.POP:
			i.pop()
		case bytecode.F64LD:
			val := binary.BigEndian.Uint64(insns[frame.ip+1:])
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
		case bytecode.F642C:
			val, _ := i.pop().(types.Float64)
			i.push(types.NewString(strconv.FormatFloat(val.Value, 'f', -1, 64)))
		case bytecode.CLD:
			offset := int(binary.BigEndian.Uint32(insns[frame.ip+1:]))
			size := int(binary.BigEndian.Uint32(insns[frame.ip+5:]))
			i.push(types.NewString(string(consts[offset : offset+size])))
			frame.ip += 8
		case bytecode.CADD:
			val2, _ := i.pop().(types.String)
			val1, _ := i.pop().(types.String)
			i.push(types.NewString(val1.Value + val2.Value))
		case bytecode.C2F64:
			val, _ := i.pop().(types.String)
			v, err := strconv.ParseFloat(val.Value, 64)
			if err != nil {
				v = math.NaN()
			}
			i.push(types.NewFloat64(v))
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

func (i *Interpreter) call(f *Frame) {
	i.frames[i.fp] = f
	i.fp++
}

func (i *Interpreter) push(val types.Value) {
	if i.sp >= len(i.stack) {
		stack := make([]types.Value, len(i.stack)*2)
		copy(stack, i.stack)
		i.stack = stack
	}
	i.stack[i.sp] = val
	i.sp++
}

func (i *Interpreter) pop() types.Value {
	if i.sp == 0 {
		return nil
	}
	i.sp--
	return i.stack[i.sp]
}
