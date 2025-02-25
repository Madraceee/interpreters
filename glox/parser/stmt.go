package parser

import (
	"github.com/madraceee/interpreters/glox/token"
)

type Stmt interface {
	Visit(VisitStmt) (interface{}, error)
}

type VisitStmt interface {
	VisitBlockStmt(*Block) (interface{}, error)
	VisitExpressionStmt(*Expression) (interface{}, error)
	VisitFunctionStmt(*Function) (interface{}, error)
	VisitIfStmt(*If) (interface{}, error)
	VisitPrintStmt(*Print) (interface{}, error)
	VisitReturnStmt(*Return) (interface{}, error)
	VisitVarStmt(*Var) (interface{}, error)
	VisitWhileStmt(*While) (interface{}, error)
}

type Block struct {
	Statements []Stmt
}

func NewBlock(statements []Stmt) Stmt {
	return &Block{
		Statements: statements,
	}
}

func (expr *Block) Visit(visitor VisitStmt) (interface{}, error) {
	return visitor.VisitBlockStmt(expr)
}

type Expression struct {
	Expression Expr
}

func NewExpression(expression Expr) Stmt {
	return &Expression{
		Expression: expression,
	}
}

func (expr *Expression) Visit(visitor VisitStmt) (interface{}, error) {
	return visitor.VisitExpressionStmt(expr)
}

type Function struct {
	Name   token.Token
	Params []token.Token
	Body   []Stmt
}

func NewFunction(name token.Token, params []token.Token, body []Stmt) Stmt {
	return &Function{
		Name:   name,
		Params: params,
		Body:   body,
	}
}

func (expr *Function) Visit(visitor VisitStmt) (interface{}, error) {
	return visitor.VisitFunctionStmt(expr)
}

type If struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func NewIf(condition Expr, thenBranch Stmt, elseBranch Stmt) Stmt {
	return &If{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}
}

func (expr *If) Visit(visitor VisitStmt) (interface{}, error) {
	return visitor.VisitIfStmt(expr)
}

type Print struct {
	Expression Expr
}

func NewPrint(expression Expr) Stmt {
	return &Print{
		Expression: expression,
	}
}

func (expr *Print) Visit(visitor VisitStmt) (interface{}, error) {
	return visitor.VisitPrintStmt(expr)
}

type Return struct {
	Keyword token.Token
	Value   Expr
}

func NewReturn(keyword token.Token, value Expr) Stmt {
	return &Return{
		Keyword: keyword,
		Value:   value,
	}
}

func (expr *Return) Visit(visitor VisitStmt) (interface{}, error) {
	return visitor.VisitReturnStmt(expr)
}

type Var struct {
	Name        token.Token
	Initializer Expr
}

func NewVar(name token.Token, initializer Expr) Stmt {
	return &Var{
		Name:        name,
		Initializer: initializer,
	}
}

func (expr *Var) Visit(visitor VisitStmt) (interface{}, error) {
	return visitor.VisitVarStmt(expr)
}

type While struct {
	Condition Expr
	Body      Stmt
}

func NewWhile(condition Expr, body Stmt) Stmt {
	return &While{
		Condition: condition,
		Body:      body,
	}
}

func (expr *While) Visit(visitor VisitStmt) (interface{}, error) {
	return visitor.VisitWhileStmt(expr)
}
