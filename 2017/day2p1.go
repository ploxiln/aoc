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
		lmin := toInt(fields[0])
		lmax := lmin
		for _, cell := range fields[1:] {
			val := toInt(cell)
			lmin = min(lmin, val)
			lmax = max(lmax, val)
		}
		ldiff := lmax - lmin
		csum += ldiff
		fmt.Printf("line len=%d min=%d max=%d diff=%d\n", len(fields), lmin, lmax, ldiff)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %s", err)
		os.Exit(1)
	}
	fmt.Printf("csum=%d\n", csum)
}
