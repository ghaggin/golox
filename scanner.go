package main

import (
	"errors"
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

func scanTokens(source string) ([]Token, error) {
	start, current, line := 0, 0, 0
	errs := []error{}
	tokens := []Token{}

	isAtEnd := func() bool {
		return current >= len(source)
	}

	addToken := func(tokenType TokenType) {
		tokens = append(tokens, Token{
			Type: tokenType,
		})
	}

	addTokenLiteral := func(tokenType TokenType, literal any) {
		tokens = append(tokens, Token{
			Type:    tokenType,
			Lexeme:  string(source[start:current]),
			Literal: literal,
			Line:    line,
		})
	}

	match := func(expected rune) bool {
		if isAtEnd() {
			return false
		}

		if rune(source[current]) != expected {
			return false
		}

		current++
		return true
	}

	peek := func() rune {
		if isAtEnd() {
			return rune(0)
		}
		return rune(source[current])
	}

	peekNext := func() rune {
		if (current + 1) > len(source) {
			return rune(0)
		}

		return rune(source[current+1])
	}

	advance := func() rune {
		c := source[current]
		current++
		return rune(c)
	}

	addStringToken := func() {
		for peek() != '"' && !isAtEnd() {
			if peek() == '\n' {
				line++
			}
			advance()
		}

		if isAtEnd() {
			errs = append(errs, fmt.Errorf("unterminated string"))
			return
		}

		// The closing "
		advance()

		v := source[start+1 : current-1]
		addTokenLiteral(STRING, v)
	}

	addNumberToken := func() {
		for unicode.IsDigit(peek()) {
			advance()
		}

		if peek() == '.' && unicode.IsDigit(peekNext()) {
			// consume the "."
			advance()

			for unicode.IsDigit(peek()) {
				advance()
			}
		}

		n, err := strconv.ParseFloat(string(source[start:current]), 64)
		if err != nil {
			errs = append(errs, err)
			return
		}
		addTokenLiteral(NUMBER, n)
	}

	isAlpha := func(c rune) bool {
		return unicode.IsLetter(c) || c == rune('_')
	}

	isAlphaNumeric := func(c rune) bool {
		return isAlpha(c) || unicode.IsDigit(c)
	}

	addIdentifierToken := func() {
		for isAlphaNumeric(peek()) {
			advance()
		}

		text := source[start:current]
		tokenType, ok := keywords[text]
		if !ok {
			tokenType = IDENTIFIER
		}
		addToken(tokenType)
	}

	for !isAtEnd() {
		c := advance()
		switch c {
		case '(':
			addToken(LEFT_PAREN)
		case ')':
			addToken(RIGHT_PAREN)
		case '{':
			addToken(LEFT_BRACE)
		case '}':
			addToken(RIGHT_BRACE)
		case ',':
			addToken(COMMA)
		case '.':
			addToken(DOT)
		case '-':
			addToken(MINUS)
		case '+':
			addToken(PLUS)
		case ';':
			addToken(SEMICOLON)
		case '*':
			addToken(STAR)
		case '!':
			if match('=') {
				addToken(BANG_EQUAL)
			} else {
				addToken(BANG)
			}
		case '=':
			if match('=') {
				addToken(EQUAL_EQUAL)
			} else {
				addToken(EQUAL)
			}
		case '<':
			if match('=') {
				addToken(LESS_EQUAL)
			} else {
				addToken(LESS)
			}
		case '>':
			if match('=') {
				addToken(GREATER_EQUAL)
			} else {
				addToken(GREATER)
			}
		case '/':
			if match('/') {
				for peek() != '\n' && !isAtEnd() {
					advance()
				}
			} else {
				addToken(SLASH)
			}
		case ' ':
		case '\r':
		case '\t':
		case '\n':
			line++
		case '"':
			// string literals
			addStringToken()
		default:
			if unicode.IsDigit(c) {
				addNumberToken()
			} else if isAlpha(c) {
				addIdentifierToken()
			} else {
				errs = append(errs, fmt.Errorf("unexpected character on line %d", line))
			}
		}
	}

	return tokens, errors.Join(errs...)
}
