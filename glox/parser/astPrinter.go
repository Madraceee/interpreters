package parser

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (ap *AstPrinter) Print(e Expr) {
	val, err := e.Visit(ap)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(val.(string))
}

func (ap *AstPrinter) VisitBinaryExpr(b *Binary) (interface{}, error) {
	return ap.parenthesize(b.operator.String(), b.left, b.right)
}

func (ap *AstPrinter) VisitGroupingExpr(g *Grouping) (interface{}, error) {
	return ap.parenthesize("group", g.expression)
}

func (ap *AstPrinter) VisitLiteralExpr(l *Literal) (interface{}, error) {
	return l.value.GetStringValue(), nil
}

func (ap *AstPrinter) VisitUnaryExpr(u *Unary) (interface{}, error) {
	return ap.parenthesize(u.operator.String(), u.right)
}

func (ap *AstPrinter) parenthesize(name string, exprs ...Expr) (string, error) {
	builder := strings.Builder{}

	builder.WriteString("(" + name)
	for _, expr := range exprs {
		val, _ := expr.Visit(ap)
		builder.WriteString(" " + val.(string))
	}
	builder.WriteString(")")

	return builder.String(), nil
}
