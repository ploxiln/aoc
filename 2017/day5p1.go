package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func abort(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(2)
}

func main() {
	var tape []int

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		val, err := strconv.Atoi(line)
		if err != nil {
			abort("Error parsing line %q: %s", line, err)
		}
		tape = append(tape, val)
	}
	if err := scanner.Err(); err != nil {
		abort("Error reading stdin: %s", err)
	}

	// ip : instruction pointer (index into tape)
	ip := 0
	step := 0
	for ; ip >= 0 && ip < len(tape); step++ {
		jump := tape[ip]
		if step < 50 {
			// just the first 50, to see it's working
			fmt.Printf("step=%4d  ip=%4d  jump=%4d\n", step, ip, jump)
		}
		tape[ip] = jump + 1
		ip = ip + jump
	}
	fmt.Printf("...\n step=%4d  ip=%4d  EXIT\n", step, ip)
}
