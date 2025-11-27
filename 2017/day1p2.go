package main

import (
	"bytes"
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
	total := int(0)

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		abort("Failed to read stdin: %w\n", err)
	}

	input = bytes.TrimSpace(input)

	for i, ch := range input {
		if !isdigit(ch) {
			abort("input at %d is not a char: %q\n", i, ch)
		}
		ri := (i + (len(input) / 2)) % len(input)
		if ch == input[ri] {
			total += int(ch - byte('0'))
		}
	}
	fmt.Printf("len=%d total=%d\n", len(input), total)
}
