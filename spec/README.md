# MLIR Go Binding Spec

This repository uses a progressive spec set instead of a single monolithic document.
Read the files in order.

## Reading Order

1. `00-vision.md`
2. `10-versioning-and-support.md`
3. `20-runtime-model.md`
4. `30-core-api.md`
5. `40-builder-and-dialects.md`
6. `50-codegen.md`
7. `60-testing.md`
8. `70-roadmap.md`

## Purpose

The goal of this spec is to define a practical path to a usable Go binding for MLIR.
It is informed by the design of [`google/mlir-hs`](https://github.com/google/mlir-hs), especially its layering between native bindings, builder-style APIs, generated dialect wrappers, and end-to-end tests.

This spec intentionally avoids machine-specific details such as local paths, local toolchain locations, or workstation assumptions.

