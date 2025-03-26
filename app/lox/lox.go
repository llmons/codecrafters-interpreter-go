package lox

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/app/scanner"
	"github.com/codecrafters-io/interpreter-starter-go/app/util"
)

func Main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	runFile(os.Args[2])
}

func runFile(path string) {
	fileContents, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	run(string(fileContents))
	if util.HadError {
		os.Exit(65)
	}
}

// func runPrompt() {
// 	for {
// 		fmt.Print("> ")
// 		var input string
// 		_, err := fmt.Scanln(&input)
// 		if err != nil {
// 			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
// 			os.Exit(1)
// 		}
// 		run(input)
// 		util.HadError = false
// 	}
// }

func run(source string) {
	scanner := scanner.NewScanner(source)
	tokens := scanner.ScanTokens()
	for _, token := range tokens {
		fmt.Println(token)
	}
}
