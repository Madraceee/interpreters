package parser

import (
	"github.com/madraceee/glox/token"
)

type Expr interface {
	Visit() string
}

type Binary struct {
	left     Expr
	operator token.Token
	right    Expr
}

func NewBinary(left Expr, operator token.Token, right Expr) *Binary {
	return &Binary{
		left:     left,
		operator: operator,
		right:    right,
	}
}

type Grouping struct {
	expression Expr
}

func NewGrouping(expression Expr) *Grouping {
	return &Grouping{
		expression: expression,
	}
}

type Literal struct {
	value token.Object
}

func NewLiteral(value token.Object) *Literal {
	return &Literal{
		value: value,
	}
}

type Unary struct {
	operator token.Token
	right    Expr
}

func NewUnary(operator token.Token, right Expr) *Unary {
	return &Unary{
		operator: operator,
		right:    right,
	}
}
