package main

import "fmt"

type Expr interface {
	Print() string
	Evaluate() any
}

func Parenthesize(name string, exprs ...Expr) string {
	s := "(" + name
	for _, expr := range exprs {
		s += " " + expr.Print()
	}
	return s + ")"
}

func isTruthy(v any) bool {
	if v == nil {
		return false
	}

	if b, ok := v.(bool); ok {
		return b
	}

	return true
}

func isEqual(l, r any) bool {
	if l == nil && r == nil {
		return true
	}

	if l == nil || r == nil {
		return false
	}

	return l == r
}

// AssignExpr ///////////////////////////////////
type AssignExpr struct {
	Name  Token
	Value Expr
}

func (expr AssignExpr) Evaluate() any {
	return nil
}

func (expr AssignExpr) Print() string {
	return "<print-not-implemented>"
}

// BinaryExpr ///////////////////////////////////
type BinaryExpr struct {
	Op    Token
	Left  Expr
	Right Expr
}

func (expr BinaryExpr) Evaluate() any {
	left := expr.Left.Evaluate()
	right := expr.Right.Evaluate()

	switch expr.Op.Type {
	case BANG_EQUAL:
		return !isEqual(left, right)
	case EQUAL_EQUAL:
		return isEqual(left, right)
	case GREATER:
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		return left.(float64) >= right.(float64)
	case LESS:
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		return left.(float64) <= right.(float64)
	case MINUS:
		return left.(float64) - right.(float64)
	case SLASH:
		return left.(float64) / right.(float64)
	case STAR:
		return left.(float64) * right.(float64)
	case PLUS:
		leftString, leftOK := left.(string)
		rightString, rightOK := right.(string)
		if leftOK && rightOK {
			return leftString + rightString
		}

		leftFloat, leftOK := left.(float64)
		rightFloat, rightOK := right.(float64)
		if leftOK && rightOK {
			return leftFloat + rightFloat
		}
	}

	// unreachable
	return nil
}

func (expr BinaryExpr) Print() string {
	return Parenthesize(expr.Op.Lexeme, expr.Left, expr.Right)
}

// CallExpr /////////////////////////////////////
type CallExpr struct {
	Callee Expr
	Paren  Token
	Args   []Expr
}

func (expr CallExpr) Evaluate() any {
	return nil
}

func (expr CallExpr) Print() string {
	return "<print-not-implemented>"
}

// GetExpr //////////////////////////////////////
type GetExpr struct {
	Object Expr
	Name   Token
}

func (expr GetExpr) Evaluate() any {
	return nil
}

func (expr GetExpr) Print() string {
	return "<print-not-implemented>"
}

// GroupingExpr /////////////////////////////////
type GroupingExpr struct {
	Expression Expr
}

func (expr GroupingExpr) Evaluate() any {
	return nil
}

func (expr GroupingExpr) Print() string {
	return Parenthesize("group", expr.Expression)
}

// LiteralExpr //////////////////////////////////
type LiteralExpr struct {
	Value any
}

func (expr LiteralExpr) Evaluate() any {
	return expr.Value
}

func (expr LiteralExpr) Print() string {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", expr.Value)
}

// LogicalExpr //////////////////////////////////
type LogicalExpr struct {
	Left  Expr
	Right Expr
	Op    Token
}

func (expr LogicalExpr) Evaluate() any {
	return nil
}

func (expr LogicalExpr) Print() string {
	return "<print-not-implemented>"
}

// SetExpr //////////////////////////////////////
type SetExpr struct {
	Object Expr
	Name   Token
	Value  Expr
}

func (expr SetExpr) Evaluate() any {
	return nil
}

func (expr SetExpr) Print() string {
	return "<print-not-implemented>"
}

// SuperExpr ////////////////////////////////////
type SuperExpr struct {
	Keyword Token
	Method  Token
}

func (expr SuperExpr) Evaluate() any {
	return nil
}

func (expr SuperExpr) Print() string {
	return "<print-not-implemented>"
}

// ThisExpr /////////////////////////////////////
type ThisExpr struct {
	Keyword Token
}

func (expr ThisExpr) Evaluate() any {
	return nil
}

func (expr ThisExpr) Print() string {
	return "<print-not-implemented>"
}

// UnaryExpr ////////////////////////////////////
type UnaryExpr struct {
	Op    Token
	Right Expr
}

func (expr UnaryExpr) Evaluate() any {
	right := expr.Right.Evaluate()

	switch expr.Op.Type {
	case MINUS:
		// TODO: what if right is not a float64?
		return -(right.(float64))
	case BANG:
		return !isTruthy(right)
	}

	// unreachable
	return nil
}

func (expr UnaryExpr) Print() string {
	return Parenthesize(expr.Op.Lexeme, expr.Right)
}

// VariableExpr /////////////////////////////////
type VariableExpr struct {
	Name Token
}

func (expr VariableExpr) Evaluate() any {
	return nil
}

func (expr VariableExpr) Print() string {
	return "<print-not-implemented>"
}
