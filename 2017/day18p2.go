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

// registers are named by single latin letter, so there are 26
// for register x, use index (x - 'a')
type RegisterFile [26]int

// set register's value
func (r *RegisterFile) setReg(reg string, val int) {
	if reg[0] < 'a' || reg[0] > 'z' {
		panic(fmt.Sprintf("ERROR invalid register %q", reg))
	}
	(*r)[reg[0]-'a'] = val
}

// get a register's value *or* an immediate value
func (r *RegisterFile) getRegImm(arg string) int {
	if arg[0] >= 'a' && arg[0] <= 'z' {
		return (*r)[arg[0]-'a']
	} else {
		val, err := strconv.Atoi(arg)
		if err != nil {
			panic(fmt.Sprintf("invalid arg %q", arg))
		}
		return val
	}
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

	// two "procs", each with its own registers, instruction-pointer, and recv chan
	p0_regs := RegisterFile{}
	p1_regs := RegisterFile{}
	p0_ip := 0
	p1_ip := 0
	p0_in := make(chan int, 999) // enough buffer to never block sends
	p1_in := make(chan int, 999)
	sends := [2]int{0, 0}

	// init "p" register to proc id
	p0_regs.setReg("p", 0)
	p1_regs.setReg("p", 1)

	// active/foreground proc id
	p := 0

	// active proc's regs, ip, etc
	ip := &p0_ip
	in := p0_in
	out := p1_in
	regs := &p0_regs

	for { // run current proc
		inst := program[*ip]
		next := *ip + 1

		fmt.Printf("%d : %3d : %v\n", p, ip, inst)
		switch inst[0] {
		case "set":
			regs.setReg(inst[1], regs.getRegImm(inst[2]))
		case "add":
			v := regs.getRegImm(inst[1]) + regs.getRegImm(inst[2])
			regs.setReg(inst[1], v)
		case "mul":
			v := regs.getRegImm(inst[1]) * regs.getRegImm(inst[2])
			regs.setReg(inst[1], v)
		case "mod":
			v := regs.getRegImm(inst[1]) % regs.getRegImm(inst[2])
			regs.setReg(inst[1], v)
		case "jgz":
			if regs.getRegImm(inst[1]) > 0 {
				next = *ip + regs.getRegImm(inst[2])
			}
		case "snd":
			v := regs.getRegImm(inst[1])
			out <- v
			sends[p]++
			fmt.Printf("SND(%d) %d->%d : %d\n", sends[p], p, 1-p, v)
		case "rcv":
			select {
			case v := <-in:
				regs.setReg(inst[1], v)
				fmt.Printf("RCV %d\n", v)
			default:
				fmt.Printf("SWITCH PROC %d -> %d\n", p, 1-p)
				if p == 0 {
					p = 1
					ip = &p1_ip
					in = p1_in
					out = p0_in
					regs = &p1_regs
				} else {
					p = 0
					ip = &p0_ip
					in = p0_in
					out = p1_in
					regs = &p0_regs
				}
				if len(in) == 0 {
					fmt.Printf("DEADLOCK p0_ip=%d p1_ip=%d\n", p0_ip, p1_ip)
					fmt.Printf("SENDS p0=%d p1=%d\n", sends[0], sends[1])
					os.Exit(0)
				}
				next = *ip
			}
		default:
			abort("ERROR invalid inst %d : %v", *ip, program[*ip])
		}
		*ip = next
	}
}
