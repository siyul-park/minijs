package ast

import (
	"github.com/siyul-park/minijs/internal/token"
)

type NullLiteral struct {
	expression
	Token token.Token
}

func NewNullLiteral(tok token.Token) *NullLiteral {
	return &NullLiteral{Token: tok}
}

func (n *NullLiteral) String() string {
	return n.Token.Literal
}

type UndefinedLiteral struct {
	expression
	Token token.Token
}

func NewUndefinedLiteral(tok token.Token) *UndefinedLiteral {
	return &UndefinedLiteral{Token: tok}
}

func (n *UndefinedLiteral) String() string {
	return n.Token.Literal
}

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
