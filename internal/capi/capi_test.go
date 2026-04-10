//go:build cgo

package capi

import (
	"strings"
	"testing"
	"unsafe"
)

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
	if !AttributeIsNull(NullAttribute()) {
		t.Fatalf("AttributeIsNull(NullAttribute()) = false, want true")
	}
	if !SymbolTableIsNull(NullSymbolTable()) {
		t.Fatalf("SymbolTableIsNull(NullSymbolTable()) = false, want true")
	}
	if !PassManagerIsNull(PassManager{}) {
		t.Fatalf("PassManagerIsNull(PassManager{}) = false, want true")
	}
	if !ExecutionEngineIsNull(ExecutionEngine{}) {
		t.Fatalf("ExecutionEngineIsNull(ExecutionEngine{}) = false, want true")
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

	attr := AttributeParseGet(ctx, "42 : i32")
	if AttributeIsNull(attr) {
		t.Fatalf("AttributeParseGet() returned null attribute")
	}
	if got := AttributeToString(attr); got != "42 : i32" {
		t.Fatalf("AttributeToString(attr) = %q, want %q", got, "42 : i32")
	}

	id := IdentifierGet(ctx, "value")
	if got := IdentifierToString(id); got != "value" {
		t.Fatalf("IdentifierToString(id) = %q, want %q", got, "value")
	}

	named := NamedAttributeGet(id, attr)
	if got := IdentifierToString(NamedAttributeName(named)); got != "value" {
		t.Fatalf("NamedAttributeName(named) = %q, want %q", got, "value")
	}
	if got := AttributeToString(NamedAttributeAttribute(named)); got != "42 : i32" {
		t.Fatalf("NamedAttributeAttribute(named) = %q, want %q", got, "42 : i32")
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

func TestConstructionPassAndExecutionCoverage(t *testing.T) {
	ctx := ContextCreate()
	if ContextIsNull(ctx) {
		t.Fatalf("ContextCreate() returned null context")
	}
	defer ContextDestroy(ctx)

	registry := DialectRegistryCreate()
	RegisterAllDialects(registry)
	ContextAppendDialectRegistry(ctx, registry)
	ContextLoadAllAvailableDialects(ctx)
	DialectRegistryDestroy(registry)

	RegisterAllPasses()

	const input = `module {
  func.func @add(%arg0: i32) -> i32 attributes {llvm.emit_c_interface} {
    return %arg0 : i32
  }
}
`

	mod := ModuleCreateParse(ctx, input)
	if ModuleIsNull(mod) {
		t.Fatalf("ModuleCreateParse() returned null module")
	}
	defer ModuleDestroy(mod)

	moduleBlock := ModuleGetBody(mod)
	funcOp := BlockGetFirstOperation(moduleBlock)
	funcBlock := RegionGetFirstBlock(OperationGetRegion(funcOp, 0))
	arg0 := BlockGetArgument(funcBlock, 0)
	ret := BlockGetFirstOperation(funcBlock)

	i32 := TypeParseGet(ctx, "i32")
	if TypeIsNull(i32) {
		t.Fatalf("TypeParseGet(i32) returned null type")
	}

	valueAttr := AttributeParseGet(ctx, "7 : i32")
	valueID := IdentifierGet(ctx, "value")
	constant := OperationCreate(
		"arith.constant",
		LocationUnknownGet(ctx),
		[]Type{i32},
		nil,
		nil,
		nil,
		[]NamedAttribute{NamedAttributeGet(valueID, valueAttr)},
		false,
	)
	if OperationIsNull(constant) {
		t.Fatalf("OperationCreate(arith.constant) returned null operation")
	}

	BlockInsertOwnedOperationBefore(funcBlock, ret, constant)
	constResult := OperationGetResult(constant, 0)
	if ValueIsNull(constResult) {
		t.Fatalf("OperationGetResult(constant, 0) returned null value")
	}
	ValueReplaceAllUsesOfWith(arg0, constResult)
	if !OperationVerify(ModuleGetOperation(mod)) {
		t.Fatalf("OperationVerify(module after rewrite) = false, want true")
	}

	argBlock := BlockCreate(nil, nil)
	if BlockIsNull(argBlock) {
		t.Fatalf("BlockCreate() returned null block")
	}
	addedArg := BlockAddArgument(argBlock, i32, LocationUnknownGet(ctx))
	if ValueIsNull(addedArg) || !ValueIsBlockArgument(addedArg) {
		t.Fatalf("BlockAddArgument() returned invalid block argument")
	}
	BlockDestroy(argBlock)

	detachedBlock := BlockCreate(nil, nil)
	if BlockIsNull(detachedBlock) {
		t.Fatalf("BlockCreate() returned null block")
	}
	detachedRegion := RegionCreate()
	if RegionIsNull(detachedRegion) {
		t.Fatalf("RegionCreate() returned null region")
	}
	RegionAppendOwnedBlock(detachedRegion, detachedBlock)
	moduleOp := OperationCreate("builtin.module", LocationUnknownGet(ctx), nil, nil, []Region{detachedRegion}, nil, nil, false)
	if OperationIsNull(moduleOp) {
		t.Fatalf("OperationCreate(builtin.module) returned null operation")
	}
	if !OperationVerify(moduleOp) {
		t.Fatalf("OperationVerify(detached module) = false, want true")
	}
	OperationDestroy(moduleOp)

	pm := PassManagerCreate(ctx)
	if PassManagerIsNull(pm) {
		t.Fatalf("PassManagerCreate() returned null pass manager")
	}
	defer PassManagerDestroy(pm)

	ok, msg := ParsePassPipeline(PassManagerGetAsOpPassManager(pm), "builtin.module(convert-arith-to-llvm,convert-func-to-llvm,reconcile-unrealized-casts)")
	if !ok {
		t.Fatalf("ParsePassPipeline() failed: %s", msg)
	}
	if got := PrintPassPipeline(PassManagerGetAsOpPassManager(pm)); got == "" {
		t.Fatalf("PrintPassPipeline() = empty, want non-empty")
	}
	if !PassManagerRunOnOp(pm, ModuleGetOperation(mod)) {
		t.Fatalf("PassManagerRunOnOp() = false, want true")
	}

	RegisterAllLLVMTranslations(ctx)
	engine := ExecutionEngineCreate(mod, 0)
	if ExecutionEngineIsNull(engine) {
		t.Fatalf("ExecutionEngineCreate() returned null engine")
	}
	defer ExecutionEngineDestroy(engine)

	arg := int32(21)
	var result int32
	args := []unsafe.Pointer{unsafe.Pointer(&arg), unsafe.Pointer(&result)}
	if !ExecutionEngineInvokePacked(engine, "add", args) {
		t.Fatalf("ExecutionEngineInvokePacked() = false, want true")
	}
	if result != 7 {
		t.Fatalf("result = %d, want 7", result)
	}
}

func TestBuiltinAndDiagnosticCoverage(t *testing.T) {
	ctx := ContextCreate()
	if ContextIsNull(ctx) {
		t.Fatalf("ContextCreate() returned null context")
	}
	defer ContextDestroy(ctx)

	i32 := IntegerTypeGet(ctx, 32)
	if TypeIsNull(i32) || !TypeIsAInteger(i32) || IntegerTypeGetWidth(i32) != 32 {
		t.Fatalf("IntegerTypeGet() returned unexpected type")
	}
	if !IntegerTypeIsSignless(i32) || IntegerTypeIsSigned(i32) || IntegerTypeIsUnsigned(i32) {
		t.Fatalf("integer signedness predicates = unexpected")
	}

	si16 := IntegerTypeSignedGet(ctx, 16)
	if TypeIsNull(si16) || !IntegerTypeIsSigned(si16) {
		t.Fatalf("IntegerTypeSignedGet() returned unexpected type")
	}

	ui8 := IntegerTypeUnsignedGet(ctx, 8)
	if TypeIsNull(ui8) || !IntegerTypeIsUnsigned(ui8) {
		t.Fatalf("IntegerTypeUnsignedGet() returned unexpected type")
	}

	index := IndexTypeGet(ctx)
	if TypeIsNull(index) || !TypeIsAIndex(index) {
		t.Fatalf("IndexTypeGet() returned unexpected type")
	}

	f32 := F32TypeGet(ctx)
	if TypeIsNull(f32) || !TypeIsAFloat(f32) || !TypeIsAF32(f32) || FloatTypeGetWidth(f32) != 32 {
		t.Fatalf("F32TypeGet() returned unexpected type")
	}

	f64 := F64TypeGet(ctx)
	if TypeIsNull(f64) || !TypeIsAF64(f64) || FloatTypeGetWidth(f64) != 64 {
		t.Fatalf("F64TypeGet() returned unexpected type")
	}

	none := NoneTypeGet(ctx)
	if TypeIsNull(none) || !TypeIsANone(none) {
		t.Fatalf("NoneTypeGet() returned unexpected type")
	}

	tensor := RankedTensorTypeGet([]int64{4, ShapedTypeGetDynamicSize()}, i32, NullAttribute())
	if TypeIsNull(tensor) || !TypeIsARankedTensor(tensor) {
		t.Fatalf("RankedTensorTypeGet() returned unexpected type")
	}
	if !TypeIsAShaped(tensor) || !ShapedTypeHasRank(tensor) || ShapedTypeGetRank(tensor) != 2 {
		t.Fatalf("shaped tensor predicates = unexpected")
	}
	if !TypeEqual(ShapedTypeGetElementType(tensor), i32) {
		t.Fatalf("ShapedTypeGetElementType(tensor) returned unexpected type")
	}
	if ShapedTypeGetDimSize(tensor, 0) != 4 || !ShapedTypeIsDynamicDim(tensor, 1) || ShapedTypeGetDimSize(tensor, 1) != ShapedTypeGetDynamicSize() {
		t.Fatalf("tensor shape APIs returned unexpected values")
	}
	if !AttributeIsNull(RankedTensorTypeGetEncoding(tensor)) {
		t.Fatalf("RankedTensorTypeGetEncoding() = non-null, want null")
	}

	memref := MemRefTypeContiguousGet(i32, []int64{ShapedTypeGetDynamicSize()}, NullAttribute())
	if TypeIsNull(memref) || !TypeIsAMemRef(memref) || !TypeIsAShaped(memref) {
		t.Fatalf("MemRefTypeContiguousGet() returned unexpected type")
	}
	if !TypeEqual(ShapedTypeGetElementType(memref), i32) {
		t.Fatalf("ShapedTypeGetElementType(memref) returned unexpected type")
	}
	if AttributeIsNull(MemRefTypeGetLayout(memref)) || AttributeToString(MemRefTypeGetLayout(memref)) == "" {
		t.Fatalf("memref layout should be materialized")
	}
	if !AttributeIsNull(MemRefTypeGetMemorySpace(memref)) {
		t.Fatalf("memref memory space should be null")
	}

	unrankedMemRef := UnrankedMemRefTypeGet(i32, NullAttribute())
	if TypeIsNull(unrankedMemRef) || !TypeIsAUnrankedMemRef(unrankedMemRef) {
		t.Fatalf("UnrankedMemRefTypeGet() returned unexpected type")
	}
	if ShapedTypeHasRank(unrankedMemRef) {
		t.Fatalf("ShapedTypeHasRank(unrankedMemRef) = true, want false")
	}

	fnType := FunctionTypeGet(ctx, []Type{i32, i32}, []Type{i32})
	if TypeIsNull(fnType) || !TypeIsAFunction(fnType) {
		t.Fatalf("FunctionTypeGet() returned unexpected type")
	}
	if got := FunctionTypeGetNumInputs(fnType); got != 2 {
		t.Fatalf("FunctionTypeGetNumInputs() = %d, want 2", got)
	}
	if got := FunctionTypeGetNumResults(fnType); got != 1 {
		t.Fatalf("FunctionTypeGetNumResults() = %d, want 1", got)
	}
	if !TypeEqual(FunctionTypeGetInput(fnType, 0), i32) || !TypeEqual(FunctionTypeGetResult(fnType, 0), i32) {
		t.Fatalf("FunctionTypeGetInput/Result() returned unexpected type")
	}

	intAttr := IntegerAttrGet(i32, 9)
	if AttributeIsNull(intAttr) || !AttributeIsAInteger(intAttr) || IntegerAttrGetValueInt(intAttr) != 9 {
		t.Fatalf("IntegerAttrGet() returned unexpected attribute")
	}

	boolAttr := BoolAttrGet(ctx, true)
	if AttributeIsNull(boolAttr) || !AttributeIsABool(boolAttr) || !BoolAttrGetValue(boolAttr) {
		t.Fatalf("BoolAttrGet() returned unexpected attribute")
	}

	stringAttr := StringAttrGet(ctx, "hello")
	if AttributeIsNull(stringAttr) || !AttributeIsAString(stringAttr) || StringAttrGetValue(stringAttr) != "hello" {
		t.Fatalf("StringAttrGet() returned unexpected attribute")
	}

	symbolAttr := FlatSymbolRefAttrGet(ctx, "callee")
	if AttributeIsNull(symbolAttr) || !AttributeIsAFlatSymbolRef(symbolAttr) || FlatSymbolRefAttrGetValue(symbolAttr) != "callee" {
		t.Fatalf("FlatSymbolRefAttrGet() returned unexpected attribute")
	}

	typeAttr := TypeAttrGet(fnType)
	if AttributeIsNull(typeAttr) || !AttributeIsAType(typeAttr) || !TypeEqual(TypeAttrGetValue(typeAttr), fnType) {
		t.Fatalf("TypeAttrGet() returned unexpected attribute")
	}

	unitAttr := UnitAttrGet(ctx)
	if AttributeIsNull(unitAttr) || !AttributeIsAUnit(unitAttr) {
		t.Fatalf("UnitAttrGet() returned unexpected attribute")
	}

	type capturedDiagnostic struct {
		severity DiagnosticSeverity
		text     string
		location string
		numNotes int
	}

	diagnostics := make([]capturedDiagnostic, 0, 1)
	id := ContextAttachDiagnosticCallback(ctx, func(diag Diagnostic) {
		diagnostics = append(diagnostics, capturedDiagnostic{
			severity: DiagnosticGetSeverity(diag),
			text:     DiagnosticToString(diag),
			location: LocationToString(DiagnosticGetLocation(diag)),
			numNotes: DiagnosticGetNumNotes(diag),
		})
	})
	EmitError(LocationUnknownGet(ctx), "boom")
	ContextDetachDiagnosticHandler(ctx, id)
	if len(diagnostics) != 1 {
		t.Fatalf("len(diagnostics) = %d, want 1", len(diagnostics))
	}
	if got := diagnostics[0].severity; got != DiagnosticSeverityError {
		t.Fatalf("DiagnosticGetSeverity() = %v, want %v", got, DiagnosticSeverityError)
	}
	if got := diagnostics[0].text; !strings.Contains(got, "boom") {
		t.Fatalf("DiagnosticToString() = %q, want substring %q", got, "boom")
	}
	if got := diagnostics[0].location; got == "" {
		t.Fatalf("DiagnosticGetLocation() = empty, want non-empty")
	}
	if got := diagnostics[0].numNotes; got != 0 {
		t.Fatalf("DiagnosticGetNumNotes() = %d, want 0", got)
	}
}

func TestIRQueryMutationCoverage(t *testing.T) {
	ctx := ContextCreate()
	if ContextIsNull(ctx) {
		t.Fatalf("ContextCreate() returned null context")
	}
	defer ContextDestroy(ctx)

	registry := DialectRegistryCreate()
	RegisterAllDialects(registry)
	ContextAppendDialectRegistry(ctx, registry)
	ContextLoadAllAvailableDialects(ctx)
	DialectRegistryDestroy(registry)

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

	mod := ModuleCreateParse(ctx, input)
	if ModuleIsNull(mod) {
		t.Fatalf("ModuleCreateParse() returned null module")
	}
	defer ModuleDestroy(mod)

	moduleOp := ModuleGetOperation(mod)
	moduleBlock := ModuleGetBody(mod)
	funcOp := BlockGetFirstOperation(moduleBlock)
	if OperationIsNull(OperationGetParentOperation(funcOp)) || IdentifierToString(OperationGetName(OperationGetParentOperation(funcOp))) != "builtin.module" {
		t.Fatalf("OperationGetParentOperation(funcOp) returned unexpected op")
	}
	if BlockIsNull(OperationGetBlock(funcOp)) {
		t.Fatalf("OperationGetBlock(funcOp) returned null block")
	}

	entry := RegionGetFirstBlock(OperationGetRegion(funcOp, 0))
	if RegionIsNull(BlockGetParentRegion(entry)) || OperationIsNull(BlockGetParentOperation(entry)) {
		t.Fatalf("block parent queries returned null")
	}
	condBr := BlockGetTerminator(entry)
	if OperationIsNull(condBr) || IdentifierToString(OperationGetName(condBr)) != "cf.cond_br" {
		t.Fatalf("BlockGetTerminator(entry) returned unexpected op")
	}

	if got := OperationGetNumOperands(condBr); got != 1 {
		t.Fatalf("OperationGetNumOperands(condBr) = %d, want 1", got)
	}
	arg0 := BlockGetArgument(entry, 0)
	if !ValueEqual(OperationGetOperand(condBr, 0), arg0) {
		t.Fatalf("OperationGetOperand(condBr, 0) should equal arg0")
	}
	OperationSetOperand(condBr, 0, arg0)
	OperationSetOperands(condBr, []Value{arg0})

	if got := OperationGetNumSuccessors(condBr); got != 2 {
		t.Fatalf("OperationGetNumSuccessors(condBr) = %d, want 2", got)
	}
	succ0 := OperationGetSuccessor(condBr, 0)
	if BlockIsNull(succ0) {
		t.Fatalf("OperationGetSuccessor(condBr, 0) returned null block")
	}
	OperationSetSuccessor(condBr, 0, succ0)

	if got := OperationGetNumAttributes(funcOp); got < 2 {
		t.Fatalf("OperationGetNumAttributes(funcOp) = %d, want >= 2", got)
	}
	firstAttr := OperationGetAttribute(funcOp, 0)
	if IdentifierToString(NamedAttributeName(firstAttr)) == "" {
		t.Fatalf("OperationGetAttribute(funcOp, 0) returned unnamed attribute")
	}
	if AttributeIsNull(OperationGetAttributeByName(funcOp, "sym_name")) {
		t.Fatalf("OperationGetAttributeByName(funcOp, sym_name) returned null")
	}
	tag := StringAttrGet(ctx, "x")
	OperationSetAttributeByName(funcOp, "test.tag", tag)
	if got := StringAttrGetValue(OperationGetAttributeByName(funcOp, "test.tag")); got != "x" {
		t.Fatalf("OperationGetAttributeByName(funcOp, test.tag) = %q, want %q", got, "x")
	}
	if !OperationRemoveAttributeByName(funcOp, "test.tag") {
		t.Fatalf("OperationRemoveAttributeByName(funcOp, test.tag) = false, want true")
	}

	if BlockIsNull(BlockArgumentGetOwner(arg0)) || BlockArgumentGetArgNumber(arg0) != 0 {
		t.Fatalf("block argument owner queries returned unexpected values")
	}
	constant := BlockGetFirstOperation(entry)
	constResult := OperationGetResult(constant, 0)
	if OperationIsNull(OpResultGetOwner(constResult)) || OpResultGetResultNumber(constResult) != 0 {
		t.Fatalf("op result owner queries returned unexpected values")
	}

	i32 := ValueGetType(constResult)
	i64 := IntegerTypeGet(ctx, 64)
	tmpBlock := BlockCreate([]Type{i32}, []Location{LocationUnknownGet(ctx)})
	if BlockIsNull(tmpBlock) {
		t.Fatalf("BlockCreate() returned null block")
	}
	tmpArg := BlockGetArgument(tmpBlock, 0)
	BlockArgumentSetType(tmpArg, i64)
	if !TypeEqual(ValueGetType(tmpArg), i64) {
		t.Fatalf("BlockArgumentSetType() did not update type")
	}
	ValueSetType(tmpArg, i32)
	if !TypeEqual(ValueGetType(tmpArg), i32) {
		t.Fatalf("ValueSetType() did not update type")
	}
	BlockDestroy(tmpBlock)

	if !OperationVerify(moduleOp) {
		t.Fatalf("OperationVerify(moduleOp) = false, want true")
	}
}

func TestSymbolTableCoverage(t *testing.T) {
	ctx := ContextCreate()
	if ContextIsNull(ctx) {
		t.Fatalf("ContextCreate() returned null context")
	}
	defer ContextDestroy(ctx)

	registry := DialectRegistryCreate()
	RegisterAllDialects(registry)
	ContextAppendDialectRegistry(ctx, registry)
	ContextLoadAllAvailableDialects(ctx)
	DialectRegistryDestroy(registry)

	const input = `module {
  func.func @callee() -> i32 {
    %0 = arith.constant 7 : i32
    return %0 : i32
  }
  func.func @caller() -> i32 {
    %0 = func.call @callee() : () -> i32
    return %0 : i32
  }
}
`

	mod := ModuleCreateParse(ctx, input)
	if ModuleIsNull(mod) {
		t.Fatalf("ModuleCreateParse() returned null module")
	}
	defer ModuleDestroy(mod)

	if got := SymbolTableGetSymbolAttributeName(); got != "sym_name" {
		t.Fatalf("SymbolTableGetSymbolAttributeName() = %q, want %q", got, "sym_name")
	}
	if got := SymbolTableGetVisibilityAttributeName(); got != "sym_visibility" {
		t.Fatalf("SymbolTableGetVisibilityAttributeName() = %q, want %q", got, "sym_visibility")
	}

	moduleOp := ModuleGetOperation(mod)
	moduleBlock := ModuleGetBody(mod)
	calleeOp := BlockGetFirstOperation(moduleBlock)
	callerOp := OperationGetNextInBlock(calleeOp)

	visibility := StringAttrGet(ctx, "private")
	OperationSetAttributeByName(callerOp, SymbolTableGetVisibilityAttributeName(), visibility)
	if got := StringAttrGetValue(OperationGetAttributeByName(callerOp, SymbolTableGetVisibilityAttributeName())); got != "private" {
		t.Fatalf("caller visibility = %q, want %q", got, "private")
	}

	symtab := SymbolTableCreate(moduleOp)
	if SymbolTableIsNull(symtab) {
		t.Fatalf("SymbolTableCreate(moduleOp) returned null symbol table")
	}
	defer SymbolTableDestroy(symtab)

	lookup := SymbolTableLookup(symtab, "callee")
	if OperationIsNull(lookup) {
		t.Fatalf("SymbolTableLookup(symtab, callee) returned null operation")
	}

	loc := LocationUnknownGet(ctx)
	i32 := IntegerTypeGet(ctx, 32)
	attr42 := IntegerAttrGet(i32, 42)
	valueID := IdentifierGet(ctx, "value")
	constant := OperationCreate(
		"arith.constant",
		loc,
		[]Type{i32},
		nil,
		nil,
		nil,
		[]NamedAttribute{NamedAttributeGet(valueID, attr42)},
		false,
	)
	if OperationIsNull(constant) {
		t.Fatalf("OperationCreate(arith.constant) returned null operation")
	}

	funcBlock := BlockCreate(nil, nil)
	if BlockIsNull(funcBlock) {
		t.Fatalf("BlockCreate() returned null block")
	}
	BlockAppendOwnedOperation(funcBlock, constant)

	ret := OperationCreate(
		"func.return",
		loc,
		nil,
		[]Value{OperationGetResult(constant, 0)},
		nil,
		nil,
		nil,
		false,
	)
	if OperationIsNull(ret) {
		t.Fatalf("OperationCreate(func.return) returned null operation")
	}
	BlockAppendOwnedOperation(funcBlock, ret)

	body := RegionCreate()
	if RegionIsNull(body) {
		t.Fatalf("RegionCreate() returned null region")
	}
	RegionAppendOwnedBlock(body, funcBlock)

	fnType := FunctionTypeGet(ctx, nil, []Type{i32})
	symNameID := IdentifierGet(ctx, SymbolTableGetSymbolAttributeName())
	functionTypeID := IdentifierGet(ctx, "function_type")
	duplicate := OperationCreate(
		"func.func",
		loc,
		nil,
		nil,
		[]Region{body},
		nil,
		[]NamedAttribute{
			NamedAttributeGet(symNameID, StringAttrGet(ctx, "callee")),
			NamedAttributeGet(functionTypeID, TypeAttrGet(fnType)),
		},
		false,
	)
	if OperationIsNull(duplicate) {
		t.Fatalf("OperationCreate(func.func) returned null operation")
	}
	BlockAppendOwnedOperation(moduleBlock, duplicate)

	renamedAttr := SymbolTableInsert(symtab, duplicate)
	if AttributeIsNull(renamedAttr) {
		t.Fatalf("SymbolTableInsert(symtab, duplicate) returned null attribute")
	}
	renamed := StringAttrGetValue(renamedAttr)
	if renamed == "" || renamed == "callee" || !strings.HasPrefix(renamed, "callee") {
		t.Fatalf("renamed symbol = %q, want non-empty unique callee-prefixed name", renamed)
	}

	if !SymbolTableReplaceAllSymbolUses("callee", renamed, moduleOp) {
		t.Fatalf("SymbolTableReplaceAllSymbolUses(callee -> %s) = false, want true", renamed)
	}
	callOp := BlockGetFirstOperation(RegionGetFirstBlock(OperationGetRegion(callerOp, 0)))
	if got := FlatSymbolRefAttrGetValue(OperationGetAttributeByName(callOp, "callee")); got != renamed {
		t.Fatalf("call target = %q, want %q", got, renamed)
	}

	var walked []string
	SymbolTableWalkSymbolTables(moduleOp, true, func(op Operation, allVisible bool) {
		if !allVisible {
			t.Fatalf("walk callback got allVisible=false, want true")
		}
		walked = append(walked, IdentifierToString(OperationGetName(op)))
	})
	if len(walked) != 1 || walked[0] != "builtin.module" {
		t.Fatalf("SymbolTableWalkSymbolTables() visited %v, want [builtin.module]", walked)
	}

	if !OperationVerify(moduleOp) {
		t.Fatalf("OperationVerify(moduleOp) after symbol rewrite = false, want true")
	}

	if !SymbolTableReplaceAllSymbolUses(renamed, "callee", moduleOp) {
		t.Fatalf("SymbolTableReplaceAllSymbolUses(%s -> callee) = false, want true", renamed)
	}
	SymbolTableErase(symtab, duplicate)
	if !OperationIsNull(SymbolTableLookup(symtab, renamed)) {
		t.Fatalf("SymbolTableLookup(symtab, %s) returned non-null op after erase", renamed)
	}
	if !OperationVerify(moduleOp) {
		t.Fatalf("OperationVerify(moduleOp) after erase = false, want true")
	}
}
