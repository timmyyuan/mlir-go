//go:build !cgo

package cf

import (
	"fmt"

	mlir "github.com/timmyyuan/mlir-go"
	"github.com/timmyyuan/mlir-go/builder"
)

func Branch(*builder.Builder, mlir.Block, ...mlir.Value) (mlir.Operation, error) {
	return mlir.Operation{}, fmt.Errorf("mlir: cgo is required")
}

func CondBranch(*builder.Builder, mlir.Value, mlir.Block, []mlir.Value, mlir.Block, []mlir.Value) (mlir.Operation, error) {
	return mlir.Operation{}, fmt.Errorf("mlir: cgo is required")
}
