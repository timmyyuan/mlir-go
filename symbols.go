//go:build cgo

package mlir

import (
	"fmt"
	"sync"

	"github.com/timmyyuan/mlir-go/internal/capi"
)

// SymbolTable owns an MLIR symbol table helper for a symbol table operation.
type SymbolTable struct {
	mu     sync.Mutex
	raw    capi.SymbolTable
	closed bool
}

// SymbolAttributeName returns the canonical attribute name used for symbols.
func SymbolAttributeName() string {
	return capi.SymbolTableGetSymbolAttributeName()
}

// VisibilityAttributeName returns the canonical attribute name used for symbol visibility.
func VisibilityAttributeName() string {
	return capi.SymbolTableGetVisibilityAttributeName()
}

// NewSymbolTable creates a symbol table helper for op.
func NewSymbolTable(op Operation) (*SymbolTable, error) {
	if op.IsNull() {
		return nil, fmt.Errorf("mlir: null operation")
	}
	raw := capi.SymbolTableCreate(op.raw)
	if capi.SymbolTableIsNull(raw) {
		return nil, fmt.Errorf("mlir: operation %q is not a symbol table", op.Name())
	}
	return &SymbolTable{raw: raw}, nil
}

// SymbolTable creates a symbol table helper for the module operation.
func (m *Module) SymbolTable() (*SymbolTable, error) {
	if m == nil {
		return nil, ErrClosed
	}
	return NewSymbolTable(m.Operation())
}

// Close destroys the underlying symbol table helper.
//
// Close is idempotent.
func (s *SymbolTable) Close() error {
	if s == nil {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}

	capi.SymbolTableDestroy(s.raw)
	s.raw = capi.NullSymbolTable()
	s.closed = true
	return nil
}

// Lookup resolves name in the symbol table and returns a null operation if missing.
func (s *SymbolTable) Lookup(name string) (Operation, error) {
	if name == "" {
		return Operation{}, fmt.Errorf("mlir: symbol name is required")
	}
	raw, err := s.requireRaw()
	if err != nil {
		return Operation{}, err
	}
	return Operation{raw: capi.SymbolTableLookup(raw, name)}, nil
}

// Insert updates the symbol table with op and returns the final symbol name attribute.
//
// The caller is responsible for placing op in the symbol table body separately.
func (s *SymbolTable) Insert(op Operation) (Attribute, error) {
	if op.IsNull() {
		return Attribute{}, fmt.Errorf("mlir: null operation")
	}
	raw, err := s.requireRaw()
	if err != nil {
		return Attribute{}, err
	}
	attr := capi.SymbolTableInsert(raw, op.raw)
	if capi.AttributeIsNull(attr) {
		return Attribute{}, fmt.Errorf("mlir: failed to insert symbol")
	}
	return Attribute{raw: attr}, nil
}

// Erase removes op from the symbol table and erases it from the IR.
func (s *SymbolTable) Erase(op Operation) error {
	if op.IsNull() {
		return fmt.Errorf("mlir: null operation")
	}
	raw, err := s.requireRaw()
	if err != nil {
		return err
	}
	capi.SymbolTableErase(raw, op.raw)
	return nil
}

// ReplaceAllSymbolUses rewrites symbol references nested within from.
func ReplaceAllSymbolUses(oldSymbol, newSymbol string, from Operation) error {
	if oldSymbol == "" || newSymbol == "" {
		return fmt.Errorf("mlir: symbol names are required")
	}
	if from.IsNull() {
		return fmt.Errorf("mlir: null operation")
	}
	if !capi.SymbolTableReplaceAllSymbolUses(oldSymbol, newSymbol, from.raw) {
		return fmt.Errorf("mlir: failed to replace symbol uses from %q to %q", oldSymbol, newSymbol)
	}
	return nil
}

// WalkSymbolTables visits the symbol table operations nested within from.
func WalkSymbolTables(from Operation, allSymUsesVisible bool, callback func(Operation, bool)) error {
	if from.IsNull() {
		return fmt.Errorf("mlir: null operation")
	}
	if callback == nil {
		return fmt.Errorf("mlir: symbol table walk callback is required")
	}
	capi.SymbolTableWalkSymbolTables(from.raw, allSymUsesVisible, func(op capi.Operation, allVisible bool) {
		callback(Operation{raw: op}, allVisible)
	})
	return nil
}

// SymbolName returns the string symbol name attached to op.
func (op Operation) SymbolName() (string, bool) {
	return op.stringAttributeValue(SymbolAttributeName())
}

// Visibility returns the string visibility attached to op.
func (op Operation) Visibility() (string, bool) {
	return op.stringAttributeValue(VisibilityAttributeName())
}

func (op Operation) stringAttributeValue(name string) (string, bool) {
	if op.IsNull() || name == "" {
		return "", false
	}
	attr := op.Attribute(name)
	if attr.IsNull() || !attr.IsString() {
		return "", false
	}
	value, err := attr.StringValue()
	if err != nil {
		return "", false
	}
	return value, true
}

func (s *SymbolTable) requireRaw() (capi.SymbolTable, error) {
	if s == nil {
		return capi.NullSymbolTable(), ErrClosed
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed || capi.SymbolTableIsNull(s.raw) {
		return capi.NullSymbolTable(), ErrClosed
	}
	return s.raw, nil
}
