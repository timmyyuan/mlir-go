//go:build cgo

package mlir

import "testing"

func TestBuiltinTypesAndAttributes(t *testing.T) {
	ctx, err := NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer func() {
		if err := ctx.Close(); err != nil {
			t.Fatalf("ctx.Close() error = %v", err)
		}
	}()

	i32, err := SignlessIntegerType(ctx, 32)
	if err != nil {
		t.Fatalf("SignlessIntegerType() error = %v", err)
	}
	if !i32.IsInteger() || !i32.IsSignlessInteger() || i32.IntegerWidth() != 32 {
		t.Fatalf("i32 predicates = unexpected")
	}

	si16, err := SignedIntegerType(ctx, 16)
	if err != nil {
		t.Fatalf("SignedIntegerType() error = %v", err)
	}
	if !si16.IsSignedInteger() || si16.IntegerWidth() != 16 {
		t.Fatalf("si16 predicates = unexpected")
	}

	ui8, err := UnsignedIntegerType(ctx, 8)
	if err != nil {
		t.Fatalf("UnsignedIntegerType() error = %v", err)
	}
	if !ui8.IsUnsignedInteger() || ui8.IntegerWidth() != 8 {
		t.Fatalf("ui8 predicates = unexpected")
	}

	index, err := IndexType(ctx)
	if err != nil {
		t.Fatalf("IndexType() error = %v", err)
	}
	if !index.IsIndex() {
		t.Fatalf("index.IsIndex() = false, want true")
	}

	f32, err := F32Type(ctx)
	if err != nil {
		t.Fatalf("F32Type() error = %v", err)
	}
	if !f32.IsF32() || f32.FloatWidth() != 32 {
		t.Fatalf("f32 predicates = unexpected")
	}

	f64, err := F64Type(ctx)
	if err != nil {
		t.Fatalf("F64Type() error = %v", err)
	}
	if !f64.IsF64() || f64.FloatWidth() != 64 {
		t.Fatalf("f64 predicates = unexpected")
	}

	none, err := NoneType(ctx)
	if err != nil {
		t.Fatalf("NoneType() error = %v", err)
	}
	if !none.IsNone() {
		t.Fatalf("none.IsNone() = false, want true")
	}

	tensor, err := RankedTensorType([]int64{4, DynamicSize()}, i32)
	if err != nil {
		t.Fatalf("RankedTensorType() error = %v", err)
	}
	if !tensor.IsRankedTensor() {
		t.Fatalf("tensor.IsRankedTensor() = false, want true")
	}
	if !tensor.IsShaped() || !tensor.HasRank() || tensor.Rank() != 2 {
		t.Fatalf("tensor shaped predicates = unexpected")
	}
	if !tensor.ElementType().Equal(i32) {
		t.Fatalf("tensor.ElementType() = %s, want %s", tensor.ElementType(), i32)
	}
	if tensor.DimSize(0) != 4 || !tensor.IsDynamicDim(1) || tensor.DimSize(1) != DynamicSize() {
		t.Fatalf("tensor shape = %v, want [4 dynamic]", tensor.Shape())
	}
	if !tensor.TensorEncoding().IsNull() {
		t.Fatalf("tensor encoding = non-null, want null")
	}

	memref, err := MemRefType([]int64{DynamicSize()}, i32)
	if err != nil {
		t.Fatalf("MemRefType() error = %v", err)
	}
	if !memref.IsMemRef() || !memref.IsShaped() {
		t.Fatalf("memref predicates = unexpected")
	}
	if !memref.ElementType().Equal(i32) {
		t.Fatalf("memref.ElementType() = %s, want %s", memref.ElementType(), i32)
	}
	if memref.MemRefLayout().IsNull() || memref.MemRefLayout().String() == "" {
		t.Fatalf("memref layout should be materialized")
	}
	if !memref.MemRefMemorySpace().IsNull() {
		t.Fatalf("memref memory space should be null")
	}

	unrankedMemRef, err := UnrankedMemRefType(i32)
	if err != nil {
		t.Fatalf("UnrankedMemRefType() error = %v", err)
	}
	if !unrankedMemRef.IsUnrankedMemRef() || unrankedMemRef.HasRank() {
		t.Fatalf("unranked memref predicates = unexpected")
	}

	fnType, err := FunctionType(ctx, []Type{i32, i32}, []Type{i32})
	if err != nil {
		t.Fatalf("FunctionType() error = %v", err)
	}
	if !fnType.IsFunction() {
		t.Fatalf("fnType.IsFunction() = false, want true")
	}
	if got := len(fnType.FunctionInputs()); got != 2 {
		t.Fatalf("len(fnType.FunctionInputs()) = %d, want 2", got)
	}
	if got := len(fnType.FunctionResults()); got != 1 {
		t.Fatalf("len(fnType.FunctionResults()) = %d, want 1", got)
	}

	intAttr, err := IntegerAttribute(i32, 42)
	if err != nil {
		t.Fatalf("IntegerAttribute() error = %v", err)
	}
	if !intAttr.IsInteger() {
		t.Fatalf("intAttr.IsInteger() = false, want true")
	}
	if got, err := intAttr.Int64Value(); err != nil || got != 42 {
		t.Fatalf("intAttr.Int64Value() = (%d, %v), want (42, nil)", got, err)
	}

	boolAttr, err := BoolAttribute(ctx, true)
	if err != nil {
		t.Fatalf("BoolAttribute() error = %v", err)
	}
	if !boolAttr.IsBool() {
		t.Fatalf("boolAttr.IsBool() = false, want true")
	}
	if got, err := boolAttr.BoolValue(); err != nil || !got {
		t.Fatalf("boolAttr.BoolValue() = (%v, %v), want (true, nil)", got, err)
	}

	strAttr, err := StringAttribute(ctx, "hello")
	if err != nil {
		t.Fatalf("StringAttribute() error = %v", err)
	}
	if !strAttr.IsString() {
		t.Fatalf("strAttr.IsString() = false, want true")
	}
	if got, err := strAttr.StringValue(); err != nil || got != "hello" {
		t.Fatalf("strAttr.StringValue() = (%q, %v), want (%q, nil)", got, err, "hello")
	}

	symAttr, err := FlatSymbolRefAttribute(ctx, "callee")
	if err != nil {
		t.Fatalf("FlatSymbolRefAttribute() error = %v", err)
	}
	if !symAttr.IsFlatSymbolRef() {
		t.Fatalf("symAttr.IsFlatSymbolRef() = false, want true")
	}
	if got, err := symAttr.FlatSymbolValue(); err != nil || got != "callee" {
		t.Fatalf("symAttr.FlatSymbolValue() = (%q, %v), want (%q, nil)", got, err, "callee")
	}

	typeAttr, err := TypeAttribute(fnType)
	if err != nil {
		t.Fatalf("TypeAttribute() error = %v", err)
	}
	if !typeAttr.IsTypeAttribute() {
		t.Fatalf("typeAttr.IsTypeAttribute() = false, want true")
	}
	if got, err := typeAttr.TypeValue(); err != nil || !got.Equal(fnType) {
		t.Fatalf("typeAttr.TypeValue() = (%v, %v), want fnType", got, err)
	}

	unitAttr, err := UnitAttribute(ctx)
	if err != nil {
		t.Fatalf("UnitAttribute() error = %v", err)
	}
	if !unitAttr.IsUnit() {
		t.Fatalf("unitAttr.IsUnit() = false, want true")
	}
}

func TestCaptureDiagnostics(t *testing.T) {
	ctx, err := NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer func() {
		if err := ctx.Close(); err != nil {
			t.Fatalf("ctx.Close() error = %v", err)
		}
	}()

	i32, err := SignlessIntegerType(ctx, 32)
	if err != nil {
		t.Fatalf("SignlessIntegerType() error = %v", err)
	}
	loc, err := FileLineColLocation(ctx, "tensor.mlir", 1, 1)
	if err != nil {
		t.Fatalf("FileLineColLocation() error = %v", err)
	}

	diags, runErr := ctx.CaptureDiagnostics(func() error {
		_, err := CheckedRankedTensorType(loc, []int64{-2}, i32, Attribute{})
		return err
	})
	if runErr == nil {
		t.Fatalf("CaptureDiagnostics() runErr = nil, want non-nil")
	}
	if got := len(diags); got == 0 {
		t.Fatalf("len(diags) = 0, want > 0")
	}
	if diags[0].Severity != DiagnosticError {
		t.Fatalf("diags[0].Severity = %v, want %v", diags[0].Severity, DiagnosticError)
	}
	if diags[0].Location == "" || diags[0].Message == "" {
		t.Fatalf("diagnostic = %+v, want non-empty location and message", diags[0])
	}
}
