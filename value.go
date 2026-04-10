//go:build cgo

package mlir

import (
	"fmt"

	"github.com/timmyyuan/mlir-go/internal/capi"
)

// Value is a borrowed handle to an MLIR SSA value.
type Value struct {
	raw capi.Value
}

// IsNull reports whether the value handle is null.
func (v Value) IsNull() bool {
	return capi.ValueIsNull(v.raw)
}

// IsBlockArgument reports whether the value is defined as a block argument.
func (v Value) IsBlockArgument() bool {
	if v.IsNull() {
		return false
	}
	return capi.ValueIsBlockArgument(v.raw)
}

// IsOpResult reports whether the value is defined as an operation result.
func (v Value) IsOpResult() bool {
	if v.IsNull() {
		return false
	}
	return capi.ValueIsOpResult(v.raw)
}

// Type returns the MLIR type of the value.
func (v Value) Type() Type {
	if v.IsNull() {
		return Type{}
	}
	return Type{raw: capi.ValueGetType(v.raw)}
}

// String returns the textual assembly form of the value.
func (v Value) String() string {
	if v.IsNull() {
		return ""
	}
	return capi.ValueToString(v.raw)
}

// Equal reports whether two values refer to the same SSA value.
func (v Value) Equal(other Value) bool {
	if v.IsNull() || other.IsNull() {
		return false
	}
	return capi.ValueEqual(v.raw, other.raw)
}

// OwnerBlock returns the defining block for a block argument, or null.
func (v Value) OwnerBlock() Block {
	if !v.IsBlockArgument() {
		return Block{}
	}
	return Block{raw: capi.BlockArgumentGetOwner(v.raw)}
}

// ArgumentNumber returns the position of a block argument.
func (v Value) ArgumentNumber() int {
	if !v.IsBlockArgument() {
		return -1
	}
	return capi.BlockArgumentGetArgNumber(v.raw)
}

// OwnerOperation returns the defining operation for an op result, or null.
func (v Value) OwnerOperation() Operation {
	if !v.IsOpResult() {
		return Operation{}
	}
	return Operation{raw: capi.OpResultGetOwner(v.raw)}
}

// ResultNumber returns the position of an op result.
func (v Value) ResultNumber() int {
	if !v.IsOpResult() {
		return -1
	}
	return capi.OpResultGetResultNumber(v.raw)
}

// SetType mutates the value type in-place.
func (v Value) SetType(typ Type) error {
	if v.IsNull() {
		return fmt.Errorf("mlir: null value")
	}
	if typ.IsNull() {
		return fmt.Errorf("mlir: null type")
	}
	capi.ValueSetType(v.raw, typ.raw)
	return nil
}
