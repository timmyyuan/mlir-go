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
	memref "github.com/timmyyuan/mlir-go/dialect/memref"
	scf "github.com/timmyyuan/mlir-go/dialect/scf"
	tensor "github.com/timmyyuan/mlir-go/dialect/tensor"
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
