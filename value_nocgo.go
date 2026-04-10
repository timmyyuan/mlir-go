//go:build !cgo

package mlir

import "fmt"

type Value struct{}

func (v Value) IsNull() bool {
	return true
}

func (v Value) IsBlockArgument() bool {
	return false
}

func (v Value) IsOpResult() bool {
	return false
}

func (v Value) Type() Type {
	return Type{}
}

func (v Value) String() string {
	return ""
}

func (v Value) Equal(Value) bool {
	return false
}

func (v Value) OwnerBlock() Block {
	return Block{}
}

func (v Value) ArgumentNumber() int {
	return -1
}

func (v Value) OwnerOperation() Operation {
	return Operation{}
}

func (v Value) ResultNumber() int {
	return -1
}

func (v Value) SetType(Type) error {
	return fmt.Errorf("mlir: cgo is required")
}
