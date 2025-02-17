package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/siyul-park/minijs/internal/ast"
	"github.com/siyul-park/minijs/internal/lexer"
	"github.com/siyul-park/minijs/internal/token"
)

type Parser struct {
	lexer  *lexer.Lexer
	tokens [3]token.Token
	prefix map[token.Type]func() (ast.Expression, error)
	infix  map[token.Type]func(ast.Expression) (ast.Expression, error)
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

var precedences = map[string]int{
	token.PLUS.Literal:       SUM,
	token.MINUS.Literal:      SUM,
	token.MULTIPLE.Literal:   PRODUCT,
	token.DIVIDE.Literal:     PRODUCT,
	token.MODULO.Literal:     MODULO,
	token.LEFT_PAREN.Literal: MODULO,
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{lexer: lexer, tokens: [3]token.Token{
		token.EOF,
		lexer.Next(),
		lexer.Next(),
	}}

	p.prefix = map[token.Type]func() (ast.Expression, error){
		token.NUMBER:            p.numberLiteral,
		token.STRING:            p.stringLiteral,
		token.BOOLEAN:           p.boolLiteral,
		token.IDENTIFIER:        p.identifierLiteral,
		token.PLUS.Kind():       p.prefixExpression,
		token.MINUS.Kind():      p.prefixExpression,
		token.LEFT_PAREN.Kind(): p.groupedExpression,
	}
	p.infix = map[token.Type]func(ast.Expression) (ast.Expression, error){
		token.PLUS.Kind():     p.infixExpression,
		token.MINUS.Kind():    p.infixExpression,
		token.MULTIPLE.Kind(): p.infixExpression,
		token.DIVIDE.Kind():   p.infixExpression,
		token.MODULO.Kind():   p.infixExpression,
	}

	return p
}

func (p *Parser) Parse() (*ast.Program, error) {
	program := &ast.Program{}
	for p.peek(CURR) != token.EOF {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		program.Statements = append(program.Statements, stmt)
	}
	return program, nil
}

func (p *Parser) statement() (ast.Statement, error) {
	if p.peek(CURR) == token.SEMICOLON {
		p.pop()
		return ast.NewEmptyStatement(), nil
	}
	if p.peek(CURR) == token.LEFT_BRACE {
		return p.blockStatement()
	}

	exp, err := p.expression(LOWEST)
	if err != nil {
		return nil, err
	}
	if p.peek(CURR) == token.SEMICOLON {
		p.pop()
	}
	return ast.NewExpressionStatement(exp), nil
}

func (p *Parser) expression(precedence int) (ast.Expression, error) {
	prefix, ok := p.prefix[p.peek(CURR).Kind()]
	if !ok {
		return nil, fmt.Errorf("no prefix expression function for %s", p.peek(CURR).Kind())
	}

	left, err := prefix()
	if err != nil {
		return nil, err
	}

	for p.peek(CURR) != token.SEMICOLON && precedence < p.precedence(CURR) {
		infix, ok := p.infix[p.peek(CURR).Kind()]
		if !ok {
			return left, nil
		}
		left, err = infix(left)
		if err != nil {
			return nil, err
		}
	}
	return left, nil
}

func (p *Parser) numberLiteral() (ast.Expression, error) {
	curr := p.peek(CURR)
	p.pop()

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

func (p *Parser) stringLiteral() (ast.Expression, error) {
	curr := p.peek(CURR)
	p.pop()
	return ast.NewStringLiteral(curr, curr.Literal), nil
}

func (p *Parser) boolLiteral() (ast.Expression, error) {
	curr := p.peek(CURR)
	p.pop()
	return ast.NewBoolLiteral(curr, curr.Literal == "true"), nil
}

func (p *Parser) identifierLiteral() (ast.Expression, error) {
	curr := p.peek(CURR)
	p.pop()
	return ast.NewIdentifierLiteral(curr, curr.Literal), nil
}

func (p *Parser) blockStatement() (ast.Statement, error) {
	p.pop()

	var statements []ast.Statement
	for p.peek(CURR) != token.RIGHT_BRACE {
		if node, err := p.statement(); err != nil {
			return nil, err
		} else {
			statements = append(statements, node)
		}
	}

	p.pop()
	return ast.NewBlockStatement(statements...), nil
}

func (p *Parser) prefixExpression() (ast.Expression, error) {
	curr := p.peek(CURR)
	p.pop()

	right, err := p.expression(PREFIX)
	if err != nil {
		return nil, err
	}
	return ast.NewPrefixExpression(curr, right), nil
}

func (p *Parser) infixExpression(left ast.Expression) (ast.Expression, error) {
	curr := p.peek(CURR)
	precedence := p.precedence(CURR)
	p.pop()

	right, err := p.expression(precedence)
	if err != nil {
		return nil, err
	}
	return ast.NewInfixExpression(curr, left, right), nil
}

func (p *Parser) groupedExpression() (ast.Expression, error) {
	p.pop()

	n, err := p.expression(LOWEST)
	if err != nil {
		return nil, err
	}

	if p.peek(NEXT) != token.RIGHT_PAREN {
		return nil, fmt.Errorf("expected next token to be %s, got %s instead", token.RIGHT_PAREN, p.peek(NEXT).Kind())
	}

	p.pop()
	return n, nil
}

func (p *Parser) precedence(i int) int {
	peek := p.peek(i)
	precedence, ok := precedences[peek.Literal]
	if !ok {
		return LOWEST
	}
	return precedence
}

func (p *Parser) peek(i int) token.Token {
	if i >= len(p.tokens) {
		return token.EOF
	}
	return p.tokens[i]
}

func (p *Parser) pop() {
	p.tokens[PREV] = p.tokens[CURR]
	p.tokens[CURR] = p.tokens[NEXT]
	p.tokens[NEXT] = p.lexer.Next()
}
