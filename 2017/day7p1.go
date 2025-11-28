package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Node struct {
	Name     string
	Weight   int
	Parent   string
	Children []string
}

func abort(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, "ERROR: "+msg+"\n", args...)
	os.Exit(1)
}

func main() {
	nodes := map[string]*Node{}

	// match example lines:
	// ktlj (57)
	// fwft (72) -> ktlj, cntj, xhth
	lineRe := regexp.MustCompile(`([a-z]+) \(([0-9]+)\)(?: -> ([a-z]+(?:, [a-z]+)*))?`)

	// parse nodes + children refs
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := lineRe.FindStringSubmatch(line)
		if len(fields) == 0 {
			abort("failed to parse line: %q", line)
		}
		name := fields[1]

		weight, err := strconv.Atoi(fields[2])
		if err != nil {
			abort("failed to parse weight from %q: %s", line, err)
		}

		var children []string
		if len(fields[3]) != 0 {
			children = strings.Split(fields[3], ", ")
		}

		nodes[name] = &Node{
			Name:     name,
			Weight:   weight,
			Children: children,
		}
	}

	// add parent refs
	for nn, node := range nodes {
		for _, cn := range node.Children {
			child, ok := nodes[cn]
			if !ok {
				abort("node %s child %s not found", nn, cn)
			}
			if child.Parent != "" {
				abort("node %s child %s already has parent %s", nn, cn, child.Parent)
			}
			child.Parent = nn
		}
	}

	// get 10 random-ish nodes' chain to root, for quick rough consistency check
	limit := 10
	for nn, node := range nodes {
		p := node
		chain := []string{nn}

		for p.Parent != "" {
			p = nodes[p.Parent]
			chain = append(chain, p.Name)
		}

		fmt.Printf("%v\n", chain)

		limit--
		if limit <= 0 {
			break
		}
	}
	// ...if all printed chains to root end with the same node, that's the answer
}
