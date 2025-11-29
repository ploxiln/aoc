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
		n    = 0
		ne   = 0
		step = 0
	)
	// s = -n
	// sw = -ne
	// nw = n - ne
	// se = ne - n
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

		switch bs {
		case "n":
			n++
		case "s":
			n--
		case "ne":
			ne++
		case "sw":
			ne--
		case "nw":
			n++
			ne--
		case "se":
			n--
			ne++
		default:
			fmt.Printf("ERROR invalid token step=%d : %q\n", step, bs)
			os.Exit(2)
		}
		if err == io.EOF {
			break
		}
	}
	fmt.Printf("raw n=%d ne=%d ...\n", n, ne)

	// Convert negative n,ne to positive s,se,sw,nw.
	// Each group of ops is net zero movement, and moves
	// n,ne counts towards zero, reducing total steps.
	var (
		s  = 0
		se = 0
		sw = 0
		nw = 0
	)
	for n < 0 {
		if ne > 0 {
			n++
			ne--
			se++
		} else {
			n++
			s++
		}
	}
	for ne < 0 {
		if n > 0 {
			ne++
			n--
			nw++
		} else {
			ne++
			sw++
		}
	}
	dist := n + ne + nw + s + se + sw
	fmt.Printf("n=%d s=%d ne=%d sw=%d nw=%d se=%d dist=%d\n", n, s, ne, sw, nw, se, dist)
}
