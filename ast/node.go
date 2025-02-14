package ast

import (
	"bytes"
	"github.com/siyul-park/minijs/token"
)

type Node interface {
	String() string
}

type Program struct {
	Statements []*Statement
}

type Statement struct {
	Node Node
}

type NumberLiteral struct {
	Token token.Token
	Value float64
}

type StringLiteral struct {
	Token token.Token
	Value string
}

type BoolLiteral struct {
	Token token.Token
	Value bool
}

type IdentifierLiteral struct {
	Token token.Token
	Value string
}

type PrefixExpression struct {
	Token token.Token
	Right Node
}

type InfixExpression struct {
	Token token.Token
	Left  Node
	Right Node
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, stmt := range p.Statements {
		out.WriteString(stmt.String())
		out.WriteString("\n")
	}
	return out.String()
}

func (n *Statement) String() string {
	return n.Node.String() + ";"
}

func (n *NumberLiteral) String() string {
	return n.Token.Literal
}

func (n *StringLiteral) String() string {
	return n.Token.Literal
}

func (n *BoolLiteral) String() string {
	return n.Token.Literal
}

func (n *IdentifierLiteral) String() string {
	return n.Value
}

func (n *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(n.Token.Literal)
	out.WriteString(n.Right.String())
	out.WriteString(")")
	return out.String()
}

func (n *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(n.Left.String())
	out.WriteString(n.Token.Literal)
	out.WriteString(n.Right.String())
	out.WriteString(")")
	return out.String()
}
