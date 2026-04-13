//go:build cgo

package cf

import (
	"fmt"

	mlir "github.com/timmyyuan/mlir-go"
	"github.com/timmyyuan/mlir-go/builder"
	lowcf "github.com/timmyyuan/mlir-go/dialect/cf"
)

func Branch(b *builder.Builder, dest mlir.Block, operands ...mlir.Value) (mlir.Operation, error) {
	if b == nil {
		return mlir.Operation{}, fmt.Errorf("mlir: builder is required")
	}
	return b.EmitWithLocation(func(loc mlir.Location) (*mlir.OwnedOperation, error) {
		return lowcf.Branch(loc, dest, operands...)
	})
}

func CondBranch(b *builder.Builder, cond mlir.Value, trueDest mlir.Block, trueOperands []mlir.Value, falseDest mlir.Block, falseOperands []mlir.Value) (mlir.Operation, error) {
	if b == nil {
		return mlir.Operation{}, fmt.Errorf("mlir: builder is required")
	}
	return b.EmitWithLocation(func(loc mlir.Location) (*mlir.OwnedOperation, error) {
		return lowcf.CondBranch(b.Context(), loc, cond, trueDest, trueOperands, falseDest, falseOperands)
	})
}
