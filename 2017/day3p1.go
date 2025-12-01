package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

func abort(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(2)
}

func getPos(idx int) (x, y int) {
	if idx < 1 {
		panic(fmt.Sprintf("Invalid index: %d", idx))
	}
	if idx == 1 {
		return 0, 0
	}
	// length of one side of the square "ring" this number lies on
	ringSide := int(math.Ceil(math.Sqrt(float64(idx)))) | 1
	// first (lowest) number in this "ring" (... shell)
	ringStart := (ringSide-2)*(ringSide-2) + 1

	// first number in this "ring" is one above bottom-right corner, so to
	// make the math below simpler pretend that the bottom-right corner
	// "sideStart1" has the number just before that (though in reality the
	// bottom right corner is "sideStart5", the largest number in this ring)
	sideStart1 := ringStart - 1
	sideStart2 := (ringStart - 1) + 1*(ringSide-1)
	sideStart3 := (ringStart - 1) + 2*(ringSide-1)
	sideStart4 := (ringStart - 1) + 3*(ringSide-1)
	sideStart5 := (ringStart - 1) + 4*(ringSide-1) // upper bound check

	if idx <= sideStart2 {
		x = ringSide/2
		y = -ringSide/2 + (idx - sideStart1)
		return
	}
	if idx <= sideStart3 {
		x = ringSide/2 - (idx - sideStart2)
		y = ringSide/2
		return
	}
	if idx <= sideStart4 {
		x = -ringSide/2
		y = ringSide/2 - (idx - sideStart3)
		return
	}
	if idx <= sideStart5 {
		x = -ringSide/2 + (idx - sideStart4)
		y = -ringSide/2
		return
	}
	// impossible!
	return -1, -1
}

func main() {
	if len(os.Args) != 2 {
		abort("Usage: %s <idx>", os.Args[0])
	}

	idx, err := strconv.Atoi(os.Args[1])
	if err != nil {
		abort("ERROR invalid value: %q", os.Args[1])
	}

	x, y := getPos(idx)

	// "manhattan" distance
	dist := max(x, -x) + max(y, -y)

	fmt.Printf("idx=%d @ (%d, %d) = %d\n", idx, x, y, dist)
}
