package main

import "fmt"

type Stmt interface {
	Execute() error
}

type ExprStmt struct {
	Expr Expr
}

func (stmt ExprStmt) Execute() error {
	_, err := stmt.Expr.Evaluate()
	return err
}

type PrintStmt struct {
	Expr Expr
}

func (stmt PrintStmt) Execute() error {
	v, err := stmt.Expr.Evaluate()
	if err != nil {
		return err
	}
	fmt.Println(v)
	return nil
}

type VarStmt struct {
	Name Token
	Expr Expr
}

func (stmt VarStmt) Execute() error {
	var v any
	if stmt.Expr != nil {
		vv, err := stmt.Expr.Evaluate()
		if err != nil {
			return err
		}
		v = vv
	}
	environment.Define(stmt.Name.Lexeme, v)
	return nil
}

type BlockStmt struct {
	Stmts []Stmt
}

func (stmt BlockStmt) Execute() error {
	return executeBlock(stmt.Stmts, NewEnvironment(environment))
}

func executeBlock(stmts []Stmt, env *Environment) error {
	prev := environment
	environment = env
	defer func() {
		environment = prev
	}()
	for _, stmt := range stmts {
		err := stmt.Execute()
		if err != nil {
			return err
		}
	}
	return nil
}
