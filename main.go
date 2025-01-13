package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var (
	environment = NewEnvironment(nil)
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: golox [script]")
		return
	} else if len(os.Args) == 2 {
		if err := runFile(os.Args[1]); err != nil {
			fmt.Println(err)
			return
		}
	} else {
		err := runPrompt()
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func runFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if err := run(string(b)); err != nil {
		return fmt.Errorf("error running file: %w", err)
	}
	return nil
}

func runPrompt() error {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("> ")
		if scanner.Scan() {
			line := scanner.Text()
			err := run(line)
			if err != nil {
				fmt.Println(err.Error())
			}
			hadError = false
		}
	}
}

func run(s string) error {
	tokens, err := NewScanner(s).scanTokens()
	if err != nil {
		return fmt.Errorf("failed to scan tokens: %w", err)
	}

	parser, err := NewParser(tokens)
	if err != nil {
		return fmt.Errorf("failed to run: %w", err)
	}

	stmts, err := parser.Parse()
	if err != nil || hadError {
		if err != nil {
			fmt.Println(err)
		}
		return fmt.Errorf("failed to parse")
	}

	interpret(stmts)

	return nil
}

func interpret(stmts []Stmt) {
	for _, stmt := range stmts {
		err := stmt.Execute()
		if err != nil {
			fmt.Println(err)
		}
	}
}
