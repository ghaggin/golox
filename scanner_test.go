package main

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsAtEnd(t *testing.T) {
	testCases := []struct {
		name     string
		source   string
		current  int
		expected bool
	}{
		{
			name:     "before end",
			source:   "1234",
			current:  3,
			expected: false,
		},
		{
			name:     "at end",
			source:   "1234",
			current:  4,
			expected: true,
		},
		{
			name:     "after end",
			source:   "1234",
			current:  5,
			expected: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s := Scanner{
				source:  tt.source,
				current: tt.current,
			}

			assert.Equal(t, tt.expected, s.isAtEnd())
		})
	}
}

func TestAdvance(t *testing.T) {
	testCases := []struct {
		name        string
		source      string
		current     int
		expected1   rune
		expected2   rune
		expectPanic bool
	}{
		{
			name:      "advance returns normally",
			source:    "1234",
			current:   0,
			expected1: rune('1'),
			expected2: rune('2'),
		},
		{
			name:        "advance panics",
			source:      "1234",
			current:     4,
			expectPanic: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scanner{
				source:  tt.source,
				current: tt.current,
			}

			if tt.expectPanic {
				assert.Panics(t, func() {
					s.advance()
				})
			} else {
				assert.Equal(t, tt.expected1, s.advance())
				assert.Equal(t, tt.expected2, s.advance())
			}
		})
	}
}

func TestPeek(t *testing.T) {
	testCases := []struct {
		name     string
		source   string
		current  int
		expected rune
	}{
		{
			name:     "peek at end",
			source:   "1234",
			current:  4,
			expected: rune(0),
		},
		{
			name:     "peek after end",
			source:   "1234",
			current:  5,
			expected: rune(0),
		},
		{
			name:     "peek at end",
			source:   "1234",
			current:  3,
			expected: rune('4'),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scanner{
				source:  tt.source,
				current: tt.current,
			}
			assert.Equal(t, tt.expected, s.peek())
		})
	}
}

func TestPeekNext(t *testing.T) {
	testCases := []struct {
		name     string
		source   string
		current  int
		expected rune
	}{
		{
			name:     "peek at end",
			source:   "1234",
			current:  3,
			expected: rune(0),
		},
		{
			name:     "peek after end",
			source:   "1234",
			current:  4,
			expected: rune(0),
		},
		{
			name:     "peek after end 2",
			source:   "1234",
			current:  5,
			expected: rune(0),
		},
		{
			name:     "peek at end",
			source:   "1234",
			current:  2,
			expected: rune('4'),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scanner{
				source:  tt.source,
				current: tt.current,
			}
			assert.Equal(t, tt.expected, s.peekNext())
		})
	}
}

func TestMatch(t *testing.T) {
	testCases := []struct {
		name     string
		source   string
		current  int
		match1   byte
		match2   byte
		expected bool
	}{
		{
			name:     "at end",
			source:   "1234",
			current:  4,
			expected: false,
		},
		{
			name:     "does match",
			source:   "1234",
			current:  2,
			match1:   '3',
			match2:   '4',
			expected: true,
		},
		{
			name:     "doesn't match",
			source:   "1234",
			current:  2,
			match1:   '0',
			match2:   '0',
			expected: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scanner{
				source:  tt.source,
				current: tt.current,
			}

			assert.Equal(t, tt.expected, s.match(rune(tt.match1)))
			assert.Equal(t, tt.expected, s.match(rune(tt.match2)))
		})
	}
}

func TestAddStringToken(t *testing.T) {
	testCases := []struct {
		name     string
		source   string
		current  int
		expected string
	}{
		{
			name:     "basic",
			source:   `var a = "test string";`,
			current:  8,
			expected: "test string",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scanner{
				source:  tt.source,
				current: tt.current,
				start:   tt.current,
				line:    3,
			}

			assert.Len(t, s.tokens, 0)
			s.scanToken()
			require.Len(t, s.tokens, 1)
			assert.Equal(t, Token{
				Type:    STRING,
				Lexeme:  fmt.Sprintf(`"%s"`, tt.expected),
				Literal: tt.expected,
				Line:    3,
			}, s.tokens[0])
			assert.Equal(t, s.current, tt.current+len(tt.expected)+2)
			assert.Equal(t, ';', s.peek())
		})
	}
}

func TestAddNumberToken(t *testing.T) {
	testCases := []struct {
		name     string
		source   string
		current  int
		expected string
	}{
		{
			name:     "int",
			source:   `var a = 1234;`,
			current:  8,
			expected: "1234",
		},
		{
			name:     "decimal",
			source:   `var a = 12.34;`,
			current:  8,
			expected: "12.34",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scanner{
				source:  tt.source,
				current: tt.current,
				start:   tt.current,
				line:    3,
			}

			expectedFloat, err := strconv.ParseFloat(tt.expected, 64)
			require.NoError(t, err)

			assert.Len(t, s.tokens, 0)
			s.scanToken()
			require.Len(t, s.tokens, 1)
			assert.Equal(t, Token{
				Type:    NUMBER,
				Lexeme:  tt.expected,
				Literal: expectedFloat,
				Line:    3,
			}, s.tokens[0])
			assert.Equal(t, s.current, tt.current+len(tt.expected))
			assert.Equal(t, ';', s.peek())
		})
	}
}

func TestAddNonKeywordIdentifier(t *testing.T) {
	testCases := []struct {
		name     string
		source   string
		current  int
		expected []Token
	}{
		{
			name:   "basic",
			source: "a",
			expected: []Token{
				{
					Type:   IDENTIFIER,
					Lexeme: "a",
				},
			},
		},
		{
			name:   "basic",
			source: "a",
			expected: []Token{
				{
					Type:   IDENTIFIER,
					Lexeme: "a",
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scanner{
				current: tt.current,
				start:   tt.current,
				source:  tt.source,
			}

			s.scanToken()
			require.Len(t, s.tokens, 1)
			assert.Equal(t, tt.expected, s.tokens)
		})
	}
}

func TestScanTokens(t *testing.T) {
	testCases := []struct {
		name     string
		source   string
		expected []Token
	}{
		{
			name:   "ands",
			source: "a and b and c",
			expected: []Token{
				{
					Type:   IDENTIFIER,
					Lexeme: "a",
				},
				{
					Type:   AND,
					Lexeme: "and",
				},
				{
					Type:   IDENTIFIER,
					Lexeme: "b",
				},
				{
					Type:   AND,
					Lexeme: "and",
				},
				{
					Type:   IDENTIFIER,
					Lexeme: "c",
				},
				{
					Type:   EOF,
					Lexeme: "",
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			s := NewScanner(tt.source)
			s.line = 0 // to simplify the tests
			tokens, err := s.scanTokens()
			require.NoError(t, err)
			assert.Equal(t, tt.expected, tokens)
		})
	}
}
