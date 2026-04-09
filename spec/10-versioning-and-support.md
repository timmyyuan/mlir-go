# 10. Versioning And Support

## Versioning Rule

One release line supports one LLVM/MLIR major version.

This project must not track upstream `main` or `HEAD` as its public compatibility target.
That strategy is acceptable for experiments, but not for a Go module that users may pin in production or tooling.

## Initial Baseline

The initial implementation should target one stable LLVM/MLIR major release, for example `23.x`.
The exact baseline can be adjusted before implementation starts, but it must be fixed explicitly and documented.

## Upgrade Policy

When LLVM/MLIR moves to a new major version:

- the repository may create a new release line, branch, or tag series
- generated code is regenerated against the new baseline
- CI images and tests move with that baseline
- older release lines are not required to compile unchanged against the new upstream release

## Supported Platforms

The first supported CI targets should be:

- `linux/amd64`
- `darwin/arm64`

Additional platforms may be added later, but support should follow passing CI rather than aspiration.

## Toolchain Expectations

The project should document these prerequisites:

- a matching `llvm-config`
- MLIR C headers
- MLIR C libraries and dependent shared libraries
- a Go toolchain version declared in `go.mod`

The build system should fail early with a readable error when the MLIR toolchain is missing or mismatched.

