//go:build cgo

package linalg

import mlir "github.com/timmyyuan/mlir-go"

func Index(ctx *mlir.Context, loc mlir.Location, dim int64) (*mlir.OwnedOperation, error) {
	return IndexOp(ctx, loc, dim)
}

func Yield(loc mlir.Location, values ...mlir.Value) (*mlir.OwnedOperation, error) {
	return YieldOp(loc, values...)
}
