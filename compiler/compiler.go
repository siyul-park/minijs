package compiler

import (
	"fmt"
	"math"

	"github.com/siyul-park/minijs/ast"
	"github.com/siyul-park/minijs/bytecode"
	"github.com/siyul-park/minijs/interpreter"
	"github.com/siyul-park/minijs/token"
)

type Compiler struct {
	code      bytecode.Bytecode
	constants map[string]int
}

var casts = map[interpreter.Kind]map[interpreter.Kind][]bytecode.Instruction{
	interpreter.KindInt32: {
		interpreter.KindInt32:   []bytecode.Instruction{},
		interpreter.KindFloat64: []bytecode.Instruction{bytecode.New(bytecode.I322F64)},
		interpreter.KindString:  []bytecode.Instruction{bytecode.New(bytecode.I322C)},
	},
	interpreter.KindFloat64: {
		interpreter.KindInt32:   []bytecode.Instruction{bytecode.New(bytecode.F64I32)},
		interpreter.KindFloat64: []bytecode.Instruction{},
		interpreter.KindString:  []bytecode.Instruction{bytecode.New(bytecode.F642C)},
	},
	interpreter.KindString: {
		interpreter.KindInt32:   []bytecode.Instruction{bytecode.New(bytecode.C2I32)},
		interpreter.KindFloat64: []bytecode.Instruction{bytecode.New(bytecode.C2F64)},
		interpreter.KindString:  []bytecode.Instruction{},
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

func (c *Compiler) compile(node ast.Node) (interpreter.Kind, error) {
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
		return interpreter.KindInvalid, fmt.Errorf("unsupported operand type: %T", node)
	}
}

func (c *Compiler) program(node *ast.Program) (interpreter.Kind, error) {
	for _, n := range node.Statements {
		if _, err := c.statement(n); err != nil {
			return interpreter.KindInvalid, err
		}
	}
	return interpreter.KindVoid, nil
}

func (c *Compiler) statement(node *ast.Statement) (interpreter.Kind, error) {
	if _, err := c.compile(node.Node); err != nil {
		return interpreter.KindInvalid, err
	}
	c.emit(bytecode.POP)
	return interpreter.KindVoid, nil
}

func (c *Compiler) number(node *ast.NumberLiteral) (interpreter.Kind, error) {
	switch node.Token.Literal {
	case "NaN":
		c.emit(bytecode.F64LOAD, math.Float64bits(math.NaN()))
	case "Infinity":
		c.emit(bytecode.F64LOAD, math.Float64bits(math.Inf(1)))
	default:
		if node.IsInteger() {
			c.emit(bytecode.I32LOAD, uint64(int32(node.Value)))
			return interpreter.KindInt32, nil
		}
		c.emit(bytecode.F64LOAD, math.Float64bits(node.Value))
	}
	return interpreter.KindFloat64, nil
}

func (c *Compiler) string(node *ast.StringLiteral) (interpreter.Kind, error) {
	offset, size := c.store(node.Value)
	c.emit(bytecode.CLOAD, uint64(offset), uint64(size))
	return interpreter.KindString, nil
}

func (c *Compiler) prefixExpression(node *ast.PrefixExpression) (interpreter.Kind, error) {
	right, err := c.compile(node.Right)
	if err != nil {
		return interpreter.KindInvalid, err
	}

	switch node.Token.Type {
	case token.PLUS, token.MINUS:
		if right != interpreter.KindInt32 && right != interpreter.KindFloat64 {
			if right, err = c.cast(right, interpreter.KindFloat64); err != nil {
				return interpreter.KindInvalid, err
			}
		}
		if node.Token.Type == token.MINUS {
			switch right {
			case interpreter.KindInt32:
				if node.Token.Type == token.MINUS {
					c.emit(bytecode.I32LOAD, uint64(0xFFFFFFFFFFFFFFFF))
					c.emit(bytecode.I32MUL)
				}
			case interpreter.KindFloat64:
				if node.Token.Type == token.MINUS {
					c.emit(bytecode.F64LOAD, math.Float64bits(-1))
					c.emit(bytecode.F64MUL)
				}
			}
		}
		return right, nil
	default:
		return interpreter.KindInvalid, fmt.Errorf("invalid token for prefix expression: %s", node.Token.Type)
	}
}

func (c *Compiler) infixExpression(node *ast.InfixExpression) (interpreter.Kind, error) {
	left := c.kind(node.Left)
	right := c.kind(node.Right)

	if left == interpreter.KindString || right == interpreter.KindString {
		left, right = interpreter.KindString, interpreter.KindString
	}

	if left == interpreter.KindFloat64 || right == interpreter.KindFloat64 || node.Token.Type == token.SLASH || node.Token.Type == token.PERCENT {
		left, right = interpreter.KindFloat64, interpreter.KindFloat64
	}

	left, err := c.compile(node.Left)
	if err != nil {
		return interpreter.KindInvalid, err
	}

	if left != right {
		if left, err = c.cast(left, right); err != nil {
			return interpreter.KindInvalid, err
		}
	}

	right, err = c.compile(node.Right)
	if err != nil {
		return interpreter.KindInvalid, err
	}

	if right != left {
		if right, err = c.cast(right, left); err != nil {
			return interpreter.KindInvalid, err
		}
	}

	if left == interpreter.KindInt32 && right == interpreter.KindInt32 {
		switch node.Token.Type {
		case token.PLUS:
			c.emit(bytecode.I32ADD)
			return interpreter.KindInt32, nil
		case token.MINUS:
			c.emit(bytecode.I32SUB)
			return interpreter.KindInt32, nil
		case token.ASTERISK:
			c.emit(bytecode.I32MUL)
			return interpreter.KindInt32, nil
		default:
			return interpreter.KindInvalid, fmt.Errorf("unsupported operator for int32: %s", node.Token.Type)
		}
	}

	if left == interpreter.KindFloat64 && right == interpreter.KindFloat64 {
		switch node.Token.Type {
		case token.PLUS:
			c.emit(bytecode.F64ADD)
			return interpreter.KindFloat64, nil
		case token.MINUS:
			c.emit(bytecode.F64SUB)
			return interpreter.KindFloat64, nil
		case token.ASTERISK:
			c.emit(bytecode.F64MUL)
			return interpreter.KindFloat64, nil
		case token.SLASH:
			c.emit(bytecode.F64DIV)
			return interpreter.KindFloat64, nil
		case token.PERCENT:
			c.emit(bytecode.F64MOD)
			return interpreter.KindFloat64, nil
		default:
			return interpreter.KindInvalid, fmt.Errorf("unsupported operator for float64: %s", node.Token.Type)
		}
	}

	if left == interpreter.KindString && right == interpreter.KindString {
		switch node.Token.Type {
		case token.PLUS:
			c.emit(bytecode.CADD)
			return interpreter.KindString, nil
		default:
			return interpreter.KindInvalid, fmt.Errorf("unsupported operator for string: %s", node.Token.Type)
		}
	}

	return interpreter.KindInvalid, fmt.Errorf("unsupported operator for types %s and %s: %s", left, right, node.Token.Type)
}

func (c *Compiler) cast(from, to interpreter.Kind) (interpreter.Kind, error) {
	queue := []interpreter.Kind{from}
	visited := map[interpreter.Kind]bool{}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if insns, ok := casts[curr][to]; ok {
			for _, insn := range insns {
				c.code.Add(insn)
			}
			return to, nil
		}

		if visited[curr] {
			continue
		}
		visited[curr] = true

		for next, insns := range casts[curr] {
			if !visited[next] {
				for _, insn := range insns {
					c.code.Add(insn)
				}
				queue = append(queue, next)
			}
		}
	}
	return interpreter.KindInvalid, fmt.Errorf("no cast path found from %s to %s", from, to)
}

func (c *Compiler) kind(node ast.Node) interpreter.Kind {
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
