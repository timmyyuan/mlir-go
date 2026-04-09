#!/usr/bin/env bash
set -euo pipefail

llvm_config="${LLVM_CONFIG:-llvm-config}"

if ! command -v "$llvm_config" >/dev/null 2>&1; then
  echo "missing llvm-config; set LLVM_CONFIG or add llvm-config to PATH" >&2
  exit 1
fi

include_dir="$("$llvm_config" --includedir)"
lib_dir="$("$llvm_config" --libdir)"
bin_dir="$("$llvm_config" --bindir)"
version="$("$llvm_config" --version)"
llvm_libs="$("$llvm_config" --libs --system-libs)"

printf 'export MLIRGO_LLVM_CONFIG=%q\n' "$llvm_config"
printf 'export MLIRGO_LLVM_VERSION=%q\n' "$version"
printf 'export PATH=%q:$PATH\n' "$bin_dir"
printf 'export CGO_CPPFLAGS=%q\n' "-I$include_dir"
printf 'export CGO_LDFLAGS=%q\n' "-L$lib_dir -Wl,-rpath,$lib_dir -lMLIR -lMLIRCAPIIR -lMLIRCAPIRegisterEverything $llvm_libs -lc++"
