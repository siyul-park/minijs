package compiler

import (
	"bytes"
	"fmt"
	"math"
	"strings"

	"github.com/siyul-park/minijs/internal/ast"
	"github.com/siyul-park/minijs/internal/bytecode"
	"github.com/siyul-park/minijs/internal/interpreter"
	"github.com/siyul-park/minijs/internal/token"
)

type Compiler struct {
	instructions []bytecode.Instruction
	constants    [][]byte
}

var casts = map[interpreter.Type]map[interpreter.Type][]bytecode.Instruction{
	interpreter.BOOL: {
		interpreter.BOOL:    {},
		interpreter.INT32:   {bytecode.New(bytecode.BLTOI32)},
		interpreter.FLOAT64: {bytecode.New(bytecode.BLTOI32), bytecode.New(bytecode.I32TOF64)},
		interpreter.STRING:  {bytecode.New(bytecode.BLTOSTR)},
	},
	interpreter.INT32: {
		interpreter.BOOL:    {bytecode.New(bytecode.I32TOBL)},
		interpreter.INT32:   {},
		interpreter.FLOAT64: {bytecode.New(bytecode.I32TOF64)},
		interpreter.STRING:  {bytecode.New(bytecode.I32TOSTR)},
	},
	interpreter.FLOAT64: {
		interpreter.BOOL:    {},
		interpreter.INT32:   {bytecode.New(bytecode.F64TOI32)},
		interpreter.FLOAT64: {},
		interpreter.STRING:  {bytecode.New(bytecode.F64TOSTR)},
	},
	interpreter.STRING: {
		interpreter.BOOL:    {},
		interpreter.INT32:   {bytecode.New(bytecode.STRTOI32)},
		interpreter.FLOAT64: {bytecode.New(bytecode.STRTOF64)},
		interpreter.STRING:  {},
	},
}

func New() *Compiler {
	return &Compiler{}
}

func (c *Compiler) Compile(node ast.Node) (bytecode.Bytecode, error) {
	if err := c.compile(node); err != nil {
		return bytecode.Bytecode{}, err
	}
	return c.bytecode(), nil
}

func (c *Compiler) compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		return c.compileProgram(node)
	case *ast.EmptyStatement:
		return c.compileEmptyStatement(node)
	case *ast.BlockStatement:
		return c.compileBlockStatement(node)
	case *ast.ExpressionStatement:
		return c.compileExpressionStatement(node)
	case *ast.VariableStatement:
		return c.compileVariableStatement(node)
	case *ast.PrefixExpression:
		return c.compilePrefixExpression(node)
	case *ast.InfixExpression:
		return c.compileInfixExpression(node)
	case *ast.AssignmentExpression:
		return c.compileAssignmentExpression(node)
	case *ast.BoolLiteral:
		return c.compileBoolLiteral(node)
	case *ast.NumberLiteral:
		return c.compileNumberLiteral(node)
	case *ast.StringLiteral:
		return c.compileStringLiteral(node)
	case *ast.IdentifierLiteral:
		return c.compileIdentifierLiteral(node)
	default:
		return fmt.Errorf("unsupported operand type: %T", node)
	}
}

func (c *Compiler) bytecode() bytecode.Bytecode {
	code := bytecode.Bytecode{}
	for _, instruction := range c.instructions {
		code.Instructions = append(code.Instructions, instruction...)
	}
	for _, constant := range c.constants {
		code.Constants = append(code.Constants, constant...)
	}

	c.instructions = nil
	c.constants = nil
	return code
}

func (c *Compiler) compileProgram(node *ast.Program) error {
	for _, n := range node.Statements {
		if err := c.compile(n); err != nil {
			return err
		}
	}
	return nil
}

func (c *Compiler) compileEmptyStatement(_ *ast.EmptyStatement) error {
	return nil
}

func (c *Compiler) compileBlockStatement(node *ast.BlockStatement) error {
	for _, n := range node.Statements {
		if err := c.compile(n); err != nil {
			return err
		}
	}
	return nil
}

func (c *Compiler) compileExpressionStatement(node *ast.ExpressionStatement) error {
	if err := c.compile(node.Expression); err != nil {
		return err
	}
	c.emit(bytecode.POP)
	return nil
}

func (c *Compiler) compileVariableStatement(node *ast.VariableStatement) error {
	switch node.Token.Type {
	case token.VAR:
		for _, n := range node.Right {
			if err := c.compile(n); err != nil {
				return err
			}
			c.emit(bytecode.POP)
		}
		return nil
	default:
		return fmt.Errorf("invalid variable token type: %s", node.Token.Type)
	}
}

func (c *Compiler) compilePrefixExpression(node *ast.PrefixExpression) error {
	typ := c.getType(node)
	right := c.getType(node.Right)

	if err := c.compile(node.Right); err != nil {
		return err
	}
	if err := c.cast(right, typ); err != nil {
		return err
	}

	switch node.Token.Type {
	case token.PLUS, token.MINUS:
		if node.Token.Type == token.MINUS {
			switch typ {
			case interpreter.INT32:
				c.emit(bytecode.I32LOAD, uint64(0xFFFFFFFFFFFFFFFF))
				c.emit(bytecode.I32MUL)
			case interpreter.FLOAT64:
				c.emit(bytecode.F64LOAD, math.Float64bits(-1))
				c.emit(bytecode.F64MUL)
			default:
			}
		}
		return nil
	}
	return fmt.Errorf("unsupported operator '%s' for types %v", node.Token.Type, right)
}

func (c *Compiler) compileInfixExpression(node *ast.InfixExpression) error {
	typ := c.getType(node)
	left := c.getType(node.Left)
	right := c.getType(node.Right)

	if err := c.compile(node.Left); err != nil {
		return err
	}
	if err := c.cast(left, typ); err != nil {
		return err
	}

	if err := c.compile(node.Right); err != nil {
		return err
	}
	if err := c.cast(right, typ); err != nil {
		return err
	}

	switch typ {
	case interpreter.INT32:
		switch node.Token.Type {
		case token.PLUS:
			c.emit(bytecode.I32ADD)
			return nil
		case token.MINUS:
			c.emit(bytecode.I32SUB)
			return nil
		case token.MULTIPLY:
			c.emit(bytecode.I32MUL)
			return nil
		}
	case interpreter.FLOAT64:
		switch node.Token.Type {
		case token.PLUS:
			c.emit(bytecode.F64ADD)
			return nil
		case token.MINUS:
			c.emit(bytecode.F64SUB)
			return nil
		case token.MULTIPLY:
			c.emit(bytecode.F64MUL)
			return nil
		case token.DIVIDE:
			c.emit(bytecode.F64DIV)
			return nil
		case token.MODULUS:
			c.emit(bytecode.F64MOD)
			return nil
		}
	case interpreter.STRING:
		switch node.Token.Type {
		case token.PLUS:
			c.emit(bytecode.STRADD)
			return nil
		}
	default:
	}
	return fmt.Errorf("unsupported operator '%s' for types %v and %v", node.Token.Type, left, right)
}

func (c *Compiler) compileAssignmentExpression(node *ast.AssignmentExpression) error {
	if err := c.compile(node.Left); err != nil {
		return err
	}
	c.instructions = c.instructions[:len(c.instructions)-1]
	if err := c.compile(node.Right); err != nil {
		return err
	}
	c.emit(bytecode.OBJSET)
	return nil
}

func (c *Compiler) compileBoolLiteral(node *ast.BoolLiteral) error {
	value := uint64(0)
	if node.Value {
		value = 1
	}
	c.emit(bytecode.BLLOAD, value)
	return nil
}

func (c *Compiler) compileNumberLiteral(node *ast.NumberLiteral) error {
	switch node.Token.Literal {
	case "NaN":
		c.emit(bytecode.F64LOAD, math.Float64bits(math.NaN()))
	case "Infinity":
		c.emit(bytecode.F64LOAD, math.Float64bits(math.Inf(1)))
	default:
		if c.getType(node) == interpreter.INT32 {
			c.emit(bytecode.I32LOAD, uint64(int32(node.Value)))
		} else {
			c.emit(bytecode.F64LOAD, math.Float64bits(node.Value))
		}
	}
	return nil
}

func (c *Compiler) compileStringLiteral(node *ast.StringLiteral) error {
	offset, size := c.store([]byte(node.Value))
	c.emit(bytecode.STRLOAD, offset, size)
	return nil
}

func (c *Compiler) compileIdentifierLiteral(node *ast.IdentifierLiteral) error {
	c.emit(bytecode.GLBLOAD)
	offset, size := c.store([]byte(node.Value))
	c.emit(bytecode.STRLOAD, offset, size)
	c.emit(bytecode.OBJGET)
	return nil
}

func (c *Compiler) cast(from, to interpreter.Type) error {
	if from == to {
		return nil
	}

	queue := []struct {
		kind         interpreter.Type
		instructions []bytecode.Instruction
	}{{from, nil}}

	visited := map[interpreter.Type]bool{}
	visited[from] = true

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if insns := casts[curr.kind][to]; len(insns) > 0 {
			c.instructions = append(c.instructions, append(curr.instructions, insns...)...)
			return nil
		}

		for next, insns := range casts[curr.kind] {
			if !visited[next] && len(insns) > 0 {
				visited[next] = true
				queue = append(queue, struct {
					kind         interpreter.Type
					instructions []bytecode.Instruction
				}{
					kind:         next,
					instructions: append(curr.instructions, insns...),
				})
			}
		}
	}

	return fmt.Errorf("no cast path found from %v to %v", from, to)
}

func (c *Compiler) getType(node ast.Expression) interpreter.Type {
	switch node := node.(type) {
	case *ast.PrefixExpression:
		return c.getPrefixExpressionType(node)
	case *ast.InfixExpression:
		return c.getInfixExpressionType(node)
	case *ast.AssignmentExpression:
		return c.getAssignmentExpression(node)
	case *ast.BoolLiteral:
		return c.getBoolLiteralType(node)
	case *ast.NumberLiteral:
		return c.getNumberLiteralType(node)
	case *ast.StringLiteral:
		return c.getStringLiteralType(node)
	default:
		return interpreter.UNKNOWN
	}
}

func (c *Compiler) getPrefixExpressionType(node *ast.PrefixExpression) interpreter.Type {
	right := c.getType(node.Right)
	switch node.Token.Type {
	case token.PLUS, token.MINUS:
		switch right {
		case interpreter.BOOL:
			return interpreter.INT32
		case interpreter.STRING:
			return interpreter.FLOAT64
		case interpreter.INT32, interpreter.FLOAT64:
			return right
		default:
			return interpreter.UNKNOWN
		}
	}
	return interpreter.UNKNOWN
}

func (c *Compiler) getInfixExpressionType(node *ast.InfixExpression) interpreter.Type {
	left := c.getType(node.Left)
	right := c.getType(node.Right)

	if left == interpreter.UNKNOWN || right == interpreter.UNKNOWN {
		return interpreter.UNKNOWN
	}

	switch node.Token.Type {
	case token.PLUS:
		if left == interpreter.STRING || right == interpreter.STRING {
			return interpreter.STRING
		} else if left == interpreter.FLOAT64 || right == interpreter.FLOAT64 {
			return interpreter.FLOAT64
		}
		return interpreter.INT32
	case token.DIVIDE, token.MODULUS:
		return interpreter.FLOAT64
	default:
		if left == interpreter.INT32 && right == interpreter.INT32 {
			return interpreter.INT32
		}
		return interpreter.FLOAT64
	}
}

func (c *Compiler) getAssignmentExpression(node *ast.AssignmentExpression) interpreter.Type {
	return c.getType(node.Right)
}

func (c *Compiler) getBoolLiteralType(_ *ast.BoolLiteral) interpreter.Type {
	return interpreter.BOOL
}

func (c *Compiler) getNumberLiteralType(node *ast.NumberLiteral) interpreter.Type {
	if strings.Contains(node.Token.Literal, ".") || strings.Contains(node.Token.Literal, "e") {
		return interpreter.FLOAT64
	} else if node.Value != float64(int32(node.Value)) {
		return interpreter.FLOAT64
	}
	return interpreter.INT32
}

func (c *Compiler) getStringLiteralType(_ *ast.StringLiteral) interpreter.Type {
	return interpreter.STRING
}

func (c *Compiler) emit(op bytecode.Opcode, operands ...uint64) {
	c.instructions = append(c.instructions, bytecode.New(op, operands...))
}

func (c *Compiler) store(val []byte) (uint64, uint64) {
	offset := 0
	for _, v := range c.constants {
		if bytes.Equal(v[:len(v)-1], val) {
			return uint64(offset), uint64(len(v))
		}
		offset += len(v)
	}
	c.constants = append(c.constants, append(val, 0))
	return uint64(offset), uint64(len(val))
}
