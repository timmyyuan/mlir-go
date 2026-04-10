//go:build !cgo

package mlir

import "fmt"

type Operation struct{}

func (op Operation) IsNull() bool {
	return true
}

func (op Operation) String() string {
	return ""
}

func (op Operation) Verify() error {
	return fmt.Errorf("mlir: cgo is required")
}

func (op Operation) Name() string {
	return ""
}

func (op Operation) Location() Location {
	return Location{}
}

func (op Operation) Block() Block {
	return Block{}
}

func (op Operation) ParentOperation() Operation {
	return Operation{}
}

func (op Operation) Regions() []Region {
	return nil
}

func (op Operation) Results() []Value {
	return nil
}

func (op Operation) Operands() []Value {
	return nil
}

func (op Operation) SetOperand(int, Value) error {
	return fmt.Errorf("mlir: cgo is required")
}

func (op Operation) SetOperands(...Value) error {
	return fmt.Errorf("mlir: cgo is required")
}

func (op Operation) Successors() []Block {
	return nil
}

func (op Operation) SetSuccessor(int, Block) error {
	return fmt.Errorf("mlir: cgo is required")
}

func (op Operation) Attributes() []NamedAttribute {
	return nil
}

func (op Operation) Attribute(string) Attribute {
	return Attribute{}
}

func (op Operation) SetAttribute(string, Attribute) error {
	return fmt.Errorf("mlir: cgo is required")
}

func (op Operation) RemoveAttribute(string) (bool, error) {
	return false, fmt.Errorf("mlir: cgo is required")
}
