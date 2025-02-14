package ast

import (
	"bytes"
	"github.com/siyul-park/minijs/token"
)

type PrefixExpression struct {
	Token token.Token
	Right Node
}

type InfixExpression struct {
	Token token.Token
	Left  Node
	Right Node
}

func NewPrefixExpression(token token.Token, right Node) *PrefixExpression {
	return &PrefixExpression{Token: token, Right: right}
}

func NewInfixExpression(token token.Token, left, right Node) *InfixExpression {
	return &InfixExpression{Token: token, Left: left, Right: right}
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
