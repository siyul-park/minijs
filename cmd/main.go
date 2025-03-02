package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/siyul-park/minijs"

	"github.com/siyul-park/minijs/internal/compiler"
	"github.com/siyul-park/minijs/internal/interpreter"
	"github.com/siyul-park/minijs/internal/lexer"
	"github.com/siyul-park/minijs/internal/parser"
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
	r := minijs.NewREPL("> ", minijs.REPLOption{PrintBytecode: printBytecode})
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

	l := lexer.New(file)
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

	o := interpreter.NewOptimizer()
	code, err = o.Optimize(code)
	if err != nil {
		log.Fatal("Error optimize program: ", err)
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
