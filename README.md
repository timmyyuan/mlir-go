# mlir-go

An early attempt at a Go-facing MLIR interface built around the MLIR C API.

The intent is modest: keep the core close to native MLIR handles, add a small builder where it helps, generate wrappers for common dialects, and keep enough low-level access available when the high-level layer is not enough.

## Status

The repository is still in the design and bootstrap stage.
The current work is centered on:

- core handle-based bindings
- a lightweight builder layer
- generated dialect wrappers
- pass management and execution support
- reproducible code generation and testing

The first version is not aiming for full dialect coverage, direct C++ bindings, or one codebase that spans multiple LLVM/MLIR major versions.

## Current Slice

The repository now includes the first bootstrap layer:

- a Go module
- a minimal `internal/capi` bridge
- explicit `Context` lifecycle management
- minimal module parse, print, and verify support
- location, region, and block traversal
- basic `Type` and `Value` inspection
- small bootstrap tests for context and module round-trips

## Design Notes

The design takes cues from existing MLIR binding projects, especially in layering, code generation, and end-to-end validation, but keeps the public API shaped around Go conventions and explicit resource ownership.

## Developer Setup

This project depends on a local MLIR toolchain that exposes the MLIR C headers and libraries.

To prepare a shell for development:

```bash
eval "$(./scripts/dev-env.sh)"
```

That script uses `llvm-config` from `PATH` by default.
If needed, point it at a specific installation:

```bash
eval "$(LLVM_CONFIG=/path/to/llvm-config ./scripts/dev-env.sh)"
```

With the environment prepared, the current bootstrap test can be run with:

```bash
go test ./...
```

The test suite uses two layers:

- regular Go tests for API behavior and ownership rules
- `FileCheck` fixtures for textual IR validation

The development environment script adds the LLVM tools directory to `PATH`, so `FileCheck` is available when present in the selected LLVM installation.

## Using From Another Module

The module path is:

```go
import mlir "github.com/timmyyuan/mlir-go"
```

Downstream builds still need the MLIR C headers and libraries visible to cgo.
Outside this repository, set equivalent flags with `llvm-config`, for example:

```bash
export CGO_CPPFLAGS="-I$(llvm-config --includedir)"
export CGO_LDFLAGS="-L$(llvm-config --libdir) -Wl,-rpath,$(llvm-config --libdir) -lMLIR -lMLIRCAPIIR -lMLIRCAPIRegisterEverything $(llvm-config --libs --system-libs) -lc++"
```

The initial public API is intentionally small:

```go
ctx, err := mlir.NewContext()
if err != nil {
    // handle error
}
defer ctx.Close()

if err := ctx.RegisterAllDialects(); err != nil {
    // handle error
}

mod, err := mlir.ParseModule(ctx, "module {\n}\n")
if err != nil {
    // handle error
}
defer mod.Close()

if err := mod.Verify(); err != nil {
    // handle error
}

top := mod.Operation()
for _, region := range top.Regions() {
    for _, block := range region.Blocks() {
        for _, arg := range block.Arguments() {
            _ = arg.Type().String()
        }
        for _, op := range block.Operations() {
            _ = op.Name()
            for _, result := range op.Results() {
                _ = result.Type().String()
            }
        }
    }
}
```

## Spec

The project spec is split into small documents under [spec/](spec/README.md).
