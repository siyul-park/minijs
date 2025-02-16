package compiler

import (
	"strings"

	"github.com/siyul-park/minijs/ast"
	"github.com/siyul-park/minijs/interpreter"
	"github.com/siyul-park/minijs/token"
)

type Analyzer struct {
	meta map[ast.Node]*Meta
}

type Meta struct {
	Kind interpreter.Kind
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{meta: make(map[ast.Node]*Meta)}
}

func (a *Analyzer) Analyze(node ast.Node) *Meta {
	if meta, found := a.meta[node]; found {
		return meta
	}

	var meta *Meta
	switch node := node.(type) {
	case *ast.Program:
		meta = a.program(node)
	case *ast.EmptyStatement:
		meta = a.emptyStatement(node)
	case *ast.BlockStatement:
		meta = a.blockStatement(node)
	case *ast.ExpressionStatement:
		meta = a.expressionStatement(node)
	case *ast.PrefixExpression:
		meta = a.prefixExpression(node)
	case *ast.InfixExpression:
		meta = a.infixExpression(node)
	case *ast.BoolLiteral:
		meta = a.boolLiteral(node)
	case *ast.NumberLiteral:
		meta = a.numberLiteral(node)
	case *ast.StringLiteral:
		meta = a.stringLiteral(node)
	default:
		meta = nil
	}

	a.meta[node] = meta
	return meta
}

func (a *Analyzer) Close() {
	a.meta = make(map[ast.Node]*Meta)
}

func (a *Analyzer) program(node *ast.Program) *Meta {
	for _, stmt := range node.Statements {
		a.Analyze(stmt)
	}
	return &Meta{Kind: interpreter.KindVoid}
}

func (a *Analyzer) emptyStatement(_ *ast.EmptyStatement) *Meta {
	return &Meta{Kind: interpreter.KindVoid}
}

func (a *Analyzer) blockStatement(node *ast.BlockStatement) *Meta {
	for _, stmt := range node.Statements {
		a.Analyze(stmt)
	}
	return &Meta{Kind: interpreter.KindVoid}
}

func (a *Analyzer) expressionStatement(node *ast.ExpressionStatement) *Meta {
	a.Analyze(node.Expression)
	return &Meta{Kind: interpreter.KindVoid}
}

func (a *Analyzer) prefixExpression(node *ast.PrefixExpression) *Meta {
	right := a.Analyze(node.Right)
	if right == nil {
		return nil
	}

	kind := interpreter.KindUnknown
	switch node.Token {
	case token.PLUS, token.MINUS:
		switch right.Kind {
		case interpreter.KindBool:
			kind = interpreter.KindInt32
		case interpreter.KindInt32, interpreter.KindFloat64:
			kind = right.Kind
		case interpreter.KindString:
			kind = interpreter.KindFloat64
		default:
		}
	default:
	}
	return &Meta{Kind: kind}
}

func (a *Analyzer) infixExpression(node *ast.InfixExpression) *Meta {
	left := a.Analyze(node.Left)
	right := a.Analyze(node.Right)

	kind := interpreter.KindUnknown
	switch node.Token {
	case token.PLUS:
		if left.Kind == interpreter.KindString || right.Kind == interpreter.KindString {
			kind = interpreter.KindString
		} else if left.Kind == interpreter.KindFloat64 || right.Kind == interpreter.KindFloat64 {
			kind = interpreter.KindFloat64
		} else {
			kind = interpreter.KindInt32
		}
	case token.DIVIDE, token.MODULO:
		kind = interpreter.KindFloat64
	default:
		if left.Kind == interpreter.KindInt32 && right.Kind == interpreter.KindInt32 {
			kind = interpreter.KindInt32
		} else {
			kind = interpreter.KindFloat64
		}
	}
	return &Meta{Kind: kind}
}

func (a *Analyzer) boolLiteral(_ *ast.BoolLiteral) *Meta {
	return &Meta{Kind: interpreter.KindBool}
}

func (a *Analyzer) numberLiteral(node *ast.NumberLiteral) *Meta {
	kind := interpreter.KindInt32
	if strings.Contains(node.Token.Literal, ".") || strings.Contains(node.Token.Literal, "e") {
		kind = interpreter.KindFloat64
	} else if node.Value != float64(int32(node.Value)) {
		kind = interpreter.KindFloat64
	}
	return &Meta{Kind: kind}
}

func (a *Analyzer) stringLiteral(_ *ast.StringLiteral) *Meta {
	return &Meta{Kind: interpreter.KindString}
}
