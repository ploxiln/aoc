// Usage: day15p2 <init-A> <init-B>
package main

import (
	"fmt"
	"os"
	"strconv"
)

const (
	FACT_A = 16807
	FACT_B = 48271
	MODDIV = 2147483647 // 2**31 - 1 ?
	ROUNDS = 5_000_000
)

func abort(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

// string of 32 bits, but take uint64 for convenience
func uintBitStr(val uint64) string {
	var output [32]byte
	for i := 0; i < 32; i++ {
		output[i] = byte('0') + byte((val >> (31 - i)) & 1)
	}
	return string(output[:])
}

func main() {
	if len(os.Args) != 3 {
		abort("USAGE: day15p1 <init-A> <init-B>")
	}
	initA, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		abort("ERROR parsing uint64 <init-A> %q: %s", os.Args[1], err)
	}
	initB, err := strconv.ParseUint(os.Args[2], 10, 64)
	if err != nil {
		abort("ERROR parsing uint64 <init-A> %q: %s", os.Args[2], err)
	}
	fmt.Printf("init: A=%d B=%d\n", initA, initB)

	// a and b values fit in 32-bit, but intermediate calculations
	// need more bits, so just use 64-bit everywhere
	a := initA
	b := initB

	matches := 0

	for i := 0; i < ROUNDS; i++ {
		// find next A which is multiple of 4
		for {
			a = (a * FACT_A) % MODDIV
			if (a & 0x3) == 0 {
				break
			}
		}
		// find next B which is multiple of 8
		for {
			b = (b * FACT_B) % MODDIV
			if (b & 0x7) == 0 {
				break
			}
		}
		diff := (a ^ b) & 0xffff

		if i < 5 { // debug print
			mstr := ""
			if diff == 0 {
				mstr = " *"
			}
			fmt.Printf("A: %10d | %s %s\n", a, uintBitStr(a), mstr)
			fmt.Printf("B: %10d | %s %s\n", b, uintBitStr(b), mstr)
		}
		if diff == 0 {
			matches++
		}
	}
	fmt.Printf("Total matches in %d generations: %d\n", ROUNDS, matches)
}
