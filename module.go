//go:build cgo

package mlir

import (
	"fmt"
	"sync"

	"github.com/timmyyuan/mlir-go/internal/capi"
)

// Module owns an MLIR module allocated through the C API.
type Module struct {
	mu     sync.Mutex
	raw    capi.Module
	closed bool
}

// ParseModule parses MLIR assembly into a module owned by the caller.
func ParseModule(ctx *Context, asm string) (*Module, error) {
	rawCtx, err := ctx.requireRaw()
	if err != nil {
		return nil, err
	}

	raw := capi.ModuleCreateParse(rawCtx, asm)
	if capi.ModuleIsNull(raw) {
		return nil, fmt.Errorf("mlir: failed to parse module")
	}
	return &Module{raw: raw}, nil
}

// CreateEmptyModule creates an empty module at the given location.
func CreateEmptyModule(loc Location) (*Module, error) {
	if loc.IsNull() {
		return nil, fmt.Errorf("mlir: null location")
	}

	raw := capi.ModuleCreateEmpty(loc.raw)
	if capi.ModuleIsNull(raw) {
		return nil, fmt.Errorf("mlir: failed to create empty module")
	}
	return &Module{raw: raw}, nil
}

// Close destroys the underlying MLIR module.
//
// Close is idempotent.
func (m *Module) Close() error {
	if m == nil {
		return nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return nil
	}

	capi.ModuleDestroy(m.raw)
	m.raw = capi.NullModule()
	m.closed = true
	return nil
}

// Operation returns the borrowed top-level operation view of the module.
func (m *Module) Operation() Operation {
	raw, err := m.requireRaw()
	if err != nil {
		return Operation{}
	}
	return Operation{raw: capi.ModuleGetOperation(raw)}
}

// Body returns the module body block.
func (m *Module) Body() Block {
	raw, err := m.requireRaw()
	if err != nil {
		return Block{}
	}
	return Block{raw: capi.ModuleGetBody(raw)}
}

// Verify checks whether the module verifies successfully.
func (m *Module) Verify() error {
	return m.Operation().Verify()
}

// String returns the textual assembly form of the module.
func (m *Module) String() string {
	return m.Operation().String()
}

func (m *Module) requireRaw() (capi.Module, error) {
	if m == nil {
		return capi.NullModule(), ErrClosed
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed || capi.ModuleIsNull(m.raw) {
		return capi.NullModule(), ErrClosed
	}
	return m.raw, nil
}
