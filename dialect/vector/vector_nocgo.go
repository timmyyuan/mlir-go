//go:build !cgo

package vector

import mlir "github.com/timmyyuan/mlir-go"

func Broadcast(loc mlir.Location, source mlir.Value, resultType mlir.Type) (*mlir.OwnedOperation, error) {
	return BroadcastOp(loc, source, resultType)
}

func ExtractElement(loc mlir.Location, vector mlir.Value, resultType mlir.Type, position ...mlir.Value) (*mlir.OwnedOperation, error) {
	return ExtractElementOp(loc, vector, resultType, position...)
}
