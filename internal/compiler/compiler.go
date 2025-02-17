package compiler

import (
	"fmt"
	"math"
	"strings"

	"github.com/siyul-park/minijs/internal/ast"
	"github.com/siyul-park/minijs/internal/bytecode"
	"github.com/siyul-park/minijs/internal/interpreter"
	"github.com/siyul-park/minijs/internal/token"
)

type Compiler struct {
	code     bytecode.Bytecode
	kinds    map[ast.Node]interpreter.Kind
	literals map[string]int
}

var casts = map[interpreter.Kind]map[interpreter.Kind][]bytecode.Instruction{
	interpreter.KindBool: {
		interpreter.KindBool:    {},
		interpreter.KindInt32:   {bytecode.New(bytecode.BTOI32)},
		interpreter.KindFloat64: {bytecode.New(bytecode.BTOI32), bytecode.New(bytecode.I32TOF64)},
		interpreter.KindString:  {bytecode.New(bytecode.BTOC)},
	},
	interpreter.KindInt32: {
		interpreter.KindBool:    {bytecode.New(bytecode.I32TOB)},
		interpreter.KindInt32:   {},
		interpreter.KindFloat64: {bytecode.New(bytecode.I32TOF64)},
		interpreter.KindString:  {bytecode.New(bytecode.I32TOC)},
	},
	interpreter.KindFloat64: {
		interpreter.KindBool:    {},
		interpreter.KindInt32:   {bytecode.New(bytecode.F64TOI32)},
		interpreter.KindFloat64: {},
		interpreter.KindString:  {bytecode.New(bytecode.F64TOC)},
	},
	interpreter.KindString: {
		interpreter.KindBool:    {},
		interpreter.KindInt32:   {bytecode.New(bytecode.CTOI32)},
		interpreter.KindFloat64: {bytecode.New(bytecode.CTOF64)},
		interpreter.KindString:  {},
	},
}

func New() *Compiler {
	c := &Compiler{
		kinds:    make(map[ast.Node]interpreter.Kind),
		literals: make(map[string]int),
	}
	return c
}

func (c *Compiler) Compile(node ast.Node) (bytecode.Bytecode, error) {
	defer func() {
		c.code = bytecode.Bytecode{}
		c.kinds = make(map[ast.Node]interpreter.Kind)
		c.literals = make(map[string]int)
	}()

	err := c.compile(node)
	return c.code, err
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
	case *ast.PrefixExpression:
		return c.compilePrefixExpression(node)
	case *ast.InfixExpression:
		return c.compileInfixExpression(node)
	case *ast.BoolLiteral:
		return c.compileBoolLiteral(node)
	case *ast.NumberLiteral:
		return c.compileNumberLiteral(node)
	case *ast.StringLiteral:
		return c.compileStringLiteral(node)
	default:
		return fmt.Errorf("unsupported operand type: %T", node)
	}
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

func (c *Compiler) compilePrefixExpression(node *ast.PrefixExpression) error {
	kind := c.getKind(node)
	right := c.getKind(node.Right)

	if err := c.compile(node.Right); err != nil {
		return err
	}
	if err := c.cast(right, kind); err != nil {
		return err
	}

	switch node.Token.Type {
	case token.PLUS, token.MINUS:
		if node.Token.Type == token.MINUS {
			switch kind {
			case interpreter.KindInt32:
				c.emit(bytecode.I32LOAD, uint64(0xFFFFFFFFFFFFFFFF))
				c.emit(bytecode.I32MUL)
			case interpreter.KindFloat64:
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
	kind := c.getKind(node)
	left := c.getKind(node.Left)
	right := c.getKind(node.Right)

	if err := c.compile(node.Left); err != nil {
		return err
	}
	if err := c.cast(left, kind); err != nil {
		return err
	}

	if err := c.compile(node.Right); err != nil {
		return err
	}
	if err := c.cast(right, kind); err != nil {
		return err
	}

	switch kind {
	case interpreter.KindInt32:
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
	case interpreter.KindFloat64:
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
	case interpreter.KindString:
		switch node.Token.Type {
		case token.PLUS:
			c.emit(bytecode.CADD)
			return nil
		}
	default:
	}
	return fmt.Errorf("unsupported operator '%s' for types %v and %v", node.Token.Type, left, right)
}

func (c *Compiler) compileBoolLiteral(node *ast.BoolLiteral) error {
	value := uint64(0)
	if node.Value {
		value = 1
	}
	c.emit(bytecode.BLOAD, value)
	return nil
}

func (c *Compiler) compileNumberLiteral(node *ast.NumberLiteral) error {
	switch node.Token.Literal {
	case "NaN":
		c.emit(bytecode.F64LOAD, math.Float64bits(math.NaN()))
	case "Infinity":
		c.emit(bytecode.F64LOAD, math.Float64bits(math.Inf(1)))
	default:
		if c.getKind(node) == interpreter.KindInt32 {
			c.emit(bytecode.I32LOAD, uint64(int32(node.Value)))
		} else {
			c.emit(bytecode.F64LOAD, math.Float64bits(node.Value))
		}
	}
	return nil
}

func (c *Compiler) compileStringLiteral(node *ast.StringLiteral) error {
	offset, size := c.intern(node.Value)
	c.emit(bytecode.CLOAD, uint64(offset), uint64(size))
	return nil
}

func (c *Compiler) getKind(node ast.Expression) interpreter.Kind {
	if kind, found := c.kinds[node]; found {
		return kind
	}

	var kind interpreter.Kind
	switch node := node.(type) {
	case *ast.PrefixExpression:
		kind = c.getPrefixKind(node)
	case *ast.InfixExpression:
		kind = c.getInfixKind(node)
	case *ast.BoolLiteral:
		kind = interpreter.KindBool
	case *ast.NumberLiteral:
		kind = c.getNumberLiteralKind(node)
	case *ast.StringLiteral:
		kind = interpreter.KindString
	default:
		kind = interpreter.KindUnknown
	}

	c.kinds[node] = kind
	return kind
}

func (c *Compiler) getPrefixKind(node *ast.PrefixExpression) interpreter.Kind {
	right := c.getKind(node.Right)
	switch node.Token.Type {
	case token.PLUS, token.MINUS:
		return c.getPrefixOperatorKind(right)
	}
	return interpreter.KindUnknown
}

func (c *Compiler) getPrefixOperatorKind(right interpreter.Kind) interpreter.Kind {
	switch right {
	case interpreter.KindBool:
		return interpreter.KindInt32
	case interpreter.KindString:
		return interpreter.KindFloat64
	case interpreter.KindInt32, interpreter.KindFloat64:
		return right
	default:
		return interpreter.KindUnknown
	}
}

func (c *Compiler) getInfixKind(node *ast.InfixExpression) interpreter.Kind {
	left := c.getKind(node.Left)
	right := c.getKind(node.Right)

	if left == interpreter.KindUnknown || right == interpreter.KindUnknown {
		return interpreter.KindUnknown
	}

	switch node.Token.Type {
	case token.PLUS:
		if left == interpreter.KindString || right == interpreter.KindString {
			return interpreter.KindString
		} else if left == interpreter.KindFloat64 || right == interpreter.KindFloat64 {
			return interpreter.KindFloat64
		}
		return interpreter.KindInt32
	case token.DIVIDE, token.MODULUS:
		return interpreter.KindFloat64
	default:
		if left == interpreter.KindInt32 && right == interpreter.KindInt32 {
			return interpreter.KindInt32
		}
		return interpreter.KindFloat64
	}
}

func (c *Compiler) getNumberLiteralKind(node *ast.NumberLiteral) interpreter.Kind {
	if strings.Contains(node.Token.Literal, ".") || strings.Contains(node.Token.Literal, "e") {
		return interpreter.KindFloat64
	} else if node.Value != float64(int32(node.Value)) {
		return interpreter.KindFloat64
	}
	return interpreter.KindInt32
}

func (c *Compiler) cast(from, to interpreter.Kind) error {
	if from == to {
		return nil
	}

	queue := []struct {
		kind         interpreter.Kind
		instructions []bytecode.Instruction
	}{{from, nil}}

	visited := map[interpreter.Kind]bool{}
	visited[from] = true

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if insns := casts[curr.kind][to]; len(insns) > 0 {
			for _, insn := range append(curr.instructions, insns...) {
				c.code.Emit(insn)
			}
			return nil
		}

		for next, insns := range casts[curr.kind] {
			if !visited[next] && len(insns) > 0 {
				visited[next] = true
				queue = append(queue, struct {
					kind         interpreter.Kind
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

func (c *Compiler) emit(op bytecode.Opcode, operands ...uint64) int {
	return c.code.Emit(bytecode.New(op, operands...))
}

func (c *Compiler) intern(val string) (int, int) {
	offset, ok := c.literals[val]
	if !ok {
		offset = c.code.Store([]byte(val + "\x00"))
		c.literals[val] = offset
	}
	return offset, len([]byte(val))
}
