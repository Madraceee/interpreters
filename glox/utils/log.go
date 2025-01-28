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
		fmt.Printf(pattern, args...)
	}
}

func DLogf(pattern string, args ...interface{}) {
	if Debug {
		log.Printf(pattern, args...)
	}
}
