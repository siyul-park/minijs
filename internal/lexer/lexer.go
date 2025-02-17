package lexer

import (
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/siyul-park/minijs/internal/token"
)

type Lexer struct {
	source io.Reader
	buf    []rune
	pos    int
	line   int
	column int
}

func New(source io.Reader) *Lexer {
	return &Lexer{
		source: source,
		line:   1,
		column: 1,
	}
}

func (l *Lexer) Next() token.Token {
	l.hidden()

	var tk token.Token
	switch ch := l.peek(0); ch {
	case rune(0):
		tk = token.New(token.EOF, "")
	case '"', '\'':
		tk = l.string()
	case '[':
		tk = token.New(token.OPEN_BRACKET, l.read(1))
	case ']':
		tk = token.New(token.CLOSE_BRACKET, l.read(1))
	case '(':
		tk = token.New(token.OPEN_PAREN, l.read(1))
	case ')':
		tk = token.New(token.CLOSE_PAREN, l.read(1))
	case '{':
		tk = token.New(token.OPEN_BRACE, l.read(1))
	case '}':
		tk = token.New(token.CLOSE_BRACE, l.read(1))
	case ';':
		tk = token.New(token.SEMICOLON, l.read(1))
	case ',':
		tk = token.New(token.COMMA, l.read(1))
	case '=':
		if l.peek(1) == '=' {
			if l.peek(2) == '=' {
				tk = token.New(token.IDENTITY_EQUAL, l.read(3))
			} else {
				tk = token.New(token.EQUAL, l.read(2))
			}
		} else {
			tk = token.New(token.ASSIGN, l.read(1))
		}
	case '?':
		tk = token.New(token.QUESTION, l.read(1))
	case ':':
		tk = token.New(token.COLON, l.read(1))
	case '.':
		tk = token.New(token.DOT, l.read(1))
	case '~':
		tk = token.New(token.BIT_NOT, l.read(1))
	case '!':
		if l.peek(1) == '=' {
			if l.peek(2) == '=' {
				tk = token.New(token.IDENTITY_NOT_EQUAL, l.read(3))
			} else {
				tk = token.New(token.NOT_EQUAL, l.read(2))
			}
		} else {
			tk = token.New(token.NOT, l.read(1))
		}
	case '+':
		if l.peek(1) == '+' {
			tk = token.New(token.PLUS_PLUS, l.read(2))
		} else if l.peek(1) == '=' {
			tk = token.New(token.PLUS_ASSIGN, l.read(2))
		} else {
			tk = token.New(token.PLUS, l.read(1))
		}
	case '-':
		if l.peek(1) == '-' {
			tk = token.New(token.MINUS_MINUS, l.read(2))
		} else if l.peek(1) == '=' {
			tk = token.New(token.MINUS_ASSIGN, l.read(2))
		} else {
			tk = token.New(token.MINUS, l.read(1))
		}
	case '*':
		if l.peek(1) == '=' {
			tk = token.New(token.MULTIPLY_ASSIGN, l.read(2))
		} else {
			tk = token.New(token.MULTIPLY, l.read(1))
		}
	case '/':
		if l.peek(1) == '=' {
			tk = token.New(token.DIVIDE_ASSIGN, l.read(2))
		} else {
			tk = token.New(token.DIVIDE, l.read(1))
		}
	case '%':
		if l.peek(1) == '=' {
			tk = token.New(token.MODULUS_ASSIGN, l.read(2))
		} else {
			tk = token.New(token.MODULUS, l.read(1))
		}
	case '&':
		if l.peek(1) == '&' {
			tk = token.New(token.AND, l.read(2))
		} else if l.peek(1) == '=' {
			tk = token.New(token.BIT_AND_ASSIGN, l.read(2))
		} else {
			tk = token.New(token.BIT_AND, l.read(1))
		}
	case '|':
		if l.peek(1) == '|' {
			tk = token.New(token.OR, l.read(2))
		} else if l.peek(1) == '=' {
			tk = token.New(token.BIT_OR_ASSIGN, l.read(2))
		} else {
			tk = token.New(token.BIT_OR, l.read(1))
		}
	case '^':
		if l.peek(1) == '=' {
			tk = token.New(token.BIT_XOR_ASSIGN, l.read(2))
		}
	case '<':
		if l.peek(1) == '=' {
			tk = token.New(token.LESS_THAN_OR_EQUAL, l.read(2))
		} else if l.peek(1) == '<' {
			if l.peek(2) == '=' {
				tk = token.New(token.LEFT_SHIFT_ARITHMETIC_ASSIGN, l.read(3))
			} else {
				tk = token.New(token.LEFT_SHIFT_ARITHMETIC, l.read(2))
			}
		} else {
			tk = token.New(token.LESS_THAN, l.read(1))
		}
	case '>':
		if l.peek(1) == '=' {
			tk = token.New(token.GREATER_THAN_OR_EQUAL, l.read(2))
		} else if l.peek(1) == '>' {
			if l.peek(2) == '>' {
				if l.peek(3) == '=' {
					tk = token.New(token.RIGHT_SHIFT_LOGICAL_ASSIGN, l.read(4))
				} else {
					tk = token.New(token.RIGHT_SHIFT_LOGICAL, l.read(3))
				}
			} else if l.peek(2) == '=' {
				tk = token.New(token.RIGHT_SHIFT_ARITHMETIC_ASSIGN, l.read(3))
			} else {
				tk = token.New(token.RIGHT_SHIFT_ARITHMETIC, l.read(2))
			}
		} else {
			tk = token.New(token.GREATER_THAN, l.read(1))
		}
	default:
		if unicode.IsLetter(ch) || ch == '_' || ch == '$' {
			return l.identifier()
		} else if unicode.IsDigit(ch) {
			tk = l.number()
		} else {
			tk = token.New(token.ILLEGAL, l.read(1))
		}
	}

	return tk
}

func (l *Lexer) number() token.Token {
	ch := l.peek(0)
	if ch == '0' && (l.peek(1) == 'x' || l.peek(1) == 'X') {
		return l.hexInteger()
	}
	if ch == '0' && (l.peek(1) == 'b' || l.peek(1) == 'B') {
		return l.binaryInteger()
	}
	if ch == '0' {
		return l.octalInteger()
	}
	if ch == '.' || unicode.IsDigit(ch) {
		return l.decimal()
	}
	return l.syntaxError("invalid number format")
}

func (l *Lexer) decimal() token.Token {
	var builder strings.Builder

	for unicode.IsDigit(l.peek(0)) {
		builder.WriteRune(l.pop())
	}

	if l.peek(0) == '.' {
		builder.WriteRune(l.pop())
		if !unicode.IsDigit(l.peek(0)) {
			return l.syntaxError("invalid decimal number: missing digits after decimal point")
		}
	}

	for unicode.IsDigit(l.peek(0)) {
		builder.WriteRune(l.pop())
	}

	if l.peek(0) == 'e' || l.peek(0) == 'E' {
		builder.WriteRune(l.pop())
		if l.peek(0) == '+' || l.peek(0) == '-' {
			builder.WriteRune(l.pop())
		}
		if !unicode.IsDigit(l.peek(0)) {
			return l.syntaxError("invalid exponent in number")
		}
		for unicode.IsDigit(l.peek(0)) {
			builder.WriteRune(l.pop())
		}
	}

	literal := builder.String()
	return token.New(token.NUMBER, literal)
}

func (l *Lexer) hexInteger() token.Token {
	var builder strings.Builder

	builder.WriteRune(l.pop())
	builder.WriteRune(l.pop())

	for ch := l.peek(0); (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F'); ch = l.peek(0) {
		builder.WriteRune(l.pop())
	}

	literal := builder.String()
	return token.New(token.NUMBER, literal)
}

func (l *Lexer) binaryInteger() token.Token {
	var builder strings.Builder

	builder.WriteRune(l.pop())
	builder.WriteRune(l.pop())

	for l.peek(0) == '0' || l.peek(0) == '1' {
		builder.WriteRune(l.pop())
	}

	if builder.Len() == 0 {
		return l.syntaxError("invalid binary literal: no digits")
	}

	literal := builder.String()
	return token.New(token.NUMBER, literal)
}

func (l *Lexer) octalInteger() token.Token {
	var builder strings.Builder

	builder.WriteRune(l.pop())
	if l.peek(0) == 'o' || l.peek(0) == 'O' {
		builder.WriteRune(l.pop())
	}

	for ch := l.peek(0); ch >= '0' && ch <= '7'; ch = l.peek(0) {
		builder.WriteRune(l.pop())
	}

	literal := builder.String()
	return token.New(token.NUMBER, literal)
}

func (l *Lexer) string() token.Token {
	quote := l.pop()

	var builder strings.Builder
	for {
		ch := l.peek(0)
		if ch == rune(0) {
			return l.syntaxError("unterminated string literal")
		}
		if ch == quote {
			l.pop()
			break
		}
		if ch == '\\' {
			l.pop()
			ch = l.peek(0)
			switch ch {
			case 'n':
				builder.WriteRune('\n')
			case 'r':
				builder.WriteRune('\r')
			case 't':
				builder.WriteRune('\t')
			case '\\':
				builder.WriteRune('\\')
			case '"':
				builder.WriteRune('"')
			case '\'':
				builder.WriteRune('\'')
			default:
				builder.WriteRune(ch)
			}
			l.pop()
		} else if ch == '\r' || ch == '\n' {
			if l.peek(1) == '\\' {
				l.pop()
				l.pop()
				continue
			} else {
				builder.WriteRune(ch)
			}
		} else {
			builder.WriteRune(l.pop())
		}
	}

	literal := builder.String()
	return token.New(token.STRING, literal)
}

func (l *Lexer) identifier() token.Token {
	var builder strings.Builder

	builder.WriteRune(l.pop())

	for {
		ch := l.peek(0)
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' && ch != '$' {
			break
		}
		builder.WriteRune(l.pop())
	}

	literal := builder.String()
	return token.New(token.TypeOf(literal), literal)
}

func (l *Lexer) hidden() {
	l.space()
	l.comment()
}

func (l *Lexer) space() {
	for unicode.IsSpace(l.peek(0)) {
		l.pop()
	}
}

func (l *Lexer) comment() {
	for {
		ch := l.peek(0)
		if ch == '/' {
			if l.peek(1) == '*' {
				l.multiLineComment()
			} else if l.peek(1) == '/' {
				l.singleLineComment()
			} else {
				break
			}
		} else {
			break
		}
	}
}

func (l *Lexer) multiLineComment() {
	l.pop()
	l.pop()

	for {
		ch := l.peek(0)
		if ch == '*' && l.peek(1) == '/' {
			l.pop()
			l.pop()
			break
		}
		if ch == rune(0) {
			break
		}
		l.pop()
	}
}

func (l *Lexer) singleLineComment() {
	l.pop()
	l.pop()

	for {
		ch := l.peek(0)
		if ch == '\n' || ch == '\r' || ch == rune(0) {
			break
		}
		l.pop()
	}
}

func (l *Lexer) read(n int) string {
	var result []rune
	for i := 0; i < n; i++ {
		result = append(result, l.pop())
	}
	return string(result)
}

func (l *Lexer) peek(offset int) rune {
	if len(l.buf) <= l.pos+offset {
		if l.fetch(offset) == rune(0) {
			return rune(0)
		}
	}
	return l.buf[l.pos+offset]
}

func (l *Lexer) pop() rune {
	if l.pos >= len(l.buf) {
		return 0
	}

	ch := l.buf[l.pos]
	l.pos++

	if ch == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
	return ch
}

func (l *Lexer) fetch(offset int) rune {
	for len(l.buf) == 0 || l.pos+offset >= len(l.buf) {
		var ch rune
		_, err := fmt.Fscanf(l.source, "%c", &ch)
		if err == io.EOF {
			return rune(0)
		} else if err != nil {
			return rune(0)
		}
		l.buf = append(l.buf, ch)
	}
	return l.buf[l.pos]
}

func (l *Lexer) syntaxError(message string) token.Token {
	return token.New(token.ILLEGAL, fmt.Sprintf("syntax error at line %d, column %d: %s", l.line, l.column, message))
}
