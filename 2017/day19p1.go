package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

func abort(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func main() {
	// 0,0  1,0   2,0 ...
	// 0,1  1,1   2,1 ...
	// 0,2  1,2   2,2 ...
	var grid [][]byte

	// copy whole input to "two-dimensional slice"
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		grid = append(grid, slices.Clone(line))
		if len(line) != len(grid[0]) {
			abort("ERROR len(line)=%d != grid-width=%d", len(line), len(grid[0]))
		}
	}
	if err := scanner.Err(); err != nil {
		abort("ERROR reading stdin: %s", err)
	}
	fmt.Printf("grid size = %dx%d\n", len(grid[0]), len(grid))

	// initial direction is down
	dir_x := 0
	dir_y := 1

	// find entrance on top row
	x := 0
	y := 0
	for i, ch := range grid[0] {
		if ch == byte('|') {
			// found entrance
			x = i
			break
		}
	}
	fmt.Printf("start @ (%d, %d)\n", x, y)

	letters := []byte{}
FOLLOW_PATH:
	for {
		x += dir_x
		y += dir_y

		if x < 0 || x > len(grid[0]) || y < 0 || y > len(grid) {
			abort("ERROR out of grid: (%d, %d) (dir=(%d,%d))", x, y, dir_x, dir_y)
		}
		ch := grid[y][x]

		switch {
		case ch == byte('|'): // continue straight (or cross other path)
		case ch == byte('-'): // continue straight (or cross other path)
		case ch == byte(' '): // just passed the end of the path
			x -= dir_x
			y -= dir_y
			break FOLLOW_PATH

		case ch == byte('+') && dir_x != 0: // turn, from left/right to up/down
			if        y   > 0         && grid[y-1][x] != byte(' ') { dir_y = -1
			} else if y+1 < len(grid) && grid[y+1][x] != byte(' ') { dir_y =  1
			} else {
				abort("ERROR bad turn @ (%d, %d) (dir=(%d,%d))", x, y, dir_x, dir_y)
			}
			dir_x = 0
			fmt.Printf("Turn @ (%d, %d)\n", x, y)

		case ch == byte('+') && dir_y != 0: // turn, from up/down to left/right
			if        x   > 0            && grid[y][x-1] != byte(' ') { dir_x = -1
			} else if x+1 < len(grid[0]) && grid[y][x+1] != byte(' ') { dir_x =  1
			} else {
				abort("ERROR bad turn @ (%d, %d) (dir=(%d,%d))", x, y, dir_x, dir_y)
			}
			dir_y = 0
			fmt.Printf("Turn @ (%d, %d)\n", x, y)

		case ch >= byte('A') && ch <= byte('Z'):
			letters = append(letters, ch)
			fmt.Printf("Letter %c @ (%d, %d)\n", ch, x, y)

		default:
			abort("ERROR unexpected char %c @ (%d, %d)", ch, x, y)
		}
	}
	fmt.Printf("Path ended @ (%d, %d)\n", x, y)
	fmt.Printf("Letters: %s\n", string(letters))
}
