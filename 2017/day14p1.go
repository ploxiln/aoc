// Usage: day14p1 <prefix>
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "USAGE: day14p1 <prefix>\n")
		os.Exit(2)
	}
	prefix := os.Args[1]

	totalbits := 0
	var grid [128][16]byte

	for i := 0; i < 128; i++ {
		data := fmt.Sprintf("%s-%d", prefix, i)
		hash := knotHash([]byte(data))
		grid[i] = hash
		bits := 0
		for _, v := range hash {
			for j := 7; j >= 0; j-- {
				bits += int((v >> j) & 1)
			}
		}
		totalbits += bits
		fmt.Printf("%-15s bits=%3d | %s\n", data, bits, toHex(hash[:]))
	}
	fmt.Printf("totalbits=%d\n", totalbits)
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
			dense[i] ^= arr[i*16 + j]
		}
	}
	return dense
}
