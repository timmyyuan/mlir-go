//go:build cgo

package mlir

import (
	"fmt"

	"github.com/timmyyuan/mlir-go/internal/capi"
)

// Attribute is a borrowed handle to an MLIR attribute.
type Attribute struct {
	raw capi.Attribute
}

// ParseAttribute parses an MLIR attribute within the given context.
func ParseAttribute(ctx *Context, asm string) (Attribute, error) {
	raw, err := ctx.requireRaw()
	if err != nil {
		return Attribute{}, err
	}

	attr := capi.AttributeParseGet(raw, asm)
	if capi.AttributeIsNull(attr) {
		return Attribute{}, fmt.Errorf("mlir: failed to parse attribute")
	}
	return Attribute{raw: attr}, nil
}

// IsNull reports whether the attribute handle is null.
func (a Attribute) IsNull() bool {
	return capi.AttributeIsNull(a.raw)
}

// String returns the textual assembly form of the attribute.
func (a Attribute) String() string {
	if a.IsNull() {
		return ""
	}
	return capi.AttributeToString(a.raw)
}
