package compiler

import (
	"encoding/binary"
	"math"

	"github.com/siyul-park/minijs/bytecode"
	"github.com/siyul-park/minijs/interpreter"
)

type Optimizer struct {
	interpreter *interpreter.Interpreter
}

func NewOptimizer() *Optimizer {
	return &Optimizer{
		interpreter: interpreter.New(),
	}
}

func (o *Optimizer) Optimize(code bytecode.Bytecode) (bytecode.Bytecode, error) {
	constants := code.Constants

	var instructions []bytecode.Instruction
	for offset := 0; offset < len(code.Instructions); {
		inst, size := code.Instruction(offset)
		instructions = append(instructions, inst)
		offset += size
	}

	instructions, constants, err := o.fusion(instructions, constants)
	if err != nil {
		return bytecode.Bytecode{}, err
	}

	instructions, constants = o.compress(instructions, constants)

	code.Instructions = nil
	code.Constants = constants
	code.Add(instructions...)
	return code, nil
}

func (o *Optimizer) fusion(instructions []bytecode.Instruction, constants []byte) ([]bytecode.Instruction, []byte, error) {
	offsets := map[string]int{}
	for i := 0; i < len(instructions); i++ {
		inst := instructions[i]
		if inst.Opcode() == bytecode.CLOAD {
			offset := int(binary.BigEndian.Uint32(inst[1:]))
			size := int(binary.BigEndian.Uint32(inst[5:]))

			literal := string(constants[offset : offset+size])
			offsets[literal] = offset
		}
	}

	for i := 0; i < len(instructions); i++ {
		inst := instructions[i]
		if i > 0 {
			operand := instructions[i-1]
			switch operand.Opcode() {
			case bytecode.BLOAD, bytecode.I32LOAD, bytecode.F64LOAD, bytecode.CLOAD:
				switch inst.Opcode() {
				case bytecode.I32TOB:
					code := bytecode.Bytecode{Constants: constants}
					code.Add(operand, inst)
					if err := o.interpreter.Execute(code); err != nil {
						return nil, nil, err
					}

					val, _ := o.interpreter.Pop().(interpreter.Bool)
					instructions[i-1] = bytecode.New(bytecode.NOP)
					instructions[i] = bytecode.New(bytecode.BLOAD, uint64(val))
				case bytecode.BTOI32, bytecode.F64TOI32, bytecode.CTOI32:
					code := bytecode.Bytecode{Constants: constants}
					code.Add(operand, inst)
					if err := o.interpreter.Execute(code); err != nil {
						return nil, nil, err
					}

					val, _ := o.interpreter.Pop().(interpreter.Int32)

					instructions[i-1] = bytecode.New(bytecode.NOP)
					instructions[i] = bytecode.New(bytecode.I32LOAD, uint64(val))
				case bytecode.I32TOF64, bytecode.CTOF64:
					code := bytecode.Bytecode{Constants: constants}
					code.Add(operand, inst)
					if err := o.interpreter.Execute(code); err != nil {
						return nil, nil, err
					}

					val, _ := o.interpreter.Pop().(interpreter.Float64)
					instructions[i-1] = bytecode.New(bytecode.NOP)
					instructions[i] = bytecode.New(bytecode.F64LOAD, math.Float64bits(float64(val)))
				case bytecode.BTOC, bytecode.I32TOC, bytecode.F64TOC:
					code := bytecode.Bytecode{Constants: constants}
					code.Add(operand, inst)
					if err := o.interpreter.Execute(code); err != nil {
						return nil, nil, err
					}

					val, _ := o.interpreter.Pop().(interpreter.String)

					offset, ok := offsets[string(val)]
					if !ok {
						offset = len(constants)
						constants = append(constants, []byte(string(val)+"\x00")...)
						offsets[string(val)] = offset
					}

					instructions[i-1] = bytecode.New(bytecode.NOP)
					instructions[i] = bytecode.New(bytecode.CLOAD, uint64(offset), uint64(len(val)))
				default:
				}
			default:
			}
		}
	}
	return instructions, constants, nil
}

func (o *Optimizer) compress(instructions []bytecode.Instruction, constants []byte) ([]bytecode.Instruction, []byte) {
	offsets := map[string]int{}
	for i := 0; i < len(instructions); i++ {
		inst := instructions[i]
		if inst.Opcode() == bytecode.CLOAD {
			offset := int(binary.BigEndian.Uint32(inst[1:]))
			size := int(binary.BigEndian.Uint32(inst[5:]))

			literal := string(constants[offset : offset+size])
			offsets[literal] = offset
		}
	}

	compressed := make([]byte, 0, len(constants))
	for literal := range offsets {
		offsets[literal] = len(compressed)
		compressed = append(compressed, []byte(literal+"\x00")...)
	}

	for i := 0; i < len(instructions); i++ {
		inst := instructions[i]
		if inst.Opcode() == bytecode.CLOAD {
			offset := int(binary.BigEndian.Uint32(inst[1:]))
			size := int(binary.BigEndian.Uint32(inst[5:]))
			instructions[i] = bytecode.New(bytecode.CLOAD, uint64(offsets[string(constants[offset:offset+size])]), uint64(size))
		}
	}

	for i := len(instructions) - 1; i >= 0; i-- {
		if instructions[i].Opcode() == bytecode.NOP {
			instructions = append(instructions[:i], instructions[i+1:]...)
		}
	}

	return instructions, compressed
}
