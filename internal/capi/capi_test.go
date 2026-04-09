//go:build cgo

package capi

import "testing"

func TestNullWrappers(t *testing.T) {
	if !ContextIsNull(NullContext()) {
		t.Fatalf("ContextIsNull(NullContext()) = false, want true")
	}
	if !ModuleIsNull(NullModule()) {
		t.Fatalf("ModuleIsNull(NullModule()) = false, want true")
	}
	if !LocationIsNull(NullLocation()) {
		t.Fatalf("LocationIsNull(NullLocation()) = false, want true")
	}
	if !ValueIsNull(NullValue()) {
		t.Fatalf("ValueIsNull(NullValue()) = false, want true")
	}
	if !TypeIsNull(NullType()) {
		t.Fatalf("TypeIsNull(NullType()) = false, want true")
	}
}

func TestWrapperCoverage(t *testing.T) {
	ctx := ContextCreate()
	if ContextIsNull(ctx) {
		t.Fatalf("ContextCreate() returned null context")
	}
	defer ContextDestroy(ctx)

	ContextSetAllowUnregisteredDialects(ctx, true)

	registry := DialectRegistryCreate()
	if DialectRegistryIsNull(registry) {
		t.Fatalf("DialectRegistryCreate() returned null registry")
	}
	RegisterAllDialects(registry)
	ContextAppendDialectRegistry(ctx, registry)
	ContextLoadAllAvailableDialects(ctx)
	if got := ContextGetNumLoadedDialects(ctx); got <= 0 {
		t.Fatalf("ContextGetNumLoadedDialects() = %d, want > 0", got)
	}
	DialectRegistryDestroy(registry)

	unknown := LocationUnknownGet(ctx)
	if LocationIsNull(unknown) {
		t.Fatalf("LocationUnknownGet() returned null location")
	}
	if got := LocationToString(unknown); got == "" {
		t.Fatalf("LocationToString(unknown) = empty, want non-empty")
	}

	fileLoc := LocationFileLineColGet(ctx, "test.mlir", 7, 11)
	if got := LocationToString(fileLoc); got != "loc(\"test.mlir\":7:11)" {
		t.Fatalf("LocationToString(fileLoc) = %q, want %q", got, "loc(\"test.mlir\":7:11)")
	}

	empty := ModuleCreateEmpty(fileLoc)
	if ModuleIsNull(empty) {
		t.Fatalf("ModuleCreateEmpty() returned null module")
	}
	if BlockIsNull(ModuleGetBody(empty)) {
		t.Fatalf("ModuleGetBody(empty) returned null block")
	}
	ModuleDestroy(empty)

	const input = `module {
  func.func @add(%arg0: i32) -> i32 {
    %0 = arith.addi %arg0, %arg0 : i32
    return %0 : i32
  }
}
`

	mod := ModuleCreateParse(ctx, input)
	if ModuleIsNull(mod) {
		t.Fatalf("ModuleCreateParse() returned null module")
	}
	defer ModuleDestroy(mod)

	moduleOp := ModuleGetOperation(mod)
	if OperationIsNull(moduleOp) {
		t.Fatalf("ModuleGetOperation() returned null operation")
	}
	if !OperationVerify(moduleOp) {
		t.Fatalf("OperationVerify(moduleOp) = false, want true")
	}
	if got := OperationToString(moduleOp); got == "" {
		t.Fatalf("OperationToString(moduleOp) = empty, want non-empty")
	}
	if got := IdentifierToString(OperationGetName(moduleOp)); got != "builtin.module" {
		t.Fatalf("module op name = %q, want %q", got, "builtin.module")
	}
	if got := LocationToString(OperationGetLocation(moduleOp)); got == "" {
		t.Fatalf("module op location = empty, want non-empty")
	}
	if got := OperationGetNumRegions(moduleOp); got != 1 {
		t.Fatalf("OperationGetNumRegions(moduleOp) = %d, want 1", got)
	}

	region := OperationGetRegion(moduleOp, 0)
	if RegionIsNull(region) {
		t.Fatalf("OperationGetRegion(moduleOp, 0) returned null region")
	}
	if !RegionIsNull(RegionGetNextInOperation(region)) {
		t.Fatalf("RegionGetNextInOperation(last) should be null")
	}

	moduleBlock := RegionGetFirstBlock(region)
	if BlockIsNull(moduleBlock) {
		t.Fatalf("RegionGetFirstBlock(module region) returned null block")
	}
	if got := BlockToString(moduleBlock); got == "" {
		t.Fatalf("BlockToString(module block) = empty, want non-empty")
	}
	if !BlockIsNull(BlockGetNextInRegion(moduleBlock)) {
		t.Fatalf("BlockGetNextInRegion(last) should be null")
	}

	funcOp := BlockGetFirstOperation(moduleBlock)
	if OperationIsNull(funcOp) {
		t.Fatalf("BlockGetFirstOperation(module block) returned null op")
	}
	if got := IdentifierToString(OperationGetName(funcOp)); got != "func.func" {
		t.Fatalf("func op name = %q, want %q", got, "func.func")
	}

	funcRegion := OperationGetRegion(funcOp, 0)
	funcBlock := RegionGetFirstBlock(funcRegion)
	if got := BlockGetNumArguments(funcBlock); got != 1 {
		t.Fatalf("BlockGetNumArguments(func block) = %d, want 1", got)
	}
	arg := BlockGetArgument(funcBlock, 0)
	if ValueIsNull(arg) {
		t.Fatalf("BlockGetArgument() returned null value")
	}
	if !ValueIsBlockArgument(arg) {
		t.Fatalf("ValueIsBlockArgument(arg) = false, want true")
	}
	if got := ValueToString(arg); got == "" {
		t.Fatalf("ValueToString(arg) = empty, want non-empty")
	}

	argType := ValueGetType(arg)
	if TypeIsNull(argType) {
		t.Fatalf("ValueGetType(arg) returned null type")
	}
	if got := TypeToString(argType); got != "i32" {
		t.Fatalf("TypeToString(argType) = %q, want %q", got, "i32")
	}

	parsedType := TypeParseGet(ctx, "i32")
	if TypeIsNull(parsedType) {
		t.Fatalf("TypeParseGet() returned null type")
	}
	if !TypeEqual(argType, parsedType) {
		t.Fatalf("TypeEqual(argType, parsedType) = false, want true")
	}

	addi := BlockGetFirstOperation(funcBlock)
	if got := IdentifierToString(OperationGetName(addi)); got != "arith.addi" {
		t.Fatalf("addi name = %q, want %q", got, "arith.addi")
	}
	result := OperationGetResult(addi, 0)
	if ValueIsNull(result) {
		t.Fatalf("OperationGetResult(addi, 0) returned null value")
	}
	if !ValueIsOpResult(result) {
		t.Fatalf("ValueIsOpResult(result) = false, want true")
	}
	if got := OperationGetNumResults(addi); got != 1 {
		t.Fatalf("OperationGetNumResults(addi) = %d, want 1", got)
	}
	if got := TypeToString(ValueGetType(result)); got != "i32" {
		t.Fatalf("result type = %q, want %q", got, "i32")
	}
	if got := ValueToString(result); got == "" {
		t.Fatalf("ValueToString(result) = empty, want non-empty")
	}

	ret := OperationGetNextInBlock(addi)
	if OperationIsNull(ret) {
		t.Fatalf("OperationGetNextInBlock(addi) returned null op")
	}
	if got := IdentifierToString(OperationGetName(ret)); got != "func.return" {
		t.Fatalf("return name = %q, want %q", got, "func.return")
	}
	if !OperationIsNull(OperationGetNextInBlock(ret)) {
		t.Fatalf("OperationGetNextInBlock(last) should be null")
	}
}
