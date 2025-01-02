package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrettyPrint(t *testing.T) {
	testCases := []struct {
		name     string
		expr     Expr
		expected string
	}{
		{
			name:     "basic",
			expected: "(* (- 123) (group 45.67))",
			expr: BinaryExpr{
				Op: Token{
					Type:   STAR,
					Lexeme: "*",
				},
				Left: UnaryExpr{
					Op: Token{
						Type:   MINUS,
						Lexeme: "-",
					},
					Right: LiteralExpr{
						Value: 123,
					},
				},
				Right: GroupingExpr{
					Expression: LiteralExpr{
						Value: 45.67,
					},
				},
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.expr.Print())
		})
	}
}
