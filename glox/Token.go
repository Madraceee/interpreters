package main

import "strconv"

type ObjectType int

const (
	NOT_ASSIGNED_TYPE ObjectType = iota
	NUMBER_TYPE
	STRING_TYPE
)

type Object struct {
	objType     ObjectType
	value_str   string
	value_float float64
}

func (o *Object) getStringValue() string {
	if o.objType == NUMBER_TYPE {
		return strconv.FormatFloat(o.value_float, 'f', -1, 64)
	} else if o.objType == STRING_TYPE {
		return o.value_str
	}
	return ""
}

type Token struct {
	tokenType int
	lexeme    string
	literal   Object
	line      int
}

func NewToken(tokenType int, lexeme string, literal Object, line int) Token {
	return Token{
		tokenType: tokenType,
		lexeme:    lexeme,
		literal:   literal,
		line:      line,
	}
}

func (t Token) String() string {
	return getTokenType(t.tokenType) + " " + t.lexeme + " " + t.literal.getStringValue()
}
