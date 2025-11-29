package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func abort(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, "ERROR "+msg+"\n", args...)
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
	nodes := map[int][]int{}
	lineCount := 0
	connCount := 0

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString(byte('\n'))
		if err != nil && err != io.EOF {
			abort("reading stdin: %s", err)
		}
		if line == "" && err == io.EOF {
			break
		}
		line = strings.TrimSpace(line)
		nodestr, connstr, ok := strings.Cut(line, " <-> ")
		if !ok {
			abort("failed to parse line %q", line)
		}
		lineCount++

		connsplit := strings.Split(connstr, ", ")
		conns := make([]int, len(connsplit))
		for i, s := range connsplit {
			conns[i] = toInt(s)
			connCount++
		}
		nodes[toInt(nodestr)] = conns

		if err == io.EOF {
			break
		}
	}
	fmt.Printf("parsed input lines=%d nodes=%d connections=%d\n", lineCount, len(nodes), connCount)

	// set of nodes which are already connected to this group
	inGroup := map[int]bool{}
	// poor-man's queue of nodes to visit/process, start with node 0
	queue := []int{0}

	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		fmt.Printf("pop queue n=%d len=%d\n", n, len(queue)) // debug

		// process this node if not already visited (processed)
		if !inGroup[n] {
			inGroup[n] = true
			queue = append(queue, nodes[n]...)
		}
	}
	fmt.Printf("group len=%d\n", len(inGroup))
}
