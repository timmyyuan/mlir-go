//go:build !cgo

package mlir

import "fmt"

type Module struct{}

func ParseModule(*Context, string) (*Module, error) {
	return nil, fmt.Errorf("mlir: cgo is required")
}

func CreateEmptyModule(Location) (*Module, error) {
	return nil, fmt.Errorf("mlir: cgo is required")
}

func (m *Module) Close() error {
	return fmt.Errorf("mlir: cgo is required")
}

func (m *Module) Operation() Operation {
	return Operation{}
}

func (m *Module) Body() Block {
	return Block{}
}

func (m *Module) Verify() error {
	return fmt.Errorf("mlir: cgo is required")
}

func (m *Module) String() string {
	return ""
}
