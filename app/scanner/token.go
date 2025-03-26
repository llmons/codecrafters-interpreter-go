package scanner

import (
	"fmt"
	"math"
)

type token struct {
	tokenType tokenType
	lexeme    string
	literal   any
	line      int
}

func (t token) String() string {
	if t.literal == nil {
		return fmt.Sprintf("%s %s null", t.tokenType, t.lexeme)
	}
	if t.tokenType == NUMBER {
		if math.Floor(t.literal.(float64)) == t.literal.(float64) {
			return fmt.Sprintf("%s %s %.1f", t.tokenType, t.lexeme, t.literal.(float64))
		}
		return fmt.Sprintf("%s %s %g", t.tokenType, t.lexeme, t.literal.(float64))
	}
	return fmt.Sprintf("%s %s %v", t.tokenType, t.lexeme, t.literal)
}
