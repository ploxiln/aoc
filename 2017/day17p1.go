package main

import (
	"fmt"
	"os"
	"strconv"
)

const INSMAX = 2017

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("USAGE: day17p1 <step>\n")
		os.Exit(2)
	}

	step, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("ERROR parsing integer step %q: %s\n", os.Args[1], err)
		os.Exit(2)
	}

	// Total length is initial 1 + INSMAX additional insertions
	ring := make([]int, INSMAX+1)
	pos := 0
	ring[pos] = 0

	for i := 1; i <= INSMAX; i++ {
		// start at previous inserted item (pos), step forward (step),
		// wrap around at current array size (i),
		// then insert value i *after* that (+1)
		pos = ((pos + step) % i) + 1
		// shift all items forward at or after (pos) to make room at (pos)
		for j := i; j > pos; j-- {
			ring[j] = ring[j-1]
		}
		ring[pos] = i

		if i < 10 { // debug print
			fmt.Printf("i=%d ring=%v\n", i, ring[:i+1])
		}
	}

	fmt.Printf("pos=%d cur=%d next=%d\n", pos, ring[pos], ring[(pos+1) % (INSMAX+1)])
}
