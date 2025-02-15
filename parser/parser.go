package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/siyul-park/minijs/ast"
	"github.com/siyul-park/minijs/lexer"
	"github.com/siyul-park/minijs/token"
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
	token.PLUS:       SUM,
	token.MINUS:      SUM,
	token.MULTIPLE:   PRODUCT,
	token.DIVIDE:     PRODUCT,
	token.MODULAR:    MODULO,
	token.PAREN_OPEN: CALL,
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{lexer: lexer, tokens: [3]token.Token{
		token.NewToken(token.EOF, ""),
		lexer.Next(),
		lexer.Next(),
	}}

	p.prefix = map[token.Type]func() (ast.Node, error){
		token.NUMBER:     p.numberLiteral,
		token.STRING:     p.stringLiteral,
		token.BOOLEAN:    p.boolLiteral,
		token.IDENTIFIER: p.identifierLiteral,
		token.PLUS:       p.prefixExpression,
		token.MINUS:      p.prefixExpression,
		token.PAREN_OPEN: p.groupedExpression,
	}
	p.infix = map[token.Type]func(ast.Node) (ast.Node, error){
		token.PLUS:     p.infixExpression,
		token.MINUS:    p.infixExpression,
		token.MULTIPLE: p.infixExpression,
		token.DIVIDE:   p.infixExpression,
		token.MODULAR:  p.infixExpression,
	}

	return p
}

func (p *Parser) Parse() (*ast.Program, error) {
	program := &ast.Program{}
	for p.peek(CURR).Type != token.EOF {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		program.Statements = append(program.Statements, stmt)
		p.pop()
	}
	return program, nil
}

func (p *Parser) statement() (*ast.Statement, error) {
	if p.peek(CURR).Type == token.SEMICOLON {
		p.pop()
		return ast.NewStatement(nil), nil
	}
	exp, err := p.expression(LOWEST)
	if err != nil {
		return nil, err
	}
	if p.peek(NEXT).Type == token.SEMICOLON {
		p.pop()
	}
	return ast.NewStatement(exp), nil
}

func (p *Parser) expression(precedence int) (ast.Node, error) {
	prefix, ok := p.prefix[p.peek(CURR).Type]
	if !ok {
		return nil, fmt.Errorf("no prefix expression function for %s", p.peek(CURR).Type)
	}
	left, err := prefix()
	if err != nil {
		return nil, err
	}

	for p.peek(CURR).Type != token.SEMICOLON && precedence < p.precedence(NEXT) {
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

	if curr.Literal == "NaN" || curr.Literal == "Infinity" {
		return ast.NewNumberLiteral(curr, 0), nil
	}

	lit := curr.Literal
	base := 10
	if strings.HasPrefix(lit, "0b") || strings.HasPrefix(lit, "0B") {
		base = 2
		lit = lit[2:]
	} else if strings.HasPrefix(lit, "0o") || strings.HasPrefix(lit, "0O") { // 8진수 (0o)
		base = 8
		lit = lit[2:]
	} else if strings.HasPrefix(lit, "0x") || strings.HasPrefix(lit, "0X") { // 16진수 (0x)
		base = 16
		lit = lit[2:]
	}

	var value float64
	if base == 10 {
		parsedValue, err := strconv.ParseFloat(lit, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid number literal: %s", curr.Literal)
		}
		value = parsedValue
	} else {
		parsedValue, err := strconv.ParseInt(lit, base, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid %d-based literal: %s", base, curr.Literal)
		}
		value = float64(parsedValue)
	}
	return ast.NewNumberLiteral(curr, value), nil
}

func (p *Parser) stringLiteral() (ast.Node, error) {
	curr := p.peek(CURR)
	return ast.NewStringLiteral(curr, curr.Literal), nil
}

func (p *Parser) boolLiteral() (ast.Node, error) {
	curr := p.peek(CURR)
	return ast.NewBoolLiteral(curr, curr.Literal == "true"), nil
}

func (p *Parser) identifierLiteral() (ast.Node, error) {
	curr := p.peek(CURR)
	return ast.NewIdentifierLiteral(curr, curr.Literal), nil
}

func (p *Parser) prefixExpression() (ast.Node, error) {
	curr := p.peek(CURR)

	p.pop()

	right, err := p.expression(PREFIX)
	if err != nil {
		return nil, err
	}
	return ast.NewPrefixExpression(curr, right), nil
}

func (p *Parser) infixExpression(left ast.Node) (ast.Node, error) {
	curr := p.peek(CURR)
	precedence := p.precedence(CURR)

	p.pop()

	right, err := p.expression(precedence)
	if err != nil {
		return nil, err
	}
	return ast.NewInfixExpression(curr, left, right), nil
}

func (p *Parser) groupedExpression() (ast.Node, error) {
	p.pop()

	n, err := p.expression(LOWEST)
	if err != nil {
		return nil, err
	}

	if err := p.assert(NEXT, token.PAREN_CLOSE); err != nil {
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
