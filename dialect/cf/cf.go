//go:build cgo

package cf

import (
	"fmt"
	"strconv"

	mlir "github.com/timmyyuan/mlir-go"
)

func Branch(loc mlir.Location, dest mlir.Block, operands ...mlir.Value) (*mlir.OwnedOperation, error) {
	if dest.IsNull() {
		return nil, fmt.Errorf("mlir: branch destination is required")
	}
	state := mlir.NewOperationState("cf.br", loc)
	state.AddOperands(operands...)
	state.AddSuccessors(dest)
	return mlir.CreateOperation(state)
}

func CondBranch(ctx *mlir.Context, loc mlir.Location, cond mlir.Value, trueDest mlir.Block, trueOperands []mlir.Value, falseDest mlir.Block, falseOperands []mlir.Value) (*mlir.OwnedOperation, error) {
	if cond.IsNull() {
		return nil, fmt.Errorf("mlir: branch condition is required")
	}
	if trueDest.IsNull() || falseDest.IsNull() {
		return nil, fmt.Errorf("mlir: both branch destinations are required")
	}
	segmentAttr, err := mlir.ParseAttribute(ctx, "array<i32: 1, "+strconv.Itoa(len(trueOperands))+", "+strconv.Itoa(len(falseOperands))+">")
	if err != nil {
		return nil, err
	}
	segmentNamedAttr, err := mlir.NamedAttributeByName(ctx, "operandSegmentSizes", segmentAttr)
	if err != nil {
		return nil, err
	}
	state := mlir.NewOperationState("cf.cond_br", loc)
	state.AddOperands(cond)
	state.AddOperands(trueOperands...)
	state.AddOperands(falseOperands...)
	state.AddSuccessors(trueDest, falseDest)
	state.AddAttributes(segmentNamedAttr)
	return mlir.CreateOperation(state)
}
