//go:build cgo

package arith

import (
	"fmt"

	mlir "github.com/timmyyuan/mlir-go"
)

func Constant(ctx *mlir.Context, loc mlir.Location, resultType mlir.Type, value mlir.Attribute) (*mlir.OwnedOperation, error) {
	if resultType.IsNull() {
		return nil, fmt.Errorf("mlir: result type is required")
	}
	if value.IsNull() {
		return nil, fmt.Errorf("mlir: value attribute is required")
	}
	valueAttr, err := mlir.NamedAttributeByName(ctx, "value", value)
	if err != nil {
		return nil, err
	}
	state := mlir.NewOperationState("arith.constant", loc)
	state.AddResults(resultType)
	state.AddAttributes(valueAttr)
	return mlir.CreateOperation(state)
}

func AddI(loc mlir.Location, lhs, rhs mlir.Value) (*mlir.OwnedOperation, error) {
	if lhs.IsNull() || rhs.IsNull() {
		return nil, fmt.Errorf("mlir: operands are required")
	}
	resultType := lhs.Type()
	if resultType.IsNull() || !resultType.Equal(rhs.Type()) {
		return nil, fmt.Errorf("mlir: addi operands must have the same non-null type")
	}
	state := mlir.NewOperationState("arith.addi", loc)
	state.AddResults(resultType)
	state.AddOperands(lhs, rhs)
	return mlir.CreateOperation(state)
}
