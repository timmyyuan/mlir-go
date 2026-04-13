#!/usr/bin/env bash

set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repo_root"

go run ./cmd/mlir-go-tblgen \
  -mode emit-dialect \
  -dialect arith \
  -output-root .

go run ./cmd/mlir-go-tblgen \
  -mode emit-dialect \
  -dialect func \
  -output-root .

go run ./cmd/mlir-go-tblgen \
  -mode emit-dialect \
  -dialect cf \
  -output-root .

go run ./cmd/mlir-go-tblgen \
  -mode emit-dialect \
  -dialect memref \
  -output-root .

go run ./cmd/mlir-go-tblgen \
  -mode emit-dialect \
  -dialect tensor \
  -output-root .

go run ./cmd/mlir-go-tblgen \
  -mode emit-dialect \
  -dialect vector \
  -output-root .

go run ./cmd/mlir-go-tblgen \
  -mode emit-dialect \
  -dialect linalg \
  -output-root .

bash ./scripts/update-skip-report.sh
