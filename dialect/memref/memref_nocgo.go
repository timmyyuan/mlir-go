//go:build !cgo

package memref

import mlir "github.com/timmyyuan/mlir-go"

func Alloca(ctx *mlir.Context, loc mlir.Location, resultType mlir.Type, dynamicSizes ...mlir.Value) (*mlir.OwnedOperation, error) {
	return AllocaOp(ctx, loc, resultType, dynamicSizes...)
}

func Load(loc mlir.Location, memref mlir.Value, indices ...mlir.Value) (*mlir.OwnedOperation, error) {
	return LoadOp(loc, memref, indices...)
}

func Store(loc mlir.Location, value mlir.Value, memref mlir.Value, indices ...mlir.Value) (*mlir.OwnedOperation, error) {
	return StoreOp(loc, value, memref, indices...)
}
