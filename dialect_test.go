//go:build cgo

package mlir_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	mlir "github.com/timmyyuan/mlir-go"
	arith "github.com/timmyyuan/mlir-go/dialect/arith"
	cf "github.com/timmyyuan/mlir-go/dialect/cf"
	funcdialect "github.com/timmyyuan/mlir-go/dialect/func"
	linalg "github.com/timmyyuan/mlir-go/dialect/linalg"
	memref "github.com/timmyyuan/mlir-go/dialect/memref"
	scf "github.com/timmyyuan/mlir-go/dialect/scf"
	tensor "github.com/timmyyuan/mlir-go/dialect/tensor"
	vector "github.com/timmyyuan/mlir-go/dialect/vector"
)

func fileCheckPath(t *testing.T) string {
	t.Helper()

	if p := os.Getenv("FILECHECK"); p != "" {
		return p
	}
	if p, err := exec.LookPath("FileCheck"); err == nil {
		return p
	}
	t.Skip("FileCheck not available")
	return ""
}

func runFileCheck(t *testing.T, input string, fixture string) {
	t.Helper()

	checkFile := filepath.Join("testdata", "filecheck", fixture)
	cmd := exec.Command(fileCheckPath(t), checkFile)
	cmd.Stdin = strings.NewReader(input)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("FileCheck failed for %s: %v\n%s\ninput:\n%s", fixture, err, out.String(), input)
	}
}

func TestDialectWrappersBuildModule(t *testing.T) {
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
	fnType, err := mlir.FunctionType(ctx, []mlir.Type{i32}, []mlir.Type{i32})
	if err != nil {
		t.Fatalf("FunctionType() error = %v", err)
	}

	body, err := mlir.NewOwnedRegion()
	if err != nil {
		t.Fatalf("NewOwnedRegion() error = %v", err)
	}
	entry, err := mlir.NewOwnedBlock([]mlir.Type{i32}, []mlir.Location{loc})
	if err != nil {
		t.Fatalf("NewOwnedBlock() error = %v", err)
	}
	entryBlock := entry.Block()
	arg0 := entryBlock.Arguments()[0]

	constAttr, err := mlir.IntegerAttribute(i32, 5)
	if err != nil {
		t.Fatalf("IntegerAttribute() error = %v", err)
	}
	constant, err := arith.Constant(ctx, loc, i32, constAttr)
	if err != nil {
		t.Fatalf("arith.Constant() error = %v", err)
	}
	constantOp, err := entry.AppendOperation(constant)
	if err != nil {
		t.Fatalf("entry.AppendOperation(constant) error = %v", err)
	}

	add, err := arith.AddI(loc, arg0, constantOp.Results()[0])
	if err != nil {
		t.Fatalf("arith.AddI() error = %v", err)
	}
	addOp, err := entry.AppendOperation(add)
	if err != nil {
		t.Fatalf("entry.AppendOperation(add) error = %v", err)
	}

	ret, err := funcdialect.Return(loc, addOp.Results()[0])
	if err != nil {
		t.Fatalf("func.Return() error = %v", err)
	}
	if _, err := entry.AppendOperation(ret); err != nil {
		t.Fatalf("entry.AppendOperation(return) error = %v", err)
	}

	if _, err := body.AppendBlock(entry); err != nil {
		t.Fatalf("body.AppendBlock() error = %v", err)
	}

	fn, err := funcdialect.Func(ctx, loc, "increment", fnType, body)
	if err != nil {
		t.Fatalf("func.Func() error = %v", err)
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

	fnOp, err := mod.Body().AppendOwnedOperation(fn)
	if err != nil {
		t.Fatalf("mod.Body().AppendOwnedOperation() error = %v", err)
	}

	argName, err := mlir.StringAttribute(ctx, "input")
	if err != nil {
		t.Fatalf("StringAttribute() error = %v", err)
	}
	if err := fnOp.SetFuncArgumentAttribute(0, "test.name", argName); err != nil {
		t.Fatalf("SetFuncArgumentAttribute() error = %v", err)
	}

	if err := mod.Verify(); err != nil {
		t.Fatalf("mod.Verify() error = %v", err)
	}
	runFileCheck(t, mod.String(), "dialect_wrappers_module.mlir")
}

func TestControlFlowAndMemRefWrappers(t *testing.T) {
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
	i1, err := mlir.SignlessIntegerType(ctx, 1)
	if err != nil {
		t.Fatalf("SignlessIntegerType(i1) error = %v", err)
	}
	i32, err := mlir.SignlessIntegerType(ctx, 32)
	if err != nil {
		t.Fatalf("SignlessIntegerType(i32) error = %v", err)
	}
	index, err := mlir.IndexType(ctx)
	if err != nil {
		t.Fatalf("IndexType() error = %v", err)
	}
	memrefType, err := mlir.MemRefType([]int64{mlir.DynamicSize()}, i32)
	if err != nil {
		t.Fatalf("MemRefType() error = %v", err)
	}
	fnType, err := mlir.FunctionType(ctx, []mlir.Type{i1, index}, []mlir.Type{i32})
	if err != nil {
		t.Fatalf("FunctionType() error = %v", err)
	}

	body, err := mlir.NewOwnedRegion()
	if err != nil {
		t.Fatalf("NewOwnedRegion() error = %v", err)
	}
	entry, err := mlir.NewOwnedBlock([]mlir.Type{i1, index}, []mlir.Location{loc, loc})
	if err != nil {
		t.Fatalf("NewOwnedBlock(entry) error = %v", err)
	}
	thenBlock, err := mlir.NewOwnedBlock(nil, nil)
	if err != nil {
		t.Fatalf("NewOwnedBlock(then) error = %v", err)
	}
	elseBlock, err := mlir.NewOwnedBlock(nil, nil)
	if err != nil {
		t.Fatalf("NewOwnedBlock(else) error = %v", err)
	}

	entryBlock := entry.Block()
	thenBorrowed := thenBlock.Block()
	elseBorrowed := elseBlock.Block()

	cond := entryBlock.Arguments()[0]
	n := entryBlock.Arguments()[1]

	c0Attr, err := mlir.IntegerAttribute(index, 0)
	if err != nil {
		t.Fatalf("IntegerAttribute(index, 0) error = %v", err)
	}
	c7Attr, err := mlir.IntegerAttribute(i32, 7)
	if err != nil {
		t.Fatalf("IntegerAttribute(i32, 7) error = %v", err)
	}

	c0, err := arith.Constant(ctx, loc, index, c0Attr)
	if err != nil {
		t.Fatalf("arith.Constant(c0) error = %v", err)
	}
	c0Op, err := entry.AppendOperation(c0)
	if err != nil {
		t.Fatalf("entry.AppendOperation(c0) error = %v", err)
	}

	c7, err := arith.Constant(ctx, loc, i32, c7Attr)
	if err != nil {
		t.Fatalf("arith.Constant(c7) error = %v", err)
	}
	c7Op, err := entry.AppendOperation(c7)
	if err != nil {
		t.Fatalf("entry.AppendOperation(c7) error = %v", err)
	}

	stack, err := memref.Alloca(ctx, loc, memrefType, n)
	if err != nil {
		t.Fatalf("memref.Alloca() error = %v", err)
	}
	stackOp, err := entry.AppendOperation(stack)
	if err != nil {
		t.Fatalf("entry.AppendOperation(alloca) error = %v", err)
	}

	store, err := memref.Store(loc, c7Op.Results()[0], stackOp.Results()[0], c0Op.Results()[0])
	if err != nil {
		t.Fatalf("memref.Store() error = %v", err)
	}
	if _, err := entry.AppendOperation(store); err != nil {
		t.Fatalf("entry.AppendOperation(store) error = %v", err)
	}

	load, err := memref.Load(loc, stackOp.Results()[0], c0Op.Results()[0])
	if err != nil {
		t.Fatalf("memref.Load() error = %v", err)
	}
	loadOp, err := entry.AppendOperation(load)
	if err != nil {
		t.Fatalf("entry.AppendOperation(load) error = %v", err)
	}

	condBr, err := cf.CondBranch(ctx, loc, cond, thenBorrowed, nil, elseBorrowed, nil)
	if err != nil {
		t.Fatalf("cf.CondBranch() error = %v", err)
	}
	if _, err := entry.AppendOperation(condBr); err != nil {
		t.Fatalf("entry.AppendOperation(cond_br) error = %v", err)
	}

	thenRet, err := funcdialect.Return(loc, loadOp.Results()[0])
	if err != nil {
		t.Fatalf("func.Return(then) error = %v", err)
	}
	if _, err := thenBlock.AppendOperation(thenRet); err != nil {
		t.Fatalf("thenBlock.AppendOperation(return) error = %v", err)
	}

	elseRet, err := funcdialect.Return(loc, c7Op.Results()[0])
	if err != nil {
		t.Fatalf("func.Return(else) error = %v", err)
	}
	if _, err := elseBlock.AppendOperation(elseRet); err != nil {
		t.Fatalf("elseBlock.AppendOperation(return) error = %v", err)
	}

	if _, err := body.AppendBlock(entry); err != nil {
		t.Fatalf("body.AppendBlock(entry) error = %v", err)
	}
	if _, err := body.AppendBlock(thenBlock); err != nil {
		t.Fatalf("body.AppendBlock(then) error = %v", err)
	}
	if _, err := body.AppendBlock(elseBlock); err != nil {
		t.Fatalf("body.AppendBlock(else) error = %v", err)
	}

	fn, err := funcdialect.Func(ctx, loc, "branchy", fnType, body)
	if err != nil {
		t.Fatalf("func.Func() error = %v", err)
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

	if _, err := mod.Body().AppendOwnedOperation(fn); err != nil {
		t.Fatalf("mod.Body().AppendOwnedOperation(fn) error = %v", err)
	}

	if err := mod.Verify(); err != nil {
		t.Fatalf("mod.Verify() error = %v", err)
	}
	runFileCheck(t, mod.String(), "cf_memref_module.mlir")
}

func TestTensorAndSCFWrappers(t *testing.T) {
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
	i1, err := mlir.SignlessIntegerType(ctx, 1)
	if err != nil {
		t.Fatalf("SignlessIntegerType(i1) error = %v", err)
	}
	i32, err := mlir.SignlessIntegerType(ctx, 32)
	if err != nil {
		t.Fatalf("SignlessIntegerType(i32) error = %v", err)
	}
	index, err := mlir.IndexType(ctx)
	if err != nil {
		t.Fatalf("IndexType() error = %v", err)
	}
	tensorType, err := mlir.RankedTensorType([]int64{mlir.DynamicSize()}, i32)
	if err != nil {
		t.Fatalf("RankedTensorType() error = %v", err)
	}
	fnType, err := mlir.FunctionType(ctx, []mlir.Type{i1, index}, []mlir.Type{index})
	if err != nil {
		t.Fatalf("FunctionType() error = %v", err)
	}

	body, err := mlir.NewOwnedRegion()
	if err != nil {
		t.Fatalf("NewOwnedRegion() error = %v", err)
	}
	entry, err := mlir.NewOwnedBlock([]mlir.Type{i1, index}, []mlir.Location{loc, loc})
	if err != nil {
		t.Fatalf("NewOwnedBlock(entry) error = %v", err)
	}

	entryBlock := entry.Block()
	cond := entryBlock.Arguments()[0]
	n := entryBlock.Arguments()[1]

	c0Attr, err := mlir.IntegerAttribute(index, 0)
	if err != nil {
		t.Fatalf("IntegerAttribute(index, 0) error = %v", err)
	}
	c0, err := arith.Constant(ctx, loc, index, c0Attr)
	if err != nil {
		t.Fatalf("arith.Constant(c0) error = %v", err)
	}
	c0Op, err := entry.AppendOperation(c0)
	if err != nil {
		t.Fatalf("entry.AppendOperation(c0) error = %v", err)
	}

	empty, err := tensor.Empty(loc, tensorType, n)
	if err != nil {
		t.Fatalf("tensor.Empty() error = %v", err)
	}
	emptyOp, err := entry.AppendOperation(empty)
	if err != nil {
		t.Fatalf("entry.AppendOperation(empty) error = %v", err)
	}

	dim, err := tensor.Dim(loc, emptyOp.Results()[0], c0Op.Results()[0])
	if err != nil {
		t.Fatalf("tensor.Dim() error = %v", err)
	}
	dimOp, err := entry.AppendOperation(dim)
	if err != nil {
		t.Fatalf("entry.AppendOperation(dim) error = %v", err)
	}

	thenRegion, err := mlir.NewOwnedRegion()
	if err != nil {
		t.Fatalf("NewOwnedRegion(then) error = %v", err)
	}
	thenBlock, err := mlir.NewOwnedBlock(nil, nil)
	if err != nil {
		t.Fatalf("NewOwnedBlock(then) error = %v", err)
	}
	thenYield, err := scf.Yield(loc, dimOp.Results()[0])
	if err != nil {
		t.Fatalf("scf.Yield(then) error = %v", err)
	}
	if _, err := thenBlock.AppendOperation(thenYield); err != nil {
		t.Fatalf("thenBlock.AppendOperation(yield) error = %v", err)
	}
	if _, err := thenRegion.AppendBlock(thenBlock); err != nil {
		t.Fatalf("thenRegion.AppendBlock() error = %v", err)
	}

	elseRegion, err := mlir.NewOwnedRegion()
	if err != nil {
		t.Fatalf("NewOwnedRegion(else) error = %v", err)
	}
	elseBlock, err := mlir.NewOwnedBlock(nil, nil)
	if err != nil {
		t.Fatalf("NewOwnedBlock(else) error = %v", err)
	}
	elseYield, err := scf.Yield(loc, n)
	if err != nil {
		t.Fatalf("scf.Yield(else) error = %v", err)
	}
	if _, err := elseBlock.AppendOperation(elseYield); err != nil {
		t.Fatalf("elseBlock.AppendOperation(yield) error = %v", err)
	}
	if _, err := elseRegion.AppendBlock(elseBlock); err != nil {
		t.Fatalf("elseRegion.AppendBlock() error = %v", err)
	}

	ifOp, err := scf.If(loc, cond, []mlir.Type{index}, thenRegion, elseRegion)
	if err != nil {
		t.Fatalf("scf.If() error = %v", err)
	}
	ifResultOp, err := entry.AppendOperation(ifOp)
	if err != nil {
		t.Fatalf("entry.AppendOperation(if) error = %v", err)
	}

	ret, err := funcdialect.Return(loc, ifResultOp.Results()[0])
	if err != nil {
		t.Fatalf("func.Return() error = %v", err)
	}
	if _, err := entry.AppendOperation(ret); err != nil {
		t.Fatalf("entry.AppendOperation(return) error = %v", err)
	}

	if _, err := body.AppendBlock(entry); err != nil {
		t.Fatalf("body.AppendBlock(entry) error = %v", err)
	}

	fn, err := funcdialect.Func(ctx, loc, "select_dim", fnType, body)
	if err != nil {
		t.Fatalf("func.Func() error = %v", err)
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

	if _, err := mod.Body().AppendOwnedOperation(fn); err != nil {
		t.Fatalf("mod.Body().AppendOwnedOperation(fn) error = %v", err)
	}

	if err := mod.Verify(); err != nil {
		t.Fatalf("mod.Verify() error = %v", err)
	}
	runFileCheck(t, mod.String(), "tensor_scf_module.mlir")
}

func TestGeneratedArithWrappers(t *testing.T) {
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
	i1, err := mlir.SignlessIntegerType(ctx, 1)
	if err != nil {
		t.Fatalf("SignlessIntegerType(i1) error = %v", err)
	}
	i32, err := mlir.SignlessIntegerType(ctx, 32)
	if err != nil {
		t.Fatalf("SignlessIntegerType(i32) error = %v", err)
	}
	i64, err := mlir.SignlessIntegerType(ctx, 64)
	if err != nil {
		t.Fatalf("SignlessIntegerType(i64) error = %v", err)
	}
	index, err := mlir.IndexType(ctx)
	if err != nil {
		t.Fatalf("IndexType() error = %v", err)
	}
	f32, err := mlir.F32Type(ctx)
	if err != nil {
		t.Fatalf("F32Type() error = %v", err)
	}
	fnType, err := mlir.FunctionType(ctx, []mlir.Type{i32, i32, index}, []mlir.Type{i1, i32, i64, i64, f32, i32, i64})
	if err != nil {
		t.Fatalf("FunctionType() error = %v", err)
	}

	body, err := mlir.NewOwnedRegion()
	if err != nil {
		t.Fatalf("NewOwnedRegion() error = %v", err)
	}
	entry, err := mlir.NewOwnedBlock([]mlir.Type{i32, i32, index}, []mlir.Location{loc, loc, loc})
	if err != nil {
		t.Fatalf("NewOwnedBlock() error = %v", err)
	}
	args := entry.Block().Arguments()

	sub, err := arith.SubI(loc, args[0], args[1])
	if err != nil {
		t.Fatalf("arith.SubI() error = %v", err)
	}
	subOp, err := entry.AppendOperation(sub)
	if err != nil {
		t.Fatalf("entry.AppendOperation(sub) error = %v", err)
	}

	mul, err := arith.MulI(loc, args[0], args[1])
	if err != nil {
		t.Fatalf("arith.MulI() error = %v", err)
	}
	mulOp, err := entry.AppendOperation(mul)
	if err != nil {
		t.Fatalf("entry.AppendOperation(mul) error = %v", err)
	}

	extsi, err := arith.ExtSI(loc, subOp.Results()[0], i64)
	if err != nil {
		t.Fatalf("arith.ExtSI() error = %v", err)
	}
	extsiOp, err := entry.AppendOperation(extsi)
	if err != nil {
		t.Fatalf("entry.AppendOperation(extsi) error = %v", err)
	}

	extui, err := arith.ExtUI(loc, subOp.Results()[0], i64)
	if err != nil {
		t.Fatalf("arith.ExtUI() error = %v", err)
	}
	extuiOp, err := entry.AppendOperation(extui)
	if err != nil {
		t.Fatalf("entry.AppendOperation(extui) error = %v", err)
	}

	trunc, err := arith.TruncI(loc, extsiOp.Results()[0], i32)
	if err != nil {
		t.Fatalf("arith.TruncI() error = %v", err)
	}
	truncOp, err := entry.AppendOperation(trunc)
	if err != nil {
		t.Fatalf("entry.AppendOperation(trunc) error = %v", err)
	}

	sitofp, err := arith.SIToFP(loc, mulOp.Results()[0], f32)
	if err != nil {
		t.Fatalf("arith.SIToFP() error = %v", err)
	}
	sitofpOp, err := entry.AppendOperation(sitofp)
	if err != nil {
		t.Fatalf("entry.AppendOperation(sitofp) error = %v", err)
	}

	fptosi, err := arith.FPToSI(loc, sitofpOp.Results()[0], i32)
	if err != nil {
		t.Fatalf("arith.FPToSI() error = %v", err)
	}
	fptosiOp, err := entry.AppendOperation(fptosi)
	if err != nil {
		t.Fatalf("entry.AppendOperation(fptosi) error = %v", err)
	}

	indexCast, err := arith.IndexCast(loc, args[2], i64)
	if err != nil {
		t.Fatalf("arith.IndexCast() error = %v", err)
	}
	indexCastOp, err := entry.AppendOperation(indexCast)
	if err != nil {
		t.Fatalf("entry.AppendOperation(index_cast) error = %v", err)
	}

	cmp, err := arith.CmpI(ctx, loc, "sgt", truncOp.Results()[0], fptosiOp.Results()[0])
	if err != nil {
		t.Fatalf("arith.CmpI() error = %v", err)
	}
	cmpOp, err := entry.AppendOperation(cmp)
	if err != nil {
		t.Fatalf("entry.AppendOperation(cmpi) error = %v", err)
	}

	ret, err := funcdialect.Return(
		loc,
		cmpOp.Results()[0],
		truncOp.Results()[0],
		extsiOp.Results()[0],
		extuiOp.Results()[0],
		sitofpOp.Results()[0],
		fptosiOp.Results()[0],
		indexCastOp.Results()[0],
	)
	if err != nil {
		t.Fatalf("func.Return() error = %v", err)
	}
	if _, err := entry.AppendOperation(ret); err != nil {
		t.Fatalf("entry.AppendOperation(return) error = %v", err)
	}

	if _, err := body.AppendBlock(entry); err != nil {
		t.Fatalf("body.AppendBlock() error = %v", err)
	}

	fn, err := funcdialect.Func(ctx, loc, "arith_generated", fnType, body)
	if err != nil {
		t.Fatalf("func.Func() error = %v", err)
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

	if _, err := mod.Body().AppendOwnedOperation(fn); err != nil {
		t.Fatalf("mod.Body().AppendOwnedOperation(fn) error = %v", err)
	}

	if err := mod.Verify(); err != nil {
		t.Fatalf("mod.Verify() error = %v", err)
	}
	runFileCheck(t, mod.String(), "generated_arith_module.mlir")
}

func TestGeneratedVectorWrappers(t *testing.T) {
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
		t.Fatalf("SignlessIntegerType(i32) error = %v", err)
	}
	index, err := mlir.IndexType(ctx)
	if err != nil {
		t.Fatalf("IndexType() error = %v", err)
	}
	vectorType, err := mlir.ParseType(ctx, "vector<4xi32>")
	if err != nil {
		t.Fatalf("ParseType(vector<4xi32>) error = %v", err)
	}
	fnType, err := mlir.FunctionType(ctx, []mlir.Type{i32, index}, []mlir.Type{i32})
	if err != nil {
		t.Fatalf("FunctionType() error = %v", err)
	}

	body, err := mlir.NewOwnedRegion()
	if err != nil {
		t.Fatalf("NewOwnedRegion() error = %v", err)
	}
	entry, err := mlir.NewOwnedBlock([]mlir.Type{i32, index}, []mlir.Location{loc, loc})
	if err != nil {
		t.Fatalf("NewOwnedBlock() error = %v", err)
	}
	args := entry.Block().Arguments()

	broadcast, err := vector.Broadcast(loc, args[0], vectorType)
	if err != nil {
		t.Fatalf("vector.Broadcast() error = %v", err)
	}
	broadcastOp, err := entry.AppendOperation(broadcast)
	if err != nil {
		t.Fatalf("entry.AppendOperation(broadcast) error = %v", err)
	}

	extract, err := vector.ExtractElement(loc, broadcastOp.Results()[0], i32, args[1])
	if err != nil {
		t.Fatalf("vector.ExtractElement() error = %v", err)
	}
	extractOp, err := entry.AppendOperation(extract)
	if err != nil {
		t.Fatalf("entry.AppendOperation(extractelement) error = %v", err)
	}

	ret, err := funcdialect.Return(loc, extractOp.Results()[0])
	if err != nil {
		t.Fatalf("func.Return() error = %v", err)
	}
	if _, err := entry.AppendOperation(ret); err != nil {
		t.Fatalf("entry.AppendOperation(return) error = %v", err)
	}

	if _, err := body.AppendBlock(entry); err != nil {
		t.Fatalf("body.AppendBlock() error = %v", err)
	}

	fn, err := funcdialect.Func(ctx, loc, "vector_generated", fnType, body)
	if err != nil {
		t.Fatalf("func.Func() error = %v", err)
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

	if _, err := mod.Body().AppendOwnedOperation(fn); err != nil {
		t.Fatalf("mod.Body().AppendOwnedOperation(fn) error = %v", err)
	}

	if err := mod.Verify(); err != nil {
		t.Fatalf("mod.Verify() error = %v", err)
	}
	runFileCheck(t, mod.String(), "vector_generated_module.mlir")
}

func TestGeneratedLinalgWrappers(t *testing.T) {
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

	indexOp, err := linalg.Index(ctx, loc, 3)
	if err != nil {
		t.Fatalf("linalg.Index() error = %v", err)
	}
	defer func() {
		if err := indexOp.Close(); err != nil {
			t.Fatalf("indexOp.Close() error = %v", err)
		}
	}()

	indexBorrowed := indexOp.Operation()
	if got := indexBorrowed.Name(); got != "linalg.index" {
		t.Fatalf("indexBorrowed.Name() = %q, want %q", got, "linalg.index")
	}
	if len(indexBorrowed.Results()) != 1 || !indexBorrowed.Results()[0].Type().IsIndex() {
		t.Fatalf("linalg.Index() returned unexpected result types")
	}
	if got := indexBorrowed.Attribute("dim").String(); !strings.Contains(got, "3") {
		t.Fatalf("indexBorrowed.Attribute(\"dim\") = %q, want dim containing 3", got)
	}

	yieldOp, err := linalg.Yield(loc, indexBorrowed.Results()[0])
	if err != nil {
		t.Fatalf("linalg.Yield() error = %v", err)
	}
	defer func() {
		if err := yieldOp.Close(); err != nil {
			t.Fatalf("yieldOp.Close() error = %v", err)
		}
	}()

	yieldBorrowed := yieldOp.Operation()
	if got := yieldBorrowed.Name(); got != "linalg.yield" {
		t.Fatalf("yieldBorrowed.Name() = %q, want %q", got, "linalg.yield")
	}
	if len(yieldBorrowed.Operands()) != 1 || !yieldBorrowed.Operands()[0].Equal(indexBorrowed.Results()[0]) {
		t.Fatalf("linalg.Yield() returned unexpected operands")
	}
}
