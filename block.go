//go:build cgo

package mlir

import "github.com/timmyyuan/mlir-go/internal/capi"

// Block is a borrowed handle to an MLIR block.
type Block struct {
	raw capi.Block
}

// IsNull reports whether the block handle is null.
func (b Block) IsNull() bool {
	return capi.BlockIsNull(b.raw)
}

// String returns the textual assembly form of the block.
func (b Block) String() string {
	if b.IsNull() {
		return ""
	}
	return capi.BlockToString(b.raw)
}

// Arguments returns the block arguments in source order.
func (b Block) Arguments() []Value {
	if b.IsNull() {
		return nil
	}

	n := capi.BlockGetNumArguments(b.raw)
	args := make([]Value, 0, n)
	for i := 0; i < n; i++ {
		args = append(args, Value{raw: capi.BlockGetArgument(b.raw, i)})
	}
	return args
}

// Operations returns the operations contained in the block in source order.
func (b Block) Operations() []Operation {
	if b.IsNull() {
		return nil
	}

	var ops []Operation
	for raw := capi.BlockGetFirstOperation(b.raw); !capi.OperationIsNull(raw); raw = capi.OperationGetNextInBlock(raw) {
		ops = append(ops, Operation{raw: raw})
	}
	return ops
}
