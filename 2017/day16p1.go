package main

import (
	"bufio"
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

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse int %q: %s", s, err))
	}
	return i
}

func main() {
	progs := []byte("abcdefghijklmnop")

	// instead of just `([sxp])([0-9]+)...` this checks if args are appropriate for move type
	moveRe, err := regexp.Compile(`(s)([0-9]+)|(x)([0-9]+)/([0-9]+)|(p)([a-p])/([a-p]),?`)
	if err != nil {
		abort("ERROR compiling regexp: %s", err)
	}
	// s=groups[1], x=groups[3], p=groups[6]

	reader := bufio.NewReader(os.Stdin)
	count := 1
	for {
		move, err := reader.ReadString(byte(','))
		if err != nil && err != io.EOF {
			abort("ERROR reading stdin: %s", err)
		}
		m := moveRe.FindStringSubmatch(move)
		if len(m) == 0 {
			abort("ERROR parsing move %d: %q", count, move)
		}

		if m[1] == "s" {
			p := toInt(m[2])
			progs = append(progs[16-p:], progs[:16-p]...)

		} else if m[3] == "x" {
			i := toInt(m[4])
			j := toInt(m[5])
			tmp     := progs[i]
			progs[i] = progs[j]
			progs[j] = tmp

		} else if m[6] == "p" {
			u := byte(m[7][0])
			v := byte(m[8][0])
			i := -1 // find idx of u
			j := -1 // find idx of v
			for x, b := range progs {
				if b == u {
					i = x
				}
				if b == v {
					j = x
				}
			}
			tmp     := progs[i]
			progs[i] = progs[j]
			progs[j] = tmp

		} else {
			abort("ERROR invalid move %d: %q", count, move)
		}
		if count < 9 { // debug print
			fmt.Printf("%d: %-7s | %q\n", count, move, string(progs))
		}
		if err == io.EOF {
			break
		}
		count++
	}
	fmt.Printf("FINAL moves=%d | %q\n", count, string(progs))
}
