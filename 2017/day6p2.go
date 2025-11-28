package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	BN = 16
)

func abort(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func main() {
	blob, err := io.ReadAll(os.Stdin)
	if err != nil {
		abort("ERROR reading stdin: %s", err)
	}
	fields := strings.Fields(string(blob))

	if len(fields) != BN {
		abort("ERROR expected %d banks, got %d", BN, len(fields))
	}

	banks := [BN]int{}

	for i, s := range fields {
		v, err := strconv.Atoi(s)
		if err != nil {
			abort("ERROR parsing field %q: %s", s, err)
		}
		banks[i] = v
	}

	seen := map[[BN]int]int{}

	for step := 0; step < 999000; step++ {
		if _, ok := seen[banks]; ok {
			fmt.Printf("cycles=%d prev=%d dist=%d : %v\n",
				step, seen[banks], step-seen[banks], banks)
			return
		}
		seen[banks] = step

		// run redistribution
		maxidx := 0
		for i, v := range banks {
			if v > banks[maxidx] {
				maxidx = i
			}
		}
		redist := banks[maxidx]
		banks[maxidx] = 0

		for i := maxidx; redist > 0; redist-- {
			i = (i + 1) % BN
			banks[i] += 1
		}
	}
	fmt.Printf("Repeat not found\n")
}
