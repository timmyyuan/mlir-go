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
- a first `builder` package for modules, functions, blocks, insertion points, and default locations
- builder-friendly dialect emit helpers under `builder/arith`, `builder/func`, and `builder/cf`
- basic IR mutation paths for regions, blocks, and SSA use replacement
- diagnostic capture around parser and checked-construction failures
- first dialect convenience wrappers under `dialect/arith`, `dialect/func`, `dialect/cf`, `dialect/memref`, `dialect/tensor`, `dialect/scf`, `dialect/vector`, and `dialect/linalg`
- a first `cmd/mlir-go-tblgen` entry point plus a checked-in `dialect_manifest.json`
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

## Generator Bootstrap

The repository now includes a first generator entry point:

```bash
go run ./cmd/mlir-go-tblgen -mode validate-manifest
go run ./cmd/mlir-go-tblgen -mode emit-skip-report
go run ./cmd/mlir-go-tblgen -mode emit-dialect -dialect arith
go run ./cmd/mlir-go-tblgen -mode emit-dialect -dialect func
go run ./cmd/mlir-go-tblgen -mode emit-dialect -dialect cf
go run ./cmd/mlir-go-tblgen -mode emit-dialect -dialect memref
go run ./cmd/mlir-go-tblgen -mode emit-dialect -dialect tensor
go run ./cmd/mlir-go-tblgen -mode emit-dialect -dialect vector
go run ./cmd/mlir-go-tblgen -mode emit-dialect -dialect linalg
bash ./scripts/update-generated.sh
bash ./scripts/update-skip-report.sh
```

By default the tool resolves the MLIR include root through `llvm-config --includedir`.
You can override that explicitly:

```bash
go run ./cmd/mlir-go-tblgen \
  -mode emit-skip-report \
  -mlir-include-root /path/to/mlir/include \
  -o /tmp/mlir-go-skip-report.json
```

The checked-in [dialect_manifest.json](dialect_manifest.json) is the contract for supported generated surfaces.
This first version validates the manifest, scans upstream TableGen operation definitions, and emits a machine-readable skip report for operations that are not generated yet.

## Using From Another Module

The module path is:

```go
import mlir "github.com/timmyyuan/mlir-go"
```

The first convenience dialect layers live under:

```go
import "github.com/timmyyuan/mlir-go/builder"
import builderarith "github.com/timmyyuan/mlir-go/builder/arith"
import buildercf "github.com/timmyyuan/mlir-go/builder/cf"
import builderfunc "github.com/timmyyuan/mlir-go/builder/func"
import arith "github.com/timmyyuan/mlir-go/dialect/arith"
import cf "github.com/timmyyuan/mlir-go/dialect/cf"
import funcdialect "github.com/timmyyuan/mlir-go/dialect/func"
import linalg "github.com/timmyyuan/mlir-go/dialect/linalg"
import memref "github.com/timmyyuan/mlir-go/dialect/memref"
import scf "github.com/timmyyuan/mlir-go/dialect/scf"
import tensor "github.com/timmyyuan/mlir-go/dialect/tensor"
import vector "github.com/timmyyuan/mlir-go/dialect/vector"
```

Downstream projects need an equivalent `cgo` environment.
A minimal setup for the current surface looks like:

```bash
export CGO_CPPFLAGS="-I$(llvm-config --includedir)"
export CGO_LDFLAGS="-L$(llvm-config --libdir) -Wl,-rpath,$(llvm-config --libdir) -lMLIR -lMLIRCAPIIR -lMLIRCAPIRegisterEverything -lMLIRCAPITransforms -lMLIRCAPIConversion -lMLIRCAPIExecutionEngine -lMLIRExecutionEngine -lMLIRCAPIFunc -lMLIRCAPIArith -lMLIRCAPILLVM $(llvm-config --libs --system-libs) -lc++"
```

The current public API is still intentionally small, but it already supports parse, traversal, operand/attribute/successor queries, symbol lookup and rewrite flows, builtin type and attribute construction, shaped and memref type helpers, generic IR building, a first stateful builder with dialect emit helpers, verification, lowering, and JIT execution:

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

For a full builder-based example, see [`examples/build_add/main.go`](examples/build_add/main.go).

## Spec

The design notes live under [spec/](spec/README.md).
