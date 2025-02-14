package main

import (
	"github.com/siyul-park/minijs/repl"
	"log"
	"os"
)

func main() {
	r := repl.New("> ")
	if err := r.Start(os.Stdin, os.Stdout); err != nil {
		log.Fatal(err)
	}
}
