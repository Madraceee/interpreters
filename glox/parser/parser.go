package parser

import (
	"errors"

	"github.com/madraceee/interpreters/glox/token"
	"github.com/madraceee/interpreters/glox/utils"
)

type Parser struct {
	Tokens  []token.Token
	Current int
}

func NewParser(tokens []token.Token) *Parser {
	return &Parser{
		Tokens:  tokens,
		Current: 0,
	}
}

func (p *Parser) Parse() []Stmt {
	statements := make([]Stmt, 0)
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			utils.DLogf("%v\n", err)
			continue
		}
		statements = append(statements, stmt)
	}

	return statements
}

// Functions for stmt.go
// Add syncrhonize
func (p *Parser) declaration() (Stmt, error) {
	if p.match(token.VAR) {
		return p.varDeclaration()
	}

	return p.statement()
}
func (p *Parser) statement() (Stmt, error) {
	if p.match(token.PRINT) {
		return p.printStatement()
	}
	if p.match(token.LEFT_BRACE) {
		stmt, err := p.block()
		return NewBlock(stmt), err
	}

	return p.expressionStatement()
}

func (p *Parser) block() ([]Stmt, error) {
	statements := make([]Stmt, 0)
	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		stmts, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmts)
	}

	p.consume(token.RIGHT_BRACE, "Expect '}' after block.")
	return statements, nil
}

func (p *Parser) printStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	p.consume(token.SEMICOLON, "Expect ';' after value")
	return &Print{
		Expression: expr,
	}, nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	p.consume(token.SEMICOLON, "Expect ';' after value")
	return &Expression{
		Expression: expr,
	}, nil
}

func (p *Parser) varDeclaration() (Stmt, error) {
	name, err := p.consume(token.IDENTIFIER, "Expecting variable name")
	if err != nil {
		return nil, err
	}

	var initializer Expr
	if p.match(token.EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(token.SEMICOLON, "Expect ';' after variable declaration")
	if err != nil {
		return nil, err
	}

	return &Var{
		Name:        *name,
		Initializer: initializer,
	}, nil
}

// Functions for expr.go
func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	if p.match(token.EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		switch expr.(type) {
		case *Variable:
			name := expr.(*Variable).Name
			return NewAssign(name, value), nil
		}

		return nil, ParserError(equals, "Invalid assignment target")

	}

	return expr, nil
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, *operator, right)
	}

	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, *operator, right)
	}

	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, *operator, right)
	}

	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, *operator, right)
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right, err := p.unary()
		return NewUnary(*operator, right), err
	}

	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.match(token.FALSE) {
		return NewLiteral(token.Object{
			ObjType:    token.BOOL_TYPE,
			Value_bool: false,
		}), nil
	}
	if p.match(token.TRUE) {
		return NewLiteral(token.Object{
			ObjType:    token.BOOL_TYPE,
			Value_bool: true,
		}), nil
	}
	if p.match(token.NIL) {
		return NewLiteral(token.Object{
			ObjType:   token.STRING_TYPE,
			Value_str: "nil",
		}), nil
	}

	if p.match(token.NUMBER, token.STRING) {
		return NewLiteral(token.Object{
			ObjType:     p.previous().Literal.ObjType,
			Value_str:   p.previous().Literal.Value_str,
			Value_float: p.previous().Literal.Value_float,
		}), nil
	}

	if p.match(token.IDENTIFIER) {
		return NewVariable(*p.previous()), nil
	}

	if p.match(token.LEFT_PARAN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(token.RIGHT_PARAN, "Expect  ')' after expression.")

		if err != nil {
			return nil, err
		}
		return NewGrouping(expr), nil
	}

	return nil, ParserError(p.peek(), "Expect expression.")
}

func (p *Parser) consume(_type token.TokenType, message string) (*token.Token, error) {
	if p.check(_type) {
		return p.advance(), nil
	}

	return nil, ParserError(p.peek(), message)
}

func (p *Parser) match(types ...token.TokenType) bool {
	for _, i := range types {
		if p.check(i) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(_type token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().TokenType == _type
}

func (p *Parser) advance() *token.Token {
	if !p.isAtEnd() {
		p.Current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return token.GetTokenType(p.peek().TokenType) == "EOF"
}

func (p *Parser) peek() *token.Token {
	return &p.Tokens[p.Current]
}

func (p *Parser) previous() *token.Token {
	return &p.Tokens[p.Current-1]
}

func (p *Parser) Synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().TokenType == token.SEMICOLON {
			return
		}

		switch p.peek().TokenType {
		case token.CLASS, token.FUN, token.VAR, token.FOR, token.IF, token.WHILE, token.PRINT, token.RETURN:
			return
		}

		p.advance()
	}
}

func ParserError(token *token.Token, message string) error {
	utils.TError(*token, message)
	return errors.New("ParserError: error while parsing")
}
