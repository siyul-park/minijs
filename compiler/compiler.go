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
	code      bytecode.Bytecode
	constants map[string]int
}

var castOpcode = map[types.Kind]map[types.Kind]bytecode.Opcode{
	types.KindString: {
		types.KindString:  bytecode.NOP,
		types.KindFloat64: bytecode.C2F64,
	},
	types.KindFloat64: {
		types.KindFloat64: bytecode.NOP,
		types.KindString:  bytecode.F642C,
	},
}

func New() *Compiler {
	return &Compiler{
		constants: make(map[string]int),
	}
}

func (c *Compiler) Compile(node ast.Node) (bytecode.Bytecode, error) {
	_, err := c.compile(node)
	code := c.code
	c.code = bytecode.Bytecode{}
	c.constants = make(map[string]int)
	return code, err
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
		c.emit(bytecode.F64LOAD, math.Float64bits(math.NaN()))
	case "Infinity":
		c.emit(bytecode.F64LOAD, math.Float64bits(math.Inf(1)))
	default:
		c.emit(bytecode.F64LOAD, math.Float64bits(node.Value))
	}
	return types.KindFloat64, nil
}

func (c *Compiler) string(node *ast.StringLiteral) (types.Kind, error) {
	offset, size := c.store(node.Value)
	c.emit(bytecode.CLOAD, uint64(offset), uint64(size))
	return types.KindString, nil
}

func (c *Compiler) prefixExpression(node *ast.PrefixExpression) (types.Kind, error) {
	right, err := c.compile(node.Right)
	if err != nil {
		return types.KindUnknown, err
	}

	switch node.Token.Type {
	case token.PLUS, token.MINUS:
		if _, err := c.cast(right, types.KindFloat64); err != nil {
			return types.KindUnknown, err
		}
		if node.Token.Type == token.MINUS {
			c.emit(bytecode.F64LOAD, math.Float64bits(-1))
			c.emit(bytecode.F64MUL)
		}
		return types.KindFloat64, nil
	default:
		return types.KindUnknown, fmt.Errorf("invalid token for prefix expression: %s", node.Token.Type)
	}
}

func (c *Compiler) infixExpression(node *ast.InfixExpression) (types.Kind, error) {
	left := c.kind(node.Left)
	right := c.kind(node.Right)

	if left == types.KindString || right == types.KindString {
		left, right = types.KindString, types.KindString
	}

	left, err := c.compile(node.Left)
	if err != nil {
		return types.KindUnknown, err
	}
	if left != right {
		left, err = c.cast(left, right)
		if err != nil {
			return types.KindUnknown, err
		}
	}

	right, err = c.compile(node.Right)
	if err != nil {
		return types.KindUnknown, err
	}
	if right != left {
		right, err = c.cast(right, left)
		if err != nil {
			return types.KindUnknown, err
		}
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

func (c *Compiler) cast(from, to types.Kind) (types.Kind, error) {
	opcode, ok := castOpcode[from][to]
	if !ok {
		return types.KindUnknown, fmt.Errorf("unsupported cast from %s to %s", from, to)
	}
	if opcode != bytecode.NOP {
		c.emit(opcode)
	}
	return to, nil
}

func (c *Compiler) kind(node ast.Node) types.Kind {
	p := New()
	kind, _ := p.compile(node)
	return kind
}

func (c *Compiler) emit(op bytecode.Opcode, operands ...uint64) int {
	return c.code.Add(bytecode.New(op, operands...))
}

func (c *Compiler) store(val string) (int, int) {
	offset, ok := c.constants[val]
	if !ok {
		offset = c.code.Store([]byte(val + "\x00"))
		c.constants[val] = offset
	}
	return offset, len([]byte(val))
}
