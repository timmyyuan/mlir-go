//go:build cgo

package funcdialect

import mlir "github.com/timmyyuan/mlir-go"

func Func(ctx *mlir.Context, loc mlir.Location, name string, functionType mlir.Type, body *mlir.OwnedRegion) (*mlir.OwnedOperation, error) {
	return FuncOp(ctx, loc, name, functionType, body)
}

func Call(ctx *mlir.Context, loc mlir.Location, callee string, results []mlir.Type, operands ...mlir.Value) (*mlir.OwnedOperation, error) {
	return CallOp(ctx, loc, callee, results, operands...)
}

func Return(loc mlir.Location, values ...mlir.Value) (*mlir.OwnedOperation, error) {
	return ReturnOp(loc, values...)
}
