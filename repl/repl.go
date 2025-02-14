package repl

import (
	"bufio"
	"fmt"
	"github.com/siyul-park/minijs/compiler"
	"github.com/siyul-park/minijs/interpreter"
	"github.com/siyul-park/minijs/lexer"
	"github.com/siyul-park/minijs/parser"
	"io"
)

type REPL struct {
	prompt string
}

func New(prompt string) *REPL {
	return &REPL{
		prompt: prompt,
	}
}

func (r *REPL) Start(reader io.Reader, writer io.Writer) error {
	scanner := bufio.NewScanner(reader)

	for {
		if _, err := fmt.Fprint(writer, r.prompt); err != nil {
			return err
		}

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}
			break
		}

		line := scanner.Text()

		l := lexer.New(line)
		p := parser.New(l)

		program, err := p.Parse()
		if err != nil {
			if _, err := fmt.Fprintln(writer, err); err != nil {
				return err
			}
			continue
		}

		stmts := program.Statements
		if len(stmts) == 0 {
			continue
		}

		c := compiler.New(stmts[0].Node)
		code, err := c.Compile()
		if err != nil {
			if _, err := fmt.Fprintln(writer, err); err != nil {
				return err
			}
			continue
		}

		i := interpreter.New(code)
		if err := i.Execute(); err != nil {
			if _, err := fmt.Fprintln(writer, err); err != nil {
				return err
			}
			continue
		}

		if _, err := fmt.Fprintln(writer, i.Peek(0)); err != nil {
			return err
		}
	}

	return nil
}
