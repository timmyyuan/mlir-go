//go:build !cgo

package mlir

import "fmt"

func (op Operation) SetFuncArgumentAttribute(int, string, Attribute) error {
	return fmt.Errorf("mlir: cgo is required")
}
