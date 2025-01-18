package parser

import (
	"strings"
)

func (b *Binary) Visit() string {
	return parenthesize(b.operator.String(), b.left, b.right)
}

func (g *Grouping) Visit() string {
	return parenthesize("group", g.expression)
}

func (l *Literal) Visit() string {
	return l.value.GetStringValue()
}

func (u *Unary) Visit() string {
	return parenthesize(u.operator.String(), u.right)
}

func parenthesize(name string, exprs ...Expr) string {
	builder := strings.Builder{}

	builder.WriteString("(" + name)
	for _, expr := range exprs {
		builder.WriteString(" " + expr.Visit())
	}
	builder.WriteString(")")

	return builder.String()
}
