# Changelog

This file records notable project changes with a short summary first and
details available on demand.

## 2026-05-10

### Downstream Consumer CI Smoke Test

Added a cgo-backed CI smoke test that proves an external Go module can consume
`github.com/timmyyuan/mlir-go` through the documented development environment.

<details>
<summary>Issue, implementation, and validation details</summary>

- Addresses GitHub issue #12.
- Added `scripts/test-downstream-consumer.sh`, which creates a temporary Go
  module outside the repository, imports the root module and builder helper
  packages, builds a minimal IR-generation program, and validates the emitted
  IR with `FileCheck`.
- Added `testdata/filecheck/downstream_consumer_module.mlir` as the textual IR
  expectation for the downstream smoke program.
- Wired the smoke test into the Dockerized cgo GitHub Actions job so it runs
  with the same `scripts/dev-env.sh` setup documented in `README.md`.
- Documented the downstream usage check in `README.md`.
- Local validation included shell syntax checking and non-cgo package tests.
  CI covers the Dockerized cgo path for the downstream smoke test.

</details>
