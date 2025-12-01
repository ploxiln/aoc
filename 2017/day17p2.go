package main

import (
	"fmt"
	"os"
	"strconv"
)

const INSMAX = 50_000_000

// Shifting in an array is n**2, that works great for 2000 items, but takes many hours
// for 50 million items. We need to avoid shifting or otherwise updating O(n) items
// every time, and also avoid needing to walk O(n) to find the place to insert.
// Thus, a binary tree! Not ordered by node value, but rather keeping track
// of left/right sub-tree sizes so we can find position/index to insert in log2(n).

// This is a minimal simplistic implementation, to make it faster:
//  * make it a red-black tree to keep more balanced
//  * put nodes in linear array where children of (i) are (i*2) and (i*2+1)
// ... but this can finish in ~ 90 secs (and 3 GB mem), so it works.

type TNode struct {
	L  *TNode
	R  *TNode
	Sz  int
	Val int
}

func insertNode(root *TNode, pos, val int) {
	// after insert, update "Sz" (tree-size) of every parent along path to root
	path := []*TNode{root}

	n := root
	x := 0 // track pos of n
	if n.L != nil {
		x += n.L.Sz
	}

	for {
		if pos <= x {
			if n.L == nil {
				break
			}
			n = n.L
			x -= 1
			if n.R != nil {
				x -= n.R.Sz
			}
		} else {
			if n.R == nil {
				break
			}
			n = n.R
			x += 1
			if n.L != nil {
				x += n.L.Sz
			}
		}
		path = append(path, n)
	}

	nn := &TNode{
		Val: val,
		Sz:  1,
	}
	if pos == x {
		// insert on left, just before
		if n.L != nil {
			fmt.Printf("ERROR n.L!=nil pos=%d x=%d path=%v\n", pos, x, path)
			os.Exit(1)
		}
		n.L = nn
	} else {
		// insert on right, just after
		if n.R != nil {
			fmt.Printf("ERROR n.R!=nil pos=%d x=%d path=%v\n", pos, x, path)
			os.Exit(1)
		}
		n.R = nn
	}

	// incr tree-size of every parent of new node, up to root
	for _, p := range path {
		p.Sz++
	}
}

// left-to-right in-order traversal print, with limit (because millions)
func printTree(n *TNode, limit int) {
	if n.L != nil {
		printTree(n.L, limit)
		limit -= n.L.Sz
		if limit <= 0 {
			return
		}
	}
	fmt.Printf(" %d ", n.Val)
	limit--
	if limit <= 0 {
		return
	}
	if n.R != nil {
		printTree(n.R, limit)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("USAGE: day17p1 <step>\n")
		os.Exit(2)
	}

	step, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("ERROR parsing integer step %q: %s\n", os.Args[1], err)
		os.Exit(2)
	}

	// start with [0]
	ring := &TNode{
		Val: 0,
		Sz:  1,
	}
	pos := 0

	for i := 1; i <= INSMAX; i++ {
		// start at previous inserted item (pos), step forward (step),
		// wrap around at current array size (i),
		// then insert value i *after* that (+1)
		pos = ((pos + step) % i) + 1

		insertNode(ring, pos, i)

		if i < 10 { // debug print
			fmt.Printf("i=%d ring=", i)
			printTree(ring, 12)
			fmt.Printf("\n")
		}
		if i & 0xfffff == 0 { // progress update
			fmt.Printf("i=%d ring=", i)
			printTree(ring, 12)
			fmt.Printf("...\n")
		}
	}
	// 0 always stays at index 0 (because insert is always *after* selected position)
	// so only print the first few items in the "circular buffer", 2nd item is the answer
	fmt.Printf("FINAL ring=")
	printTree(ring, 12)
	fmt.Printf("\n")
}
