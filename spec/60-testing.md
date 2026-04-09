# 60. Testing

## Testing Principle

A binding is only credible if it proves real workflows.
Parse-and-print alone is not enough.

## Required Test Layers

### 1. Resource And Handle Tests

Cover:

- context creation and destruction
- module creation and destruction
- null handle behavior
- borrowed handle access on valid objects

### 2. IR Functionality Tests

Cover:

- parsing valid IR
- rejecting invalid IR
- printing modules and operations
- verification success and failure
- traversing regions, blocks, and operations
- locations, types, and attributes

### 3. Builder Tests

Cover:

- building a simple function
- block arguments
- branching and multiple blocks
- default location propagation
- missing terminator behavior

### 4. Pass And Lowering Tests

Cover:

- pass-manager creation
- pass pipeline execution
- at least one lowering path from high-level dialects to LLVM-compatible IR

### 5. Execution Engine Tests

Cover:

- creating an execution engine
- invoking a simple function
- validating returned results

If JIT support is unstable on one platform, the test may live in a specialized CI job, but it must exist.

## Generated Surface Tests

Each enabled dialect should have at least one generated validation strategy:

- a smoke test
- a golden print test
- or a round-trip construction test

The goal is not to prove every op semantic exhaustively.
The goal is to catch signature drift, attribute mapping mistakes, variadic handling bugs, and naming conflicts early.

## Negative Tests

The suite must include negative cases for:

- invalid parse input
- verification failure
- pass pipeline failure
- missing symbol invocation
- use after close, whether that is documented as error or panic

## CI Requirements

The first CI matrix should cover:

- `linux/amd64`
- `darwin/arm64`

CI should perform at least:

1. generated code freshness checks
2. package build
3. unit tests
4. integration tests
5. at least one lowering plus execution-engine path

