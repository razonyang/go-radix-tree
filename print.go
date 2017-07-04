package radix

import (
	"fmt"
	"strings"
)

func Print(r *Radix) {
	printNode(r.node, 0)
}

func printNode(n *node, intent int) {
	if n == nil {
		return
	}
	for _, e := range n.edges {
		printEdge(e, intent)
	}
}

func printEdge(e *edge, intent int) {
	fmt.Printf("  %s%s :%v\n", strings.Repeat("--", intent), e.label, e.node.leaf)
	printNode(e.node, intent+2)
}
