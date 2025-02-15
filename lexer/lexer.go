package lexer

import (
	"unicode"

	"github.com/siyul-park/minijs/token"
)

type Lexer struct {
	source []rune
	pos    int
}

func New(source string) *Lexer {
	return &Lexer{source: []rune(source)}
}

func (l *Lexer) Next() token.Token {
	l.trim()

	var tk token.Token
	switch l.peek(0) {
	case rune(0):
		tk = token.NewToken(token.EOF, "")
	case '"', '\'':
		tk = l.string(l.peek(0))
	case '+':
		tk = token.NewToken(token.PLUS, string(l.pop()))
	case '-':
		tk = token.NewToken(token.MINUS, string(l.pop()))
	case '*':
		tk = token.NewToken(token.ASTERISK, string(l.pop()))
	case '/':
		tk = token.NewToken(token.SLASH, string(l.pop()))
	case '%':
		tk = token.NewToken(token.PERCENT, string(l.pop()))
	case '.':
		tk = token.NewToken(token.DOT, string(l.pop()))
	case '(':
		tk = token.NewToken(token.LPAREN, string(l.pop()))
	case ')':
		tk = token.NewToken(token.RPAREN, string(l.pop()))
	case ';':
		tk = token.NewToken(token.SEMICOLON, string(l.pop()))
	default:
		if unicode.IsDigit(l.peek(0)) {
			tk = l.number()
		} else {
			tk = l.identifier()
		}
	}
	return tk
}

func (l *Lexer) number() token.Token {
	integer := l.integer()
	if l.peek(0) == '.' && unicode.IsDigit(l.peek(1)) {
		l.pop()
		fraction := l.integer()
		integer.Literal += "." + fraction.Literal
	}

	if l.peek(0) == 'e' || l.peek(0) == 'E' {
		l.pop()
		sign := ""
		if l.peek(0) == '+' || l.peek(0) == '-' {
			sign = string(l.pop())
		}

		if !unicode.IsDigit(l.peek(0)) {
			return token.NewToken(token.ILLEGAL, "Malformed exponent")
		}
		exponent := l.integer()
		return token.NewToken(token.NUMBER, integer.Literal+"e"+sign+exponent.Literal)
	}

	return integer
}

func (l *Lexer) integer() token.Token {
	var literal []rune
	prev := rune(0)

	if l.peek(0) == '0' && (l.peek(1) == 'b' || l.peek(1) == 'B') {
		l.pop()
		l.pop()

		for ch := l.peek(0); ch == '0' || ch == '1' || ch == '_'; ch = l.peek(0) {
			if ch == '_' {
				if prev == '_' {
					return token.NewToken(token.ILLEGAL, "Consecutive underscores")
				}
				l.pop()
				continue
			}
			literal = append(literal, l.pop())
			prev = literal[len(literal)-1]
		}

		if len(literal) == 0 {
			return token.NewToken(token.ILLEGAL, "Malformed binary")
		}
		return token.NewToken(token.NUMBER, "0b"+string(literal))
	}

	literal = append(literal, l.pop())

	for unicode.IsDigit(l.peek(0)) || l.peek(0) == '_' {
		if l.peek(0) == '_' && prev != '_' {
			l.pop()
		} else if unicode.IsDigit(l.peek(0)) {
			literal = append(literal, l.pop())
		} else {
			return token.NewToken(token.ILLEGAL, "Consecutive underscores")
		}
		prev = literal[len(literal)-1]
	}

	if len(literal) > 0 && literal[len(literal)-1] == '_' {
		return token.NewToken(token.ILLEGAL, "Trailing underscore")
	}

	return token.NewToken(token.NUMBER, string(literal))
}

func (l *Lexer) identifier() token.Token {
	var literal []rune
	for unicode.IsLetter(l.peek(0)) || unicode.IsDigit(l.peek(0)) {
		literal = append(literal, l.pop())
	}
	return token.NewToken(token.TypeOf(string(literal)), string(literal))
}

func (l *Lexer) string(delim rune) token.Token {
	if l.peek(0) != delim {
		return token.NewToken(token.ILLEGAL, "Unterminated string")
	}
	l.pop()

	var literal []rune
	for {
		ch := l.peek(0)
		if ch == rune(0) {
			return token.NewToken(token.ILLEGAL, "Unterminated string")
		}
		if ch == delim {
			l.pop()
			break
		}
		literal = append(literal, l.pop())
	}
	return token.NewToken(token.KindString, string(literal))
}

func (l *Lexer) trim() {
	for unicode.IsSpace(l.peek(0)) {
		l.pop()
	}
}

func (l *Lexer) peek(i int) rune {
	if l.pos+i >= len(l.source) {
		return rune(0)
	}
	return l.source[l.pos+i]
}

func (l *Lexer) pop() rune {
	if l.pos >= len(l.source) {
		return rune(0)
	}
	ch := l.source[l.pos]
	l.pos++
	return ch
}
