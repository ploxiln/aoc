package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

func abort(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func toInt(b []byte) int {
	i, err := strconv.Atoi(string(b))
	if err != nil {
		panic(fmt.Sprintf("failed to parse int %q: %s", string(b), err))
	}
	return i
}

type Move struct {
	Type rune
	X    int
	Y    int
	U    byte
	V    byte
}

func main() {
	progs := []byte("abcdefghijklmnop")

	// instead of just `([sxp])([0-9]+)...` this checks if args are appropriate for move type
	moveRe, err := regexp.Compile(`^(s)([0-9]+)|(x)([0-9]+)/([0-9]+)|(p)([a-p])/([a-p])$`)
	if err != nil {
		abort("ERROR compiling regexp: %s", err)
	}
	// s=groups[1], x=groups[3], p=groups[6]

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		abort("ERROR reading stdin: %s", err)
	}

	var dance []Move

	for i, movestr := range bytes.Split(input, []byte(",")) {
		m := moveRe.FindSubmatch(movestr)
		if len(m) == 0 {
			abort("ERROR parsing move %d: %q", i+1, string(movestr))
		}
		var move Move

		if string(m[1]) == "s" {
			move.Type = 's'
			move.X = toInt(m[2])

		} else if string(m[3]) == "x" {
			move.Type = 'x'
			move.X = toInt(m[4])
			move.Y = toInt(m[5])

		} else if string(m[6]) == "p" {
			move.Type = 'p'
			move.U = m[7][0]
			move.V = m[8][0]

		} else {
			abort("ERROR invalid move %d: %q", i+1, string(movestr))
		}
		dance = append(dance, move)
	}
	fmt.Printf("Parsed dance = %d moves\n", len(dance))
	fmt.Printf("After %d : %s\n", 0, string(progs))

	seen := map[string]int{}
	seen[string(progs)] = 0

	// just 200 rounds for a start, maybe there is a cycle...
	for r := 1; r <= 200; r++ {
		for _, m := range dance {
			switch m.Type {
			case 's':
				progs = append(progs[16 - m.X :], progs[:16 - m.X]...)
			case 'x':
				tmp       := progs[m.X]
				progs[m.X] = progs[m.Y]
				progs[m.Y] = tmp
			case 'p':
				i := -1 // find idx of m.U
				j := -1 // find idx of m.V
				for x, b := range progs {
					if b == m.U {
						i = x
					}
					if b == m.V {
						j = x
					}
				}
				tmp     := progs[i]
				progs[i] = progs[j]
				progs[j] = tmp
			}
		}
		sp := string(progs)

		fmt.Printf("After %d : %s\n", r, sp)
		if i, ok := seen[sp]; ok {
			fmt.Printf("seen round %d == %d\n", i, r)
			fmt.Printf("After 1,000,000,000 rounds, state should match round %d\n",
				(1_000_000_000 - i) % (r - i) + i)
			// Simple case: i=0 (round 0 seen first), final pattern matches round (1 million % r)
			break
		}
		seen[sp] = r
	}
}
