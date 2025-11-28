package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func abort(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, "ERROR: "+msg+"\n", args...)
	os.Exit(1)
}

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("failed to parse int %q: %s", s, err))
	}
	return i
}

func condition(l int, tst string, r int) bool {
	switch tst {
	case "<":
		return l <  r
	case ">":
		return l >  r
	case "<=":
		return l <= r
	case ">=":
		return l >= r
	case "==":
		return l == r
	case "!=":
		return l != r
	default:
		panic(fmt.Sprintf("invalid condition %q", tst))
	}
}

func main() {
	lineRe, err := regexp.Compile(`^([a-z]+) (inc|dec) (-?[0-9]+) if ([a-z]+) (<|<=|>|>=|==|!=) (-?[0-9]+)`)
	if err != nil {
		abort("invalid line regex: %s", err)
	}

	// all registers start at 0 ; map missing default is 0
	registers := map[string]int{}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		fields := lineRe.FindStringSubmatch(line)
		if len(fields) == 0 {
			abort("failed to parse line: %q", line)
		}
		var (
			treg = fields[1] // target
			top  = fields[2]
			tval = fields[3]
			creg = fields[4] // condition
			ctst = fields[5]
			cval = fields[6]
		)
		if condition(registers[creg], ctst, toInt(cval)) {
			if top == "inc" {
				registers[treg] += toInt(tval)
			} else { // "dec"
				registers[treg] -= toInt(tval)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		abort("stdin scanner read failed: %s", err)
	}

	fMaxReg := "?"
	fMaxVal := -987654321

	for reg, val := range registers {
		if val > fMaxVal {
			fMaxReg = reg
			fMaxVal = val
		}
		fmt.Printf("%s : %d\n", reg, val)
	}
	fmt.Printf("FINAL-MAX: %s : %d\n", fMaxReg, fMaxVal)
}
