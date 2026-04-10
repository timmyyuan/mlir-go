//go:build !cgo

package arith

import (
	"fmt"

	mlir "github.com/timmyyuan/mlir-go"
)

func Constant(*mlir.Context, mlir.Location, mlir.Type, mlir.Attribute) (*mlir.OwnedOperation, error) {
	return nil, fmt.Errorf("mlir: cgo is required")
}

func AddI(mlir.Location, mlir.Value, mlir.Value) (*mlir.OwnedOperation, error) {
	return nil, fmt.Errorf("mlir: cgo is required")
}
