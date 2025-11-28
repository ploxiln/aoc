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

func mInt(s string) int {
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

	maxVal := 0
	maxReg := "?"
	maxStp := -1

	counter := 0
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
		regVal, ok := registers[treg]

		if condition(registers[creg], ctst, mInt(cval)) {
			if top == "inc" {
				regVal += mInt(tval)
			} else { // "dec"
				regVal -= mInt(tval)
			}
			// check if max value seen in any reg
			if regVal > maxVal {
				maxVal = regVal
				maxReg = treg
				maxStp = counter
			}
			registers[treg] = regVal
		} else if !ok {
			// this target reg mentioned but not modified yet,
			// remember zero value
			registers[treg] = regVal
		}

		counter++
	}
	if err := scanner.Err(); err != nil {
		abort("stdin scanner read failed: %s", err)
	}

	// calculate "final" max again also, why not
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
	fmt.Printf("INNER-MAX: %s : %d (step=%d)\n", maxReg, maxVal, maxStp)
}
