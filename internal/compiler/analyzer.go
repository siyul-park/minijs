package compiler

import (
	"strings"

	"github.com/siyul-park/minijs/internal/ast"
	"github.com/siyul-park/minijs/internal/interpreter"
	"github.com/siyul-park/minijs/internal/token"
)

type Analyzer struct {
	symbols map[ast.Node]*Symbol
}

type Symbol struct {
	Kind interpreter.Kind
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{symbols: make(map[ast.Node]*Symbol)}
}

func (a *Analyzer) Analyze(node ast.Node) *Symbol {
	return a.analyze(node)
}

func (a *Analyzer) Clear() {
	a.symbols = make(map[ast.Node]*Symbol)
}

func (a *Analyzer) analyze(node ast.Node) *Symbol {
	if sb, found := a.symbols[node]; found {
		return sb
	}

	var sb *Symbol
	switch node := node.(type) {
	case *ast.Program:
		sb = a.program(node)
	case *ast.EmptyStatement:
		sb = a.emptyStatement(node)
	case *ast.BlockStatement:
		sb = a.blockStatement(node)
	case *ast.ExpressionStatement:
		sb = a.expressionStatement(node)
	case *ast.PrefixExpression:
		sb = a.prefixExpression(node)
	case *ast.InfixExpression:
		sb = a.infixExpression(node)
	case *ast.BoolLiteral:
		sb = a.boolLiteral(node)
	case *ast.NumberLiteral:
		sb = a.numberLiteral(node)
	case *ast.StringLiteral:
		sb = a.stringLiteral(node)
	default:
		sb = nil
	}

	a.symbols[node] = sb
	return sb
}

func (a *Analyzer) program(node *ast.Program) *Symbol {
	sb := &Symbol{Kind: interpreter.KindVoid}
	for _, stmt := range node.Statements {
		a.analyze(stmt)
	}
	return sb
}

func (a *Analyzer) emptyStatement(_ *ast.EmptyStatement) *Symbol {
	return &Symbol{Kind: interpreter.KindVoid}
}

func (a *Analyzer) blockStatement(node *ast.BlockStatement) *Symbol {
	sb := &Symbol{Kind: interpreter.KindVoid}
	for _, stmt := range node.Statements {
		a.analyze(stmt)
	}
	return sb
}

func (a *Analyzer) expressionStatement(node *ast.ExpressionStatement) *Symbol {
	a.analyze(node.Expression)
	return &Symbol{Kind: interpreter.KindVoid}
}

func (a *Analyzer) prefixExpression(node *ast.PrefixExpression) *Symbol {
	right := a.analyze(node.Right)
	if right == nil {
		return nil
	}

	sb := &Symbol{}
	switch node.Token.Type {
	case token.PLUS, token.MINUS:
		switch right.Kind {
		case interpreter.KindBool:
			sb.Kind = interpreter.KindInt32
		case interpreter.KindInt32, interpreter.KindFloat64:
			sb.Kind = right.Kind
		case interpreter.KindString:
			sb.Kind = interpreter.KindFloat64
		default:
		}
	default:
	}
	return sb
}

func (a *Analyzer) infixExpression(node *ast.InfixExpression) *Symbol {
	left := a.analyze(node.Left)
	if left == nil {
		return nil
	}

	right := a.analyze(node.Right)
	if right == nil {
		return nil
	}

	sb := &Symbol{}
	switch node.Token.Type {
	case token.PLUS:
		if left.Kind == interpreter.KindString || right.Kind == interpreter.KindString {
			sb.Kind = interpreter.KindString
		} else if left.Kind == interpreter.KindFloat64 || right.Kind == interpreter.KindFloat64 {
			sb.Kind = interpreter.KindFloat64
		} else {
			sb.Kind = interpreter.KindInt32
		}
	case token.DIVIDE, token.MODULUS:
		sb.Kind = interpreter.KindFloat64
	default:
		if left.Kind == interpreter.KindInt32 && right.Kind == interpreter.KindInt32 {
			sb.Kind = interpreter.KindInt32
		} else {
			sb.Kind = interpreter.KindFloat64
		}
	}
	return sb
}

func (a *Analyzer) boolLiteral(_ *ast.BoolLiteral) *Symbol {
	return &Symbol{Kind: interpreter.KindBool}
}

func (a *Analyzer) numberLiteral(node *ast.NumberLiteral) *Symbol {
	kind := interpreter.KindInt32
	if strings.Contains(node.Token.Literal, ".") || strings.Contains(node.Token.Literal, "e") {
		kind = interpreter.KindFloat64
	} else if node.Value != float64(int32(node.Value)) {
		kind = interpreter.KindFloat64
	}
	return &Symbol{Kind: kind}
}

func (a *Analyzer) stringLiteral(_ *ast.StringLiteral) *Symbol {
	return &Symbol{Kind: interpreter.KindString}
}
