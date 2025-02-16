package ast

import "bytes"

type Program struct {
	Statements []Statement
}

func NewProgram(statements ...Statement) *Program {
	return &Program{Statements: statements}
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, stmt := range p.Statements {
		out.WriteString(stmt.String())
		out.WriteString("\n")
	}
	return out.String()
}
