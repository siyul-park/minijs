package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/siyul-park/minijs/compiler"
	"github.com/siyul-park/minijs/interpreter"
	"github.com/siyul-park/minijs/lexer"
	"github.com/siyul-park/minijs/parser"
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

		l := lexer.New(line)
		p := parser.New(l)

		program, err := p.Parse()
		if err != nil {
			if err := r.error(writer, err); err != nil {
				return err
			}
			continue
		}

		stmts := program.Statements
		if len(stmts) == 0 {
			continue
		}

		code, err := c.Compile(stmts[len(stmts)-1].Node)
		if err != nil {
			if err := r.error(writer, err); err != nil {
				return err
			}
			continue
		}

		if err := i.Execute(code); err != nil {
			if err := r.error(writer, err); err != nil {
				return err
			}
			continue
		}

		if _, err := fmt.Fprintln(writer, i.Top()); err != nil {
			return err
		}
	}

	return nil
}

func (r *REPL) error(writer io.Writer, err error) error {
	if _, writeErr := fmt.Fprintln(writer, err); writeErr != nil {
		return writeErr
	}
	return err
}
