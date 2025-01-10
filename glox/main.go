package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	debug = false
)

var (
	hadError = false
)

func dLogf(pattern string, args ...string) {
	if debug {
		log.Printf(pattern, args)
	}
}

func main() {
	args := os.Args
	if len(args) > 2 {
		fmt.Printf("%s\n", "Usage: glox [script]")
		os.Exit(1)
	} else if len(args) == 2 {
		runFile(args[1])
	} else {
		runPrompt()
	}
}

func run(source string) {
	scanner := NewScanner(source)
	tokens := scanner.scanTokens()
	for _, v := range tokens {
		fmt.Println(v)
	}
}

func runFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		dLogf("Cannot open file %s: %v", fileName, err.Error())
		fmt.Printf("Cannot open file %s\n", fileName)
		os.Exit(1)
	}

	content, err := io.ReadAll(file)
	if err != nil {
		dLogf("Error reading file%s: %v", fileName, err.Error())
		fmt.Printf("Error reading file %s\n", fileName)
		os.Exit(1)
	}
	run(string(content))
	if hadError {
		os.Exit(1)
	}
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")
		text, _ := reader.ReadString('\n')
		if len(text) == 0 {
			break
		}
		run(text)
		hadError = false
	}
}

func error(line int, message string) {
	report(line, "", message)
}

func report(line int, where, message string) {
	fmt.Printf("[Line %d] Error %s: %s\n", line, where, message)
	hadError = true
}
