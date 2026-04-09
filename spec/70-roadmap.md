# 70. Roadmap

## M0: Repository Bootstrap

Deliverables:

- `go.mod`
- initial `README.md`
- `internal/capi` skeleton
- MLIR toolchain detection
- initial CI skeleton

Exit criteria:

- a tiny program can create and destroy a context successfully

## M1: Core Handles

Deliverables:

- context, module, operation, region, block, type, attribute, value handles
- parse, print, verify
- diagnostic-backed error model

Exit criteria:

- a simple module can round-trip through parse and print
- a function body can be traversed through operations and blocks

## M2: Builder And Base Dialects

Deliverables:

- first builder package
- first hand-written dialect helpers for `func`, `arith`, and `cf`
- location propagation support

Exit criteria:

- a simple `add` function can be built and verified
- a control-flow example with multiple blocks can be expressed

## M3: Generator And Extended Dialects

Deliverables:

- `cmd/mlir-go-tblgen`
- dialect manifest
- generated wrappers for at least `memref`, `tensor`, `vector`, `linalg`, and `llvm`
- generated smoke or golden tests

Exit criteria:

- generated code is checked in
- unsupported operations are reported explicitly
- generated dialect tests pass in CI

## M4: Passes And Execution

Deliverables:

- `pass` package
- `exec` package
- example lowering pipelines

Exit criteria:

- a module can be lowered to LLVM-compatible IR
- a simple function can be invoked through the execution engine

## M5: Stabilization

Deliverables:

- examples
- package documentation
- version support matrix
- maintainer docs for regeneration and upgrades

Exit criteria:

- a new user can follow the README and complete the V1 workflow
- a maintainer can regenerate code and pass CI using only repository docs

## Definition Of Done

Version 1 is done when all of the following are true:

- the module has a documented compatibility policy
- core handle APIs are usable
- the builder is usable
- common dialect wrappers exist
- pass execution exists
- execution-engine support exists
- generated code is reproducible and checked in
- the test suite proves end-to-end workflows
- ownership and error rules are clearly documented

