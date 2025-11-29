package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	score := 0
	depth := 0
	gar := false // inside garbage section
	esc := false // cancel next garbage char ("escape")

	counter := 0
	reader := bufio.NewReader(os.Stdin)
	for {
		ch, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(fmt.Sprintf("ERROR reading stdin: %s", err))
		}
		counter++

		switch {
		case esc:
			esc = false
		case gar && ch == '!':
			esc = true
		case gar && ch == '>':
			gar = false
		case gar:
			gar = true // no-op
		case ch == '<':
			gar = true
		case ch == '}':
			depth -= 1
		case ch == '{':
			depth += 1
			score += depth
			if score < 100 { // debug-print first 20 or so
				fmt.Printf("Group Start pos=%d depth=%d\n", counter, depth)
			}
		}
	}

	fmt.Printf("Total len=%d score=%d\n", counter, score)

	if depth != 0 {
		fmt.Printf("ERROR final depth=%d\n", depth)
		os.Exit(1)
	}
}
