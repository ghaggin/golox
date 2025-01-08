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

func TestIsTruthy(t *testing.T) {
	testCases := []struct {
		name     string
		v        any
		expected bool
	}{
		{
			name:     "nil",
			v:        nil,
			expected: false,
		},
		{
			name:     "true",
			v:        true,
			expected: true,
		},
		{
			name:     "false",
			v:        false,
			expected: false,
		},
		{
			name:     "zero",
			v:        0.,
			expected: true,
		},
		{
			name:     "non_zero",
			v:        1.,
			expected: true,
		},
		{
			name:     "empty_string",
			v:        "",
			expected: true,
		},
		{
			name:     "string",
			v:        "hello world",
			expected: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isTruthy(tt.v))
		})
	}
}

func TestIsEqual(t *testing.T) {
	testCases := []struct {
		name     string
		l, r     any
		expected bool
	}{
		{
			name:     "both nil",
			l:        nil,
			r:        nil,
			expected: true,
		},
		{
			name:     "one nil",
			l:        nil,
			r:        1.,
			expected: false,
		},
		{
			name:     "one nil",
			l:        1.,
			r:        nil,
			expected: false,
		},
		{
			name:     "floats_equal",
			l:        1.,
			r:        1.,
			expected: true,
		},
		{
			name:     "floats_not_equal",
			l:        1.,
			r:        2.,
			expected: false,
		},
		{
			name:     "strings_equal",
			l:        "s1",
			r:        "s1",
			expected: true,
		},
		{
			name:     "strings_not_equal",
			l:        "s1",
			r:        "s2",
			expected: false,
		},
		{
			name:     "string_not_equal_float",
			l:        "s1",
			r:        1.,
			expected: false,
		},
		{
			name:     "boolean_equal",
			l:        true,
			r:        true,
			expected: true,
		},
		{
			name:     "boolean_equal_2",
			l:        false,
			r:        false,
			expected: true,
		},
		{
			name:     "boolean_not_equal",
			l:        true,
			r:        false,
			expected: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, isEqual(tt.l, tt.r))
		})
	}
}

func TestBinaryExpression(t *testing.T) {
	testCases := []struct {
		name     string
		l        any
		r        any
		op       TokenType
		expected any
	}{
		{
			name:     "equal_equal",
			l:        3.,
			r:        5.,
			op:       EQUAL_EQUAL,
			expected: false,
		},
		{
			name:     "bang_equal",
			l:        3.,
			r:        5.,
			op:       BANG_EQUAL,
			expected: true,
		},
		{
			name:     "greater",
			l:        3.,
			r:        5.,
			op:       GREATER,
			expected: false,
		},
		{
			name:     "greater_equal",
			l:        3.,
			r:        5.,
			op:       GREATER_EQUAL,
			expected: false,
		},
		{
			name:     "less",
			l:        3.,
			r:        5.,
			op:       LESS,
			expected: true,
		},
		{
			name:     "less_equal",
			l:        3.,
			r:        5.,
			op:       LESS_EQUAL,
			expected: true,
		},
		{
			name:     "minus",
			l:        3.,
			r:        5.,
			op:       MINUS,
			expected: -2.,
		},
		{
			name:     "slash",
			l:        3.,
			r:        5.,
			op:       SLASH,
			expected: 3. / 5.,
		},
		{
			name:     "star",
			l:        3.,
			r:        5.,
			op:       STAR,
			expected: 3. * 5.,
		},
		{
			name:     "plus_floats",
			l:        3.,
			r:        5.,
			op:       PLUS,
			expected: 3. + 5.,
		},
		{
			name:     "plus_strings",
			l:        "hello ",
			r:        "world",
			op:       PLUS,
			expected: "hello world",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, BinaryExpr{
				Op: Token{
					Type: tt.op,
				},
				Left: LiteralExpr{
					Value: tt.l,
				},
				Right: LiteralExpr{
					Value: tt.r,
				},
			}.Evaluate())
		})
	}
}

func TestUnaryEvaluate(t *testing.T) {
	testCases := []struct {
		name     string
		op       TokenType
		r        any
		expected any
	}{
		{
			name:     "minus",
			op:       MINUS,
			r:        3.,
			expected: -3.,
		},
		{
			name:     "bang",
			op:       BANG,
			r:        3.,
			expected: false,
		},
		{
			name:     "bang",
			op:       BANG,
			r:        false,
			expected: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, UnaryExpr{
				Op: Token{
					Type: tt.op,
				},
				Right: LiteralExpr{
					Value: tt.r,
				},
			}.Evaluate())
		})
	}
}
