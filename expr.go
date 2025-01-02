package main

import "fmt"

type Expr interface {
	Print() string
}

type AssignExpr struct {
	Name  Token
	Value Expr
}

func (a AssignExpr) Print() string {
	return "<print-not-implemented>"
}

type BinaryExpr struct {
	Op    Token
	Left  Expr
	Right Expr
}

func (expr BinaryExpr) Print() string {
	return Parenthesize(expr.Op.Lexeme, expr.Left, expr.Right)
}

type CallExpr struct {
	Callee Expr
	Paren  Token
	Args   []Expr
}

func (expr CallExpr) Print() string {
	return "<print-not-implemented>"
}

type GetExpr struct {
	Object Expr
	Name   Token
}

func (expr GetExpr) Print() string {
	return "<print-not-implemented>"
}

type GroupingExpr struct {
	Expression Expr
}

func (expr GroupingExpr) Print() string {
	return Parenthesize("group", expr.Expression)
}

type LiteralExpr struct {
	Value any
}

func (expr LiteralExpr) Print() string {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value)
}

type LogicalExpr struct {
	Left  Expr
	Right Expr
	Op    Token
}

func (expr LogicalExpr) Print() string {
	return "<print-not-implemented>"
}

type SetExpr struct {
	Object Expr
	Name   Token
	Value  Expr
}

func (expr SetExpr) Print() string {
	return "<print-not-implemented>"
}

type SuperExpr struct {
	Keyword Token
	Method  Token
}

func (expr SuperExpr) Print() string {
	return "<print-not-implemented>"
}

type ThisExpr struct {
	Keyword Token
}

func (expr ThisExpr) Print() string {
	return "<print-not-implemented>"
}

type UnaryExpr struct {
	Op    Token
	Right Expr
}

func (expr UnaryExpr) Print() string {
	return Parenthesize(expr.Op.Lexeme, expr.Right)
}

type VariableExpr struct {
	Name Token
}

func (expr VariableExpr) Print() string {
	return "<print-not-implemented>"
}

func Parenthesize(name string, exprs ...Expr) string {
	s := "(" + name
	for _, expr := range exprs {
		s += " " + expr.Print()
	}
	return s + ")"
}
