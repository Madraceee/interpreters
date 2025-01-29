package utils

import (
	"fmt"

	"github.com/madraceee/interpreters/glox/token"
)

var (
	HadError        = false
	HadRunTimeError = false
)

func Error(line int, message string) {
	Ereport(line, "", message)
}

func TError(_token token.Token, message string) {
	if _token.TokenType == token.EOF {
		Ereport(_token.Line, "at end", message)
	} else {
		Ereport(_token.Line, "at '"+_token.Lexeme+"'", message)
	}

}

func Ereport(line int, where, message string) {
	fmt.Printf("[Line %d] Error %s: %s\n", line, where, message)
	HadError = true
}
