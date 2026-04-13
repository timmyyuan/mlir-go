//go:build cgo

package arith

import (
	"fmt"

	mlir "github.com/timmyyuan/mlir-go"
	"github.com/timmyyuan/mlir-go/builder"
	lowarith "github.com/timmyyuan/mlir-go/dialect/arith"
)

func Constant(b *builder.Builder, resultType mlir.Type, value mlir.Attribute) (mlir.Operation, error) {
	if b == nil {
		return mlir.Operation{}, fmt.Errorf("mlir: builder is required")
	}
	return b.EmitWithLocation(func(loc mlir.Location) (*mlir.OwnedOperation, error) {
		return lowarith.Constant(b.Context(), loc, resultType, value)
	})
}

func AddI(b *builder.Builder, lhs, rhs mlir.Value) (mlir.Operation, error) {
	if b == nil {
		return mlir.Operation{}, fmt.Errorf("mlir: builder is required")
	}
	return b.EmitWithLocation(func(loc mlir.Location) (*mlir.OwnedOperation, error) {
		return lowarith.AddI(loc, lhs, rhs)
	})
}
