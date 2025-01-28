package token

import "strconv"

type ObjectType int

const (
	NOT_ASSIGNED_TYPE ObjectType = iota
	NUMBER_TYPE
	STRING_TYPE
	BOOL_TYPE
)

type Object struct {
	ObjType     ObjectType
	Value_str   string
	Value_float float64
	Value_bool  bool
}

func GetStringValue(o Object) string {
	if o.ObjType == NUMBER_TYPE {
		return strconv.FormatFloat(o.Value_float, 'f', -1, 64)
	} else if o.ObjType == STRING_TYPE {
		return o.Value_str
	} else if o.ObjType == BOOL_TYPE {
		if o.Value_bool {
			return "true"
		} else {
			return "false"
		}
	}
	return ""
}

type Token struct {
	TokenType TokenType
	Lexeme    string
	Literal   Object
	Line      int
}

func NewToken(tokenType TokenType, lexeme string, literal Object, line int) Token {
	return Token{
		TokenType: tokenType,
		Lexeme:    lexeme,
		Literal:   literal,
		Line:      line,
	}
}

func (t Token) String() string {
	return GetTokenType(t.TokenType) + " " + t.Lexeme + " " + GetStringValue(t.Literal)
}
