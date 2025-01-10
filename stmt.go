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
