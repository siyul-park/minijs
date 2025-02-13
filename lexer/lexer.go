package lexer

import (
	"github.com/siyul-park/miniscript/token"
	"unicode"
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
	case '"':
		tk = l.string(l.peek(0))
	case '+':
		tk = token.NewToken(token.PLUS, string(l.peek(0)))
		l.pop(1)
	case '-':
		tk = token.NewToken(token.MINUS, string(l.peek(0)))
		l.pop(1)
	case '*':
		tk = token.NewToken(token.MULTIPLY, string(l.peek(0)))
		l.pop(1)
	case '/':
		tk = token.NewToken(token.DIVIDE, string(l.peek(0)))
		l.pop(1)
	case '.':
		tk = token.NewToken(token.PERIOD, string(l.peek(0)))
		l.pop(1)
	case '(':
		tk = token.NewToken(token.LPAREN, string(l.peek(0)))
		l.pop(1)
	case ')':
		tk = token.NewToken(token.RPAREN, string(l.peek(0)))
		l.pop(1)
	default:
		if unicode.IsDigit(l.peek(0)) {
			tk = l.decimal()
		} else {
			tk = l.identifier()
		}
	}
	return tk
}

func (l *Lexer) string(delim rune) token.Token {
	if l.peek(0) != delim {
		return token.NewToken(token.STRING, "")
	}

	var literal []rune
	for {
		l.pop(1)

		ch := l.peek(0)
		if ch == rune(0) {
			return token.NewToken(token.ILLEGAL, "")
		}
		if ch == delim {
			l.pop(1)
			break
		}

		literal = append(literal, ch)
	}
	return token.NewToken(token.STRING, string(literal))
}

func (l *Lexer) decimal() token.Token {
	integer := l.digits()
	if l.peek(0) == '.' && unicode.IsDigit(l.peek(1)) {
		l.pop(1)
		fraction := l.digits()
		return token.NewToken(token.FLOAT, integer.Literal+"."+fraction.Literal)
	}
	return integer
}

func (l *Lexer) digits() token.Token {
	var literal []rune
	for unicode.IsDigit(l.peek(0)) {
		literal = append(literal, l.peek(0))
		l.pop(1)
	}
	return token.NewToken(token.INT, string(literal))
}

func (l *Lexer) identifier() token.Token {
	var literal []rune
	for {
		ch := l.peek(0)
		if ch == rune(0) || (!unicode.IsLetter(ch) && !unicode.IsDigit(ch)) {
			break
		}

		literal = append(literal, ch)
		l.pop(1)
	}

	return token.NewToken(token.TypeOf(string(literal)), string(literal))
}

func (l *Lexer) trim() {
	for unicode.IsSpace(l.peek(0)) {
		l.pop(1)
	}
}

func (l *Lexer) peek(i int) rune {
	if l.pos+i >= len(l.source) {
		return rune(0)
	}
	return l.source[l.pos+i]
}

func (l *Lexer) pop(i int) {
	l.pos += i
}
