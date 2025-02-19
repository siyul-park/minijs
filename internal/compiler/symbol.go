package compiler

import (
	"github.com/siyul-park/minijs/internal/interpreter"
)

type Symbol struct {
	Name  string
	Index int
	Type  interpreter.Type
}

type SymbolTable struct {
	symbols map[string]*Symbol
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		symbols: make(map[string]*Symbol),
	}
}

func (s *SymbolTable) Define(name string) *Symbol {
	sym := &Symbol{Name: name, Index: len(s.symbols)}
	s.symbols[name] = sym
	return sym
}

func (s *SymbolTable) Resolve(name string) (*Symbol, bool) {
	sym, ok := s.symbols[name]
	return sym, ok
}
