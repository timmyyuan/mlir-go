# 20. Runtime Model

## Why This Matters

Bindings fail in practice when ownership is vague.
MLIR objects are lightweight handles into native state, but their validity still depends on larger owning objects such as contexts and modules.

This repository must define those rules up front.

## Two Kinds Of Public Types

### Owning Types

Owning types manage native resources and must expose `Close()`:

- `*Context`
- `*Module`
- `*PassManager`
- `*ExecutionEngine`

Rules:

- zero values are not valid instances
- construction happens through constructors
- `Close()` must be idempotent
- docs must say which borrowed handles become invalid after closing

### Borrowed Handle Types

Borrowed handles are cheap value types without `Close()`:

- `Operation`
- `Region`
- `Block`
- `Type`
- `Attribute`
- `Value`
- `Location`
- `Identifier`
- `AffineExpr`
- `AffineMap`

Rules:

- zero value means null handle
- each type exposes `IsNull() bool`
- handles can be copied by value
- borrowed handles never outlive the owning native state they depend on

## Invalidation Rules

- Closing a `Context` invalidates all dependent objects.
- Closing a `Module` invalidates module-derived `Operation`, `Region`, `Block`, and `Value` handles.
- Borrowed handles are not independently reference-counted.
- Public docs must mark returned values as owned or borrowed.

## Error Model

Operations that can fail must return `error`.
Do not encode meaningful failure only as `bool`.

Required categories:

- parse errors
- verification errors
- pass execution errors
- execution-engine invocation errors

The project should define a diagnostic-backed error type that can carry:

- stage name
- rendered diagnostic text
- optional structured diagnostic fields

## Concurrency Model

The initial policy should be conservative:

- public objects are not goroutine-safe unless explicitly documented
- concurrent mutation through one context is unsupported
- read-only safety may be relaxed later, but only after validation

## cgo Safety Rules

The `internal/capi` layer must follow these rules:

- use official `C.Mlir*` declarations from MLIR headers
- do not re-declare MLIR struct layouts manually
- do not persist Go pointers inside C-owned state
- centralize string and slice bridging helpers
- use `runtime.KeepAlive` when object lifetime depends on cgo call boundaries
- keep `unsafe.Pointer` out of the normal user-facing path

