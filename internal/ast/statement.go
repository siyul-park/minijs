package ast

import (
	"strings"
)

type Statement interface {
	Node
	_statement()
}

type statement struct {
}

func (statement) _statement() {
}

type EmptyStatement struct {
	statement
}

func NewEmptyStatement() *EmptyStatement {
	return &EmptyStatement{}
}

func (n *EmptyStatement) String() string {
	return ";"
}

type BlockStatement struct {
	statement
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

type ExpressionStatement struct {
	statement
	Expression Expression
}

func NewExpressionStatement(expression Expression) *ExpressionStatement {
	return &ExpressionStatement{Expression: expression}
}

func (n *ExpressionStatement) String() string {
	return n.Expression.String() + ";"
}
