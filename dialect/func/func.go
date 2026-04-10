//go:build cgo

package funcdialect

import (
	"fmt"

	mlir "github.com/timmyyuan/mlir-go"
)

func Func(ctx *mlir.Context, loc mlir.Location, name string, functionType mlir.Type, body *mlir.OwnedRegion) (*mlir.OwnedOperation, error) {
	if name == "" {
		return nil, fmt.Errorf("mlir: function name is required")
	}
	if functionType.IsNull() {
		return nil, fmt.Errorf("mlir: function type is required")
	}
	if body == nil {
		return nil, fmt.Errorf("mlir: function body region is required")
	}

	symName, err := mlir.StringAttribute(ctx, name)
	if err != nil {
		return nil, err
	}
	functionTypeAttr, err := mlir.TypeAttribute(functionType)
	if err != nil {
		return nil, err
	}
	nameAttr, err := mlir.NamedAttributeByName(ctx, "sym_name", symName)
	if err != nil {
		return nil, err
	}
	typeAttr, err := mlir.NamedAttributeByName(ctx, "function_type", functionTypeAttr)
	if err != nil {
		return nil, err
	}

	state := mlir.NewOperationState("func.func", loc)
	state.AddAttributes(nameAttr, typeAttr)
	state.AddOwnedRegions(body)
	return mlir.CreateOperation(state)
}

func Call(ctx *mlir.Context, loc mlir.Location, callee string, results []mlir.Type, operands ...mlir.Value) (*mlir.OwnedOperation, error) {
	if callee == "" {
		return nil, fmt.Errorf("mlir: callee is required")
	}
	calleeAttr, err := mlir.FlatSymbolRefAttribute(ctx, callee)
	if err != nil {
		return nil, err
	}
	namedCallee, err := mlir.NamedAttributeByName(ctx, "callee", calleeAttr)
	if err != nil {
		return nil, err
	}
	state := mlir.NewOperationState("func.call", loc)
	state.AddResults(results...)
	state.AddOperands(operands...)
	state.AddAttributes(namedCallee)
	return mlir.CreateOperation(state)
}

func Return(loc mlir.Location, values ...mlir.Value) (*mlir.OwnedOperation, error) {
	state := mlir.NewOperationState("func.return", loc)
	state.AddOperands(values...)
	return mlir.CreateOperation(state)
}
