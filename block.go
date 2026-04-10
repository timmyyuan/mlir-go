//go:build cgo

package mlir

import (
	"fmt"

	"github.com/timmyyuan/mlir-go/internal/capi"
)

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

// ParentOperation returns the closest surrounding operation.
func (b Block) ParentOperation() Operation {
	if b.IsNull() {
		return Operation{}
	}
	return Operation{raw: capi.BlockGetParentOperation(b.raw)}
}

// ParentRegion returns the containing region.
func (b Block) ParentRegion() Region {
	if b.IsNull() {
		return Region{}
	}
	return Region{raw: capi.BlockGetParentRegion(b.raw)}
}

// Terminator returns the block terminator, or null if missing.
func (b Block) Terminator() Operation {
	if b.IsNull() {
		return Operation{}
	}
	return Operation{raw: capi.BlockGetTerminator(b.raw)}
}

// AppendOwnedOperation transfers ownership of a detached operation into b.
func (b Block) AppendOwnedOperation(op *OwnedOperation) (Operation, error) {
	if b.IsNull() {
		return Operation{}, fmt.Errorf("mlir: null block")
	}
	rawOp, err := op.takeRaw()
	if err != nil {
		return Operation{}, err
	}
	capi.BlockAppendOwnedOperation(b.raw, rawOp)
	return Operation{raw: rawOp}, nil
}

// InsertOwnedOperationBefore inserts a detached operation before ref.
func (b Block) InsertOwnedOperationBefore(ref Operation, op *OwnedOperation) (Operation, error) {
	if b.IsNull() {
		return Operation{}, fmt.Errorf("mlir: null block")
	}
	if ref.IsNull() {
		return Operation{}, fmt.Errorf("mlir: null reference operation")
	}
	rawOp, err := op.takeRaw()
	if err != nil {
		return Operation{}, err
	}
	capi.BlockInsertOwnedOperationBefore(b.raw, ref.raw, rawOp)
	return Operation{raw: rawOp}, nil
}

// AddArgument appends an argument to the block.
func (b Block) AddArgument(typ Type, loc Location) (Value, error) {
	if b.IsNull() {
		return Value{}, fmt.Errorf("mlir: null block")
	}
	if typ.IsNull() {
		return Value{}, fmt.Errorf("mlir: null type")
	}
	if loc.IsNull() {
		return Value{}, fmt.Errorf("mlir: null location")
	}
	return Value{raw: capi.BlockAddArgument(b.raw, typ.raw, loc.raw)}, nil
}
