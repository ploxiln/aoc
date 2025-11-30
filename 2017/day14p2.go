// Usage: day14p2 <prefix>
package main

import (
	"fmt"
	"os"
)

var hexLookup = [16]rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}

func toHex(data []byte) string {
	out := make([]rune, len(data)*2)
	for i, d := range data {
		out[i*2  ] = hexLookup[(d>>4) & 0xf]
		out[i*2+1] = hexLookup[ d     & 0xf]
	}
	return string(out)
}

func toBinary(data []byte) string {
	out := make([]byte, len(data)*8)
	for i, d := range data {
		for j := 0; j < 8; j++ {
			out[i*8+j] = byte('0') + (d >> (7 - j) & 1)
		}
	}
	return string(out)
}

type Point struct {
	X int
	Y int
}

func bitCoords(x, y int) (row int, col int, bit byte) {
	row = y
	col = x / 8
	bit = byte(1) << (7 - (x % 8))
	return
}

func isBitSet(grid *[128][16]byte, p Point) bool {
	row, col, bit := bitCoords(p.X, p.Y)
	return (grid[row][col] & bit) != 0
}
func setBit(grid *[128][16]byte, p Point, v bool) {
	row, col, bit := bitCoords(p.X, p.Y)
	if v {
		grid[row][col] |= bit
	} else {
		grid[row][col] &= ^bit
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "USAGE: day14p1 <prefix>\n")
		os.Exit(2)
	}
	prefix := os.Args[1]

	var grid [128][16]byte

	for i := 0; i < 128; i++ {
		data := fmt.Sprintf("%s-%d", prefix, i)
		hash := knotHash([]byte(data))
		grid[i] = hash
		fmt.Printf("%s | %s\n", toHex(hash[:]), toBinary(hash[:]))
	}

	groups := 0
	var seen [128][16]byte // already "visited" bits

	for y := 0; y < 128; y++ {
		for x := 0; x < 128; x++ {
			// find not-yet-visited bit to start new group
			p := Point{x, y}
			if !isBitSet(&grid, p) {
				continue
			}
			if isBitSet(&seen, p) {
				continue
			}
			fmt.Printf("found group start=(%d,%d)", x, y)
			groups++
			groupSize := 0

			// poor-man's queue of bit coords to visit
			queue := []Point{p}
			queueAppend := func(u, v int) {
				queue = append(queue, Point{u, v})
			}

			for len(queue) > 0 {
				np := queue[0]
				queue = queue[1:]

				if !isBitSet(&grid, np) {
					continue // no bit set here
				}
				if isBitSet(&seen, np) {
					continue // already visited
				}
				setBit(&seen, np, true)
				groupSize++

				// queue neighboring bit coords to visit,
				// but avoid going over the edge
				if np.X > 0 {
					queueAppend(np.X-1, np.Y)
				}
				if np.X < 127 {
					queueAppend(np.X+1, np.Y)
				}
				if np.Y > 0 {
					queueAppend(np.X, np.Y-1)
				}
				if np.Y < 127 {
					queueAppend(np.X, np.Y+1)
				}
			}
			fmt.Printf(" size=%d\n", groupSize)
		}
	}
	fmt.Printf("total-groups=%d\n", groups)
}

func knotHash(data []byte) [16]byte {
	var arr [256]byte
	for i := range arr {
		arr[i] = byte(i)
	}
	data = append(data, []byte{17, 31, 73, 47, 23}...)
	pos := byte(0)
	skip := byte(0)

	for r := 0; r < 64; r++ {
		for _, length := range data {
			for i := byte(0); i < length/2; i++ {
				a := (pos /* ...... */ + i) // % SIZE
				b := (pos + length - 1 - i) // % SIZE

				tmp   := arr[a]
				arr[a] = arr[b]
				arr[b] = tmp
			}
			pos = (pos + length + skip) // % SIZE
			skip++
		}
	}
	var dense [16]byte
	for i := range dense {
		for j := 0; j < 16; j++ {
			dense[i] ^= arr[i * 16 + j]
		}
	}
	return dense
}
