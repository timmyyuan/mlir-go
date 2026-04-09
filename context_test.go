//go:build cgo

package mlir

import "testing"

func TestContextLifecycle(t *testing.T) {
	ctx, err := NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}

	if err := ctx.RegisterAllDialects(); err != nil {
		t.Fatalf("RegisterAllDialects() error = %v", err)
	}

	n, err := ctx.NumLoadedDialects()
	if err != nil {
		t.Fatalf("NumLoadedDialects() error = %v", err)
	}
	if n <= 0 {
		t.Fatalf("NumLoadedDialects() = %d, want > 0", n)
	}

	if err := ctx.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if err := ctx.Close(); err != nil {
		t.Fatalf("second Close() error = %v", err)
	}

	if _, err := ctx.NumLoadedDialects(); err == nil {
		t.Fatalf("NumLoadedDialects() after Close() = nil error, want non-nil")
	}
}

func TestParseModuleRoundTrip(t *testing.T) {
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

	const input = "module {\n}\n"

	mod, err := ParseModule(ctx, input)
	if err != nil {
		t.Fatalf("ParseModule() error = %v", err)
	}
	defer func() {
		if err := mod.Close(); err != nil {
			t.Fatalf("mod.Close() error = %v", err)
		}
	}()

	if err := mod.Verify(); err != nil {
		t.Fatalf("Verify() error = %v", err)
	}

	if got := mod.String(); got != input {
		t.Fatalf("String() = %q, want %q", got, input)
	}

	if op := mod.Operation(); op.IsNull() {
		t.Fatalf("Operation() returned a null handle")
	}
}

func TestIRTraversal(t *testing.T) {
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
  func.func @add(%arg0: i32) -> i32 {
    %0 = arith.addi %arg0, %arg0 : i32
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

	top := mod.Operation()
	if got := top.Name(); got != "builtin.module" {
		t.Fatalf("top.Name() = %q, want %q", got, "builtin.module")
	}
	if got := top.Location().String(); got == "" {
		t.Fatalf("top.Location().String() = empty, want non-empty")
	}

	regions := top.Regions()
	if len(regions) != 1 {
		t.Fatalf("len(top.Regions()) = %d, want 1", len(regions))
	}

	blocks := regions[0].Blocks()
	if len(blocks) != 1 {
		t.Fatalf("len(regions[0].Blocks()) = %d, want 1", len(blocks))
	}

	ops := blocks[0].Operations()
	if len(ops) != 1 {
		t.Fatalf("len(blocks[0].Operations()) = %d, want 1", len(ops))
	}
	if got := ops[0].Name(); got != "func.func" {
		t.Fatalf("ops[0].Name() = %q, want %q", got, "func.func")
	}

	funcRegions := ops[0].Regions()
	if len(funcRegions) != 1 {
		t.Fatalf("len(funcRegions) = %d, want 1", len(funcRegions))
	}

	funcBlocks := funcRegions[0].Blocks()
	if len(funcBlocks) != 1 {
		t.Fatalf("len(funcBlocks) = %d, want 1", len(funcBlocks))
	}

	args := funcBlocks[0].Arguments()
	if len(args) != 1 {
		t.Fatalf("len(args) = %d, want 1", len(args))
	}
	if !args[0].IsBlockArgument() {
		t.Fatalf("args[0].IsBlockArgument() = false, want true")
	}
	if got := args[0].Type().String(); got != "i32" {
		t.Fatalf("args[0].Type().String() = %q, want %q", got, "i32")
	}

	innerOps := funcBlocks[0].Operations()
	if len(innerOps) != 2 {
		t.Fatalf("len(innerOps) = %d, want 2", len(innerOps))
	}
	if got := innerOps[0].Name(); got != "arith.addi" {
		t.Fatalf("innerOps[0].Name() = %q, want %q", got, "arith.addi")
	}
	results := innerOps[0].Results()
	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
	if !results[0].IsOpResult() {
		t.Fatalf("results[0].IsOpResult() = false, want true")
	}
	if got := results[0].Type().String(); got != "i32" {
		t.Fatalf("results[0].Type().String() = %q, want %q", got, "i32")
	}
	if got := innerOps[1].Name(); got != "func.return" {
		t.Fatalf("innerOps[1].Name() = %q, want %q", got, "func.return")
	}
}

func TestCreateEmptyModuleWithLocation(t *testing.T) {
	ctx, err := NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer func() {
		if err := ctx.Close(); err != nil {
			t.Fatalf("Close() error = %v", err)
		}
	}()

	loc, err := FileLineColLocation(ctx, "test.mlir", 3, 9)
	if err != nil {
		t.Fatalf("FileLineColLocation() error = %v", err)
	}

	mod, err := CreateEmptyModule(loc)
	if err != nil {
		t.Fatalf("CreateEmptyModule() error = %v", err)
	}
	defer func() {
		if err := mod.Close(); err != nil {
			t.Fatalf("mod.Close() error = %v", err)
		}
	}()

	if got := mod.Operation().Location().String(); got != "loc(\"test.mlir\":3:9)" {
		t.Fatalf("module location = %q, want %q", got, "loc(\"test.mlir\":3:9)")
	}
	if body := mod.Body(); body.IsNull() {
		t.Fatalf("mod.Body() returned a null block")
	}
}

func TestParseType(t *testing.T) {
	ctx, err := NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer func() {
		if err := ctx.Close(); err != nil {
			t.Fatalf("Close() error = %v", err)
		}
	}()

	typ, err := ParseType(ctx, "tensor<4xi32>")
	if err != nil {
		t.Fatalf("ParseType() error = %v", err)
	}
	if got := typ.String(); got != "tensor<4xi32>" {
		t.Fatalf("typ.String() = %q, want %q", got, "tensor<4xi32>")
	}
}
