package token

import "strconv"

type ObjectType int

const (
	NOT_ASSIGNED_TYPE ObjectType = iota
	NUMBER_TYPE
	STRING_TYPE
)

type Object struct {
	ObjType     ObjectType
	Value_str   string
	Value_float float64
}

func (o *Object) GetStringValue() string {
	if o.ObjType == NUMBER_TYPE {
		return strconv.FormatFloat(o.Value_float, 'f', -1, 64)
	} else if o.ObjType == STRING_TYPE {
		return o.Value_str
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
	return getTokenType(t.tokenType) + " " + t.lexeme + " " + t.literal.GetStringValue()
}
