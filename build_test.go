//go:build cgo

package mlir

import "testing"

func TestBuildEmptyModuleOp(t *testing.T) {
	ctx, err := NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer func() {
		if err := ctx.Close(); err != nil {
			t.Fatalf("Close() error = %v", err)
		}
	}()

	loc, err := UnknownLocation(ctx)
	if err != nil {
		t.Fatalf("UnknownLocation() error = %v", err)
	}

	region, err := NewOwnedRegion()
	if err != nil {
		t.Fatalf("NewOwnedRegion() error = %v", err)
	}
	block, err := NewOwnedBlock(nil, nil)
	if err != nil {
		t.Fatalf("NewOwnedBlock() error = %v", err)
	}
	if _, err := region.AppendBlock(block); err != nil {
		t.Fatalf("region.AppendBlock() error = %v", err)
	}

	state := NewOperationState("builtin.module", loc)
	state.AddOwnedRegions(region)
	moduleOp, err := CreateOperation(state)
	if err != nil {
		t.Fatalf("CreateOperation() error = %v", err)
	}
	defer func() {
		if err := moduleOp.Close(); err != nil {
			t.Fatalf("moduleOp.Close() error = %v", err)
		}
	}()

	op := moduleOp.Operation()
	if err := op.Verify(); err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	runFileCheck(t, op.String(), "empty_module.mlir")
}

func TestGenericConstructionAndMutation(t *testing.T) {
	ctx, err := NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer func() {
		if err := ctx.Close(); err != nil {
			t.Fatalf("Close() error = %v", err)
		}
	}()

	if err := ctx.RegisterAllDialects(); err != nil {
		t.Fatalf("RegisterAllDialects() error = %v", err)
	}

	const input = `module {
  func.func @f(%arg0: i32) -> i32 {
    return %arg0 : i32
  }
}
`

	mod, err := ParseModule(ctx, input)
	if err != nil {
		t.Fatalf("ParseModule() error = %v", err)
	}
	defer func() {
		if err := mod.Close(); err != nil {
			t.Fatalf("mod.Close() error = %v", err)
		}
	}()

	funcBlock := mod.Body().Operations()[0].Regions()[0].Blocks()[0]
	arg0 := funcBlock.Arguments()[0]
	ret := funcBlock.Operations()[0]
	i32, err := ParseType(ctx, "i32")
	if err != nil {
		t.Fatalf("ParseType() error = %v", err)
	}
	attr, err := ParseAttribute(ctx, "42 : i32")
	if err != nil {
		t.Fatalf("ParseAttribute() error = %v", err)
	}
	id, err := InternIdentifier(ctx, "value")
	if err != nil {
		t.Fatalf("InternIdentifier() error = %v", err)
	}
	loc, err := UnknownLocation(ctx)
	if err != nil {
		t.Fatalf("UnknownLocation() error = %v", err)
	}

	state := NewOperationState("arith.constant", loc)
	state.AddResults(i32)
	state.AddAttributes(NewNamedAttribute(id, attr))
	constantOwned, err := CreateOperation(state)
	if err != nil {
		t.Fatalf("CreateOperation(constant) error = %v", err)
	}

	ownedBlock, err := NewOwnedBlock(nil, nil)
	if err == nil {
		v, err := ownedBlock.AddArgument(i32, loc)
		if err != nil {
			t.Fatalf("ownedBlock.AddArgument() error = %v", err)
		}
		if v.IsNull() {
			t.Fatalf("ownedBlock.AddArgument() returned null value")
		}
		if err := ownedBlock.Close(); err != nil {
			t.Fatalf("ownedBlock.Close() error = %v", err)
		}
	} else {
		t.Fatalf("NewOwnedBlock() error = %v", err)
	}

	insertedConstant, err := funcBlock.InsertOwnedOperationBefore(ret, constantOwned)
	if err != nil {
		t.Fatalf("InsertOperationBefore() error = %v", err)
	}
	constResult := insertedConstant.Results()[0]
	if err := arg0.ReplaceAllUsesWith(constResult); err != nil {
		t.Fatalf("ReplaceAllUsesWith() error = %v", err)
	}

	if err := mod.Verify(); err != nil {
		t.Fatalf("Verify() after mutation error = %v", err)
	}
	runFileCheck(t, mod.String(), "constant_replace_module.mlir")
}
