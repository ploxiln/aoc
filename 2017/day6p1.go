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

	// can use a fixed-length array of ints as a map key
	// (otherwise would use string, comma-joined numbers, or just hex)
	seen := map[[BN]int]bool{}

	for step := 0; step < 999000; step++ {
		if seen[banks] {
			fmt.Printf("SEEN BEFORE cycles=%d : %v\n", step, banks)
			return
		}
		seen[banks] = true

		if step < 50 {
			// just to see it's working ...
			fmt.Printf("step=%3d  banks: %v\n", step, banks)
			if step == 49 {
				fmt.Printf("...\n")
			}
		}

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
