package utils

import "fmt"

var (
	HadError = false
)

func Error(line int, message string) {
	Eeport(line, "", message)
}

func Eeport(line int, where, message string) {
	fmt.Printf("[Line %d] Error %s: %s\n", line, where, message)
	HadError = true
}
