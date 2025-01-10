package main

import "fmt"

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
		stmt, err := p.statement()
		if err != nil {
			// synchronize??
			return nil, err
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

func (p *Parser) statement() (Stmt, error) {
	if p.match(PRINT) {
		return p.printStatement()
	}

	return p.expressionStatement()
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
	return p.equality()
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
