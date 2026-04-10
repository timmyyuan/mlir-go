//go:build !cgo

package funcdialect

import (
	"fmt"

	mlir "github.com/timmyyuan/mlir-go"
)

func Func(*mlir.Context, mlir.Location, string, mlir.Type, *mlir.OwnedRegion) (*mlir.OwnedOperation, error) {
	return nil, fmt.Errorf("mlir: cgo is required")
}

func Call(*mlir.Context, mlir.Location, string, []mlir.Type, ...mlir.Value) (*mlir.OwnedOperation, error) {
	return nil, fmt.Errorf("mlir: cgo is required")
}

func Return(mlir.Location, ...mlir.Value) (*mlir.OwnedOperation, error) {
	return nil, fmt.Errorf("mlir: cgo is required")
}
