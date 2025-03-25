package main

import (
	"fmt"
	"os"
)

func main() {
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
		'/': "SLASH",
	}
	line, hasErr := 1, false
	for i := 0; i < len(fileContents); i++ {
		token := fileContents[i]
		if token == '\n' {
			line++
			continue
		}

		if token == ' ' || token == '\t' || token == '\r' {
			continue
		}

		if i < len(fileContents)-1 && token == '/' && fileContents[i+1] == '/' {
			for i < len(fileContents) && fileContents[i] != '\n' {
				i++
			}
			line++
			continue
		}

		if token == '"' {
			j := i + 1
			for ; j < len(fileContents) && fileContents[j] != '"'; j++ {
				if fileContents[j] == '\n' {
					line++
				}
			}
			if j == len(fileContents) {
				fmt.Fprintf(os.Stderr, "[line %d] Error: Unterminated string.\n", line)
				hasErr = true
				break
			}
			fmt.Printf("STRING %s %s\n", string(fileContents[i:j+1]), string(fileContents[i+1:j]))
			i = j
			continue
		}

		if token == '=' {
			if i < len(fileContents)-1 && fileContents[i+1] == '=' {
				fmt.Printf("EQUAL_EQUAL == null\n")
				i++
			} else {
				fmt.Printf("EQUAL = null\n")
			}
			continue
		}

		if token == '!' {
			if i < len(fileContents)-1 && fileContents[i+1] == '=' {
				fmt.Printf("BANG_EQUAL != null\n")
				i++
			} else {
				fmt.Printf("BANG ! null\n")
			}
			continue
		}

		if token == '<' {
			if i < len(fileContents)-1 && fileContents[i+1] == '=' {
				fmt.Printf("LESS_EQUAL <= null\n")
				i++
			} else {
				fmt.Printf("LESS < null\n")
			}
			continue
		}

		if token == '>' {
			if i < len(fileContents)-1 && fileContents[i+1] == '=' {
				fmt.Printf("GREATER_EQUAL >= null\n")
				i++
			} else {
				fmt.Printf("GREATER > null\n")
			}
			continue
		}

		if val, ok := mp[token]; ok {
			fmt.Printf("%s %c null\n", val, token)
		} else {
			fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %c\n", line, token)
			hasErr = true
		}
	}

	fmt.Println("EOF  null")
	if hasErr {
		os.Exit(65)
	}
}
