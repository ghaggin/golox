package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: golox [script]")
		return
	} else if len(os.Args) == 2 {
		if err := runFile(os.Args[0]); err != nil {
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

	parser := NewParser(tokens)
	expr := parser.Parse()
	if hadError {
		return fmt.Errorf("failed to parse")
	}

	fmt.Println(expr.Print())
	return nil
}
