//go:build cgo

package mlir_test

import (
	"strings"
	"testing"

	mlir "github.com/timmyyuan/mlir-go"
	arith "github.com/timmyyuan/mlir-go/dialect/arith"
	funcdialect "github.com/timmyyuan/mlir-go/dialect/func"
)

func buildConstantFunction(t *testing.T, ctx *mlir.Context, loc mlir.Location, name string, fnType mlir.Type, resultType mlir.Type, value int64) *mlir.OwnedOperation {
	t.Helper()

	body, err := mlir.NewOwnedRegion()
	if err != nil {
		t.Fatalf("NewOwnedRegion() error = %v", err)
	}
	entry, err := mlir.NewOwnedBlock(nil, nil)
	if err != nil {
		t.Fatalf("NewOwnedBlock() error = %v", err)
	}
	attr, err := mlir.IntegerAttribute(resultType, value)
	if err != nil {
		t.Fatalf("IntegerAttribute() error = %v", err)
	}
	constant, err := arith.Constant(ctx, loc, resultType, attr)
	if err != nil {
		t.Fatalf("arith.Constant() error = %v", err)
	}
	constantOp, err := entry.AppendOperation(constant)
	if err != nil {
		t.Fatalf("entry.AppendOperation(constant) error = %v", err)
	}
	ret, err := funcdialect.Return(loc, constantOp.Results()[0])
	if err != nil {
		t.Fatalf("func.Return() error = %v", err)
	}
	if _, err := entry.AppendOperation(ret); err != nil {
		t.Fatalf("entry.AppendOperation(return) error = %v", err)
	}
	if _, err := body.AppendBlock(entry); err != nil {
		t.Fatalf("body.AppendBlock() error = %v", err)
	}
	fn, err := funcdialect.Func(ctx, loc, name, fnType, body)
	if err != nil {
		t.Fatalf("func.Func(%q) error = %v", name, err)
	}
	return fn
}

func buildCallerFunction(t *testing.T, ctx *mlir.Context, loc mlir.Location, name string, fnType mlir.Type, resultType mlir.Type, callee string) *mlir.OwnedOperation {
	t.Helper()

	body, err := mlir.NewOwnedRegion()
	if err != nil {
		t.Fatalf("NewOwnedRegion() error = %v", err)
	}
	entry, err := mlir.NewOwnedBlock(nil, nil)
	if err != nil {
		t.Fatalf("NewOwnedBlock() error = %v", err)
	}
	call, err := funcdialect.Call(ctx, loc, callee, []mlir.Type{resultType})
	if err != nil {
		t.Fatalf("func.Call() error = %v", err)
	}
	callOp, err := entry.AppendOperation(call)
	if err != nil {
		t.Fatalf("entry.AppendOperation(call) error = %v", err)
	}
	ret, err := funcdialect.Return(loc, callOp.Results()[0])
	if err != nil {
		t.Fatalf("func.Return() error = %v", err)
	}
	if _, err := entry.AppendOperation(ret); err != nil {
		t.Fatalf("entry.AppendOperation(return) error = %v", err)
	}
	if _, err := body.AppendBlock(entry); err != nil {
		t.Fatalf("body.AppendBlock() error = %v", err)
	}
	fn, err := funcdialect.Func(ctx, loc, name, fnType, body)
	if err != nil {
		t.Fatalf("func.Func(%q) error = %v", name, err)
	}
	return fn
}

func TestSymbolTableAPIs(t *testing.T) {
	ctx, err := mlir.NewContext()
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

	loc, err := mlir.UnknownLocation(ctx)
	if err != nil {
		t.Fatalf("UnknownLocation() error = %v", err)
	}
	i32, err := mlir.SignlessIntegerType(ctx, 32)
	if err != nil {
		t.Fatalf("SignlessIntegerType() error = %v", err)
	}
	fnType, err := mlir.FunctionType(ctx, nil, []mlir.Type{i32})
	if err != nil {
		t.Fatalf("FunctionType() error = %v", err)
	}

	mod, err := mlir.CreateEmptyModule(loc)
	if err != nil {
		t.Fatalf("CreateEmptyModule() error = %v", err)
	}
	defer func() {
		if err := mod.Close(); err != nil {
			t.Fatalf("mod.Close() error = %v", err)
		}
	}()

	calleeOp, err := mod.Body().AppendOwnedOperation(buildConstantFunction(t, ctx, loc, "callee", fnType, i32, 7))
	if err != nil {
		t.Fatalf("AppendOwnedOperation(callee) error = %v", err)
	}
	callerOp, err := mod.Body().AppendOwnedOperation(buildCallerFunction(t, ctx, loc, "caller", fnType, i32, "callee"))
	if err != nil {
		t.Fatalf("AppendOwnedOperation(caller) error = %v", err)
	}

	if got := mlir.SymbolAttributeName(); got != "sym_name" {
		t.Fatalf("SymbolAttributeName() = %q, want %q", got, "sym_name")
	}
	if got := mlir.VisibilityAttributeName(); got != "sym_visibility" {
		t.Fatalf("VisibilityAttributeName() = %q, want %q", got, "sym_visibility")
	}
	if name, ok := calleeOp.SymbolName(); !ok || name != "callee" {
		t.Fatalf("calleeOp.SymbolName() = (%q, %t), want (%q, true)", name, ok, "callee")
	}

	privateAttr, err := mlir.StringAttribute(ctx, "private")
	if err != nil {
		t.Fatalf("StringAttribute(private) error = %v", err)
	}
	if err := callerOp.SetAttribute(mlir.VisibilityAttributeName(), privateAttr); err != nil {
		t.Fatalf("callerOp.SetAttribute(visibility) error = %v", err)
	}
	if visibility, ok := callerOp.Visibility(); !ok || visibility != "private" {
		t.Fatalf("callerOp.Visibility() = (%q, %t), want (%q, true)", visibility, ok, "private")
	}

	symtab, err := mod.SymbolTable()
	if err != nil {
		t.Fatalf("mod.SymbolTable() error = %v", err)
	}
	defer func() {
		if err := symtab.Close(); err != nil {
			t.Fatalf("symtab.Close() error = %v", err)
		}
	}()

	lookup, err := symtab.Lookup("callee")
	if err != nil {
		t.Fatalf("symtab.Lookup(callee) error = %v", err)
	}
	if lookup.IsNull() {
		t.Fatalf("symtab.Lookup(callee) returned null operation")
	}
	if name, ok := lookup.SymbolName(); !ok || name != "callee" {
		t.Fatalf("lookup.SymbolName() = (%q, %t), want (%q, true)", name, ok, "callee")
	}

	duplicateOp, err := mod.Body().AppendOwnedOperation(buildConstantFunction(t, ctx, loc, "callee", fnType, i32, 42))
	if err != nil {
		t.Fatalf("AppendOwnedOperation(duplicate) error = %v", err)
	}
	renamedAttr, err := symtab.Insert(duplicateOp)
	if err != nil {
		t.Fatalf("symtab.Insert(duplicate) error = %v", err)
	}
	renamed, err := renamedAttr.StringValue()
	if err != nil {
		t.Fatalf("renamedAttr.StringValue() error = %v", err)
	}
	if renamed == "" || renamed == "callee" || !strings.HasPrefix(renamed, "callee") {
		t.Fatalf("renamed symbol = %q, want non-empty unique callee-prefixed name", renamed)
	}
	if duplicateName, ok := duplicateOp.SymbolName(); !ok || duplicateName != renamed {
		t.Fatalf("duplicateOp.SymbolName() = (%q, %t), want (%q, true)", duplicateName, ok, renamed)
	}

	if err := mlir.ReplaceAllSymbolUses("callee", renamed, mod.Operation()); err != nil {
		t.Fatalf("ReplaceAllSymbolUses(callee -> %s) error = %v", renamed, err)
	}
	callAttr := callerOp.Regions()[0].Blocks()[0].Operations()[0].Attribute("callee")
	callTarget, err := callAttr.FlatSymbolValue()
	if err != nil {
		t.Fatalf("callAttr.FlatSymbolValue() error = %v", err)
	}
	if callTarget != renamed {
		t.Fatalf("call target = %q, want %q", callTarget, renamed)
	}

	var walked []string
	if err := mlir.WalkSymbolTables(mod.Operation(), true, func(op mlir.Operation, allVisible bool) {
		walked = append(walked, op.Name())
		if !allVisible {
			t.Fatalf("WalkSymbolTables callback got allVisible=false, want true")
		}
	}); err != nil {
		t.Fatalf("WalkSymbolTables() error = %v", err)
	}
	if len(walked) != 1 || walked[0] != "builtin.module" {
		t.Fatalf("WalkSymbolTables() visited %v, want [builtin.module]", walked)
	}

	if err := mod.Verify(); err != nil {
		t.Fatalf("mod.Verify() error = %v", err)
	}
	runFileCheck(t, mod.String(), "symbol_table_module.mlir")

	if err := mlir.ReplaceAllSymbolUses(renamed, "callee", mod.Operation()); err != nil {
		t.Fatalf("ReplaceAllSymbolUses(%s -> callee) error = %v", renamed, err)
	}
	if err := symtab.Erase(duplicateOp); err != nil {
		t.Fatalf("symtab.Erase(duplicate) error = %v", err)
	}
	erased, err := symtab.Lookup(renamed)
	if err != nil {
		t.Fatalf("symtab.Lookup(%s) error = %v", renamed, err)
	}
	if !erased.IsNull() {
		t.Fatalf("symtab.Lookup(%s) returned non-null operation after erase", renamed)
	}
	if err := mod.Verify(); err != nil {
		t.Fatalf("mod.Verify() after erase error = %v", err)
	}
}
