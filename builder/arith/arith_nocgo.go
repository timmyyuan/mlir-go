//go:build !cgo

package arith

import (
	"fmt"

	mlir "github.com/timmyyuan/mlir-go"
	"github.com/timmyyuan/mlir-go/builder"
)

func Constant(*builder.Builder, mlir.Type, mlir.Attribute) (mlir.Operation, error) {
	return mlir.Operation{}, fmt.Errorf("mlir: cgo is required")
}

func AddI(*builder.Builder, mlir.Value, mlir.Value) (mlir.Operation, error) {
	return mlir.Operation{}, fmt.Errorf("mlir: cgo is required")
}
