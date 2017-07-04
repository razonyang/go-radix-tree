package radix

func newNode(leaf *leaf, edges ...*edge) *node {
	return &node{
		leaf:  leaf,
		edges: edges,
	}
}

type node struct {
	leaf  *leaf
	edges []*edge
}

func (n *node) isLeaf() bool {
	return n.leaf != nil
}
