package minijs

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/siyul-park/minijs/internal/bytecode"
	"github.com/siyul-park/minijs/internal/compiler"
	"github.com/siyul-park/minijs/internal/interpreter"
	"github.com/siyul-park/minijs/internal/lexer"
	"github.com/siyul-park/minijs/internal/parser"
)

type REPLOption struct {
	PrintBytecode bool
}

type REPL struct {
	prompt        string
	printBytecode bool
}

func NewREPL(prompt string, opts ...REPLOption) *REPL {
	repl := &REPL{prompt: prompt}

	for _, opt := range opts {
		repl.printBytecode = opt.PrintBytecode
	}

	return repl
}

func (r *REPL) Start(reader io.Reader, writer io.Writer) error {
	scanner := bufio.NewScanner(reader)

	c := compiler.New()
	i := interpreter.New()

	for {
		if r.prompt != "" {
			if _, err := fmt.Fprint(writer, r.prompt); err != nil {
				return err
			}
		}

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}
			break
		}

		line := scanner.Text()

		l := lexer.New(strings.NewReader(line))
		p := parser.New(l)

		program, err := p.Parse()
		if err != nil {
			if err := r.error(writer, err); err != nil {
				return err
			}
			continue
		}

		code, err := c.Compile(program)
		if err != nil {
			if err := r.error(writer, err); err != nil {
				return err
			}
			continue
		}

		if r.printBytecode {
			if _, err := fmt.Fprintln(writer, code.String()); err != nil {
				return err
			}
		}

		var insts []bytecode.Instruction
		for offset := 0; offset < len(code.Instructions); {
			inst, size := code.Instruction(offset)
			insts = append(insts, inst)
			offset += size
		}
		if len(insts) > 0 {
			if insts[len(insts)-1].Opcode() == bytecode.POP {
				insts = insts[:len(insts)-1]
			}

			code.Instructions = nil
			code.Emit(insts...)
		}

		if err := i.Execute(code); err != nil {
			if err := r.error(writer, err); err != nil {
				return err
			}
			continue
		}

		if _, err := fmt.Fprintln(writer, i.Pop()); err != nil {
			return err
		}
	}

	return nil
}

func (r *REPL) error(writer io.Writer, err error) error {
	_, err = fmt.Fprintln(writer, err)
	return err
}
