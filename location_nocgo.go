//go:build !cgo

package mlir

import "fmt"

type Location struct{}

func UnknownLocation(*Context) (Location, error) {
	return Location{}, fmt.Errorf("mlir: cgo is required")
}

func FileLineColLocation(*Context, string, uint, uint) (Location, error) {
	return Location{}, fmt.Errorf("mlir: cgo is required")
}

func (loc Location) IsNull() bool {
	return true
}

func (loc Location) String() string {
	return ""
}
