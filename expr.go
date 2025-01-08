package main

import "fmt"

type Expr interface {
	Print() string
	Evaluate() (any, error)
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

func (expr AssignExpr) Evaluate() (any, error) {
	return nil, nil
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

func (expr BinaryExpr) Evaluate() (any, error) {
	left, err := expr.Left.Evaluate()
	if err != nil {
		return nil, err
	}

	right, err := expr.Right.Evaluate()
	if err != nil {
		return nil, err
	}

	var leftFloat, rightFloat float64
	var leftString, rightString string
	var plusFloat bool

	// Type Checking
	switch expr.Op.Type {
	case BANG_EQUAL, EQUAL_EQUAL, GREATER, GREATER_EQUAL, LESS, LESS_EQUAL, MINUS, SLASH, STAR:
		// should be floats

		if leftCast, ok := left.(float64); ok {
			leftFloat = leftCast
		} else {
			return nil, fmt.Errorf("left operand of binary '%s' expression should be number: %v", expr.Op.Lexeme, left)
		}

		if rightCast, ok := right.(float64); ok {
			rightFloat = rightCast
		} else {
			return nil, fmt.Errorf("right operand of binary '%s' expression should be number: %v", expr.Op.Lexeme, left)
		}
	case PLUS:
		// should be either be floats or strings

		leftFloatCast, leftFloatOK := left.(float64)
		rightFloatCast, rightFloatOK := right.(float64)

		if leftFloatOK && rightFloatOK {
			// use float for plus if left and right both casted to a float
			plusFloat = true
			leftFloat = leftFloatCast
			rightFloat = rightFloatCast
		}

		leftStringCast, leftStringOK := left.(string)
		rightStringCast, rightStringOK := right.(string)
		if leftStringOK && rightStringOK {
			plusFloat = false
			leftString = leftStringCast
			rightString = rightStringCast
		}

		// If left and right are not both floats and they are not both strings,
		// return a type error
		if !(leftFloatOK && rightFloatOK) && !(leftStringOK && rightStringOK) {
			return nil, fmt.Errorf("left and right operant of '+' expression should both be numbers or both be strings: %v, %v", left, right)
		}
	}

	switch expr.Op.Type {
	case BANG_EQUAL:
		return !isEqual(left, right), nil
	case EQUAL_EQUAL:
		return isEqual(left, right), nil
	case GREATER:
		return leftFloat > rightFloat, nil
	case GREATER_EQUAL:
		return leftFloat >= rightFloat, nil
	case LESS:
		return leftFloat < rightFloat, nil
	case LESS_EQUAL:
		return leftFloat <= rightFloat, nil
	case MINUS:
		return leftFloat - rightFloat, nil
	case SLASH:
		return leftFloat / rightFloat, nil
	case STAR:
		return leftFloat * rightFloat, nil
	case PLUS:
		if plusFloat {
			return leftFloat + rightFloat, nil
		}
		return leftString + rightString, nil
	}

	// unreachable
	return nil, nil
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

func (expr CallExpr) Evaluate() (any, error) {
	return nil, nil
}

func (expr CallExpr) Print() string {
	return "<print-not-implemented>"
}

// GetExpr //////////////////////////////////////
type GetExpr struct {
	Object Expr
	Name   Token
}

func (expr GetExpr) Evaluate() (any, error) {
	return nil, nil
}

func (expr GetExpr) Print() string {
	return "<print-not-implemented>"
}

// GroupingExpr /////////////////////////////////
type GroupingExpr struct {
	Expression Expr
}

func (expr GroupingExpr) Evaluate() (any, error) {
	return nil, nil
}

func (expr GroupingExpr) Print() string {
	return Parenthesize("group", expr.Expression)
}

// LiteralExpr //////////////////////////////////
type LiteralExpr struct {
	Value any
}

func (expr LiteralExpr) Evaluate() (any, error) {
	return expr.Value, nil
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

func (expr LogicalExpr) Evaluate() (any, error) {
	return nil, nil
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

func (expr SetExpr) Evaluate() (any, error) {
	return nil, nil
}

func (expr SetExpr) Print() string {
	return "<print-not-implemented>"
}

// SuperExpr ////////////////////////////////////
type SuperExpr struct {
	Keyword Token
	Method  Token
}

func (expr SuperExpr) Evaluate() (any, error) {
	return nil, nil
}

func (expr SuperExpr) Print() string {
	return "<print-not-implemented>"
}

// ThisExpr /////////////////////////////////////
type ThisExpr struct {
	Keyword Token
}

func (expr ThisExpr) Evaluate() (any, error) {
	return nil, nil
}

func (expr ThisExpr) Print() string {
	return "<print-not-implemented>"
}

// UnaryExpr ////////////////////////////////////
type UnaryExpr struct {
	Op    Token
	Right Expr
}

func (expr UnaryExpr) Evaluate() (any, error) {
	right, err := expr.Right.Evaluate()
	if err != nil {
		return nil, err
	}

	switch expr.Op.Type {
	case MINUS:
		// TODO: what if right is not a float64?
		rightFloat, ok := right.(float64)
		if !ok {
			return nil, fmt.Errorf("operand for unary '-' expression should be a number: %v", right)
		}
		return -rightFloat, nil
	case BANG:
		return !isTruthy(right), nil
	}

	// unreachable
	return nil, nil
}

func (expr UnaryExpr) Print() string {
	return Parenthesize(expr.Op.Lexeme, expr.Right)
}

// VariableExpr /////////////////////////////////
type VariableExpr struct {
	Name Token
}

func (expr VariableExpr) Evaluate() (any, error) {
	return nil, nil
}

func (expr VariableExpr) Print() string {
	return "<print-not-implemented>"
}
