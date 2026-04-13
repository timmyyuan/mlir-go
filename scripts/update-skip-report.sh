#!/usr/bin/env bash

set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repo_root"

go run ./cmd/mlir-go-tblgen \
  -mode emit-skip-report \
  -o internal/gen/skip_report.json
