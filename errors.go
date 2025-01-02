package main

import (
	"fmt"
	"os"
)

var hadError bool

type ParseError struct{}

func (p ParseError) Error() string {
	return "parse error"
}

func Error(line int, message string) {
	Report(line, "", message)
}

func Report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s", line, where, message)
	hadError = true
}

func TokenError(token Token, message string) {
	if token.Type == EOF {
		Report(token.Line, " at end", message)
	} else {
		Report(token.Line, fmt.Sprintf(" at '%s'", token.Lexeme), message)
	}
}
