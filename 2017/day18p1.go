package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func abort(msg string, args ...any) {
	fmt.Printf(msg+"\n", args...)
	os.Exit(1)
}

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		abort("ERROR reading stdin: %s", err)
	}

	lines := strings.Split(string(input), "\n")
	program := [][]string{}
	for _, l := range lines {
		if len(l) == 0 {
			continue
		}
		// lazy parse: just split each instruction into tokens
		program = append(program, strings.Fields(l))
	}

	// registers are named by single latin letter, so there are 26
	// for register x, use index (x - 'a')
	registers := [26]int{}
	lastplay := 0
	ip := 0

	// set register's value
	setReg := func(reg string, val int) {
		if reg[0] < 'a' || reg[0] > 'z' {
			abort("ERROR invalid register %q ip=%d", reg, ip)
		}
		registers[reg[0]-'a'] = val
	}

	// get a register's value *or* an immediate value
	getRegImm := func(s string) int {
		if s[0] >= 'a' && s[0] <= 'z' {
			return registers[s[0]-'a']
		} else {
			val, err := strconv.Atoi(s)
			if err != nil {
				panic(fmt.Sprintf("invalid arg %q ip=%d", s, ip))
			}
			return val
		}
	}

	for { // execute program
		inst := program[ip]
		next := ip + 1

		fmt.Printf("%3d : %v\n", ip, inst)
		switch inst[0] {
		case "snd":
			lastplay = getRegImm(inst[1])
			fmt.Printf("PLAY : %d\n", lastplay)
		case "set":
			setReg(inst[1], getRegImm(inst[2]))
		case "add":
			v := getRegImm(inst[1]) + getRegImm(inst[2])
			setReg(inst[1], v)
		case "mul":
			v := getRegImm(inst[1]) * getRegImm(inst[2])
			setReg(inst[1], v)
		case "mod":
			v := getRegImm(inst[1]) % getRegImm(inst[2])
			setReg(inst[1], v)
		case "rcv":
			if getRegImm(inst[1]) != 0 {
				// program done (for part 1)
				fmt.Printf("LASTPLAY : %d\n", lastplay)
				return
			}
		case "jgz":
			if getRegImm(inst[1]) > 0 {
				next = ip + getRegImm(inst[2])
			}
		default:
			abort("ERROR invalid inst %d : %v", ip, program[ip])
		}
		ip = next
	}
}
