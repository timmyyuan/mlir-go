//go:build !cgo

package mlir

import "fmt"

type Identifier struct{}

func InternIdentifier(*Context, string) (Identifier, error) {
	return Identifier{}, fmt.Errorf("mlir: cgo is required")
}

func (id Identifier) String() string {
	return ""
}

type NamedAttribute struct{}

func NamedAttributeByName(*Context, string, Attribute) (NamedAttribute, error) {
	return NamedAttribute{}, fmt.Errorf("mlir: cgo is required")
}

func NewNamedAttribute(Identifier, Attribute) NamedAttribute {
	return NamedAttribute{}
}

func (na NamedAttribute) Name() Identifier {
	return Identifier{}
}

func (na NamedAttribute) Attribute() Attribute {
	return Attribute{}
}
