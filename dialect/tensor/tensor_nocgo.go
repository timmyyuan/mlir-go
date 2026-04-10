//go:build !cgo

package tensor

import (
	"fmt"

	mlir "github.com/timmyyuan/mlir-go"
)

func Empty(mlir.Location, mlir.Type, ...mlir.Value) (*mlir.OwnedOperation, error) {
	return nil, fmt.Errorf("mlir: cgo is required")
}

func Dim(mlir.Location, mlir.Value, mlir.Value) (*mlir.OwnedOperation, error) {
	return nil, fmt.Errorf("mlir: cgo is required")
}

func Cast(mlir.Location, mlir.Value, mlir.Type) (*mlir.OwnedOperation, error) {
	return nil, fmt.Errorf("mlir: cgo is required")
}
