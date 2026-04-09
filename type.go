//go:build cgo

package mlir

import (
	"fmt"

	"github.com/timmyyuan/mlir-go/internal/capi"
)

// Type is a borrowed handle to an MLIR type.
type Type struct {
	raw capi.Type
}

// ParseType parses an MLIR type within the given context.
func ParseType(ctx *Context, asm string) (Type, error) {
	raw, err := ctx.requireRaw()
	if err != nil {
		return Type{}, err
	}

	typ := capi.TypeParseGet(raw, asm)
	if capi.TypeIsNull(typ) {
		return Type{}, fmt.Errorf("mlir: failed to parse type")
	}
	return Type{raw: typ}, nil
}

// IsNull reports whether the type handle is null.
func (t Type) IsNull() bool {
	return capi.TypeIsNull(t.raw)
}

// Equal reports whether two type handles refer to the same MLIR type.
func (t Type) Equal(other Type) bool {
	if t.IsNull() || other.IsNull() {
		return false
	}
	return capi.TypeEqual(t.raw, other.raw)
}

// String returns the textual assembly form of the type.
func (t Type) String() string {
	if t.IsNull() {
		return ""
	}
	return capi.TypeToString(t.raw)
}
