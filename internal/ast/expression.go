package ast

import (
	"bytes"

	"github.com/siyul-park/minijs/internal/token"
)

type Expression interface {
	Node
	_expression()
}

type expression struct {
}

func (expression) _expression() {
}

type PrefixExpression struct {
	expression
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

type InfixExpression struct {
	expression
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

type AssignmentExpression struct {
	expression
	Token token.Token
	Left  Expression
	Right Expression
}

func NewAssignmentExpression(token token.Token, left Expression, right Expression) *AssignmentExpression {
	return &AssignmentExpression{Token: token, Left: left, Right: right}
}

func (n *AssignmentExpression) String() string {
	var out bytes.Buffer
	out.WriteString(n.Left.String())
	out.WriteString(n.Token.Literal)
	out.WriteString(n.Right.String())
	return out.String()
}
