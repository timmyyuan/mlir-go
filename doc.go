// Package mlir provides Go bindings built on top of the MLIR C API.
//
// The repository is at an early stage. The public surface currently combines
// explicit resource ownership and direct access to native MLIR handles with a
// first builder package, dialect helpers, pass execution, and JIT support.
package mlir
