package compiler

import (
	"errors"
	"fmt"
	"math"

	"github.com/siyul-park/minijs/ast"
	"github.com/siyul-park/minijs/bytecode"
	"github.com/siyul-park/minijs/token"
	"github.com/siyul-park/minijs/types"
)

type Compiler struct {
	node ast.Node
	code bytecode.Bytecode
}

func New(node ast.Node) *Compiler {
	return &Compiler{node: node}
}

func (c *Compiler) Compile() (bytecode.Bytecode, error) {
	if err := c.compile(c.node); err != nil {
		return bytecode.Bytecode{}, err
	}
	return c.code, nil
}

func (c *Compiler) compile(node ast.Node) error {
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
		return fmt.Errorf("unsupported operand type: %T", node)
	}
}

func (c *Compiler) program(node *ast.Program) error {
	for _, n := range node.Statements {
		if err := c.statement(n); err != nil {
			return err
		}
	}
	return nil
}

func (c *Compiler) statement(node *ast.Statement) error {
	if err := c.compile(node.Node); err != nil {
		return err
	}
	c.emit(bytecode.POP)
	return nil
}

func (c *Compiler) number(node *ast.NumberLiteral) error {
	switch node.Token.Literal {
	case "NaN":
		c.emit(bytecode.F64LD, math.Float64bits(math.NaN()))
	case "Infinity":
		c.emit(bytecode.F64LD, math.Float64bits(math.Inf(1)))
	default:
		c.emit(bytecode.F64LD, math.Float64bits(node.Value))
	}
	return nil
}

func (c *Compiler) string(node *ast.StringLiteral) error {
	offset, size := c.store([]byte(node.Value))
	c.emit(bytecode.CLD, uint64(offset), uint64(size))
	return nil
}

func (c *Compiler) prefixExpression(node *ast.PrefixExpression) error {
	if err := c.compile(node.Right); err != nil {
		return err
	}
	switch node.Token.Type {
	case token.PLUS:
	case token.MINUS:
		c.emit(bytecode.F64LD, math.Float64bits(-1))
		c.emit(bytecode.F64MUL)
	default:
		return errors.New("invalid token")
	}
	return nil
}

func (c *Compiler) infixExpression(node *ast.InfixExpression) error {
	if err := c.compile(node.Left); err != nil {
		return err
	}
	if err := c.compile(node.Right); err != nil {
		return err
	}

	switch c.kind(node.Left) {
	case types.KindFloat64:
		switch node.Token.Type {
		case token.PLUS:
			c.emit(bytecode.F64ADD)
		case token.MINUS:
			c.emit(bytecode.F64SUB)
		case token.ASTERISK:
			c.emit(bytecode.F64MUL)
		case token.SLASH:
			c.emit(bytecode.F64DIV)
		case token.PERCENT:
			c.emit(bytecode.F64MOD)
		default:
			return fmt.Errorf("invalid operator for float64: %s", node.Token.Type)
		}
	default:
		return fmt.Errorf("unsupported operand type: %s", c.kind(node.Left))
	}
	return nil
}

func (c *Compiler) kind(node ast.Node) types.Kind {
	switch node := node.(type) {
	case *ast.Program, *ast.Statement:
		return types.KindVoid
	case *ast.NumberLiteral:
		return types.KindFloat64
	case *ast.StringLiteral:
		return types.KindString
	case *ast.PrefixExpression:
		return c.kind(node.Right)
	case *ast.InfixExpression:
		left := c.kind(node.Left)
		right := c.kind(node.Right)
		if left == right {
			return left
		}
		return types.KindUnknown
	}
	return types.KindUnknown
}

func (c *Compiler) emit(op bytecode.Opcode, operands ...uint64) {
	c.code.Append(bytecode.New(op, operands...))
}

func (c *Compiler) store(val []byte) (int, int) {
	return c.code.Store(val)
}
