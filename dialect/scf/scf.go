//go:build cgo

package scf

import (
	"fmt"

	mlir "github.com/timmyyuan/mlir-go"
)

func If(loc mlir.Location, condition mlir.Value, results []mlir.Type, thenRegion *mlir.OwnedRegion, elseRegion *mlir.OwnedRegion) (*mlir.OwnedOperation, error) {
	if condition.IsNull() {
		return nil, fmt.Errorf("mlir: condition is required")
	}
	if thenRegion == nil {
		return nil, fmt.Errorf("mlir: then region is required")
	}
	state := mlir.NewOperationState("scf.if", loc)
	state.AddResults(results...)
	state.AddOperands(condition)
	state.AddOwnedRegions(thenRegion)
	if elseRegion != nil {
		state.AddOwnedRegions(elseRegion)
	}
	return mlir.CreateOperation(state)
}

func Yield(loc mlir.Location, values ...mlir.Value) (*mlir.OwnedOperation, error) {
	state := mlir.NewOperationState("scf.yield", loc)
	state.AddOperands(values...)
	return mlir.CreateOperation(state)
}
