package main

import (
	"fmt"
	"strconv"
	"unicode"
)

var (
	keywords = map[string]TokenType{
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
)

type Scanner struct {
	start, current, line int
	source               string
	tokens               []Token
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source: source,
		tokens: []Token{},
		line:   1,
	}
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() rune {
	c := s.source[s.current]
	s.current++
	return rune(c)
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return rune(0)
	}
	return rune(s.source[s.current])
}

func (s *Scanner) peekNext() rune {
	if (s.current + 1) >= len(s.source) {
		return rune(0)
	}

	return rune(s.source[s.current+1])
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	if rune(s.source[s.current]) != expected {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) addToken(tokenType TokenType) {
	s.addTokenLiteral(tokenType, nil)
}

func (s *Scanner) addTokenLiteral(tokenType TokenType, literal any) {
	s.tokens = append(s.tokens, Token{
		Type:    tokenType,
		Lexeme:  string(s.source[s.start:s.current]),
		Literal: literal,
		Line:    s.line,
	})
}

func (s *Scanner) addStringToken() {
	startLine := s.line

	for s.peek() != '"' && !s.isAtEnd() {
		if s.advance() == '\n' {
			s.line++
		}
	}

	if s.isAtEnd() {
		Error(startLine, "Unterminated string.")
		return
	}

	// consume the closing "
	s.advance()

	s.addTokenLiteral(STRING, s.source[s.start+1:s.current-1])
}

func (s *Scanner) addNumberToken() {
	for unicode.IsDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && unicode.IsDigit(s.peekNext()) {
		// consume the "."
		s.advance()

		for unicode.IsDigit(s.peek()) {
			s.advance()
		}
	}

	n, err := strconv.ParseFloat(string(s.source[s.start:s.current]), 64)
	if err != nil {
		Error(s.line, fmt.Sprintf("error parsing float: %s", err.Error()))
		return
	}
	s.addTokenLiteral(NUMBER, n)
}

func (s *Scanner) isAlpha(c rune) bool {
	return unicode.IsLetter(c) || c == rune('_')
}

func (s *Scanner) isAlphaNumeric(c rune) bool {
	return s.isAlpha(c) || unicode.IsDigit(c)
}

func (s *Scanner) addIdentifierToken() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	tokenType, ok := keywords[text]
	if !ok {
		tokenType = IDENTIFIER
	}
	s.addToken(tokenType)
}

// advance (consume and return token) and check for basic
// single character lexemes.  For multi-character lexemes
// peek and/or comume more characters to match the lexeme
func (s *Scanner) scanToken() {
	c := s.advance()

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
			for s.peek() != rune('\n') && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH)
		}
	case ' ', '\r', '\t':
		// noop, ignore whitespace
	case '\n':
		// ignore newline, but track it
		s.line++
	case '"':
		// string literals
		s.addStringToken()
	default:
		if unicode.IsDigit(c) {
			s.addNumberToken()
		} else if s.isAlpha(c) {
			s.addIdentifierToken()
		} else {
			Error(s.line, "Unexpected character.")
		}
	}
}

func (s *Scanner) scanTokens() ([]Token, error) {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, Token{
		Type:   EOF,
		Lexeme: "",
		Line:   s.line,
	})

	return s.tokens, nil
}
