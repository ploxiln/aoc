package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	var (
		n  = 0
		ne = 0
		nw = 0

		step = 0
		dist = 0

		maxDist = 0
		maxStep = 0
	)
	// s = -n
	// se = -nw = ne - n
	// sw = -ne = nw - n
	// ...

	reader := bufio.NewReader(os.Stdin)
	for {
		// read up to next ',' (or EOF)
		bs, err := reader.ReadString(byte(','))
		if err != nil && err != io.EOF {
			fmt.Printf("ERROR reading stdin: %s\n", err)
			os.Exit(1)
		}
		bs = strings.TrimRight(bs, ",\n")
		step++

		// To find true distance while stepping, we need to
		// reduce total steps by moving some value towards
		// zero whenever possible.
		switch bs {
		case "n":
			if n >= 0 && (ne < 0 || nw < 0) {
				ne++
				nw++
			} else {
				n++
			}
		case "s":
			if n <= 0 && (ne > 0 || nw > 0) {
				ne--
				nw--
			} else {
				n--
			}
		case "ne":
			if ne >= 0 && (n < 0 || nw > 0) {
				n++
				nw--
			} else {
				ne++
			}
		case "sw":
			if ne <= 0 && (n > 0 || nw < 0) {
				n--
				nw++
			} else {
				ne--
			}
		case "nw":
			if nw >= 0 && (n < 0 || ne > 0) {
				n++
				ne--
			} else {
				nw++
			}
		case "se":
			if nw <= 0 && (n > 0 || ne < 0) {
				n--
				ne++
			} else {
				nw--
			}
		default:
			fmt.Printf("ERROR invalid token step=%d : %q\n", step, bs)
			os.Exit(2)
		}
		dist = max(n, -n) + max(ne, -ne) + max(nw, -nw)
		if dist > maxDist {
			maxDist = dist
			maxStep = step
		}
		if step < 50 { // debug print first 50 steps
			fmt.Printf("step=%d op=%s dist=%d n=%d ne=%d nw=%d\n", step, bs, dist, n, ne, nw)
		}

		if err == io.EOF {
			break
		}
	}
	fmt.Printf("dist=%d step=%d n=%d ne=%d nw=%d\n", dist, step, n, ne, nw)
	fmt.Printf("max-dist=%d max-step=%d\n", maxDist, maxStep)
}
