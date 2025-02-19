package ast

import (
	"bytes"
	"strings"

	"github.com/siyul-park/minijs/internal/token"
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
	var out strings.Builder
	out.WriteString("{\n")
	for _, node := range n.Statements {
		out.WriteString(node.String())
		out.WriteString(";")
	}
	out.WriteString("}")
	return out.String()
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

type VariableStatement struct {
	statement
	Token token.Token
	Right []*AssignmentExpression
}

func NewVariableStatement(token token.Token, right ...*AssignmentExpression) *VariableStatement {
	return &VariableStatement{Token: token, Right: right}
}

func (n *VariableStatement) String() string {
	var out bytes.Buffer
	out.WriteString(n.Token.Literal)
	out.WriteString(" ")
	for i, node := range n.Right {
		out.WriteString(node.String())
		if i < len(n.Right)-1 {
			out.WriteString(",")
		}
	}
	out.WriteString(";")
	return out.String()
}
