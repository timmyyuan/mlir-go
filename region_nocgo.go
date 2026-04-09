//go:build !cgo

package mlir

type Region struct{}

func (r Region) IsNull() bool {
	return true
}

func (r Region) Blocks() []Block {
	return nil
}
