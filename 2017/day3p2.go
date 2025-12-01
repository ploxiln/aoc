package main

import (
	"fmt"
	"os"
	"strconv"
)

const dim = 65

func abort(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(2)
}

func main() {
	if len(os.Args) != 2 {
		abort("Usage: %s <val>", os.Args[0])
	}

	val, err := strconv.Atoi(os.Args[1])
	if err != nil {
		abort("ERROR invalid value: %q", os.Args[1])
	}

	// calculate and store grid cells sequentially
	// (values initialize to zero by default)
	grid := [dim][dim]int{}

	side := 1
	x := dim/2
	y := dim/2
	grid[y][x] = 1

	calcCell := func() {
		grid[y][x] = (grid[y  ][x+1] +
		              grid[y+1][x+1] +
		              grid[y+1][x  ] +
		              grid[y+1][x-1] +
		              grid[y  ][x-1] +
		              grid[y-1][x-1] +
		              grid[y-1][x  ] +
		              grid[y-1][x+1] )

		fmt.Printf("(%3d,%3d) = %6d\n", x-dim/2, y-dim/2, grid[y][x])

		if grid[y][x] > val {
			fmt.Printf("Found greater value ^^^\n")
			os.Exit(0)
		}
	}

	// walk around each square "ring" / "shell"
	for x < dim-2 && y < dim-2 {
		side += 2
		x    += 1
		calcCell()
		for i := 2; i < side; i++ {
			y += 1
			calcCell()
		}
		for i := 1; i < side; i++ {
			x -= 1
			calcCell()
		}
		for i := 1; i < side; i++ {
			y -= 1
			calcCell()
		}
		for i := 1; i < side; i++ {
			x += 1
			calcCell()
		}
	}
}
