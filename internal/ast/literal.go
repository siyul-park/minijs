package ast

import (
	"github.com/siyul-park/minijs/internal/token"
)

type BoolLiteral struct {
	expression
	Token token.Token
	Value bool
}

func NewBoolLiteral(tok token.Token, value bool) *BoolLiteral {
	return &BoolLiteral{Token: tok, Value: value}
}

func (n *BoolLiteral) String() string {
	return n.Token.Literal
}

type NumberLiteral struct {
	expression
	Token token.Token
	Value float64
}

func NewNumberLiteral(tok token.Token, value float64) *NumberLiteral {
	return &NumberLiteral{Token: tok, Value: value}
}

func (n *NumberLiteral) String() string {
	return n.Token.Literal
}

type StringLiteral struct {
	expression
	Token token.Token
	Value string
}

func NewStringLiteral(tok token.Token, value string) *StringLiteral {
	return &StringLiteral{Token: tok, Value: value}
}

func (n *StringLiteral) String() string {
	return "\"" + n.Token.Literal + "\""
}

type IdentifierLiteral struct {
	expression
	Token token.Token
	Value string
}

func NewIdentifierLiteral(tok token.Token, value string) *IdentifierLiteral {
	return &IdentifierLiteral{Token: tok, Value: value}
}

func (n *IdentifierLiteral) String() string {
	return n.Value
}
