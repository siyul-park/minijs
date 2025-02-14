package interpreter

import (
	"fmt"
	"github.com/siyul-park/minijs/bytecode"
	"github.com/siyul-park/minijs/types"
	"math"
)

type Interpreter struct {
	bytecode bytecode.Bytecode
	stack    []types.Value
	pc       int
}

func New(code bytecode.Bytecode) *Interpreter {
	return &Interpreter{bytecode: code}
}

func (i *Interpreter) Execute() error {
	for i.pc < len(i.bytecode.Instructions) {
		instruction, read := i.bytecode.Instruction(i.pc)
		i.pc += read

		opcode := instruction.Opcode()
		operands := instruction.Operands()

		switch opcode {
		case bytecode.NOP:
		case bytecode.POP:
			if len(i.stack) > 0 {
				i.stack = i.stack[:len(i.stack)-1]
			}
		case bytecode.F64LOAD:
			i.stack = append(i.stack, types.NewFloat64(math.Float64frombits(operands[0])))
		case bytecode.F64ADD:
			val1, _ := i.stack[len(i.stack)-2].(types.Float64)
			val2, _ := i.stack[len(i.stack)-1].(types.Float64)
			i.stack = i.stack[:len(i.stack)-1]
			i.stack[len(i.stack)-1] = types.NewFloat64(val1.Value + val2.Value)
		case bytecode.F64SUB:
			val1, _ := i.stack[len(i.stack)-2].(types.Float64)
			val2, _ := i.stack[len(i.stack)-1].(types.Float64)
			i.stack = i.stack[:len(i.stack)-1]
			i.stack[len(i.stack)-1] = types.NewFloat64(val1.Value - val2.Value)
		case bytecode.F64MUL:
			val1, _ := i.stack[len(i.stack)-2].(types.Float64)
			val2, _ := i.stack[len(i.stack)-1].(types.Float64)
			i.stack = i.stack[:len(i.stack)-1]
			i.stack[len(i.stack)-1] = types.NewFloat64(val1.Value * val2.Value)
		case bytecode.F64DIV:
			val1, _ := i.stack[len(i.stack)-2].(types.Float64)
			val2, _ := i.stack[len(i.stack)-1].(types.Float64)
			i.stack = i.stack[:len(i.stack)-1]
			i.stack[len(i.stack)-1] = types.NewFloat64(val1.Value / val2.Value)
		case bytecode.F64MOD:
			val1, _ := i.stack[len(i.stack)-2].(types.Float64)
			val2, _ := i.stack[len(i.stack)-1].(types.Float64)
			i.stack = i.stack[:len(i.stack)-1]
			i.stack[len(i.stack)-1] = types.NewFloat64(math.Mod(val1.Value, val2.Value))
		default:
			return fmt.Errorf("unsupported opcode: %d", opcode)
		}
	}
	return nil
}

func (i *Interpreter) Peek(offset int) types.Value {
	if offset >= len(i.stack) {
		return nil
	}
	return i.stack[len(i.stack)-1-offset]
}
