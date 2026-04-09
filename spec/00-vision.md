# 00. Vision

## Problem

MLIR already exposes a C API, but that is not enough to make MLIR pleasant to use from Go.
Users still need:

- a Go-friendly ownership model
- predictable error handling
- a workable builder layer
- generated wrappers for common dialects
- tests that prove real workflows, not only raw handle creation

## Goal

This repository should provide a Go binding for MLIR that lets users:

- create and destroy MLIR contexts and modules
- parse, print, and verify IR
- traverse operations, regions, blocks, types, attributes, and values
- construct IR using a Go-style builder
- run common pass pipelines
- invoke simple JIT-compiled functions through the execution engine
- use generated wrappers for common dialect operations without losing access to lower-level escape hatches

## Success Criteria For V1

Version 1 is successful when a Go user can complete this path:

1. Create a context.
2. Register dialects.
3. Parse or build a module containing `func` and `arith`.
4. Verify and print that module.
5. Run a lowering pipeline.
6. Invoke a simple function through the execution engine.

## Guiding Principles

- Use the MLIR C API as the stable boundary.
- Prefer live IR handles over a full shadow AST.
- Keep a low-level escape hatch available at all times.
- Make ownership and invalidation rules explicit.
- Generate mechanical dialect wrappers instead of hand-writing large surfaces.
- Check generated code into the repository.
- Validate behavior with end-to-end tests, not only unit tests.

## Non-Goals For V1

- binding the MLIR C++ API directly
- covering every dialect or every extension point
- supporting multiple LLVM/MLIR major versions in one release line
- relying on Go finalizers for correctness
- shipping a full rewrite-pattern DSL in the first milestone
- promising Windows support in the first milestone

