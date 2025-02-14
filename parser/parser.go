package parser

import (
	"fmt"
	"github.com/siyul-park/minijs/ast"
	"github.com/siyul-park/minijs/lexer"
	"github.com/siyul-park/minijs/token"
	"strconv"
)

type Parser struct {
	lexer  *lexer.Lexer
	tokens [3]token.Token
	prefix map[token.Type]func() (ast.Node, error)
	infix  map[token.Type]func(ast.Node) (ast.Node, error)
}

const (
	PREV = iota
	CURR
	NEXT
)

const (
	_ int = iota
	LOWEST
	SUM
	PRODUCT
	MODULO
	PREFIX
	CALL
	HIGHEST
)

// each token precedence
var precedences = map[token.Type]int{
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.MULTIPLY: PRODUCT,
	token.DIVIDE:   PRODUCT,
	token.MODULO:   MODULO,
	token.LPAREN:   CALL,
	token.PERIOD:   CALL,
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{lexer: lexer, tokens: [3]token.Token{
		token.NewToken(token.EOF, ""),
		lexer.Next(),
		lexer.Next(),
	}}

	p.prefix = map[token.Type]func() (ast.Node, error){
		token.DECIMAL:     p.numberLiteral,
		token.EXPONENTIAL: p.numberLiteral,
		token.BINARY:      p.numberLiteral,
		token.NAN:         p.numberLiteral,
		token.INFINITY:    p.numberLiteral,
		token.STRING:      p.stringLiteral,
		token.TRUE:        p.boolLiteral,
		token.FALSE:       p.boolLiteral,
		token.IDENTIFIER:  p.identifierLiteral,
		token.PLUS:        p.prefixExpression,
		token.MINUS:       p.prefixExpression,
		token.LPAREN:      p.groupedExpression,
	}
	p.infix = map[token.Type]func(ast.Node) (ast.Node, error){
		token.PLUS:     p.infixExpression,
		token.MINUS:    p.infixExpression,
		token.MULTIPLY: p.infixExpression,
		token.DIVIDE:   p.infixExpression,
		token.MODULO:   p.infixExpression,
	}

	return p
}

func (p *Parser) Parse() (*ast.Program, error) {
	program := &ast.Program{}
	for p.peek(CURR).Type != token.EOF {
		node, err := p.parse(LOWEST)
		if err != nil {
			return nil, err
		}
		program.Nodes = append(program.Nodes, node)
		p.pop()
	}
	return program, nil
}

func (p *Parser) parse(precedence int) (ast.Node, error) {
	prefix, ok := p.prefix[p.peek(CURR).Type]
	if !ok {
		return nil, fmt.Errorf("no prefix parse function for %s", p.peek(CURR).Type)
	}
	left, err := prefix()
	if err != nil {
		return nil, err
	}

	for precedence < p.precedence(NEXT) {
		infix, ok := p.infix[p.peek(NEXT).Type]
		if !ok {
			return left, nil
		}

		p.pop()
		left, err = infix(left)
		if err != nil {
			return nil, err
		}
	}
	return left, nil
}

func (p *Parser) numberLiteral() (ast.Node, error) {
	curr := p.peek(CURR)

	switch curr.Type {
	case token.BINARY:
		binaryValue := curr.Literal[2:]
		value, err := strconv.ParseInt(binaryValue, 2, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid binary literal: %s", curr.Literal)
		}
		return &ast.NumberLiteral{
			Token: curr,
			Value: float64(value),
		}, nil
	case token.DECIMAL, token.EXPONENTIAL:
		value, err := strconv.ParseFloat(curr.Literal, 64)
		if err != nil {
			return nil, err
		}
		return &ast.NumberLiteral{
			Token: curr,
			Value: value,
		}, nil
	case token.NAN, token.INFINITY:
		return &ast.NumberLiteral{Token: curr}, nil
	default:
		return nil, fmt.Errorf("unexpected token type for number literal: %s", curr.Type)
	}
}

func (p *Parser) stringLiteral() (ast.Node, error) {
	curr := p.peek(CURR)
	return &ast.StringLiteral{
		Token: curr,
		Value: curr.Literal,
	}, nil
}

func (p *Parser) boolLiteral() (ast.Node, error) {
	curr := p.peek(CURR)
	return &ast.BoolLiteral{
		Token: curr,
		Value: curr.Literal == "true",
	}, nil
}

func (p *Parser) identifierLiteral() (ast.Node, error) {
	curr := p.peek(CURR)
	return &ast.IdentifierLiteral{
		Token: curr,
		Value: curr.Literal,
	}, nil
}

func (p *Parser) prefixExpression() (ast.Node, error) {
	curr := p.peek(CURR)

	p.pop()
	right, err := p.parse(PREFIX)
	if err != nil {
		return nil, err
	}

	return &ast.PrefixExpression{
		Token: curr,
		Right: right,
	}, nil
}

func (p *Parser) infixExpression(left ast.Node) (ast.Node, error) {
	curr := p.peek(CURR)
	precedence := p.precedence(CURR)

	p.pop()
	right, err := p.parse(precedence)
	if err != nil {
		return nil, err
	}

	return &ast.InfixExpression{
		Token: curr,
		Left:  left,
		Right: right,
	}, nil
}

func (p *Parser) groupedExpression() (ast.Node, error) {
	p.pop()
	n, err := p.parse(LOWEST)
	if err != nil {
		return nil, err
	}
	if err := p.assert(NEXT, token.RPAREN); err != nil {
		return nil, err
	}
	p.pop()
	return n, nil
}

func (p *Parser) assert(i int, typ token.Type) error {
	if p.peek(i).Type != typ {
		return fmt.Errorf("expected next token to be %s, got %s instead", typ, p.peek(i).Type)
	}
	return nil
}

func (p *Parser) precedence(i int) int {
	peek := p.peek(i)
	precedence, ok := precedences[peek.Type]
	if !ok {
		return LOWEST
	}
	return precedence
}

func (p *Parser) peek(i int) token.Token {
	if i >= len(p.tokens) {
		return token.NewToken(token.EOF, "")
	}
	return p.tokens[i]
}

func (p *Parser) pop() {
	p.tokens[PREV] = p.tokens[CURR]
	p.tokens[CURR] = p.tokens[NEXT]
	p.tokens[NEXT] = p.lexer.Next()
}
