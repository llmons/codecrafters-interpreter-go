package main

import (
	"fmt"
	"math"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/app/scanner"
)

var hadError = false

func runFile(path string) {
	fileContents, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	run(string(fileContents))
	if hadError {
		os.Exit(65)
	}
}

func runPrompt() {
	for {
		fmt.Print("> ")
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}
		run(input)
		hadError = false
	}
}

func run(source string) {
	scanner := scanner.NewScanner(source)
	tokens := scanner.ScanTokens()
	for _, token := range tokens {
		fmt.Println(token)
	}
}

func error(line int, message string) {
	report(line, "", message)
}

func report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error %s: %s\n", line, where, message)
	hadError = true
}

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
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			j := i + 1
			var val float64
			val = float64(token - '0')
			for ; j < len(fileContents) && (fileContents[j] >= '0' && fileContents[j] <= '9'); j++ {
				val = val*10 + float64(fileContents[j]-'0')
			}
			if j < len(fileContents) && fileContents[j] == '.' {
				j++
				for k := 10; j < len(fileContents) && (fileContents[j] >= '0' && fileContents[j] <= '9'); j++ {
					val += float64(fileContents[j]-'0') * 1.0 / float64(k)
					k *= 10
				}
			}
			if math.Floor(val) == val {
				fmt.Printf("NUMBER %s %.1f\n", string(fileContents[i:j]), val)
			} else {
				fmt.Printf("NUMBER %s %.4f\n", string(fileContents[i:j]), val)
			}
			// now j points to the character after the number
			i = j - 1
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
			// now j points to the closing "
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
		case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
			'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
			'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
			'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '_':
			j := i + 1
			for ; j < len(fileContents) && (fileContents[j] >= 'a' && fileContents[j] <= 'z') || (fileContents[j] >= 'A' && fileContents[j] <= 'Z') || fileContents[j] == '_'; j++ {
			}
			fmt.Printf("IDENTIFIER %s null\n", string(fileContents[i:j]))
			// now j points to the character after the identifier
			i = j - 1
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
