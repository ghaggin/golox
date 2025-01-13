package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefineAndGetNoEnclosing(t *testing.T) {
	testCases := []struct {
		name    string
		varName string
		v       any
		want    func(any, error)
	}{
		{
			name:    "number",
			v:       3.,
			varName: "test",
			want: func(v any, err error) {
				require.NoError(t, err)
				v, ok := v.(float64)
				require.True(t, ok)
				assert.Equal(t, 3., v)
			},
		},
		{
			name:    "string",
			v:       "value",
			varName: "test",
			want: func(v any, err error) {
				require.NoError(t, err)
				v, ok := v.(string)
				require.True(t, ok)
				assert.Equal(t, "value", v)
			},
		},
		{
			name:    "bool",
			v:       true,
			varName: "test",
			want: func(v any, err error) {
				require.NoError(t, err)
				v, ok := v.(bool)
				require.True(t, ok)
				assert.Equal(t, true, v)
			},
		},
		{
			name:    "error",
			v:       true,
			varName: "invalid",
			want: func(_ any, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEnvironment(nil)
			e.Define("test", tt.v)
			tt.want(e.Get(Token{Lexeme: tt.varName}))
		})
	}
}

func TestAssignNoEnclosing(t *testing.T) {
	e := NewEnvironment(nil)

	e.Define("a", 3.)
	v, err := e.Get(Token{Lexeme: "a"})
	require.NoError(t, err)
	vv, ok := v.(float64)
	require.True(t, ok)
	assert.Equal(t, 3., vv)

	err = e.Assign(Token{Lexeme: "b"}, nil)
	assert.Error(t, err)

	err = e.Assign(Token{Lexeme: "a"}, 4.)
	require.NoError(t, err)
	v, err = e.Get(Token{Lexeme: "a"})
	require.NoError(t, err)
	vv, ok = v.(float64)
	require.True(t, ok)
	assert.Equal(t, 4., vv)
}

func TestEnclosing(t *testing.T) {
	e1 := NewEnvironment(nil)
	e2 := NewEnvironment(e1)

	// a is in both
	e1.Define("a", 1.)
	e2.Define("a", 1.1)

	// b is only in e1
	e1.Define("b", 2.)

	// c is only in e2
	e2.Define("c", 3.)

	assertGet := func(e *Environment, name string, expected float64) {
		v, err := e.Get(Token{Lexeme: name})
		require.NoError(t, err)
		vv, ok := v.(float64)
		require.True(t, ok)
		assert.Equal(t, expected, vv)
	}

	assertGet(e1, "a", 1.)
	assertGet(e2, "a", 1.1)
	assertGet(e1, "b", 2.)
	assertGet(e2, "b", 2.)
	assertGet(e2, "c", 3.)

	_, err := e1.Get(Token{Lexeme: "c"})
	assert.Error(t, err)

	err = e2.Assign(Token{Lexeme: "b"}, 2.2)
	require.NoError(t, err)
	assertGet(e1, "b", 2.2)
	assertGet(e2, "b", 2.2)
}
