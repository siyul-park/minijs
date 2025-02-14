package ast

import "bytes"

type Program struct {
	Statements []*Statement
}

func NewProgram(stmts ...*Statement) *Program {
	return &Program{Statements: stmts}
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, stmt := range p.Statements {
		out.WriteString(stmt.String())
		out.WriteString("\n")
	}
	return out.String()
}
