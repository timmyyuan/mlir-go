//go:build !cgo

package cf

import mlir "github.com/timmyyuan/mlir-go"

func Branch(loc mlir.Location, dest mlir.Block, operands ...mlir.Value) (*mlir.OwnedOperation, error) {
	return BranchOp(loc, dest, operands...)
}

func CondBranch(ctx *mlir.Context, loc mlir.Location, cond mlir.Value, trueDest mlir.Block, trueOperands []mlir.Value, falseDest mlir.Block, falseOperands []mlir.Value) (*mlir.OwnedOperation, error) {
	return CondBranchOp(ctx, loc, cond, trueDest, trueOperands, falseDest, falseOperands)
}
