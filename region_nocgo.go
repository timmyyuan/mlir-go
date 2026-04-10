//go:build !cgo

package mlir

import "fmt"

type Region struct{}

func (r Region) IsNull() bool {
	return true
}

func (r Region) Blocks() []Block {
	return nil
}

func (r Region) AppendOwnedBlock(*OwnedBlock) (Block, error) {
	return Block{}, fmt.Errorf("mlir: cgo is required")
}
