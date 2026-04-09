# 40. Builder And Dialects

## Builder Goals

The builder layer should make common IR construction natural in Go while staying honest about MLIR concepts.
It should help users build valid IR, not pretend MLIR has disappeared.

## Required Builder Capabilities

The first builder milestone should support:

- module construction
- function construction
- block creation
- block arguments
- SSA result value flow
- default location propagation
- terminator-aware block finalization

## Builder Rules

- SSA result names are internal by default.
- Builder-returned `Value` handles must compose directly with dialect helpers.
- Missing terminators should fail as early as practical.
- Default locations must affect subsequent emitted operations.
- The builder should not require users to manage `MlirOperationState` directly for common cases.

## Dialect Coverage Priority

The first generated or hand-written dialect set should prioritize:

- `builtin`
- `func`
- `arith`
- `cf`
- `memref`
- `tensor`
- `vector`
- `linalg`
- `llvm`

## Dialect API Shape

Dialect APIs should have two forms:

1. low-level wrappers around generic operation construction
2. builder-friendly helpers that emit into the current block

That keeps generated code useful both inside and outside the builder package.

## Go Naming Rules

Some MLIR dialect names collide with Go keywords.
The project should use these rules:

- import paths stay close to upstream dialect naming
- package names use valid Go identifiers
- examples use a single repository-wide alias convention

For example, an import path under `dialect/func` may expose package name `funcdialect`.

## Why Not A Full Shadow AST

The reference Haskell project uses a higher-level AST-centric style.
That is reasonable in Haskell, but this Go project should not make a full shadow AST the center of the design.

Reasons:

- it duplicates MLIR structure
- it creates extra invalid intermediate states
- it weakens the clarity of native object lifetime
- it is not required to provide a good builder experience in Go

If an offline AST is ever added, it should live in an experimental package, not define the entire public architecture.

