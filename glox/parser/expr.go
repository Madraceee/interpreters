package parser

import (
	"github.com/madraceee/interpreters/glox/token"
)

type Expr interface {
	Visit(VisitExpr) (interface{}, error)
}

type VisitExpr interface {
	VisitAssignExpr(*Assign) (interface{}, error)
	VisitBinaryExpr(*Binary) (interface{}, error)
	VisitCallExpr(*Call) (interface{}, error)
	VisitGroupingExpr(*Grouping) (interface{}, error)
	VisitLiteralExpr(*Literal) (interface{}, error)
	VisitLogicalExpr(*Logical) (interface{}, error)
	VisitUnaryExpr(*Unary) (interface{}, error)
	VisitVariableExpr(*Variable) (interface{}, error)
}

type Assign struct {
	Name  token.Token
	Value Expr
}

func NewAssign(name token.Token, value Expr) Expr {
	return &Assign{
		Name:  name,
		Value: value,
	}
}

func (expr *Assign) Visit(visitor VisitExpr) (interface{}, error) {
	return visitor.VisitAssignExpr(expr)
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func NewBinary(left Expr, operator token.Token, right Expr) Expr {
	return &Binary{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func (expr *Binary) Visit(visitor VisitExpr) (interface{}, error) {
	return visitor.VisitBinaryExpr(expr)
}

type Call struct {
	Callee    Expr
	Paren     token.Token
	Arguments []Expr
}

func NewCall(callee Expr, paren token.Token, arguments []Expr) Expr {
	return &Call{
		Callee:    callee,
		Paren:     paren,
		Arguments: arguments,
	}
}

func (expr *Call) Visit(visitor VisitExpr) (interface{}, error) {
	return visitor.VisitCallExpr(expr)
}

type Grouping struct {
	Expression Expr
}

func NewGrouping(expression Expr) Expr {
	return &Grouping{
		Expression: expression,
	}
}

func (expr *Grouping) Visit(visitor VisitExpr) (interface{}, error) {
	return visitor.VisitGroupingExpr(expr)
}

type Literal struct {
	Value token.Object
}

func NewLiteral(value token.Object) Expr {
	return &Literal{
		Value: value,
	}
}

func (expr *Literal) Visit(visitor VisitExpr) (interface{}, error) {
	return visitor.VisitLiteralExpr(expr)
}

type Logical struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func NewLogical(left Expr, operator token.Token, right Expr) Expr {
	return &Logical{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

func (expr *Logical) Visit(visitor VisitExpr) (interface{}, error) {
	return visitor.VisitLogicalExpr(expr)
}

type Unary struct {
	Operator token.Token
	Right    Expr
}

func NewUnary(operator token.Token, right Expr) Expr {
	return &Unary{
		Operator: operator,
		Right:    right,
	}
}

func (expr *Unary) Visit(visitor VisitExpr) (interface{}, error) {
	return visitor.VisitUnaryExpr(expr)
}

type Variable struct {
	Name token.Token
}

func NewVariable(name token.Token) Expr {
	return &Variable{
		Name: name,
	}
}

func (expr *Variable) Visit(visitor VisitExpr) (interface{}, error) {
	return visitor.VisitVariableExpr(expr)
}
