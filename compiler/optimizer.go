package compiler

import (
	"bytes"
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
	consts := code.Constants

	var insts []bytecode.Instruction
	for offset := 0; offset < len(code.Instructions); {
		inst, size := code.Instruction(offset)
		insts = append(insts, inst)
		offset += size
	}

	for {
		prevInsts := append([]bytecode.Instruction{}, insts...)
		prevConsts := append([]byte{}, consts...)

		var err error
		insts, consts, err = o.fusion(insts, consts)
		if err != nil {
			return bytecode.Bytecode{}, err
		}
		insts, consts = o.compress(insts, consts)

		if len(prevInsts) != len(insts) {
			continue
		}

		same := true
		for i := range prevInsts {
			if !bytes.Equal(prevInsts[i], insts[i]) {
				same = false
				break
			}
		}

		if same && bytes.Equal(prevConsts, consts) {
			break
		}
	}

	code.Instructions = nil
	code.Constants = consts
	code.Add(insts...)
	return code, nil
}

func (o *Optimizer) fusion(insts []bytecode.Instruction, consts []byte) ([]bytecode.Instruction, []byte, error) {
	offsets := map[string]int{}
	for i := 0; i < len(insts); i++ {
		inst := insts[i]
		if inst.Opcode() == bytecode.CLOAD {
			offset := int(binary.BigEndian.Uint32(inst[1:]))
			size := int(binary.BigEndian.Uint32(inst[5:]))

			literal := string(consts[offset : offset+size])
			offsets[literal] = offset
		}
	}

	for i := 0; i < len(insts); i++ {
		inst := insts[i]
		if i > 0 {
			operand := insts[i-1]
			switch operand.Opcode() {
			case bytecode.BLOAD, bytecode.I32LOAD, bytecode.F64LOAD, bytecode.CLOAD:
				switch inst.Opcode() {
				case bytecode.I32TOB:
					code := bytecode.Bytecode{Constants: consts}
					code.Add(operand, inst)
					if err := o.interpreter.Execute(code); err != nil {
						return nil, nil, err
					}

					val, _ := o.interpreter.Pop().(interpreter.Bool)
					insts[i-1] = bytecode.New(bytecode.NOP)
					insts[i] = bytecode.New(bytecode.BLOAD, uint64(val))
				case bytecode.BTOI32, bytecode.F64TOI32, bytecode.CTOI32:
					code := bytecode.Bytecode{Constants: consts}
					code.Add(operand, inst)
					if err := o.interpreter.Execute(code); err != nil {
						return nil, nil, err
					}

					val, _ := o.interpreter.Pop().(interpreter.Int32)

					insts[i-1] = bytecode.New(bytecode.NOP)
					insts[i] = bytecode.New(bytecode.I32LOAD, uint64(val))
				case bytecode.I32TOF64, bytecode.CTOF64:
					code := bytecode.Bytecode{Constants: consts}
					code.Add(operand, inst)
					if err := o.interpreter.Execute(code); err != nil {
						return nil, nil, err
					}

					val, _ := o.interpreter.Pop().(interpreter.Float64)

					insts[i-1] = bytecode.New(bytecode.NOP)
					insts[i] = bytecode.New(bytecode.F64LOAD, math.Float64bits(float64(val)))
				case bytecode.BTOC, bytecode.I32TOC, bytecode.F64TOC:
					code := bytecode.Bytecode{Constants: consts}
					code.Add(operand, inst)
					if err := o.interpreter.Execute(code); err != nil {
						return nil, nil, err
					}

					val, _ := o.interpreter.Pop().(interpreter.String)

					offset, ok := offsets[string(val)]
					if !ok {
						offset = len(consts)
						consts = append(consts, []byte(string(val)+"\x00")...)
						offsets[string(val)] = offset
					}

					insts[i-1] = bytecode.New(bytecode.NOP)
					insts[i] = bytecode.New(bytecode.CLOAD, uint64(offset), uint64(len(val)))
				default:
				}
			default:
			}
		}

		if i > 1 {
			operand1 := insts[i-1]
			operand2 := insts[i-2]
			if operand1.Opcode() == operand2.Opcode() {
				switch operand1.Opcode() {
				case bytecode.BLOAD, bytecode.I32LOAD, bytecode.F64LOAD, bytecode.CLOAD:
					switch inst.Opcode() {
					case bytecode.I32ADD, bytecode.I32SUB, bytecode.I32MUL, bytecode.I32DIV, bytecode.I32MOD:
						code := bytecode.Bytecode{Constants: consts}
						code.Add(operand2, operand1, inst)
						if err := o.interpreter.Execute(code); err != nil {
							return nil, nil, err
						}

						val, _ := o.interpreter.Pop().(interpreter.Int32)

						insts[i-2] = bytecode.New(bytecode.NOP)
						insts[i-1] = bytecode.New(bytecode.NOP)
						insts[i] = bytecode.New(bytecode.I32LOAD, uint64(val))
					case bytecode.F64ADD, bytecode.F64SUB, bytecode.F64MUL, bytecode.F64DIV, bytecode.F64MOD:
						code := bytecode.Bytecode{Constants: consts}
						code.Add(operand2, operand1, inst)
						if err := o.interpreter.Execute(code); err != nil {
							return nil, nil, err
						}

						val, _ := o.interpreter.Pop().(interpreter.Float64)

						insts[i-2] = bytecode.New(bytecode.NOP)
						insts[i-1] = bytecode.New(bytecode.NOP)
						insts[i] = bytecode.New(bytecode.F64LOAD, math.Float64bits(float64(val)))
					case bytecode.CADD:
						code := bytecode.Bytecode{Constants: consts}
						code.Add(operand2, operand1, inst)
						if err := o.interpreter.Execute(code); err != nil {
							return nil, nil, err
						}

						val, _ := o.interpreter.Pop().(interpreter.String)

						offset, ok := offsets[string(val)]
						if !ok {
							offset = len(consts)
							consts = append(consts, []byte(string(val)+"\x00")...)
							offsets[string(val)] = offset
						}

						insts[i-2] = bytecode.New(bytecode.NOP)
						insts[i-1] = bytecode.New(bytecode.NOP)
						insts[i] = bytecode.New(bytecode.CLOAD, uint64(offset), uint64(len(val)))
					default:
					}
				default:
				}
			}
		}
	}
	return insts, consts, nil
}

func (o *Optimizer) compress(insts []bytecode.Instruction, consts []byte) ([]bytecode.Instruction, []byte) {
	offsets := map[string]int{}
	for i := 0; i < len(insts); i++ {
		inst := insts[i]
		if inst.Opcode() == bytecode.CLOAD {
			offset := int(binary.BigEndian.Uint32(inst[1:]))
			size := int(binary.BigEndian.Uint32(inst[5:]))

			literal := string(consts[offset : offset+size])
			offsets[literal] = offset
		}
	}

	compressed := make([]byte, 0, len(consts))
	for literal := range offsets {
		offsets[literal] = len(compressed)
		compressed = append(compressed, []byte(literal+"\x00")...)
	}

	for i := 0; i < len(insts); i++ {
		inst := insts[i]
		if inst.Opcode() == bytecode.CLOAD {
			offset := int(binary.BigEndian.Uint32(inst[1:]))
			size := int(binary.BigEndian.Uint32(inst[5:]))
			insts[i] = bytecode.New(bytecode.CLOAD, uint64(offsets[string(consts[offset:offset+size])]), uint64(size))
		}
	}

	for i := len(insts) - 1; i >= 0; i-- {
		if insts[i].Opcode() == bytecode.NOP {
			insts = append(insts[:i], insts[i+1:]...)
		}
	}

	return insts, compressed
}
