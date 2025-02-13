package ast

import (
	"bytes"
	"github.com/siyul-park/miniscript/token"
)

type Node interface {
	Literal() string
	String() string
}

type Program struct {
	Nodes []Node
}

type IntLiteral struct {
	Token token.Token
	Value int
}

type FloatLiteral struct {
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

func (p *Program) Literal() string {
	if len(p.Nodes) > 0 {
		return p.Nodes[0].Literal()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, n := range p.Nodes {
		out.WriteString(n.String())
	}
	return out.String()
}

func (n *IntLiteral) Literal() string {
	return n.Token.Literal
}

func (n *IntLiteral) String() string {
	return n.Token.Literal
}

func (n *FloatLiteral) Literal() string {
	return n.Token.Literal
}

func (n *FloatLiteral) String() string {
	return n.Token.Literal
}

func (n *StringLiteral) Literal() string {
	return n.Token.Literal
}

func (n *StringLiteral) String() string {
	return n.Token.Literal
}

func (n *BoolLiteral) Literal() string {
	return n.Token.Literal
}

func (n *BoolLiteral) String() string {
	return n.Token.Literal
}

func (n *IdentifierLiteral) Literal() string {
	return n.Token.Literal
}

func (n *IdentifierLiteral) String() string {
	return n.Value
}

func (n *PrefixExpression) Literal() string {
	return n.Token.Literal
}

func (n *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(n.Token.Literal)
	out.WriteString(n.Right.String())
	out.WriteString(")")
	return out.String()
}

func (n *InfixExpression) Literal() string {
	return n.Token.Literal
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
