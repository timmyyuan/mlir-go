//go:build !cgo

package mlir

import "fmt"

func RegisterAllPasses() {}

func (c *Context) RegisterAllLLVMTranslations() error {
	return fmt.Errorf("mlir: cgo is required")
}

type PassManager struct{}

func NewPassManager(*Context) (*PassManager, error) { return nil, fmt.Errorf("mlir: cgo is required") }
func (pm *PassManager) Close() error                { return fmt.Errorf("mlir: cgo is required") }
func (pm *PassManager) ParsePipeline(string) error  { return fmt.Errorf("mlir: cgo is required") }
func (pm *PassManager) String() string              { return "" }
func (pm *PassManager) Run(Operation) error         { return fmt.Errorf("mlir: cgo is required") }
