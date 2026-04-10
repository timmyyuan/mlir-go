//go:build cgo

package mlir

import (
	"fmt"
	"sync"

	"github.com/timmyyuan/mlir-go/internal/capi"
)

// RegisterAllPasses registers all upstream MLIR passes.
func RegisterAllPasses() {
	capi.RegisterAllPasses()
}

// RegisterAllLLVMTranslations registers upstream LLVM IR translations in the context.
func (c *Context) RegisterAllLLVMTranslations() error {
	raw, err := c.requireRaw()
	if err != nil {
		return err
	}
	capi.RegisterAllLLVMTranslations(raw)
	return nil
}

// PassManager owns an MLIR pass manager.
type PassManager struct {
	mu     sync.Mutex
	raw    capi.PassManager
	closed bool
}

func NewPassManager(ctx *Context) (*PassManager, error) {
	rawCtx, err := ctx.requireRaw()
	if err != nil {
		return nil, err
	}
	pm := capi.PassManagerCreate(rawCtx)
	if capi.PassManagerIsNull(pm) {
		return nil, fmt.Errorf("mlir: failed to create pass manager")
	}
	return &PassManager{raw: pm}, nil
}

func (pm *PassManager) Close() error {
	if pm == nil {
		return nil
	}
	pm.mu.Lock()
	defer pm.mu.Unlock()
	if pm.closed {
		return nil
	}
	capi.PassManagerDestroy(pm.raw)
	pm.closed = true
	return nil
}

func (pm *PassManager) ParsePipeline(pipeline string) error {
	raw, err := pm.requireRaw()
	if err != nil {
		return err
	}
	ok, msg := capi.ParsePassPipeline(capi.PassManagerGetAsOpPassManager(raw), pipeline)
	if !ok {
		if msg == "" {
			msg = "failed to parse pipeline"
		}
		return fmt.Errorf("mlir: %s", msg)
	}
	return nil
}

func (pm *PassManager) String() string {
	raw, err := pm.requireRaw()
	if err != nil {
		return ""
	}
	return capi.PrintPassPipeline(capi.PassManagerGetAsOpPassManager(raw))
}

func (pm *PassManager) Run(op Operation) error {
	raw, err := pm.requireRaw()
	if err != nil {
		return err
	}
	if op.IsNull() {
		return fmt.Errorf("mlir: null operation")
	}
	if !capi.PassManagerRunOnOp(raw, op.raw) {
		return fmt.Errorf("mlir: pass manager execution failed")
	}
	return nil
}

func (pm *PassManager) requireRaw() (capi.PassManager, error) {
	if pm == nil {
		return capi.PassManager{}, ErrClosed
	}
	pm.mu.Lock()
	defer pm.mu.Unlock()
	if pm.closed || capi.PassManagerIsNull(pm.raw) {
		return capi.PassManager{}, ErrClosed
	}
	return pm.raw, nil
}
