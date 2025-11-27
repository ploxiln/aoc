package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

func validPassphrase(phrase []string) bool {
	// sort letters within each "word"
	sorted := make([][]byte, len(phrase))
	for i, w := range phrase {
		bw := []byte(w)
		slices.Sort(bw)
		sorted[i] = bw
	}
	// then check if any words match
	for i := 0; i < len(phrase); i++ {
		for j := i + 1; j < len(phrase); j++ {
			if slices.Compare(sorted[i], sorted[j]) == 0 {
				return false
			}
		}
	}
	return true
}

func main() {
	valid := 0
	inval := 0

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 0 {
			continue
		}
		if validPassphrase(fields) {
			valid++
		} else {
			inval++
		}

	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %s", err)
		os.Exit(1)
	}

	fmt.Printf("valid=%d invalid=%d\n", valid, inval)
}
