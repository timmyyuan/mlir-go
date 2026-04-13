//go:build !cgo

package funcdialect

import (
	"fmt"

	mlir "github.com/timmyyuan/mlir-go"
	"github.com/timmyyuan/mlir-go/builder"
)

func Call(*builder.Builder, string, []mlir.Type, ...mlir.Value) (mlir.Operation, error) {
	return mlir.Operation{}, fmt.Errorf("mlir: cgo is required")
}

func Return(*builder.Builder, ...mlir.Value) (mlir.Operation, error) {
	return mlir.Operation{}, fmt.Errorf("mlir: cgo is required")
}
