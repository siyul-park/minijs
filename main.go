package main

import (
	"bufio"
	"fmt"
	"github.com/siyul-park/miniscript/lexer"
	"github.com/siyul-park/miniscript/parser"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()

		l := lexer.New(line)
		p := parser.New(l)

		program, err := p.Parse()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(program.String())
	}
}
