//go:build !cgo

package scf

import (
	"fmt"

	mlir "github.com/timmyyuan/mlir-go"
)

func If(mlir.Location, mlir.Value, []mlir.Type, *mlir.OwnedRegion, *mlir.OwnedRegion) (*mlir.OwnedOperation, error) {
	return nil, fmt.Errorf("mlir: cgo is required")
}

func Yield(mlir.Location, ...mlir.Value) (*mlir.OwnedOperation, error) {
	return nil, fmt.Errorf("mlir: cgo is required")
}
