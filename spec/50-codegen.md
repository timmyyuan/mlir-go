# 50. Code Generation

## Why Code Generation Is Required

Common dialect surfaces are too large to maintain by hand.
If this repository relies only on manual wrappers, coverage will remain narrow and fragile.

The project therefore needs a generator for dialect wrappers.

## Generator Responsibilities

The generator should:

- read MLIR TableGen definitions for selected dialects
- generate Go wrappers for operations, types, and attributes where feasible
- generate smoke tests or golden tests for generated surfaces
- produce a machine-readable skip report for unsupported features

## Generator Non-Responsibilities

The generator should not:

- run implicitly during `go build`
- hide unsupported cases silently
- attempt to encode every complex semantic rule in one pass

## Checked-In Generated Code

Generated `.go` files should be committed to the repository.

Reasons:

- users should not need the full generator toolchain to consume the module
- builds become more deterministic
- diffs from MLIR upgrades become reviewable
- CI can verify freshness of generated output

## Dialect Manifest

The repository should maintain an explicit manifest that declares:

- dialect name
- source `.td` file or source group
- Go import path
- package name
- prefix-stripping rules
- whether test generation is enabled

That manifest is the contract for supported generated surfaces.

## Unsupported Operations

When the generator cannot emit a usable wrapper, it must record:

- dialect
- operation name
- reason for skipping
- whether the limitation comes from generator logic or public API shape

Unsupported operations may be skipped in V1, but they must be visible and auditable.

