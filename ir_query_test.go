//go:build cgo

package mlir

import "testing"

func TestIRQueryAndMutation(t *testing.T) {
	ctx, err := NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer func() {
		if err := ctx.Close(); err != nil {
			t.Fatalf("ctx.Close() error = %v", err)
		}
	}()

	if err := ctx.RegisterAllDialects(); err != nil {
		t.Fatalf("RegisterAllDialects() error = %v", err)
	}

	const input = `module {
  func.func @branch(%arg0: i1) -> i32 {
    %0 = arith.constant 7 : i32
    cf.cond_br %arg0, ^bb1, ^bb2
  ^bb1:
    return %0 : i32
  ^bb2:
    return %0 : i32
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

	moduleBody := mod.Body()
	if moduleBody.ParentOperation().Name() != "builtin.module" {
		t.Fatalf("moduleBody.ParentOperation().Name() = %q, want %q", moduleBody.ParentOperation().Name(), "builtin.module")
	}

	funcOp := moduleBody.Operations()[0]
	if funcOp.ParentOperation().Name() != "builtin.module" {
		t.Fatalf("funcOp.ParentOperation().Name() = %q, want %q", funcOp.ParentOperation().Name(), "builtin.module")
	}
	if funcOp.Block().IsNull() || !funcOp.Block().ParentOperation().IsNull() && funcOp.Block().ParentOperation().Name() != "builtin.module" {
		t.Fatalf("funcOp.Block() returned unexpected block")
	}

	entry := funcOp.Regions()[0].Blocks()[0]
	thenBlock := funcOp.Regions()[0].Blocks()[1]
	elseBlock := funcOp.Regions()[0].Blocks()[2]
	arg0 := entry.Arguments()[0]
	if !arg0.OwnerBlock().IsNull() && arg0.OwnerBlock().String() == "" {
		t.Fatalf("arg0.OwnerBlock() returned invalid block")
	}
	if got := arg0.ArgumentNumber(); got != 0 {
		t.Fatalf("arg0.ArgumentNumber() = %d, want 0", got)
	}

	constant := entry.Operations()[0]
	condBr := entry.Terminator()
	if condBr.Name() != "cf.cond_br" {
		t.Fatalf("entry.Terminator().Name() = %q, want %q", condBr.Name(), "cf.cond_br")
	}
	if got := len(condBr.Operands()); got != 1 {
		t.Fatalf("len(condBr.Operands()) = %d, want 1", got)
	}
	if !condBr.Operands()[0].Equal(arg0) {
		t.Fatalf("condBr operand 0 should equal arg0")
	}
	if got := len(condBr.Successors()); got != 2 {
		t.Fatalf("len(condBr.Successors()) = %d, want 2", got)
	}
	if condBr.Successors()[0].IsNull() || condBr.Successors()[1].IsNull() {
		t.Fatalf("condBr successors should be non-null")
	}
	if err := condBr.SetSuccessor(0, thenBlock); err != nil {
		t.Fatalf("condBr.SetSuccessor() error = %v", err)
	}
	if err := condBr.SetOperand(0, arg0); err != nil {
		t.Fatalf("condBr.SetOperand() error = %v", err)
	}
	if err := condBr.SetOperands(arg0); err != nil {
		t.Fatalf("condBr.SetOperands() error = %v", err)
	}

	foundSymName := false
	for _, attr := range funcOp.Attributes() {
		if attr.Name().String() == "sym_name" {
			foundSymName = true
		}
	}
	if !foundSymName {
		t.Fatalf("funcOp.Attributes() missing sym_name")
	}

	tag, err := StringAttribute(ctx, "x")
	if err != nil {
		t.Fatalf("StringAttribute() error = %v", err)
	}
	if err := funcOp.SetAttribute("test.tag", tag); err != nil {
		t.Fatalf("funcOp.SetAttribute() error = %v", err)
	}
	if got, err := funcOp.Attribute("test.tag").StringValue(); err != nil || got != "x" {
		t.Fatalf("funcOp.Attribute(\"test.tag\").StringValue() = (%q, %v), want (%q, nil)", got, err, "x")
	}
	removed, err := funcOp.RemoveAttribute("test.tag")
	if err != nil {
		t.Fatalf("funcOp.RemoveAttribute() error = %v", err)
	}
	if !removed {
		t.Fatalf("funcOp.RemoveAttribute() = false, want true")
	}

	constResult := constant.Results()[0]
	if constResult.OwnerOperation().Name() != "arith.constant" {
		t.Fatalf("constResult.OwnerOperation().Name() = %q, want %q", constResult.OwnerOperation().Name(), "arith.constant")
	}
	if got := constResult.ResultNumber(); got != 0 {
		t.Fatalf("constResult.ResultNumber() = %d, want 0", got)
	}

	retThen := thenBlock.Terminator()
	if got := len(retThen.Operands()); got != 1 {
		t.Fatalf("len(retThen.Operands()) = %d, want 1", got)
	}
	if err := retThen.SetOperand(0, constResult); err != nil {
		t.Fatalf("retThen.SetOperand() error = %v", err)
	}
	retElse := elseBlock.Terminator()
	if err := retElse.SetOperands(constResult); err != nil {
		t.Fatalf("retElse.SetOperands() error = %v", err)
	}

	i64, err := SignlessIntegerType(ctx, 64)
	if err != nil {
		t.Fatalf("SignlessIntegerType(64) error = %v", err)
	}
	loc, err := UnknownLocation(ctx)
	if err != nil {
		t.Fatalf("UnknownLocation() error = %v", err)
	}
	block, err := NewOwnedBlock([]Type{arg0.Type()}, []Location{loc})
	if err != nil {
		t.Fatalf("NewOwnedBlock() error = %v", err)
	}
	defer func() {
		if err := block.Close(); err != nil {
			t.Fatalf("block.Close() error = %v", err)
		}
	}()
	blockArg := block.Block().Arguments()[0]
	if err := blockArg.SetType(i64); err != nil {
		t.Fatalf("blockArg.SetType() error = %v", err)
	}
	if !blockArg.Type().Equal(i64) {
		t.Fatalf("blockArg.Type() = %s, want %s", blockArg.Type(), i64)
	}

	if err := mod.Verify(); err != nil {
		t.Fatalf("mod.Verify() error = %v", err)
	}
}
