package main

import (
	"fmt"
)

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) (*Parser, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("cannot parse empty list of tokens")
	}

	lastToken := tokens[len(tokens)-1]
	if lastToken.Type != EOF {
		tokens = append(tokens, Token{
			Type:   EOF,
			Lexeme: "",
			Line:   lastToken.Line,
		})
	}

	return &Parser{
		tokens: tokens,
	}, nil
}

func (p *Parser) Parse() ([]Stmt, error) {
	stmts := []Stmt{}
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			p.Synchronize()
			continue
		}

		stmts = append(stmts, stmt)
	}
	return stmts, nil
}

func (p *Parser) match(tokenTypes ...TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(tokenType TokenType, message string) (Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}

	return Token{}, p.error(p.peek(), message)
}

func (p *Parser) error(token Token, message string) error {
	TokenError(token, message)
	return ParseError{}
}

func (p *Parser) declaration() (Stmt, error) {
	if p.match(VAR) {
		return p.varDeclaration()
	} else {
		return p.statement()
	}
}

func (p *Parser) varDeclaration() (Stmt, error) {
	name, err := p.consume(IDENTIFIER, "Expect variable name")
	if err != nil {
		return nil, err
	}

	var initializer Expr
	if p.match(EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	p.consume(SEMICOLON, "Expect ';' after variable declaration.")
	return VarStmt{
		Name: name,
		Expr: initializer,
	}, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(PRINT) {
		return p.printStatement()
	}

	if p.match(LEFT_BRACE) {
		stmts, err := p.block()
		return BlockStmt{
			Stmts: stmts,
		}, err
	}

	return p.expressionStatement()
}

func (p *Parser) block() ([]Stmt, error) {
	stmts := []Stmt{}
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		declaration, err := p.declaration()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, declaration)
	}

	p.consume(RIGHT_BRACE, "Expect '}' after block.")
	return stmts, nil
}

func (p *Parser) printStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	p.consume(SEMICOLON, "Expect ';' after value.")
	return PrintStmt{
		Expr: expr,
	}, nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	p.consume(SEMICOLON, "Expect ';' after expression.")
	return ExprStmt{
		Expr: expr,
	}, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	if p.match(EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if varExpr, ok := expr.(VariableExpr); ok {
			return AssignExpr{
				Name:  varExpr.Name,
				Value: value,
			}, nil
		}

		Error(equals.Line, "Invalid assignment target.")
	}

	return expr, nil
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpr{
			Left:  expr,
			Op:    operator,
			Right: right,
		}
	}
	return expr, err
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpr{
			Op:    operator,
			Left:  expr,
			Right: right,
		}
	}
	return expr, err
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpr{
			Op:    operator,
			Left:  expr,
			Right: right,
		}
	}

	return expr, err
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = BinaryExpr{
			Op:    operator,
			Left:  expr,
			Right: right,
		}
	}

	return expr, err
}

func (p *Parser) unary() (Expr, error) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, err := p.unary()
		return UnaryExpr{
			Op:    operator,
			Right: right,
		}, err
	}
	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.match(FALSE) {
		return LiteralExpr{
			Value: false,
		}, nil
	}
	if p.match(TRUE) {
		return LiteralExpr{
			Value: true,
		}, nil
	}
	if p.match(NIL) {
		return LiteralExpr{
			Value: nil,
		}, nil
	}

	if p.match(NUMBER, STRING) {
		return LiteralExpr{
			Value: p.previous().Literal,
		}, nil
	}

	if p.match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return GroupingExpr{
			Expression: expr,
		}, err
	}

	if p.match(IDENTIFIER) {
		return VariableExpr{
			Name: p.previous(),
		}, nil
	}

	return nil, p.error(p.peek(), "Expect expression.")
}

func (p *Parser) Synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == SEMICOLON {
			return
		}

		switch p.peek().Type {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
			return
		}

		p.advance()
	}
}
