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
