package scanner

import (
	"strconv"

	"github.com/codecrafters-io/interpreter-starter-go/app/util"
)

var keywords = map[string]tokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
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

type Scanner struct {
	source  string
	tokens  []token
	start   int
	current int
	line    int
}

func NewScanner(source string) Scanner {
	return Scanner{
		source:  source,
		tokens:  []token{},
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) ScanTokens() []token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.addToken(EOF)
	return s.tokens
}

func (s Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() {
	c := s.advance() // now indexof(c) and current are different points
	switch c {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
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
		} else {
			s.addToken(SLASH)
		}
	case ' ', '\r', '\t':
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
			util.Error(s.line, "Unexpected character: "+string(c))
		}
	}
}

func (s *Scanner) advance() byte {
	ret := s.source[s.current]
	s.current++ // current increases after returning the character
	return ret
}

func (s *Scanner) addToken(tokenType tokenType) {
	s.addTokenWithLiteral(tokenType, nil)
}

func (s *Scanner) addTokenWithLiteral(tokenType tokenType, literal any) {
	var text string
	if tokenType == EOF {
		text = ""
	} else {
		text = s.source[s.start:s.current]
	}
	s.tokens = append(s.tokens, token{
		tokenType: tokenType,
		lexeme:    text,
		literal:   literal,
		line:      s.line,
	})
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	// advance has already been called in scanToken, so current points to the next character
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s Scanner) peek() byte {
	if s.isAtEnd() {
		return '\000'
	}
	return s.source[s.current]
}

func (s Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return '\000'
	}
	return s.source[s.current+1]
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.isAtEnd() {
		util.Error(s.line, "Unterminated string.")
		return
	}

	s.advance()
	// now start points to the opening quote, current points to the closing quote
	value := s.source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(STRING, value)
}

func (s *Scanner) number() {
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
	s.addTokenWithLiteral(NUMBER, value)
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current]

	if tokenType, ok := keywords[text]; !ok {
		s.addToken(IDENTIFIER)
	} else {
		s.addToken(tokenType)
	}
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}
