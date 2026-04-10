//go:build cgo

package mlir

import (
	"fmt"

	"github.com/timmyyuan/mlir-go/internal/capi"
)

// Identifier is a borrowed handle to an MLIR identifier.
type Identifier struct {
	raw capi.Identifier
}

// InternIdentifier returns an identifier owned by the context.
func InternIdentifier(ctx *Context, name string) (Identifier, error) {
	raw, err := ctx.requireRaw()
	if err != nil {
		return Identifier{}, err
	}
	id := capi.IdentifierGet(raw, name)
	if capi.IdentifierToString(id) == "" && name != "" {
		return Identifier{}, fmt.Errorf("mlir: failed to intern identifier")
	}
	return Identifier{raw: id}, nil
}

// String returns the string value of the identifier.
func (id Identifier) String() string {
	return capi.IdentifierToString(id.raw)
}

// NamedAttribute is a name-attribute pair.
type NamedAttribute struct {
	raw capi.NamedAttribute
}

// NamedAttributeByName interns name in ctx and constructs a named attribute.
func NamedAttributeByName(ctx *Context, name string, attr Attribute) (NamedAttribute, error) {
	id, err := InternIdentifier(ctx, name)
	if err != nil {
		return NamedAttribute{}, err
	}
	return NewNamedAttribute(id, attr), nil
}

// NewNamedAttribute constructs a named attribute pair.
func NewNamedAttribute(name Identifier, attr Attribute) NamedAttribute {
	return NamedAttribute{raw: capi.NamedAttributeGet(name.raw, attr.raw)}
}

// Name returns the attribute name.
func (na NamedAttribute) Name() Identifier {
	return Identifier{raw: capi.NamedAttributeName(na.raw)}
}

// Attribute returns the attribute payload.
func (na NamedAttribute) Attribute() Attribute {
	return Attribute{raw: capi.NamedAttributeAttribute(na.raw)}
}
