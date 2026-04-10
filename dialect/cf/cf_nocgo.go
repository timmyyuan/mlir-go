//go:build !cgo

package cf

import (
	"fmt"

	mlir "github.com/timmyyuan/mlir-go"
)

func Branch(mlir.Location, mlir.Block, ...mlir.Value) (*mlir.OwnedOperation, error) {
	return nil, fmt.Errorf("mlir: cgo is required")
}

func CondBranch(*mlir.Context, mlir.Location, mlir.Value, mlir.Block, []mlir.Value, mlir.Block, []mlir.Value) (*mlir.OwnedOperation, error) {
	return nil, fmt.Errorf("mlir: cgo is required")
}
