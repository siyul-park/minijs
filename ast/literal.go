package ast

import (
	"strings"

	"github.com/siyul-park/minijs/token"
)

type BoolLiteral struct {
	Token token.Token
	Value bool
}

func NewBoolLiteral(tok token.Token, value bool) *BoolLiteral {
	return &BoolLiteral{Token: tok, Value: value}
}

func (n *BoolLiteral) String() string {
	return n.Token.Literal
}

func (n *BoolLiteral) expression() {
}

type NumberLiteral struct {
	Token token.Token
	Value float64
}

func NewNumberLiteral(tok token.Token, value float64) *NumberLiteral {
	return &NumberLiteral{Token: tok, Value: value}
}

func (n *NumberLiteral) IsInteger() bool {
	if strings.Contains(n.Token.Literal, ".") || strings.Contains(n.Token.Literal, "e") {
		return false
	}
	return n.Value == float64(int32(n.Value))
}

func (n *NumberLiteral) String() string {
	return n.Token.Literal
}

func (n *NumberLiteral) expression() {
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func NewStringLiteral(tok token.Token, value string) *StringLiteral {
	return &StringLiteral{Token: tok, Value: value}
}

func (n *StringLiteral) String() string {
	return n.Token.Literal
}

func (n *StringLiteral) expression() {
}

type IdentifierLiteral struct {
	Token token.Token
	Value string
}

func NewIdentifierLiteral(tok token.Token, value string) *IdentifierLiteral {
	return &IdentifierLiteral{Token: tok, Value: value}
}

func (n *IdentifierLiteral) String() string {
	return n.Value
}

func (n *IdentifierLiteral) expression() {
}
