package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/madraceee/interpreters/glox/interpreter"
	"github.com/madraceee/interpreters/glox/parser"
	"github.com/madraceee/interpreters/glox/scanner"
	"github.com/madraceee/interpreters/glox/utils"
)

func main() {
	args := os.Args
	utils.Debug = false
	if len(args) > 3 {
		fmt.Printf("%s\n", "Usage: glox [debug] [script]")
		os.Exit(1)
	} else if len(args) == 3 {
		if args[1] == "debug" {
			utils.Debug = true
		}
		runFile(args[2])
	} else if len(args) == 2 {
		runFile(args[1])
	} else {
		runPrompt()
	}
}

func run(source string) {
	utils.DPrintf("%s\n", "----Scanning----")
	scanner := scanner.NewScanner(source)
	tokens := scanner.ScanTokens()
	for _, v := range tokens {
		utils.DPrintf("%s\n", v)
	}

	utils.DPrintf("%s\n", "----Parsing----")
	// ap := parser.NewAstPrinter()
	parser := parser.NewParser(tokens)
	statements := parser.Parse()
	// if utils.Debug {
	// 	ap.Print(statements)
	// }
	if utils.HadError {
		return
	}

	utils.DPrintf("%s\n", "----Interpreter----")
	gloxInterpreter := interpreter.NewInterpreter()
	gloxInterpreter.Interpret(statements)
}

func runFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		utils.DLogf("Cannot open file %s: %v", fileName, err.Error())
		fmt.Printf("Cannot open file %s\n", fileName)
		os.Exit(1)
	}

	content, err := io.ReadAll(file)
	if err != nil {
		utils.DLogf("Error reading file%s: %v", fileName, err.Error())
		fmt.Printf("Error reading file %s\n", fileName)
		os.Exit(1)
	}
	run(string(content))
	if utils.HadError {
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
		fmt.Println("")
		utils.HadError = false
	}
}
