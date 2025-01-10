package main

import "strconv"

var (
	keywords map[string]int = map[string]int{
		"and":    AND,
		"class":  CLASS,
		"else":   ELSE,
		"false":  FALSE,
		"fun":    FUN,
		"for":    FOR,
		"if":     IF,
		"nil":    NIL,
		"or":     OR,
		"print":  PRINT,
		"return": RETURN,
		"super":  SUPER,
		"this":   THIS,
		"true":   TRUE,
		"var":    VAR,
		"while":  WHILE,
	}
)

type Scanner interface {
	scanTokens() []Token
}

type Scan struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

func NewScanner(source string) *Scan {
	return &Scan{
		source:  source,
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scan) scanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, NewToken(EOF, "", Object{}, s.line))
	return s.tokens
}

func (s *Scan) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scan) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PARAN)
	case ')':
		s.addToken(RIGHT_PARAN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ',':
		s.addToken(COMMA)
	case '.':
		s.addToken(DOT)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case ';':
		s.addToken(SEMICOLON)
	case '*':
		s.addToken(STAR)
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL)
		} else {
			s.addToken(BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else if s.match('*') {
			count := 1
			for count > 0 {
				if s.isAtEnd() {
					error(s.line, "Multi line comment is not closed")
					break
				}
				if s.peek() == '*' && s.peekNext() == '/' {
					s.advance()
					count--
				}

				if s.peek() == '/' && s.peekNext() == '*' {
					s.advance()
					count++
				}

				if s.peek() == '\n' {
					s.line++
				}
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}
	case ' ', '\r', '\t':
		break
	case '\n':
		s.line++
	case '"':
		s.string()
	default:
		if isDigit(c) {
			s.number()
		} else if isAlpha(c) {
			s.identifier()
		} else {
			error(s.line, "Unexpected character.")
		}
	}
}

// advance Return the character at current position and advance to next position
func (s *Scan) advance() rune {
	defer func() { s.current += 1 }()
	return rune(s.source[s.current])
}

// match If the expected character is same as the character present at s.current
// then true is returned
// else false
func (s *Scan) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	if rune(s.source[s.current]) != expected {
		return false
	}

	s.current += 1
	return true
}

// peek Returns the characeter that is present at s.current
func (s *Scan) peek() rune {
	if s.isAtEnd() {
		return '\n'
	}
	return rune(s.source[s.current])
}

// peekNext Returns the character present at current+1
func (s *Scan) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return '\n'
	}
	return rune(s.source[s.current+1])
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func isAlphaNumeric(c rune) bool {
	return isAlpha(c) || isDigit(c)
}

func (s *Scan) number() {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}

	value, _ := strconv.ParseFloat(s.source[s.start:s.current], 64)
	s.addTokenObj(NUMBER, Object{
		objType:     NUMBER_TYPE,
		value_float: value,
	})
}

func (s *Scan) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	value := s.source[s.start:s.current]
	tokenType, ok := keywords[value]
	if !ok {
		tokenType = IDENTIFIER
	}
	s.addToken(tokenType)
}

func (s *Scan) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line += 1
		}
		s.advance()
	}

	if s.isAtEnd() {
		error(s.line, "Unterminated String.")
		return
	}

	s.advance()

	value := s.source[s.start+1 : s.current-1]
	s.addTokenObj(STRING, Object{
		objType:   STRING_TYPE,
		value_str: value,
	})
}

func (s *Scan) addToken(tokenType int) {
	s.addTokenObj(tokenType, Object{})
}

func (s *Scan) addTokenObj(tokenType int, object Object) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, NewToken(tokenType, text, object, s.line))
}
