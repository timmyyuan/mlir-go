//go:build cgo

package mlir

import (
	"fmt"

	"github.com/timmyyuan/mlir-go/internal/capi"
)

// SetFuncArgumentAttribute sets an argument attribute on a func.func operation.
func (op Operation) SetFuncArgumentAttribute(pos int, name string, attr Attribute) error {
	if op.IsNull() {
		return fmt.Errorf("mlir: null operation")
	}
	if attr.IsNull() {
		return fmt.Errorf("mlir: null attribute")
	}
	if name == "" {
		return fmt.Errorf("mlir: attribute name is required")
	}
	capi.FuncSetArgAttr(op.raw, pos, name, attr.raw)
	return nil
}
