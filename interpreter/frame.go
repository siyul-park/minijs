package interpreter

import "github.com/siyul-park/minijs/bytecode"

type Frame struct {
	bytecode bytecode.Bytecode
	ip       int
	bp       int
}

func NewFrame(bytecode bytecode.Bytecode, bp int) *Frame {
	return &Frame{bytecode: bytecode, ip: -1, bp: bp}
}

func (f *Frame) Next() bytecode.Instruction {
	instruction, read := f.bytecode.Instruction(f.ip + 1)
	f.ip += read
	return instruction
}
