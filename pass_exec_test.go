//go:build cgo

package mlir

import (
	"testing"
	"unsafe"
)

func addModuleSource() string {
	return `module {
  func.func @add(%arg0: i32) -> i32 attributes {llvm.emit_c_interface} {
    %0 = arith.addi %arg0, %arg0 : i32
    return %0 : i32
  }
}
`
}

func lowerAddModule(t *testing.T) (*Context, *Module) {
	t.Helper()

	ctx, err := NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}

	if err := ctx.RegisterAllDialects(); err != nil {
		t.Fatalf("RegisterAllDialects() error = %v", err)
	}
	RegisterAllPasses()

	mod, err := ParseModule(ctx, addModuleSource())
	if err != nil {
		t.Fatalf("ParseModule() error = %v", err)
	}

	pm, err := NewPassManager(ctx)
	if err != nil {
		t.Fatalf("NewPassManager() error = %v", err)
	}
	defer func() {
		if err := pm.Close(); err != nil {
			t.Fatalf("pm.Close() error = %v", err)
		}
	}()

	const pipeline = "builtin.module(convert-arith-to-llvm,convert-func-to-llvm,reconcile-unrealized-casts)"
	if err := pm.ParsePipeline(pipeline); err != nil {
		t.Fatalf("ParsePipeline() error = %v", err)
	}
	if got := pm.String(); got == "" {
		t.Fatalf("pm.String() = empty, want non-empty")
	}
	if err := pm.Run(mod.Operation()); err != nil {
		t.Fatalf("pm.Run() error = %v", err)
	}
	return ctx, mod
}

func TestPassManagerLowering(t *testing.T) {
	ctx, mod := lowerAddModule(t)
	defer func() {
		if err := mod.Close(); err != nil {
			t.Fatalf("mod.Close() error = %v", err)
		}
		if err := ctx.Close(); err != nil {
			t.Fatalf("ctx.Close() error = %v", err)
		}
	}()

	runFileCheck(t, mod.String(), "lowered_add_module.mlir")
}

func TestExecutionEngineInvokePacked(t *testing.T) {
	ctx, mod := lowerAddModule(t)
	defer func() {
		if err := mod.Close(); err != nil {
			t.Fatalf("mod.Close() error = %v", err)
		}
		if err := ctx.Close(); err != nil {
			t.Fatalf("ctx.Close() error = %v", err)
		}
	}()

	if err := ctx.RegisterAllLLVMTranslations(); err != nil {
		t.Fatalf("RegisterAllLLVMTranslations() error = %v", err)
	}

	engine, err := NewExecutionEngine(mod, 0)
	if err != nil {
		t.Fatalf("NewExecutionEngine() error = %v", err)
	}
	defer func() {
		if err := engine.Close(); err != nil {
			t.Fatalf("engine.Close() error = %v", err)
		}
	}()

	arg := int32(123)
	var result int32
	if err := engine.InvokePacked("add", unsafe.Pointer(&result), unsafe.Pointer(&arg)); err != nil {
		t.Fatalf("InvokePacked() error = %v", err)
	}
	if result != 246 {
		t.Fatalf("result = %d, want 246", result)
	}
}
