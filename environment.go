package main

import "fmt"

type Environment struct {
	enclosing *Environment
	values    map[string]any
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		values:    make(map[string]any),
	}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Get(name Token) (any, error) {
	v, ok := e.values[name.Lexeme]
	if ok {
		return v, nil
	}

	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}

	return nil, fmt.Errorf("undefined variable '%s'", name.Lexeme)
}

func (e *Environment) Assign(name Token, v any) error {
	if _, ok := e.values[name.Lexeme]; ok {
		e.values[name.Lexeme] = v
		return nil
	}

	if e.enclosing != nil {
		return e.enclosing.Assign(name, v)
	}

	return fmt.Errorf("undefined variable '%s'", name.Lexeme)
}
