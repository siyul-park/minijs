package compiler

import (
	"fmt"
	"math"

	"github.com/siyul-park/minijs/internal/ast"
	"github.com/siyul-park/minijs/internal/bytecode"
	"github.com/siyul-park/minijs/internal/interpreter"
	"github.com/siyul-park/minijs/internal/token"
)

type Compiler struct {
	analyzer *Analyzer
	context  *compileContext
}

type compileContext struct {
	code     bytecode.Bytecode
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
		analyzer: NewAnalyzer(),
		context:  &compileContext{},
	}
	c.context.reset()
	return c
}

func (c *Compiler) Compile(node ast.Node) (bytecode.Bytecode, error) {
	defer c.analyzer.Close()
	defer c.context.reset()

	err := c.compile(node)
	code := c.context.code
	return code, err
}

func (c *Compiler) compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		return c.program(node)
	case *ast.EmptyStatement:
		return c.emptyStatement(node)
	case *ast.BlockStatement:
		return c.blockStatement(node)
	case *ast.ExpressionStatement:
		return c.expressionStatement(node)
	case *ast.PrefixExpression:
		return c.prefixExpression(node)
	case *ast.InfixExpression:
		return c.infixExpression(node)
	case *ast.BoolLiteral:
		return c.boolLiteral(node)
	case *ast.NumberLiteral:
		return c.numberLiteral(node)
	case *ast.StringLiteral:
		return c.stringLiteral(node)
	default:
		return fmt.Errorf("unsupported operand type: %T", node)
	}
}

func (c *Compiler) program(node *ast.Program) error {
	for _, n := range node.Statements {
		if err := c.compile(n); err != nil {
			return err
		}
	}
	return nil
}

func (c *Compiler) emptyStatement(_ *ast.EmptyStatement) error {
	return nil
}

func (c *Compiler) blockStatement(node *ast.BlockStatement) error {
	for _, n := range node.Statements {
		if err := c.compile(n); err != nil {
			return err
		}
	}
	return nil
}

func (c *Compiler) expressionStatement(node *ast.ExpressionStatement) error {
	if err := c.compile(node.Expression); err != nil {
		return err
	}
	c.emit(bytecode.POP)
	return nil
}

func (c *Compiler) prefixExpression(node *ast.PrefixExpression) error {
	meta := c.analyzer.Analyze(node)
	right := c.analyzer.Analyze(node.Right)

	if err := c.compile(node.Right); err != nil {
		return err
	}
	if err := c.cast(right.Kind, meta.Kind); err != nil {
		return err
	}

	switch node.Token.Type {
	case token.PLUS, token.MINUS:
		if node.Token.Type == token.MINUS {
			switch meta.Kind {
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
	return fmt.Errorf("unsupported operator '%s' for types %v", node.Token.Type, right.Kind)
}

func (c *Compiler) infixExpression(node *ast.InfixExpression) error {
	meta := c.analyzer.Analyze(node)
	left := c.analyzer.Analyze(node.Left)
	right := c.analyzer.Analyze(node.Right)

	if err := c.compile(node.Left); err != nil {
		return err
	}
	if err := c.cast(left.Kind, meta.Kind); err != nil {
		return err
	}

	if err := c.compile(node.Right); err != nil {
		return err
	}
	if err := c.cast(right.Kind, meta.Kind); err != nil {
		return err
	}

	switch meta.Kind {
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
	return fmt.Errorf("unsupported operator '%s' for types %v and %v", node.Token.Type, left.Kind, right.Kind)
}

func (c *Compiler) boolLiteral(node *ast.BoolLiteral) error {
	value := uint64(0)
	if node.Value {
		value = 1
	}
	c.emit(bytecode.BLOAD, value)
	return nil
}

func (c *Compiler) numberLiteral(node *ast.NumberLiteral) error {
	switch node.Token.Literal {
	case "NaN":
		c.emit(bytecode.F64LOAD, math.Float64bits(math.NaN()))
	case "Infinity":
		c.emit(bytecode.F64LOAD, math.Float64bits(math.Inf(1)))
	default:
		meta := c.analyzer.Analyze(node)
		if meta.Kind == interpreter.KindInt32 {
			c.emit(bytecode.I32LOAD, uint64(int32(node.Value)))
		} else {
			c.emit(bytecode.F64LOAD, math.Float64bits(node.Value))
		}
	}
	return nil
}

func (c *Compiler) stringLiteral(node *ast.StringLiteral) error {
	offset, size := c.store(node.Value)
	c.emit(bytecode.CLOAD, uint64(offset), uint64(size))
	return nil
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
				c.context.code.Add(insn)
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
	return c.context.code.Add(bytecode.New(op, operands...))
}

func (c *Compiler) store(val string) (int, int) {
	offset, ok := c.context.literals[val]
	if !ok {
		offset = c.context.code.Store([]byte(val + "\x00"))
		c.context.literals[val] = offset
	}
	return offset, len([]byte(val))
}

func (c *compileContext) reset() {
	c.literals = make(map[string]int)
	c.code = bytecode.Bytecode{}
}
