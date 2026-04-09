//go:build cgo

package mlir

import (
	"fmt"

	"github.com/timmyyuan/mlir-go/internal/capi"
)

// Operation is a borrowed handle to an MLIR operation.
type Operation struct {
	raw capi.Operation
}

// IsNull reports whether the operation handle is null.
func (op Operation) IsNull() bool {
	return capi.OperationIsNull(op.raw)
}

// String returns the textual assembly form of the operation.
func (op Operation) String() string {
	if op.IsNull() {
		return ""
	}
	return capi.OperationToString(op.raw)
}

// Verify checks whether the operation verifies successfully.
func (op Operation) Verify() error {
	if op.IsNull() {
		return fmt.Errorf("mlir: null operation")
	}
	if !capi.OperationVerify(op.raw) {
		return fmt.Errorf("mlir: operation verification failed")
	}
	return nil
}

// Name returns the operation name, such as "func.func" or "arith.addi".
func (op Operation) Name() string {
	if op.IsNull() {
		return ""
	}
	return capi.IdentifierToString(capi.OperationGetName(op.raw))
}

// Location returns the operation location.
func (op Operation) Location() Location {
	if op.IsNull() {
		return Location{}
	}
	return Location{raw: capi.OperationGetLocation(op.raw)}
}

// Regions returns the regions attached to the operation in source order.
func (op Operation) Regions() []Region {
	if op.IsNull() {
		return nil
	}

	n := capi.OperationGetNumRegions(op.raw)
	regions := make([]Region, 0, n)
	for i := 0; i < n; i++ {
		regions = append(regions, Region{raw: capi.OperationGetRegion(op.raw, i)})
	}
	return regions
}

// Results returns the operation results in source order.
func (op Operation) Results() []Value {
	if op.IsNull() {
		return nil
	}

	n := capi.OperationGetNumResults(op.raw)
	results := make([]Value, 0, n)
	for i := 0; i < n; i++ {
		results = append(results, Value{raw: capi.OperationGetResult(op.raw, i)})
	}
	return results
}
