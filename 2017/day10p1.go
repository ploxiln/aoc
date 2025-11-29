package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const SIZE = 256

func main() {
	var arr [SIZE]int
	for i := range arr {
		arr[i] = i
	}

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf("Error reading stdin: %s\n", err)
		os.Exit(1)
	}
	instr := strings.TrimSuffix(string(input), "\n")
	sequence := strings.Split(instr, ",")

	pos := 0

	for skip, lenstr := range sequence {
		length, err := strconv.Atoi(lenstr)
		if err != nil || length < 0 || length > SIZE {
			fmt.Printf("ERROR length #%d not valid: %q\n", skip, lenstr)
			os.Exit(1)
		}

		for i := 0; i < length/2; i++ {
			a := (pos /* ...... */ + i) % SIZE
			b := (pos + length - 1 - i) % SIZE

			tmp   := arr[a]
			arr[a] = arr[b]
			arr[b] = tmp
		}

		fmt.Printf("pos=%d len=%d arr=%v...\n", pos, length, arr[:28])

		pos = (pos + length + skip) % SIZE
	}

	fmt.Printf("arr=[%d %d %d...] prod(arr[:2])=%d\n", arr[0], arr[1], arr[2], arr[0]*arr[1])
}
