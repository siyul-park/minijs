package ast

import (
	"strings"
)

type Statement interface {
	Node
	statement()
}

type EmptyStatement struct {
}

func NewEmptyStatement() *EmptyStatement {
	return &EmptyStatement{}
}

func (n *EmptyStatement) String() string {
	return ";"
}

func (n *EmptyStatement) statement() {
}

type BlockStatement struct {
	Statements []Statement
}

func NewBlockStatement(statements ...Statement) *BlockStatement {
	return &BlockStatement{Statements: statements}
}

func (n *BlockStatement) String() string {
	indent := "    "

	var builder strings.Builder
	builder.WriteString("{\n")
	for _, node := range n.Statements {
		builder.WriteString(indent)
		builder.WriteString(node.String())
		builder.WriteString("\n")
	}
	builder.WriteString("}")
	return builder.String()
}

func (n *BlockStatement) statement() {
}

type ExpressionStatement struct {
	Expression Expression
}

func NewExpressionStatement(expression Expression) *ExpressionStatement {
	return &ExpressionStatement{Expression: expression}
}

func (n *ExpressionStatement) String() string {
	return n.Expression.String() + ";"
}

func (n *ExpressionStatement) statement() {
}
