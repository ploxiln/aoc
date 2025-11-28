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
	TowerWgh int
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

	// check 10 random-ish nodes chain to same root
	root := ""
	limit := 10
	for nn, node := range nodes {
		p := node
		for p.Parent != "" {
			p = nodes[p.Parent]
		}

		if p.Name != root {
			if root == "" {
				root = p.Name
			} else {
				abort("inconsistent root nodes %s != %s (on child %s after %d checks)",
					root, p.Name, nn, 10-limit)
			}
		}

		limit--
		if limit <= 0 {
			break
		}
	}

	calcTowerWeight(nodes, nodes[root], 0)
	// Manually examine "Unbalanced children" printouts, to find deepest
	// child with inconsistency. It's easy enough to see it, and by how much.
	// (Yeah, it's lazy to skip writing the code to automate this last bit ...)
}

// recursively calculate "tower" weight of each subtree
// (and check for unbalanced children)
func calcTowerWeight(nodes map[string]*Node, node *Node, depth int) {
	if node.TowerWgh != 0 {
		abort("node visited twice: %s", node.Name)
	}
	if len(node.Children) == 0 {
		node.TowerWgh = node.Weight
		return
	}

	children := make([]int, len(node.Children))
	total := 0
	for i, cn := range node.Children {
		child := nodes[cn]
		calcTowerWeight(nodes, child, depth+1)
		children[i] = child.TowerWgh
		total += child.TowerWgh
	}
	for _, wgh := range children[1:] {
		if wgh != children[0] {
			fmt.Printf("Unbalanced children for %s (%d):\n", node.Name, node.Weight)
			for i, cn := range node.Children {
				weight := nodes[cn].Weight
				fmt.Printf("	%s (%d): %d\n", cn, weight, children[i])
			}
			break
		}
	}
	node.TowerWgh = node.Weight + total
}
