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
	ASSIGN
	SUM
	PRODUCT
	MODULUS
	PREFIX
	CALL
	HIGHEST
)

var precedences = map[token.Type]int{
	token.ASSIGN:     ASSIGN,
	token.PLUS:       SUM,
	token.MINUS:      SUM,
	token.MULTIPLY:   PRODUCT,
	token.DIVIDE:     PRODUCT,
	token.MODULUS:    MODULUS,
	token.OPEN_PAREN: MODULUS,
}

func New(lexer *lexer.Lexer) *Parser {
	p := &Parser{
		lexer: lexer,
		tokens: [3]token.Token{
			token.New(token.EOF, ""),
			lexer.Next(),
			lexer.Next(),
		},
	}
	p.prefix = map[token.Type]func() (ast.Expression, error){
		token.NUMBER:     p.numberLiteral,
		token.STRING:     p.stringLiteral,
		token.TRUE:       p.boolLiteral,
		token.FALSE:      p.boolLiteral,
		token.IDENTIFIER: p.identifierLiteral,
		token.PLUS:       p.prefixExpression,
		token.MINUS:      p.prefixExpression,
		token.OPEN_PAREN: p.groupedExpression,
	}
	p.infix = map[token.Type]func(ast.Expression) (ast.Expression, error){
		token.PLUS:     p.infixExpression,
		token.MINUS:    p.infixExpression,
		token.MULTIPLY: p.infixExpression,
		token.DIVIDE:   p.infixExpression,
		token.MODULUS:  p.infixExpression,
		token.ASSIGN:   p.assignmentExpression,
	}
	return p
}

func (p *Parser) Parse() (*ast.Program, error) {
	var statements []ast.Statement
	for p.peek(CURR).Type != token.EOF {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}
	return ast.NewProgram(statements...), nil
}

func (p *Parser) statement() (ast.Statement, error) {
	switch p.peek(CURR).Type {
	case token.SEMICOLON:
		return p.emptyStatement()
	case token.OPEN_BRACE:
		return p.blockStatement()
	case token.VAR:
		return p.variableStatement()
	default:
		return p.expressionStatement()
	}
}

func (p *Parser) expression(precedence int) (ast.Expression, error) {
	prefix, ok := p.prefix[p.peek(CURR).Type]
	if !ok {
		return nil, fmt.Errorf("no prefix expression function for %s", p.peek(CURR).Type)
	}

	left, err := prefix()
	if err != nil {
		return nil, err
	}

	for p.precedence(CURR) > precedence {
		infix, ok := p.infix[p.peek(CURR).Type]
		if !ok {
			break
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

	lit := curr.Literal
	base := 10
	if strings.HasPrefix(lit, "0b") || strings.HasPrefix(lit, "0B") {
		base = 2
		lit = lit[2:]
	} else if strings.HasPrefix(lit, "0o") || strings.HasPrefix(lit, "0O") {
		base = 8
		lit = lit[2:]
	} else if strings.HasPrefix(lit, "0x") || strings.HasPrefix(lit, "0X") {
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

func (p *Parser) emptyStatement() (ast.Statement, error) {
	p.pop()
	return ast.NewEmptyStatement(), nil
}

func (p *Parser) blockStatement() (ast.Statement, error) {
	p.pop()

	var statements []ast.Statement
	for p.peek(CURR).Type != token.CLOSE_BRACE {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}

	p.pop()
	return ast.NewBlockStatement(statements...), nil
}

func (p *Parser) expressionStatement() (ast.Statement, error) {
	exp, err := p.expression(LOWEST)
	if err != nil {
		return nil, err
	}
	if p.peek(CURR).Type == token.SEMICOLON {
		p.pop()
	}
	return ast.NewExpressionStatement(exp), nil
}

func (p *Parser) variableStatement() (ast.Statement, error) {
	curr := p.peek(CURR)
	p.pop()

	var expressions []*ast.AssignmentExpression
	for {
		exp, err := p.expression(LOWEST)
		if err != nil {
			return nil, err
		}
		right, ok := exp.(*ast.AssignmentExpression)
		if !ok {
			return nil, fmt.Errorf("expected assignment expressions, got %s", p.peek(CURR).Literal)
		}
		expressions = append(expressions, right)

		if p.peek(CURR).Type != token.COMMA {
			break
		}
		p.pop()
	}
	return ast.NewVariableStatement(curr, expressions...), nil
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

	if p.peek(CURR).Type != token.CLOSE_PAREN {
		return nil, fmt.Errorf("expected next token to be %s, got %s instead", token.CLOSE_PAREN, p.peek(CURR).Type)
	}
	p.pop()
	return n, nil
}

func (p *Parser) assignmentExpression(left ast.Expression) (ast.Expression, error) {
	curr := p.peek(CURR)
	p.pop()

	right, err := p.expression(LOWEST)
	if err != nil {
		return nil, err
	}
	return ast.NewAssignmentExpression(curr, left, right), nil
}

func (p *Parser) precedence(i int) int {
	peek := p.peek(i)
	if precedence, ok := precedences[peek.Type]; ok {
		return precedence
	}
	return LOWEST
}

func (p *Parser) peek(i int) token.Token {
	if i >= len(p.tokens) {
		return token.New(token.EOF, "")
	}
	return p.tokens[i]
}

func (p *Parser) pop() {
	p.tokens[PREV] = p.tokens[CURR]
	p.tokens[CURR] = p.tokens[NEXT]
	p.tokens[NEXT] = p.lexer.Next()
}
