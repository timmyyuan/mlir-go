//go:build cgo

package arith

import mlir "github.com/timmyyuan/mlir-go"

func Constant(ctx *mlir.Context, loc mlir.Location, resultType mlir.Type, value mlir.Attribute) (*mlir.OwnedOperation, error) {
	return ConstantOp(ctx, loc, resultType, value)
}

func AddI(loc mlir.Location, lhs, rhs mlir.Value) (*mlir.OwnedOperation, error) {
	return AddIOp(loc, lhs, rhs)
}

func SubI(loc mlir.Location, lhs, rhs mlir.Value) (*mlir.OwnedOperation, error) {
	return SubIOp(loc, lhs, rhs)
}

func MulI(loc mlir.Location, lhs, rhs mlir.Value) (*mlir.OwnedOperation, error) {
	return MulIOp(loc, lhs, rhs)
}

func ExtUI(loc mlir.Location, input mlir.Value, resultType mlir.Type) (*mlir.OwnedOperation, error) {
	return ExtUIOp(loc, input, resultType)
}

func ExtSI(loc mlir.Location, input mlir.Value, resultType mlir.Type) (*mlir.OwnedOperation, error) {
	return ExtSIOp(loc, input, resultType)
}

func TruncI(loc mlir.Location, input mlir.Value, resultType mlir.Type) (*mlir.OwnedOperation, error) {
	return TruncIOp(loc, input, resultType)
}

func SIToFP(loc mlir.Location, input mlir.Value, resultType mlir.Type) (*mlir.OwnedOperation, error) {
	return SIToFPOp(loc, input, resultType)
}

func FPToSI(loc mlir.Location, input mlir.Value, resultType mlir.Type) (*mlir.OwnedOperation, error) {
	return FPToSIOp(loc, input, resultType)
}

func IndexCast(loc mlir.Location, input mlir.Value, resultType mlir.Type) (*mlir.OwnedOperation, error) {
	return IndexCastOp(loc, input, resultType)
}

func CmpI(ctx *mlir.Context, loc mlir.Location, predicate string, lhs, rhs mlir.Value) (*mlir.OwnedOperation, error) {
	return CmpIOp(ctx, loc, predicate, lhs, rhs)
}
