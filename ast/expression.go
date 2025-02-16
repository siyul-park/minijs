package ast

import (
	"bytes"

	"github.com/siyul-park/minijs/token"
)

type Expression interface {
	Node
	expression()
}

type PrefixExpression struct {
	Token token.Token
	Right Expression
}

func NewPrefixExpression(token token.Token, right Expression) *PrefixExpression {
	return &PrefixExpression{Token: token, Right: right}
}

func (n *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(n.Token.Literal)
	out.WriteString(n.Right.String())
	out.WriteString(")")
	return out.String()
}

func (n *PrefixExpression) expression() {
}

type InfixExpression struct {
	Token token.Token
	Left  Expression
	Right Expression
}

func NewInfixExpression(token token.Token, left, right Expression) *InfixExpression {
	return &InfixExpression{Token: token, Left: left, Right: right}
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

func (n *InfixExpression) expression() {
}
