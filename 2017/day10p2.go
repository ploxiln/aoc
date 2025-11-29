package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

const SIZE = 256

func main() {
	// less casting between int / byte if we
	// switch everything to byte (more efficient too)
	var arr [SIZE]byte
	for i := range arr {
		arr[i] = byte(i)
	}

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf("Error reading stdin: %s\n", err)
		os.Exit(1)
	}
	input = bytes.TrimSpace(input)
	input = append(input, []byte{17, 31, 73, 47, 23}...)

	pos := byte(0)
	skip := byte(0)

	for r := 0; r < 64; r++ {
		for _, length := range input {
			for i := byte(0); i < length/2; i++ {
				// byte or uint8 automatically wraps 255 -> 0
				a := (pos /* ...... */ + i) // % SIZE
				b := (pos + length - 1 - i) // % SIZE

				tmp   := arr[a]
				arr[a] = arr[b]
				arr[b] = tmp
			}

			if r == 0 && skip < 40 { // debug print first 40
				fmt.Printf("pos=%d len=%d arr=%v...\n", pos, length, arr[:32])
			}
			pos = (pos + length + skip) // % SIZE
			skip++
		}
		fmt.Printf("done round %d | pos=%d skip=%d\n", r, pos, skip)
	}

	fmt.Printf("Result: ")
	var dense [SIZE / 16]byte
	for i := range dense {
		for j := 0; j < 16; j++ {
			dense[i] ^= arr[i * 16 + j]
		}
		// output "dense" hash in hex
		fmt.Printf("%02x", dense[i])
	}
	fmt.Printf("\n")
}
