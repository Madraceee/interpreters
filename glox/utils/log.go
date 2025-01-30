package utils

import (
	"fmt"
	"log"
)

var (
	Debug bool = false
)

func DPrintf(pattern string, args ...interface{}) {
	if Debug {
		Print(pattern, args...)
	}
}

func DLogf(pattern string, args ...interface{}) {
	if Debug {
		log.Printf(pattern, args...)
	}
}

func Print(pattern string, args ...interface{}) {
	fmt.Printf(pattern, args...)
}
