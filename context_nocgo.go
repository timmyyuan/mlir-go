//go:build !cgo

package mlir

import "fmt"

type Context struct{}

func NewContext() (*Context, error) {
	return nil, fmt.Errorf("mlir: cgo is required")
}

func (c *Context) Close() error {
	return fmt.Errorf("mlir: cgo is required")
}

func (c *Context) AllowUnregisteredDialects(bool) error {
	return fmt.Errorf("mlir: cgo is required")
}

func (c *Context) RegisterAllDialects() error {
	return fmt.Errorf("mlir: cgo is required")
}

func (c *Context) NumLoadedDialects() (int, error) {
	return 0, fmt.Errorf("mlir: cgo is required")
}
