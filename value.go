//go:build cgo

package mlir

import "github.com/timmyyuan/mlir-go/internal/capi"

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
