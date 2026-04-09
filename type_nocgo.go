//go:build !cgo

package mlir

import "fmt"

type Type struct{}

func ParseType(*Context, string) (Type, error) {
	return Type{}, fmt.Errorf("mlir: cgo is required")
}

func (t Type) IsNull() bool {
	return true
}

func (t Type) Equal(Type) bool {
	return false
}

func (t Type) String() string {
	return ""
}
