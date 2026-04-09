//go:build !cgo

package mlir

type Block struct{}

func (b Block) IsNull() bool {
	return true
}

func (b Block) String() string {
	return ""
}

func (b Block) Arguments() []Value {
	return nil
}

func (b Block) Operations() []Operation {
	return nil
}
