package main

import (
	"fmt"
	"os"
)

func getFileContents() []byte {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	return fileContents
}

func main() {
	fileContents := getFileContents()

	mp := map[byte]string{
		'(': "LEFT_PAREN",
		')': "RIGHT_PAREN",
		'{': "LEFT_BRACE",
		'}': "RIGHT_BRACE",
		',': "COMMA",
		'.': "DOT",
		'-': "MINUS",
		'+': "PLUS",
		';': "SEMICOLON",
		'*': "STAR",
	}
	line, hasErr := 1, false
outer:
	for i := 0; i < len(fileContents); i++ {
		token := fileContents[i]
		switch token {
		case '\n':
			line++
		case ' ', '\t', '\r':
		case '/':
			if i < len(fileContents)-1 && fileContents[i+1] == '/' {
				for i < len(fileContents) && fileContents[i] != '\n' {
					i++
				}
				line++
			} else {
				fmt.Printf("SLASH / null\n")
			}
		case '"':
			j := i + 1
			for ; j < len(fileContents) && fileContents[j] != '"'; j++ {
				if fileContents[j] == '\n' {
					line++
				}
			}
			if j == len(fileContents) {
				fmt.Fprintf(os.Stderr, "[line %d] Error: Unterminated string.\n", line)
				hasErr = true
				break outer
			}
			fmt.Printf("STRING %s %s\n", string(fileContents[i:j+1]), string(fileContents[i+1:j]))
			i = j
		case '=':
			if i < len(fileContents)-1 && fileContents[i+1] == '=' {
				fmt.Printf("EQUAL_EQUAL == null\n")
				i++
			} else {
				fmt.Printf("EQUAL = null\n")
			}
		case '!':
			if i < len(fileContents)-1 && fileContents[i+1] == '=' {
				fmt.Printf("BANG_EQUAL != null\n")
				i++
			} else {
				fmt.Printf("BANG ! null\n")
			}
		case '<':
			if i < len(fileContents)-1 && fileContents[i+1] == '=' {
				fmt.Printf("LESS_EQUAL <= null\n")
				i++
			} else {
				fmt.Printf("LESS < null\n")
			}
		case '>':
			if i < len(fileContents)-1 && fileContents[i+1] == '=' {
				fmt.Printf("GREATER_EQUAL >= null\n")
				i++
			} else {
				fmt.Printf("GREATER > null\n")
			}
		default:
			if val, ok := mp[token]; ok {
				fmt.Printf("%s %c null\n", val, token)
			} else {
				fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %c\n", line, token)
				hasErr = true
			}
		}
	}

	fmt.Println("EOF  null")
	if hasErr {
		os.Exit(65)
	}
}
