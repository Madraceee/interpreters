package main

import (
	"github.com/madraceee/interpreters/glox/parser"
	"github.com/madraceee/interpreters/glox/token"
)

func astPrinter() {
	_ = parser.NewBinary(
		parser.NewUnary(
			token.NewToken(token.MINUS, "-", token.Object{}, 1),
			parser.NewLiteral(token.Object{
				ObjType:     token.NUMBER_TYPE,
				Value_float: 123,
			}),
		),
		token.NewToken(token.STAR, "*", token.Object{}, 1),
		parser.NewGrouping(
			parser.NewLiteral(token.Object{
				ObjType:     token.NUMBER_TYPE,
				Value_float: 45.67,
			}),
		),
	)

	// ap := parser.NewAstPrinter()
	// ap.Print(expr)
}
