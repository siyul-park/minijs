package lexer

import (
	"fmt"
	"strings"
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
	l.space()

	var tk token.Token
	switch l.peek(0) {
	case rune(0):
		tk = token.NewToken(token.EOF, "")

	case '"', '\'':
		tk = l.string(l.peek(0))
	case '`':
		tk = token.NewToken(token.TEMPLATE, string(l.pop()))

	case '(':
		tk = token.NewToken(token.PAREN_OPEN, string(l.pop()))
	case ')':
		tk = token.NewToken(token.PAREN_CLOSE, string(l.pop()))
	case '[':
		tk = token.NewToken(token.BRACKET_OPEN, string(l.pop()))
	case ']':
		tk = token.NewToken(token.BRACKET_CLOSE, string(l.pop()))
	case '{':
		tk = token.NewToken(token.CURLY_OPEN, string(l.pop()))
	case '}':
		tk = token.NewToken(token.CURLY_CLOSE, string(l.pop()))
	case ',':
		tk = token.NewToken(token.COMMA, string(l.pop()))
	case '.':
		tk = token.NewToken(token.PERIOD, string(l.pop()))
	case ':':
		tk = token.NewToken(token.COLON, string(l.pop()))
	case ';':
		tk = token.NewToken(token.SEMICOLON, string(l.pop()))

	case '+':
		l.pop()
		if l.peek(0) == '=' {
			l.pop()
			tk = token.NewToken(token.PLUS_ASSIGN, "+=")
		} else {
			tk = token.NewToken(token.PLUS, "+")
		}
	case '-':
		l.pop()
		if l.peek(0) == '=' {
			l.pop()
			tk = token.NewToken(token.MINUS_ASSIGN, "-=")
		} else {
			tk = token.NewToken(token.MINUS, "-")
		}
	case '*':
		l.pop()
		if l.peek(0) == '=' {
			l.pop()
			tk = token.NewToken(token.MULTIPLE_ASSIGN, "*=")
		} else {
			tk = token.NewToken(token.MULTIPLE, "*")
		}
	case '/':
		l.pop()
		if l.peek(0) == '/' {
			l.lineComment()
			tk = l.Next()
		} else if l.peek(0) == '*' {
			l.blockComment()
			tk = l.Next()
		} else if l.peek(0) == '=' {
			l.pop()
			tk = token.NewToken(token.DIVIDE_ASSIGN, "/=")
		} else {
			tk = token.NewToken(token.DIVIDE, "/")
		}
	case '%':
		l.pop()
		if l.peek(0) == '=' {
			l.pop()
			tk = token.NewToken(token.MODULAR_ASSIGN, "%=")
		} else {
			tk = token.NewToken(token.MODULAR, "%")
		}
	case '=':
		l.pop()
		if l.peek(0) == '=' {
			l.pop()
			tk = token.NewToken(token.EQUAL, "==")
		} else if l.peek(0) == '>' {
			l.pop()
			tk = token.NewToken(token.ARROW, "=>")
		} else {
			tk = token.NewToken(token.ASSIGN, "=")
		}
	case '!':
		l.pop()
		if l.peek(0) == '=' {
			l.pop()
			tk = token.NewToken(token.NOT_EQUAL, "!=")
		} else {
			tk = token.NewToken(token.NOT, "!")
		}
	case '&':
		l.pop()
		if l.peek(0) == '=' {
			l.pop()
			tk = token.NewToken(token.BIT_AND_ASSIGN, "&=")
		} else if l.peek(0) == '&' {
			l.pop()
			tk = token.NewToken(token.AND, "&&")
		} else {
			tk = token.NewToken(token.BIT_AND, "&")
		}
	case '|':
		l.pop()
		if l.peek(0) == '=' {
			l.pop()
			tk = token.NewToken(token.BIT_OR_ASSIGN, "|=")
		} else if l.peek(0) == '|' {
			l.pop()
			tk = token.NewToken(token.OR, "||")
		} else {
			tk = token.NewToken(token.BIT_OR, "|")
		}
	case '^':
		l.pop()
		if l.peek(0) == '=' {
			l.pop()
			tk = token.NewToken(token.BIT_XOR_ASSIGN, "^=")
		} else {
			tk = token.NewToken(token.BIT_XOR, "^")
		}
	case '~':
		tk = token.NewToken(token.BIT_NOT, string(l.pop()))
	case '<':
		l.pop()
		if l.peek(0) == '<' {
			l.pop()
			if l.peek(0) == '=' {
				l.pop()
				tk = token.NewToken(token.LEFT_SHIFT_ASSIGN, "<<=")
			} else {
				tk = token.NewToken(token.LEFT_SHIFT, "<<")
			}
		} else if l.peek(0) == '=' {
			l.pop()
			tk = token.NewToken(token.LESS_THAN_EQUAL, "<=")
		} else {
			tk = token.NewToken(token.LESS_THAN, "<")
		}
	case '>':
		l.pop()
		if l.peek(0) == '>' {
			l.pop()
			if l.peek(0) == '>' {
				l.pop()
				if l.peek(0) == '=' {
					l.pop()
					tk = token.NewToken(token.UNSIGNED_RIGHT_SHIFT_ASSIGN, ">>>=")
				} else {
					tk = token.NewToken(token.UNSIGNED_RIGHT_SHIFT, ">>>")
				}
			} else if l.peek(0) == '=' {
				l.pop()
				tk = token.NewToken(token.RIGHT_SHIFT_ASSIGN, ">>=")
			} else {
				tk = token.NewToken(token.RIGHT_SHIFT, ">>")
			}
		} else if l.peek(0) == '=' {
			l.pop()
			tk = token.NewToken(token.GREATER_THAN_EQUAL, ">=")
		} else {
			tk = token.NewToken(token.GREATER_THAN, ">")
		}

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
			return l.syntaxError("invalid exponent syntax")
		}
		exponent := l.integer()
		return token.NewToken(token.NUMBER, integer.Literal+"e"+sign+exponent.Literal)
	}

	return integer
}

func (l *Lexer) integer() token.Token {
	var literal []rune
	prev := rune(0)

	if l.peek(0) == '0' {
		switch l.peek(1) {
		case 'b', 'B':
			l.pop()
			l.pop()
			for ch := l.peek(0); ch == '0' || ch == '1' || ch == '_'; ch = l.peek(0) {
				if ch == '_' {
					if prev == '_' {
						return l.syntaxError("unexpected underscore")
					}
					l.pop()
					continue
				}
				literal = append(literal, l.pop())
				prev = literal[len(literal)-1]
			}
			if len(literal) == 0 {
				return l.syntaxError("invalid binary number")
			}
			return token.NewToken(token.NUMBER, "0b"+string(literal))

		case 'o', 'O':
			l.pop()
			l.pop()
			for ch := l.peek(0); '0' <= ch && ch <= '7' || ch == '_'; ch = l.peek(0) {
				if ch == '_' {
					if prev == '_' {
						return l.syntaxError("unexpected underscore")
					}
					l.pop()
					continue
				}
				literal = append(literal, l.pop())
				prev = literal[len(literal)-1]
			}
			if len(literal) == 0 {
				return l.syntaxError("invalid octal number")
			}
			return token.NewToken(token.NUMBER, "0o"+string(literal))

		case 'x', 'X':
			l.pop()
			l.pop()
			for ch := l.peek(0); unicode.IsDigit(ch) || ('a' <= ch && ch <= 'f') || ('A' <= ch && ch <= 'F') || ch == '_'; ch = l.peek(0) {
				if ch == '_' {
					if prev == '_' {
						return l.syntaxError("unexpected underscore")
					}
					l.pop()
					continue
				}
				literal = append(literal, l.pop())
				prev = literal[len(literal)-1]
			}
			if len(literal) == 0 {
				return l.syntaxError("invalid hexadecimal number")
			}
			return token.NewToken(token.NUMBER, "0x"+string(literal))
		}
	}

	literal = append(literal, l.pop())

	for unicode.IsDigit(l.peek(0)) || l.peek(0) == '_' {
		if l.peek(0) == '_' {
			if prev == '_' {
				return l.syntaxError("unexpected underscore")
			}
			l.pop()
			continue
		}
		literal = append(literal, l.pop())
		prev = literal[len(literal)-1]
	}

	if len(literal) > 0 && literal[len(literal)-1] == '_' {
		return l.syntaxError("unexpected trailing underscore")
	}
	return token.NewToken(token.NUMBER, string(literal))
}

func (l *Lexer) identifier() token.Token {
	var builder strings.Builder

	ch := l.peek(0)
	if !unicode.IsLetter(ch) && ch != '_' && ch != '$' {
		return l.syntaxError("invalid identifier start")
	}

	builder.WriteRune(l.pop())

	for {
		ch := l.peek(0)
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' && ch != '$' {
			break
		}
		builder.WriteRune(l.pop())
	}

	keyword := builder.String()
	return token.NewToken(token.TypeOf(keyword), keyword)
}

func (l *Lexer) string(delim rune) token.Token {
	if l.peek(0) != delim {
		return l.syntaxError("unterminated string literal")
	}
	l.pop()

	var literal []rune
	for {
		ch := l.peek(0)
		if ch == rune(0) {
			return l.syntaxError("unterminated string literal")
		}
		if ch == delim {
			l.pop()
			break
		}
		literal = append(literal, l.pop())
	}
	return token.NewToken(token.STRING, string(literal))
}

func (l *Lexer) space() {
	for unicode.IsSpace(l.peek(0)) {
		l.pop()
	}
}

func (l *Lexer) lineComment() {
	for ch := l.peek(0); ch != rune(0) && ch != '\n'; ch = l.peek(0) {
		l.pop()
	}
}

func (l *Lexer) blockComment() {
	l.pop()
	l.pop()

	for {
		ch := l.peek(0)
		if ch == rune(0) {
			break
		}
		if ch == '*' && l.peek(1) == '/' {
			l.pop()
			l.pop()
			break
		}
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

func (l *Lexer) syntaxError(message string) token.Token {
	return token.NewToken(token.ILLEGAL, fmt.Sprintf("syntax error: %s at position %d", message, l.pos))
}
