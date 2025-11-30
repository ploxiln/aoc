package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func abort(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, "ERROR "+msg+"\n", args...)
	os.Exit(1)
}

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse int %q : %s", s, err))
	}
	return i
}

func main() {
	layers := []int{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		depStr, rngStr, ok := strings.Cut(line, ": ")
		if !ok {
			abort("parsing line %q", line)
		}
		dep := toInt(depStr)
		rng := toInt(rngStr)
		if len(layers) > dep {
			abort("out of order: layers=%d new-depth=%d", len(layers), dep)
		}
		for len(layers) < dep {
			layers = append(layers, 0)
		}
		layers = append(layers, rng)
	}
	if err := scanner.Err(); err != nil {
		abort("scanning stdin: %s\n", err)
	}

	severity := 0

	// p is picosecond, and layer/depth that packet enters
	for p, rng := range layers {
		if rng == 0 {
			// no scanner
			continue
		}
		if rng == 1 {
			// scanner always at top
			severity += p * rng
			fmt.Printf("HIT depth=%d range=%d\n", p, rng)
			continue
		}
		// "scanner" moves (rng-1) down, then (rng-1) up, then repeats
		pos := p % ((rng-1) * 2)
		if pos == 0 {
			// scanner hits packet
			severity += p * rng
			fmt.Printf("HIT depth=%d range=%d\n", p, rng)
		}
	}
	fmt.Printf("total severity: %d\n", severity)
}
