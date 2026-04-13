//go:build cgo

package tensor

import mlir "github.com/timmyyuan/mlir-go"

func Empty(loc mlir.Location, resultType mlir.Type, dynamicSizes ...mlir.Value) (*mlir.OwnedOperation, error) {
	return EmptyOp(loc, resultType, dynamicSizes...)
}

func Dim(loc mlir.Location, source mlir.Value, index mlir.Value) (*mlir.OwnedOperation, error) {
	return DimOp(loc, source, index)
}

func Cast(loc mlir.Location, source mlir.Value, resultType mlir.Type) (*mlir.OwnedOperation, error) {
	return CastOp(loc, source, resultType)
}
