package interpreter

import (
	"encoding/binary"
	"math"

	"github.com/siyul-park/minijs/internal/bytecode"
)

type Optimizer struct {
	interpreter *Interpreter
}

func NewOptimizer() *Optimizer {
	return &Optimizer{
		interpreter: New(),
	}
}

func (o *Optimizer) Optimize(code bytecode.Bytecode) (bytecode.Bytecode, error) {
	constants := code.Constants

	var instructions []bytecode.Instruction
	for offset := 0; offset < len(code.Instructions); {
		inst, size := code.Fetch(offset)
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
	code.Emit(instructions...)
	return code, nil
}

func (o *Optimizer) fusion(instructions []bytecode.Instruction, constants []byte) ([]bytecode.Instruction, []byte, error) {
	literals := map[string]int{}
	for i := 0; i < len(instructions); i++ {
		inst := instructions[i]
		if inst.Opcode() == bytecode.STRLOAD {
			offset := int(binary.BigEndian.Uint32(inst[1:]))
			size := int(binary.BigEndian.Uint32(inst[5:]))

			literal := string(constants[offset : offset+size])
			literals[literal] = offset
		}
	}

	for i := 0; i < len(instructions); i++ {
		inst := instructions[i]
		if i > 0 {
			j := i - 1
			for ; j > 0; j-- {
				if instructions[j].Opcode() != bytecode.NOP {
					break
				}
			}

			operand := instructions[j]
			switch operand.Opcode() {
			case bytecode.BOOLLOAD, bytecode.I32LOAD, bytecode.F64LOAD, bytecode.STRLOAD:
				switch inst.Opcode() {
				case bytecode.I32TOBOOL:
					code := bytecode.Bytecode{Constants: constants}
					code.Emit(operand, inst)
					if err := o.interpreter.Execute(code); err != nil {
						return nil, nil, err
					}

					val, _ := o.interpreter.Pop().(Bool)

					instructions[j] = bytecode.New(bytecode.NOP)
					instructions[i] = bytecode.New(bytecode.BOOLLOAD, uint64(val))
				case bytecode.BOOLTOI32, bytecode.F64TOI32, bytecode.STRTOI32:
					code := bytecode.Bytecode{Constants: constants}
					code.Emit(operand, inst)
					if err := o.interpreter.Execute(code); err != nil {
						return nil, nil, err
					}

					val, _ := o.interpreter.Pop().(Int32)

					instructions[j] = bytecode.New(bytecode.NOP)
					instructions[i] = bytecode.New(bytecode.I32LOAD, uint64(val))
				case bytecode.I32TOF64, bytecode.STRTOF64:
					code := bytecode.Bytecode{Constants: constants}
					code.Emit(operand, inst)
					if err := o.interpreter.Execute(code); err != nil {
						return nil, nil, err
					}

					val, _ := o.interpreter.Pop().(Float64)

					instructions[j] = bytecode.New(bytecode.NOP)
					instructions[i] = bytecode.New(bytecode.F64LOAD, math.Float64bits(float64(val)))
				case bytecode.BOOLTOSTR, bytecode.I32TOSTR, bytecode.F64TOSTR:
					code := bytecode.Bytecode{Constants: constants}
					code.Emit(operand, inst)
					if err := o.interpreter.Execute(code); err != nil {
						return nil, nil, err
					}

					val, _ := o.interpreter.Pop().(String)

					offset, ok := literals[string(val)]
					if !ok {
						offset = len(constants)
						constants = append(constants, []byte(string(val)+"\x00")...)
						literals[string(val)] = offset
					}

					instructions[j] = bytecode.New(bytecode.NOP)
					instructions[i] = bytecode.New(bytecode.STRLOAD, uint64(offset), uint64(len(val)))
				default:
				}
			default:
			}
		}

		if i > 1 {
			j := i - 1
			for ; j > 0; j-- {
				if instructions[j].Opcode() != bytecode.NOP {
					break
				}
			}

			k := j - 1
			for ; k > 0; k-- {
				if instructions[k].Opcode() != bytecode.NOP {
					break
				}
			}

			operand1 := instructions[j]
			operand2 := instructions[k]
			if operand1.Opcode() == operand2.Opcode() {
				switch operand1.Opcode() {
				case bytecode.BOOLLOAD, bytecode.I32LOAD, bytecode.F64LOAD, bytecode.STRLOAD:
					switch inst.Opcode() {
					case bytecode.I32ADD, bytecode.I32SUB, bytecode.I32MUL, bytecode.I32DIV, bytecode.I32MOD:
						code := bytecode.Bytecode{Constants: constants}
						code.Emit(operand2, operand1, inst)
						if err := o.interpreter.Execute(code); err != nil {
							return nil, nil, err
						}

						val, _ := o.interpreter.Pop().(Int32)

						instructions[k] = bytecode.New(bytecode.NOP)
						instructions[j] = bytecode.New(bytecode.NOP)
						instructions[i] = bytecode.New(bytecode.I32LOAD, uint64(val))
					case bytecode.F64ADD, bytecode.F64SUB, bytecode.F64MUL, bytecode.F64DIV, bytecode.F64MOD:
						code := bytecode.Bytecode{Constants: constants}
						code.Emit(operand2, operand1, inst)
						if err := o.interpreter.Execute(code); err != nil {
							return nil, nil, err
						}

						val, _ := o.interpreter.Pop().(Float64)

						instructions[k] = bytecode.New(bytecode.NOP)
						instructions[j] = bytecode.New(bytecode.NOP)
						instructions[i] = bytecode.New(bytecode.F64LOAD, math.Float64bits(float64(val)))
					case bytecode.STRADD:
						code := bytecode.Bytecode{Constants: constants}
						code.Emit(operand2, operand1, inst)
						if err := o.interpreter.Execute(code); err != nil {
							return nil, nil, err
						}

						val, _ := o.interpreter.Pop().(String)

						offset, ok := literals[string(val)]
						if !ok {
							offset = len(constants)
							constants = append(constants, []byte(string(val)+"\x00")...)
							literals[string(val)] = offset
						}

						instructions[k] = bytecode.New(bytecode.NOP)
						instructions[j] = bytecode.New(bytecode.NOP)
						instructions[i] = bytecode.New(bytecode.STRLOAD, uint64(offset), uint64(len(val)))
					default:
					}
				default:
				}
			}
		}
	}
	return instructions, constants, nil
}

func (o *Optimizer) compress(instructions []bytecode.Instruction, constants []byte) ([]bytecode.Instruction, []byte) {
	literals := map[string]int{}
	for i := 0; i < len(instructions); i++ {
		inst := instructions[i]
		if inst.Opcode() == bytecode.STRLOAD {
			offset := int(binary.BigEndian.Uint32(inst[1:]))
			size := int(binary.BigEndian.Uint32(inst[5:]))

			literal := string(constants[offset : offset+size])
			literals[literal] = offset
		}
	}

	compressed := make([]byte, 0, len(constants))
	for literal := range literals {
		literals[literal] = len(compressed)
		compressed = append(compressed, []byte(literal+"\x00")...)
	}

	for i := 0; i < len(instructions); i++ {
		inst := instructions[i]
		if inst.Opcode() == bytecode.STRLOAD {
			offset := int(binary.BigEndian.Uint32(inst[1:]))
			size := int(binary.BigEndian.Uint32(inst[5:]))
			instructions[i] = bytecode.New(bytecode.STRLOAD, uint64(literals[string(constants[offset:offset+size])]), uint64(size))
		}
	}

	for i := len(instructions) - 1; i >= 0; i-- {
		if instructions[i].Opcode() == bytecode.NOP {
			instructions = append(instructions[:i], instructions[i+1:]...)
		}
	}

	return instructions, compressed
}
