package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrimaryNoGroupings(t *testing.T) {
	testCases := []struct {
		name     string
		tokens   []Token
		expected Expr
	}{
		{
			name: "false",
			tokens: []Token{
				{
					Type:   FALSE,
					Lexeme: "false",
				},
			},
			expected: LiteralExpr{
				Value: false,
			},
		},
		{
			name: "true",
			tokens: []Token{
				{
					Type:   TRUE,
					Lexeme: "true",
				},
			},
			expected: LiteralExpr{
				Value: true,
			},
		},
		{
			name: "nil",
			tokens: []Token{
				{
					Type:   NIL,
					Lexeme: "nil",
				},
			},
			expected: LiteralExpr{
				Value: nil,
			},
		},
		{
			name: "number",
			tokens: []Token{
				{
					Type:    NUMBER,
					Lexeme:  "52.25",
					Literal: float64(52.25),
				},
			},
			expected: LiteralExpr{
				Value: float64(52.25),
			},
		},
		{
			name: "string",
			tokens: []Token{
				{
					Type:    STRING,
					Lexeme:  `"hello world"`,
					Literal: "hello world",
				},
			},
			expected: LiteralExpr{
				Value: "hello world",
			},
		},
		{
			name: "group",
			tokens: []Token{
				{
					Type:   LEFT_PAREN,
					Lexeme: "(",
				},
				{
					Type:   NIL,
					Lexeme: "nil",
				},
				{
					Type:   RIGHT_PAREN,
					Lexeme: ")",
				},
			},
			expected: GroupingExpr{
				Expression: LiteralExpr{
					Value: nil,
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewParser(tt.tokens)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, p.Parse())
		})
	}
}

func TestUnary(t *testing.T) {
	testCases := []struct {
		name     string
		tokens   []Token
		expected Expr
	}{
		{
			name: "bang",
			tokens: []Token{
				{
					Type:   BANG,
					Lexeme: "!",
				},
				{
					Type:   TRUE,
					Lexeme: "true",
				},
			},
			expected: UnaryExpr{
				Op: Token{
					Type:   BANG,
					Lexeme: "!",
				},
				Right: LiteralExpr{
					Value: true,
				},
			},
		},
		{
			name: "minus",
			tokens: []Token{
				{
					Type:   MINUS,
					Lexeme: "-",
				},
				{
					Type:    NUMBER,
					Lexeme:  "52.25",
					Literal: float64(52.25),
				},
			},
			expected: UnaryExpr{
				Op: Token{
					Type:   MINUS,
					Lexeme: "-",
				},
				Right: LiteralExpr{
					Value: 52.25,
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewParser(tt.tokens)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, p.Parse())
		})
	}
}

func TestFactor(t *testing.T) {
	testCases := []struct {
		name     string
		tokens   []Token
		expected Expr
	}{
		{
			name: "slash",
			tokens: []Token{
				{
					Type:    NUMBER,
					Lexeme:  "3",
					Literal: float64(3),
				},
				{
					Type:   SLASH,
					Lexeme: "/",
				},
				{
					Type:    NUMBER,
					Lexeme:  "3",
					Literal: float64(3),
				},
			},
			expected: BinaryExpr{
				Op: Token{
					Type:   SLASH,
					Lexeme: "/",
				},
				Left: LiteralExpr{
					Value: 3.,
				},
				Right: LiteralExpr{
					Value: 3.,
				},
			},
		},
		{
			name: "star",
			tokens: []Token{
				{
					Type:    NUMBER,
					Lexeme:  "3",
					Literal: float64(3),
				},
				{
					Type:   STAR,
					Lexeme: "*",
				},
				{
					Type:    NUMBER,
					Lexeme:  "3",
					Literal: float64(3),
				},
			},
			expected: BinaryExpr{
				Op: Token{
					Type:   STAR,
					Lexeme: "*",
				},
				Left: LiteralExpr{
					Value: 3.,
				},
				Right: LiteralExpr{
					Value: 3.,
				},
			},
		},
		{
			name: "associativity",
			tokens: []Token{
				{
					Type:    NUMBER,
					Lexeme:  "3",
					Literal: float64(3),
				},
				{
					Type:   STAR,
					Lexeme: "*",
				},
				{
					Type:    NUMBER,
					Lexeme:  "3",
					Literal: float64(3),
				},
				{
					Type:   STAR,
					Lexeme: "*",
				},
				{
					Type:    NUMBER,
					Lexeme:  "3",
					Literal: float64(3),
				},
			},
			expected: BinaryExpr{
				Op: Token{
					Type:   STAR,
					Lexeme: "*",
				},
				Left: BinaryExpr{
					Op: Token{
						Type:   STAR,
						Lexeme: "*",
					},
					Left: LiteralExpr{
						Value: 3.,
					},
					Right: LiteralExpr{
						Value: 3.,
					},
				},
				Right: LiteralExpr{
					Value: 3.,
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewParser(tt.tokens)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, p.Parse())
		})
	}
}

func TestTerm(t *testing.T) {
	testCases := []struct {
		name     string
		tokens   []Token
		expected Expr
	}{
		{
			name: "minus",
			tokens: []Token{
				{
					Type:    NUMBER,
					Lexeme:  "3",
					Literal: float64(3),
				},
				{
					Type:   MINUS,
					Lexeme: "-",
				},
				{
					Type:    NUMBER,
					Lexeme:  "3",
					Literal: float64(3),
				},
			},
			expected: BinaryExpr{
				Op: Token{
					Type:   MINUS,
					Lexeme: "-",
				},
				Left: LiteralExpr{
					Value: 3.,
				},
				Right: LiteralExpr{
					Value: 3.,
				},
			},
		},
		{
			name: "plus",
			tokens: []Token{
				{
					Type:    NUMBER,
					Lexeme:  "3",
					Literal: float64(3),
				},
				{
					Type:   PLUS,
					Lexeme: "+",
				},
				{
					Type:    NUMBER,
					Lexeme:  "3",
					Literal: float64(3),
				},
			},
			expected: BinaryExpr{
				Op: Token{
					Type:   PLUS,
					Lexeme: "+",
				},
				Left: LiteralExpr{
					Value: 3.,
				},
				Right: LiteralExpr{
					Value: 3.,
				},
			},
		},
		{
			name: "associativity",
			tokens: []Token{
				{
					Type:    NUMBER,
					Lexeme:  "3",
					Literal: float64(3),
				},
				{
					Type:   PLUS,
					Lexeme: "+",
				},
				{
					Type:    NUMBER,
					Lexeme:  "3",
					Literal: float64(3),
				},
				{
					Type:   PLUS,
					Lexeme: "+",
				},
				{
					Type:    NUMBER,
					Lexeme:  "3",
					Literal: float64(3),
				},
			},
			expected: BinaryExpr{
				Op: Token{
					Type:   PLUS,
					Lexeme: "+",
				},
				Left: BinaryExpr{
					Op: Token{
						Type:   PLUS,
						Lexeme: "+",
					},
					Left: LiteralExpr{
						Value: 3.,
					},
					Right: LiteralExpr{
						Value: 3.,
					},
				},
				Right: LiteralExpr{
					Value: 3.,
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewParser(tt.tokens)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, p.Parse())
		})
	}
}

func TestComparison(t *testing.T) {
	testCases := []struct {
		name     string
		tokens   []Token
		expected Expr
	}{
		{
			name: "greater",
			tokens: []Token{
				{
					Type:    NUMBER,
					Lexeme:  "3",
					Literal: 3.,
				},
				{
					Type:   GREATER,
					Lexeme: ">",
				},
				{
					Type:    NUMBER,
					Lexeme:  "3",
					Literal: 3.,
				},
			},
			expected: BinaryExpr{
				Op: Token{
					Type:   GREATER,
					Lexeme: ">",
				},
				Left: LiteralExpr{
					Value: 3.,
				},
				Right: LiteralExpr{
					Value: 3.,
				},
			},
		},
		{
			name: "greater_equal",
			tokens: []Token{
				{
					Type:    NUMBER,
					Lexeme:  "3",
					Literal: 3.,
				},
				{
					Type:   GREATER_EQUAL,
					Lexeme: ">=",
				},
				{
					Type:    NUMBER,
					Lexeme:  "3",
					Literal: 3.,
				},
			},
			expected: BinaryExpr{
				Op: Token{
					Type:   GREATER_EQUAL,
					Lexeme: ">=",
				},
				Left: LiteralExpr{
					Value: 3.,
				},
				Right: LiteralExpr{
					Value: 3.,
				},
			},
		},
		{
			name: "less",
			tokens: []Token{
				{
					Type:    NUMBER,
					Lexeme:  "3",
					Literal: 3.,
				},
				{
					Type:   LESS,
					Lexeme: "<",
				},
				{
					Type:    NUMBER,
					Lexeme:  "3",
					Literal: 3.,
				},
			},
			expected: BinaryExpr{
				Op: Token{
					Type:   LESS,
					Lexeme: "<",
				},
				Left: LiteralExpr{
					Value: 3.,
				},
				Right: LiteralExpr{
					Value: 3.,
				},
			},
		},
		{
			name: "less",
			tokens: []Token{
				{
					Type:    NUMBER,
					Lexeme:  "3",
					Literal: 3.,
				},
				{
					Type:   LESS_EQUAL,
					Lexeme: "<=",
				},
				{
					Type:    NUMBER,
					Lexeme:  "3",
					Literal: 3.,
				},
			},
			expected: BinaryExpr{
				Op: Token{
					Type:   LESS_EQUAL,
					Lexeme: "<=",
				},
				Left: LiteralExpr{
					Value: 3.,
				},
				Right: LiteralExpr{
					Value: 3.,
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewParser(tt.tokens)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, p.Parse())
		})
	}
}

func TestEquality(t *testing.T) {
	testCases := []struct {
		name     string
		tokens   []Token
		expected Expr
	}{
		{
			name: "bang_equal",
			tokens: []Token{
				{
					Type:    NUMBER,
					Lexeme:  "3",
					Literal: 3.,
				},
				{
					Type:   BANG_EQUAL,
					Lexeme: "!=",
				},
				{
					Type:    NUMBER,
					Lexeme:  "3",
					Literal: 3.,
				},
			},
			expected: BinaryExpr{
				Op: Token{
					Type:   BANG_EQUAL,
					Lexeme: "!=",
				},
				Left: LiteralExpr{
					Value: 3.,
				},
				Right: LiteralExpr{
					Value: 3.,
				},
			},
		},
		{
			name: "equal_equal",
			tokens: []Token{
				{
					Type:    NUMBER,
					Lexeme:  "3",
					Literal: 3.,
				},
				{
					Type:   EQUAL_EQUAL,
					Lexeme: "==",
				},
				{
					Type:    NUMBER,
					Lexeme:  "3",
					Literal: 3.,
				},
			},
			expected: BinaryExpr{
				Op: Token{
					Type:   EQUAL_EQUAL,
					Lexeme: "==",
				},
				Left: LiteralExpr{
					Value: 3.,
				},
				Right: LiteralExpr{
					Value: 3.,
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewParser(tt.tokens)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, p.Parse())
		})
	}
}
