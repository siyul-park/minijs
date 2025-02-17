package compiler

import (
	"strings"

	"github.com/siyul-park/minijs/internal/ast"
	"github.com/siyul-park/minijs/internal/interpreter"
	"github.com/siyul-park/minijs/internal/token"
)

type Analyzer struct {
	semantic map[ast.Node]*Semantic
}

type Semantic struct {
	Kind interpreter.Kind
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{semantic: make(map[ast.Node]*Semantic)}
}

func (a *Analyzer) Analyze(node ast.Node) *Semantic {
	if sem, found := a.semantic[node]; found {
		return sem
	}

	var meta *Semantic
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

	a.semantic[node] = meta
	return meta
}

func (a *Analyzer) Close() {
	a.semantic = make(map[ast.Node]*Semantic)
}

func (a *Analyzer) program(node *ast.Program) *Semantic {
	for _, stmt := range node.Statements {
		a.Analyze(stmt)
	}
	return &Semantic{Kind: interpreter.KindVoid}
}

func (a *Analyzer) emptyStatement(_ *ast.EmptyStatement) *Semantic {
	return &Semantic{Kind: interpreter.KindVoid}
}

func (a *Analyzer) blockStatement(node *ast.BlockStatement) *Semantic {
	for _, stmt := range node.Statements {
		a.Analyze(stmt)
	}
	return &Semantic{Kind: interpreter.KindVoid}
}

func (a *Analyzer) expressionStatement(node *ast.ExpressionStatement) *Semantic {
	a.Analyze(node.Expression)
	return &Semantic{Kind: interpreter.KindVoid}
}

func (a *Analyzer) prefixExpression(node *ast.PrefixExpression) *Semantic {
	right := a.Analyze(node.Right)
	if right == nil {
		return nil
	}

	kind := interpreter.KindUnknown
	switch node.Token.Type {
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
	return &Semantic{Kind: kind}
}

func (a *Analyzer) infixExpression(node *ast.InfixExpression) *Semantic {
	left := a.Analyze(node.Left)
	right := a.Analyze(node.Right)

	kind := interpreter.KindUnknown
	switch node.Token.Type {
	case token.PLUS:
		if left.Kind == interpreter.KindString || right.Kind == interpreter.KindString {
			kind = interpreter.KindString
		} else if left.Kind == interpreter.KindFloat64 || right.Kind == interpreter.KindFloat64 {
			kind = interpreter.KindFloat64
		} else {
			kind = interpreter.KindInt32
		}
	case token.DIVIDE, token.MODULUS:
		kind = interpreter.KindFloat64
	default:
		if left.Kind == interpreter.KindInt32 && right.Kind == interpreter.KindInt32 {
			kind = interpreter.KindInt32
		} else {
			kind = interpreter.KindFloat64
		}
	}
	return &Semantic{Kind: kind}
}

func (a *Analyzer) boolLiteral(_ *ast.BoolLiteral) *Semantic {
	return &Semantic{Kind: interpreter.KindBool}
}

func (a *Analyzer) numberLiteral(node *ast.NumberLiteral) *Semantic {
	kind := interpreter.KindInt32
	if strings.Contains(node.Token.Literal, ".") || strings.Contains(node.Token.Literal, "e") {
		kind = interpreter.KindFloat64
	} else if node.Value != float64(int32(node.Value)) {
		kind = interpreter.KindFloat64
	}
	return &Semantic{Kind: kind}
}

func (a *Analyzer) stringLiteral(_ *ast.StringLiteral) *Semantic {
	return &Semantic{Kind: interpreter.KindString}
}
