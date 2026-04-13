//go:build cgo

package funcdialect

import (
	"fmt"

	mlir "github.com/timmyyuan/mlir-go"
	"github.com/timmyyuan/mlir-go/builder"
	lowfunc "github.com/timmyyuan/mlir-go/dialect/func"
)

func Call(b *builder.Builder, callee string, results []mlir.Type, operands ...mlir.Value) (mlir.Operation, error) {
	if b == nil {
		return mlir.Operation{}, fmt.Errorf("mlir: builder is required")
	}
	return b.EmitWithLocation(func(loc mlir.Location) (*mlir.OwnedOperation, error) {
		return lowfunc.Call(b.Context(), loc, callee, results, operands...)
	})
}

func Return(b *builder.Builder, values ...mlir.Value) (mlir.Operation, error) {
	if b == nil {
		return mlir.Operation{}, fmt.Errorf("mlir: builder is required")
	}
	return b.EmitWithLocation(func(loc mlir.Location) (*mlir.OwnedOperation, error) {
		return lowfunc.Return(loc, values...)
	})
}
