package parser

import (
	"github.com/madraceee/glox/token"
)

type Expr interface {
	Visit(VisitExpr) (interface{}, error)
}

type VisitExpr interface {
	VisitBinaryExpr(*Binary) (interface{}, error)
	VisitGroupingExpr(*Grouping) (interface{}, error)
	VisitLiteralExpr(*Literal) (interface{}, error)
	VisitUnaryExpr(*Unary) (interface{}, error)
}

type Binary struct {
	left     Expr
	operator token.Token
	right    Expr
}

func NewBinary(left Expr, operator token.Token, right Expr) Expr {
	return &Binary{
		left:     left,
		operator: operator,
		right:    right,
	}
}

func (expr *Binary) Visit(visitor VisitExpr) (interface{}, error) {
	return visitor.VisitBinaryExpr(expr)
}

type Grouping struct {
	expression Expr
}

func NewGrouping(expression Expr) Expr {
	return &Grouping{
		expression: expression,
	}
}

func (expr *Grouping) Visit(visitor VisitExpr) (interface{}, error) {
	return visitor.VisitGroupingExpr(expr)
}

type Literal struct {
	value token.Object
}

func NewLiteral(value token.Object) Expr {
	return &Literal{
		value: value,
	}
}

func (expr *Literal) Visit(visitor VisitExpr) (interface{}, error) {
	return visitor.VisitLiteralExpr(expr)
}

type Unary struct {
	operator token.Token
	right    Expr
}

func NewUnary(operator token.Token, right Expr) Expr {
	return &Unary{
		operator: operator,
		right:    right,
	}
}

func (expr *Unary) Visit(visitor VisitExpr) (interface{}, error) {
	return visitor.VisitUnaryExpr(expr)
}
