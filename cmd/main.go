package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/siyul-park/minijs/compiler"
	"github.com/siyul-park/minijs/interpreter"
	"github.com/siyul-park/minijs/lexer"
	"github.com/siyul-park/minijs/parser"
	"github.com/siyul-park/minijs/repl"
)

func main() {
	printBytecode := flag.Bool("print-bytecode", false, "")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		runREPL(*printBytecode)
		return
	}
	runFile(args[0], *printBytecode)
}

func runREPL(printBytecode bool) {
	r := repl.New("> ", repl.Option{PrintBytecode: printBytecode})
	if err := r.Start(os.Stdin, os.Stdout); err != nil {
		log.Fatal("Error starting REPL: ", err)
	}
}

func runFile(filePath string, printBytecode bool) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Error opening file: ", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("Error reading file: ", err)
	}

	l := lexer.New(string(data))
	p := parser.New(l)

	program, err := p.Parse()
	if err != nil {
		log.Fatal("Error parsing program: ", err)
	}

	c := compiler.New()
	code, err := c.Compile(program)
	if err != nil {
		log.Fatal("Error compiling program: ", err)
	}

	if printBytecode {
		fmt.Println(code.String())
	} else {
		i := interpreter.New()
		if err := i.Execute(code); err != nil {
			log.Fatal("Error executing code: ", err)
		}
	}
}
