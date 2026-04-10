//go:build !cgo

package mlir

import "fmt"

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

func (b Block) ParentOperation() Operation {
	return Operation{}
}

func (b Block) ParentRegion() Region {
	return Region{}
}

func (b Block) Terminator() Operation {
	return Operation{}
}

func (b Block) AppendOwnedOperation(*OwnedOperation) (Operation, error) {
	return Operation{}, fmt.Errorf("mlir: cgo is required")
}

func (b Block) InsertOwnedOperationBefore(Operation, *OwnedOperation) (Operation, error) {
	return Operation{}, fmt.Errorf("mlir: cgo is required")
}

func (b Block) AddArgument(Type, Location) (Value, error) {
	return Value{}, fmt.Errorf("mlir: cgo is required")
}
