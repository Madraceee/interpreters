package main

import (
	"fmt"

	"github.com/madraceee/glox/parser"
	"github.com/madraceee/glox/token"
)

func main() {
	expr := parser.NewBinary(
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

	fmt.Println(expr.Visit())
}
