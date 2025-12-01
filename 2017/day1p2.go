package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func abort(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(2)
}

func isdigit(ch byte) bool {
	return ch >= byte('0') && ch <= byte('9')
}

func main() {
	total := int(0)

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		abort("ERROR reading stdin: %s", err)
	}

	input = bytes.TrimSpace(input)

	for i, ch := range input {
		if !isdigit(ch) {
			abort("ERROR not digit at position %d: %c", i, ch)
		}
		ri := (i + (len(input) / 2)) % len(input)
		if ch == input[ri] {
			total += int(ch - byte('0'))
		}
	}
	fmt.Printf("len=%d total=%d\n", len(input), total)
}
