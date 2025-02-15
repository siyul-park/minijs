package main

import (
	"log"
	"os"

	"github.com/siyul-park/minijs/repl"
)

func main() {
	r := repl.New("> ")
	if err := r.Start(os.Stdin, os.Stdout); err != nil {
		log.Fatal(err)
	}
}
