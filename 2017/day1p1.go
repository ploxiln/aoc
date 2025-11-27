package main

import (
	"fmt"
	"io"
	"os"
)

func abort(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg, args...)
	os.Exit(2)
}

func isdigit(ch byte) bool {
	return ch >= byte('0') && ch <= byte('9')
}

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		abort("Failed to read stdin: %w", err)
	}

	total := int(0)
	var last byte

	// find last digit, ignoring newline/space/etc
	for i := len(input) - 1; i >= 0; i-- {
		if isdigit(input[i]) {
			last = input[i]
			input = input[:i+1]
			break
		}
	}
	for i, ch := range input {
		if !isdigit(ch) {
			abort("ERROR not digit at position %d: %c\n", i, ch)
		}
		if ch == last {
			total += int(ch - byte('0'))
		}
		last = ch
	}
	fmt.Printf("len=%d total=%d\n", len(input), total)
}
