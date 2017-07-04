package radix

func New() *Radix {
	return &Radix{}
}

type Radix struct {
	node *node
}

func (r *Radix) Lookup(label string) interface{} {
	if r.node == nil {
		panic("empty node")
	}

	traverseNode := r.node
walk:
	for _, e := range traverseNode.edges {
		if len(label) < len(e.label) {
			continue
		}
		if label[:len(e.label)] != e.label {
			continue
		}
		if len(label) == len(e.label) {
			// label is equal to e.label
			if e.node == nil || e.node.leaf == nil {
				panic("not found")
			}
			return e.node.leaf.val
		}

		label = label[len(e.label):]
		traverseNode = e.node
		goto walk
	}

	panic("not found")
}

func (r *Radix) Insert(label string, val interface{}) {
	if r.node == nil {
		r.node = new(node)
	}

	traverseNode := r.node
walk:
	if len(traverseNode.edges) == 0 {
		// create a new edge if the current node has no edge.
		traverseNode.edges = []*edge{newEdge(label, newNode(newLeaf(val)))}
		return
	}

	for _, e := range traverseNode.edges {
		i := 0
		max := min(len(label), len(e.label))
		// find the longest common prefix.
		for i < max && label[i] == e.label[i] {
			i++
		}

		if i == 0 {
			continue
		}

		// a label contains the other label.
		if i == max {
			if len(e.label) == len(label) {
				// label is equal to e.label
				if e.node != nil && e.node.leaf != nil {
					panic("duplicate node")
				}
				if e.node == nil {
					e.node = newNode(newLeaf(val))
					return
				}
				// in order to keep node's original edges, update node's leaf
				// if node is non-nil, not create a new node.
				e.node.leaf = newLeaf(val)
				return
			} else if i == len(e.label) {
				// e.label is a substring of label.
				label = label[i:]
				// set current edge's node as traverse node and walk.
				traverseNode = e.node
				goto walk
			} else if i == len(label) {
				// label is a substring of e.label.
				// save current edge temporarily and update it's label.
				originalEdge := newEdge(e.label[i:], e.node)
				// change current edge's label and node
				e.label = label
				e.node = newNode(newLeaf(val), originalEdge)
				return
			}
		}

		// create a new node to store two edges.
		newNode := newNode(
			nil,
			newEdge(e.label[i:], e.node), // original edge
			newEdge(
				label[i:],
				newNode(newLeaf(val)),
			), // new edge
		)

		// update current edge's label and node
		e.label = label[:i]
		e.node = newNode
		return
	}

	// create a new edge if no matched edge is found.
	traverseNode.edges = append(traverseNode.edges, newEdge(label, newNode(newLeaf(val))))
}

func (r *Radix) Delete(label string) {
	var parentNode *node
	var idxOfParent int
	traverseNode := r.node
walk:
	if traverseNode == nil {
		return
	}
	for idx, e := range traverseNode.edges {
		if len(label) < len(e.label) || label[:len(e.label)] != e.label {
			continue
		}
		if len(label) == len(e.label) {
			if e.node != nil && e.node.leaf != nil {
				if len(e.node.edges) == 0 {
					// delete current edge from traverse node.
					traverseNode.edges = append(traverseNode.edges[:idx], traverseNode.edges[idx+1:]...)
				} else if len(e.node.edges) == 1 {
					// replace current edge with it's node's only edge.
					e.node.edges[0].label = e.label + e.node.edges[0].label
					traverseNode.edges[idx] = e.node.edges[0]
				} else {
					// set current edge's node's leaf as nil.
					e.node.leaf = nil
				}
				if len(traverseNode.edges) == 1 && traverseNode.leaf == nil && parentNode != nil {
					parentNode.edges[idxOfParent].label = parentNode.edges[idxOfParent].label + traverseNode.edges[0].label
					parentNode.edges[idxOfParent].node = traverseNode.edges[0].node
				}
			}
			return
		}
		label = label[len(e.label):]
		parentNode = traverseNode
		idxOfParent = idx
		traverseNode = e.node
		goto walk
	}
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
