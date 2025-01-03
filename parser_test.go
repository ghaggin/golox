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
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewParser(tt.tokens)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, p.Parse())
		})
	}
}
