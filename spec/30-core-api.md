# 30. Core API

## Package Shape

The repository should converge on a layout like this:

```text
/go.mod
/README.md
/spec/
/context.go
/module.go
/operation.go
/type.go
/attribute.go
/value.go
/location.go
/diagnostic.go
/builder/
/pass/
/exec/
/dialect/
/internal/capi/
/internal/gen/
/cmd/mlir-go-tblgen/
/examples/
/testdata/
```

## Layering

The public API should have three layers.

### Layer 1: Core Handle API

This is the minimum API needed to model MLIR through the C API.

Expected capabilities:

- create and destroy contexts
- parse modules from assembly
- print modules and operations
- verify modules and operations
- traverse regions, blocks, and operations
- query types, attributes, operands, results, and locations

### Layer 2: Generic Operation Construction

This is the low-level escape hatch used by both humans and generated code.

The API should include a generic operation-state builder, for example:

```go
state := mlir.NewOperationState("arith.addf", loc)
state.AddOperands(x, y)
state.AddResults(x.Type())
op, err := mlir.NewOperation(ctx, state)
```

This layer is required even if dialect-specific wrappers exist.

### Layer 3: Higher-Level Convenience APIs

This includes:

- `builder`
- dialect wrappers
- pass helpers
- execution-engine helpers

Higher-level layers must never remove access to lower-level construction paths.

## Minimum V1 API Surface

The following names are illustrative, not frozen, but equivalent capabilities are required:

- `NewContext() (*Context, error)`
- `(*Context).Close() error`
- `(*Context).RegisterAllDialects() error`
- `ParseModule(ctx *Context, asm string) (*Module, error)`
- `(*Module).Close() error`
- `(*Module).Verify() error`
- `(*Module).String() string`
- `(*Module).Operation() Operation`
- `Operation.Verify() error`
- `Operation.String() string`
- `Operation.Regions() []Region`
- `Region.Blocks() []Block`
- `Block.Operations() []Operation`

## Diagnostics

Diagnostics should not be an afterthought.
The core API should be designed so parse, verify, and pass failures can surface readable MLIR diagnostics without forcing users to scrape stderr.

That means:

- installable diagnostic hooks at the C API boundary
- error values that preserve the rendered message
- tests that validate actual diagnostic text on failure paths

