//go:build !cgo

package mlir

import "fmt"

type SymbolTable struct{}

func SymbolAttributeName() string {
	return ""
}

func VisibilityAttributeName() string {
	return ""
}

func NewSymbolTable(Operation) (*SymbolTable, error) {
	return nil, fmt.Errorf("mlir: cgo is required")
}

func (m *Module) SymbolTable() (*SymbolTable, error) {
	return nil, fmt.Errorf("mlir: cgo is required")
}

func (s *SymbolTable) Close() error {
	return fmt.Errorf("mlir: cgo is required")
}

func (s *SymbolTable) Lookup(string) (Operation, error) {
	return Operation{}, fmt.Errorf("mlir: cgo is required")
}

func (s *SymbolTable) Insert(Operation) (Attribute, error) {
	return Attribute{}, fmt.Errorf("mlir: cgo is required")
}

func (s *SymbolTable) Erase(Operation) error {
	return fmt.Errorf("mlir: cgo is required")
}

func ReplaceAllSymbolUses(string, string, Operation) error {
	return fmt.Errorf("mlir: cgo is required")
}

func WalkSymbolTables(Operation, bool, func(Operation, bool)) error {
	return fmt.Errorf("mlir: cgo is required")
}

func (op Operation) SymbolName() (string, bool) {
	return "", false
}

func (op Operation) Visibility() (string, bool) {
	return "", false
}
