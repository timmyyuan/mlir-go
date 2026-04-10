//go:build !cgo

package memref

import (
	"fmt"

	mlir "github.com/timmyyuan/mlir-go"
)

func Alloca(*mlir.Context, mlir.Location, mlir.Type, ...mlir.Value) (*mlir.OwnedOperation, error) {
	return nil, fmt.Errorf("mlir: cgo is required")
}

func Load(mlir.Location, mlir.Value, ...mlir.Value) (*mlir.OwnedOperation, error) {
	return nil, fmt.Errorf("mlir: cgo is required")
}

func Store(mlir.Location, mlir.Value, mlir.Value, ...mlir.Value) (*mlir.OwnedOperation, error) {
	return nil, fmt.Errorf("mlir: cgo is required")
}
