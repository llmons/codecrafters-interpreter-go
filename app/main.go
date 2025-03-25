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
	line, unkownChar := 1, false
	for i := 0; i < len(fileContents); i++ {
		token := fileContents[i]
		if token == '\n' {
			line++
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
			unkownChar = true
		}
	}
	fmt.Println("EOF  null")
	if unkownChar {
		os.Exit(65)
	}

}
