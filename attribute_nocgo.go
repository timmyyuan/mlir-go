//go:build !cgo

package mlir

import "fmt"

type Attribute struct{}

func ParseAttribute(*Context, string) (Attribute, error) {
	return Attribute{}, fmt.Errorf("mlir: cgo is required")
}

func (a Attribute) IsNull() bool {
	return true
}

func (a Attribute) String() string {
	return ""
}
