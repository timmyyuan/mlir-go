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

func (op Operation) Regions() []Region {
	return nil
}

func (op Operation) Results() []Value {
	return nil
}
