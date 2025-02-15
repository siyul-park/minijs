package compiler

import (
	"fmt"
	"math"

	"github.com/siyul-park/minijs/ast"
	"github.com/siyul-park/minijs/bytecode"
	"github.com/siyul-park/minijs/token"
	"github.com/siyul-park/minijs/types"
)

type Compiler struct {
	node      ast.Node
	code      bytecode.Bytecode
	constants map[string]int
}

func New(node ast.Node) *Compiler {
	return &Compiler{
		node:      node,
		constants: make(map[string]int),
	}
}

func (c *Compiler) Compile() (bytecode.Bytecode, error) {
	if _, err := c.compile(c.node); err != nil {
		return bytecode.Bytecode{}, err
	}
	return c.code, nil
}

func (c *Compiler) compile(node ast.Node) (types.Kind, error) {
	switch node := node.(type) {
	case *ast.Program:
		return c.program(node)
	case *ast.Statement:
		return c.statement(node)
	case *ast.NumberLiteral:
		return c.number(node)
	case *ast.StringLiteral:
		return c.string(node)
	case *ast.PrefixExpression:
		return c.prefixExpression(node)
	case *ast.InfixExpression:
		return c.infixExpression(node)
	default:
		return types.KindUnknown, fmt.Errorf("unsupported operand type: %T", node)
	}
}

func (c *Compiler) program(node *ast.Program) (types.Kind, error) {
	for _, n := range node.Statements {
		if _, err := c.statement(n); err != nil {
			return types.KindUnknown, err
		}
	}
	return types.KindVoid, nil
}

func (c *Compiler) statement(node *ast.Statement) (types.Kind, error) {
	if _, err := c.compile(node.Node); err != nil {
		return types.KindUnknown, err
	}
	c.emit(bytecode.POP)
	return types.KindVoid, nil
}

func (c *Compiler) number(node *ast.NumberLiteral) (types.Kind, error) {
	switch node.Token.Literal {
	case "NaN":
		c.emit(bytecode.F64LD, math.Float64bits(math.NaN()))
	case "Infinity":
		c.emit(bytecode.F64LD, math.Float64bits(math.Inf(1)))
	default:
		c.emit(bytecode.F64LD, math.Float64bits(node.Value))
	}
	return types.KindFloat64, nil
}

func (c *Compiler) string(node *ast.StringLiteral) (types.Kind, error) {
	offset := c.store(node.Value)
	c.emit(bytecode.CLD, uint64(offset), uint64(len(node.Value)))
	return types.KindString, nil
}

func (c *Compiler) prefixExpression(node *ast.PrefixExpression) (types.Kind, error) {
	right, err := c.compile(node.Right)
	if err != nil {
		return types.KindUnknown, err
	}
	if right == types.KindFloat64 {
		switch node.Token.Type {
		case token.PLUS:
			return types.KindFloat64, nil
		case token.MINUS:
			c.emit(bytecode.F64LD, math.Float64bits(-1))
			c.emit(bytecode.F64MUL)
			return types.KindFloat64, nil
		default:
			return types.KindUnknown, fmt.Errorf("invalid token for prefix expression: %s", node.Token.Type)
		}
	}
	return types.KindUnknown, fmt.Errorf("invalid operand type for prefix expression: %s", right)
}

func (c *Compiler) infixExpression(node *ast.InfixExpression) (types.Kind, error) {
	left, err := c.compile(node.Left)
	if err != nil {
		return types.KindUnknown, err
	}
	right, err := c.compile(node.Right)
	if err != nil {
		return types.KindUnknown, err
	}

	if left == types.KindFloat64 && right == types.KindFloat64 {
		switch node.Token.Type {
		case token.PLUS:
			c.emit(bytecode.F64ADD)
			return types.KindFloat64, nil
		case token.MINUS:
			c.emit(bytecode.F64SUB)
			return types.KindFloat64, nil
		case token.ASTERISK:
			c.emit(bytecode.F64MUL)
			return types.KindFloat64, nil
		case token.SLASH:
			c.emit(bytecode.F64DIV)
			return types.KindFloat64, nil
		case token.PERCENT:
			c.emit(bytecode.F64MOD)
			return types.KindFloat64, nil
		default:
			return types.KindUnknown, fmt.Errorf("unsupported operator for float64: %s", node.Token.Type)
		}
	}
	if left == types.KindString && right == types.KindString {
		switch node.Token.Type {
		case token.PLUS:
			c.emit(bytecode.CADD)
			return types.KindString, nil
		default:
			return types.KindUnknown, fmt.Errorf("unsupported operator for string: %s", node.Token.Type)
		}
	}
	return types.KindUnknown, fmt.Errorf("unsupported operator for types %s and %s: %s", left, right, node.Token.Type)
}

func (c *Compiler) emit(op bytecode.Opcode, operands ...uint64) int {
	return c.code.Add(bytecode.New(op, operands...))
}

func (c *Compiler) store(val string) int {
	offset, ok := c.constants[val]
	if !ok {
		offset = c.code.Store([]byte(val))
		c.constants[val] = offset
	}
	return offset
}
