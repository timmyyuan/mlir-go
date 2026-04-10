# mlir-go

A small Go-facing MLIR interface built on the MLIR C API.

The project stays close to upstream handles, adds construction helpers where the raw API gets awkward, and keeps low-level access available for cases where a future builder or dialect layer is not enough.

## Status

This repository is still early, but it is no longer documentation-only.
The current slice covers:

- core handle-based bindings for context, module, operation, location, region, block, type, value, attribute, and identifier
- symbol table helpers for lookup, uniquing, visibility, and symbol-use rewrites
- builtin type and attribute constructors for common MLIR forms
- generic operation construction through `OperationState`
- basic IR mutation paths for regions, blocks, and SSA use replacement
- diagnostic capture around parser and checked-construction failures
- first dialect convenience wrappers under `dialect/arith`, `dialect/func`, `dialect/cf`, `dialect/memref`, `dialect/tensor`, and `dialect/scf`
- pass manager setup and pipeline execution
- execution engine creation and packed invocation
- Go tests plus `FileCheck`-based textual IR checks

The first version is not trying to expose all upstream C APIs at once, bind directly to C++, or support multiple LLVM/MLIR major versions from one code line.

## Developer Setup

Real bindings require `cgo` and a local MLIR installation that exposes the MLIR C headers, libraries, and tools.
The repository also includes non-`cgo` fallback files so unsupported builds fail with explicit runtime errors instead of opaque linker failures.

To prepare a development shell:

```bash
eval "$(./scripts/dev-env.sh)"
```

That script uses `llvm-config` from `PATH` by default.
To select a different installation:

```bash
eval "$(LLVM_CONFIG=/path/to/llvm-config ./scripts/dev-env.sh)"
```

It exports:

- `CGO_CPPFLAGS`
- `CGO_LDFLAGS`
- `PATH` with the selected LLVM tool directory, so `FileCheck` is available when present

Run the full test suite with:

```bash
go test -count=1 ./...
```

The test suite has two layers:

- regular Go tests for API behavior, ownership, mutation, passes, and execution
- `FileCheck` fixtures for textual IR validation

## Using From Another Module

The module path is:

```go
import mlir "github.com/timmyyuan/mlir-go"
```

The first convenience dialect layers live under:

```go
import arith "github.com/timmyyuan/mlir-go/dialect/arith"
import cf "github.com/timmyyuan/mlir-go/dialect/cf"
import funcdialect "github.com/timmyyuan/mlir-go/dialect/func"
import memref "github.com/timmyyuan/mlir-go/dialect/memref"
import scf "github.com/timmyyuan/mlir-go/dialect/scf"
import tensor "github.com/timmyyuan/mlir-go/dialect/tensor"
```

Downstream projects need an equivalent `cgo` environment.
A minimal setup for the current surface looks like:

```bash
export CGO_CPPFLAGS="-I$(llvm-config --includedir)"
export CGO_LDFLAGS="-L$(llvm-config --libdir) -Wl,-rpath,$(llvm-config --libdir) -lMLIR -lMLIRCAPIIR -lMLIRCAPIRegisterEverything -lMLIRCAPITransforms -lMLIRCAPIConversion -lMLIRCAPIExecutionEngine -lMLIRExecutionEngine -lMLIRCAPIFunc -lMLIRCAPIArith -lMLIRCAPILLVM $(llvm-config --libs --system-libs) -lc++"
```

The current public API is still intentionally small, but it already supports parse, traversal, operand/attribute/successor queries, symbol lookup and rewrite flows, builtin type and attribute construction, shaped and memref type helpers, generic IR building, verification, lowering, and JIT execution:

```go
ctx, err := mlir.NewContext()
if err != nil {
    // handle error
}
defer ctx.Close()

if err := ctx.RegisterAllDialects(); err != nil {
    // handle error
}

i32, err := mlir.SignlessIntegerType(ctx, 32)
if err != nil {
    // handle error
}

loc, err := mlir.UnknownLocation(ctx)
if err != nil {
    // handle error
}

attr, err := mlir.IntegerAttribute(i32, 5)
if err != nil {
    // handle error
}

constOp, err := arith.Constant(ctx, loc, i32, attr)
if err != nil {
    // handle error
}
defer constOp.Close()

diags, err := ctx.CaptureDiagnostics(func() error {
    _, err := mlir.ParseModule(ctx, "module {\n  func.func @broken(\n}\n")
    return err
})
if err != nil && len(diags) > 0 {
    _ = diags[0].String()
}

mod, err := mlir.ParseModule(ctx, "module {\n}\n")
if err != nil {
    // handle error
}
defer mod.Close()
```

## Spec

The design notes live under [spec/](spec/README.md).
