package radix

func newEdge(label string, node *node) *edge {
	return &edge{
		label: label,
		node:  node,
	}
}

type edge struct {
	label string
	node  *node
}
