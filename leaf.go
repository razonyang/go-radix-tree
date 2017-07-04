package radix

func newLeaf(val interface{}) *leaf {
	return &leaf{val: val}
}

type leaf struct {
	val interface{}
}
