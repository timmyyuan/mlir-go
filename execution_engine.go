//go:build cgo

package mlir

import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/timmyyuan/mlir-go/internal/capi"
)

// ExecutionEngine owns an MLIR execution engine instance.
type ExecutionEngine struct {
	mu     sync.Mutex
	raw    capi.ExecutionEngine
	closed bool
}

func NewExecutionEngine(module *Module, optLevel int) (*ExecutionEngine, error) {
	rawModule, err := module.requireRaw()
	if err != nil {
		return nil, err
	}
	engine := capi.ExecutionEngineCreate(rawModule, optLevel)
	if capi.ExecutionEngineIsNull(engine) {
		return nil, fmt.Errorf("mlir: failed to create execution engine")
	}
	return &ExecutionEngine{raw: engine}, nil
}

func (e *ExecutionEngine) Close() error {
	if e == nil {
		return nil
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return nil
	}
	capi.ExecutionEngineDestroy(e.raw)
	e.closed = true
	return nil
}

// InvokePacked calls a function through the packed execution engine ABI.
// args contains pointers to argument storage; result points to result storage.
func (e *ExecutionEngine) InvokePacked(name string, result unsafe.Pointer, args ...unsafe.Pointer) error {
	raw, err := e.requireRaw()
	if err != nil {
		return err
	}

	packed := make([]unsafe.Pointer, 0, len(args)+1)
	packed = append(packed, args...)
	if result != nil {
		packed = append(packed, result)
	}
	if !capi.ExecutionEngineInvokePacked(raw, name, packed) {
		return fmt.Errorf("mlir: execution failed for %q", name)
	}
	return nil
}

func (e *ExecutionEngine) requireRaw() (capi.ExecutionEngine, error) {
	if e == nil {
		return capi.ExecutionEngine{}, ErrClosed
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed || capi.ExecutionEngineIsNull(e.raw) {
		return capi.ExecutionEngine{}, ErrClosed
	}
	return e.raw, nil
}
