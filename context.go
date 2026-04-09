//go:build cgo

package mlir

import (
	"fmt"
	"sync"

	"github.com/timmyyuan/mlir-go/internal/capi"
)

// Context owns an MLIR context allocated through the C API.
//
// A Context must be closed explicitly. Borrowed handles derived from the
// context become invalid after Close returns.
type Context struct {
	mu     sync.Mutex
	raw    capi.Context
	closed bool
}

// NewContext creates a new MLIR context.
func NewContext() (*Context, error) {
	raw := capi.ContextCreate()
	if capi.ContextIsNull(raw) {
		return nil, fmt.Errorf("mlir: failed to create context")
	}
	return &Context{raw: raw}, nil
}

// Close destroys the underlying MLIR context.
//
// Close is idempotent.
func (c *Context) Close() error {
	if c == nil {
		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	capi.ContextDestroy(c.raw)
	c.raw = capi.NullContext()
	c.closed = true
	return nil
}

// AllowUnregisteredDialects sets whether the context accepts operations from
// dialects that are not registered in the context.
func (c *Context) AllowUnregisteredDialects(allow bool) error {
	raw, err := c.requireRaw()
	if err != nil {
		return err
	}
	capi.ContextSetAllowUnregisteredDialects(raw, allow)
	return nil
}

// RegisterAllDialects appends the upstream "register everything" dialect
// registry and eagerly loads all available dialects into the context.
func (c *Context) RegisterAllDialects() error {
	raw, err := c.requireRaw()
	if err != nil {
		return err
	}

	registry := capi.DialectRegistryCreate()
	if capi.DialectRegistryIsNull(registry) {
		return fmt.Errorf("mlir: failed to create dialect registry")
	}
	defer capi.DialectRegistryDestroy(registry)

	capi.RegisterAllDialects(registry)
	capi.ContextAppendDialectRegistry(raw, registry)
	capi.ContextLoadAllAvailableDialects(raw)
	return nil
}

// NumLoadedDialects returns the number of dialects currently loaded in the
// context.
func (c *Context) NumLoadedDialects() (int, error) {
	raw, err := c.requireRaw()
	if err != nil {
		return 0, err
	}
	return capi.ContextGetNumLoadedDialects(raw), nil
}

func (c *Context) requireRaw() (capi.Context, error) {
	if c == nil {
		return capi.NullContext(), ErrClosed
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed || capi.ContextIsNull(c.raw) {
		return capi.NullContext(), ErrClosed
	}
	return c.raw, nil
}
