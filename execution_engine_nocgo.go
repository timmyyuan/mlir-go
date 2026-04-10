//go:build !cgo

package mlir

import (
	"fmt"
	"unsafe"
)

type ExecutionEngine struct{}

func NewExecutionEngine(*Module, int) (*ExecutionEngine, error) {
	return nil, fmt.Errorf("mlir: cgo is required")
}

func (e *ExecutionEngine) Close() error {
	return fmt.Errorf("mlir: cgo is required")
}

func (e *ExecutionEngine) InvokePacked(string, unsafe.Pointer, ...unsafe.Pointer) error {
	return fmt.Errorf("mlir: cgo is required")
}
