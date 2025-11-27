package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("ERROR parsing %q - %s", s, err))
	}
	return i
}

func main() {
	csum := 0

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 0 {
			continue
		}
		vals := make([]int, len(fields))
		for i, cell := range fields {
			vals[i] = toInt(cell)
		}
		i := 0
		j := 0
	SEARCH:
		for i = 0; i < len(vals); i++ {
			for j = 0; j < len(vals); j++ {
				if vals[i] > vals[j] && vals[i] % vals[j] == 0 {
					break SEARCH
				}
			}
		}
		if i == len(vals) || j == len(vals) {
			fmt.Fprintf(os.Stderr, "ERROR: no divisible pair in %v\n", vals)
		}
		q := vals[i] / vals[j]
		csum += q
		fmt.Printf("line len=%d num=%d div=%d q=%d\n", len(fields), vals[i], vals[j], q)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %s", err)
		os.Exit(1)
	}
	fmt.Printf("csum=%d\n", csum)
}
