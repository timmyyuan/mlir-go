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

// Block returns the containing block of the operation, or a null block.
func (op Operation) Block() Block {
	if op.IsNull() {
		return Block{}
	}
	return Block{raw: capi.OperationGetBlock(op.raw)}
}

// ParentOperation returns the closest enclosing operation, or null.
func (op Operation) ParentOperation() Operation {
	if op.IsNull() {
		return Operation{}
	}
	return Operation{raw: capi.OperationGetParentOperation(op.raw)}
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

// Operands returns the operation operands in source order.
func (op Operation) Operands() []Value {
	if op.IsNull() {
		return nil
	}

	n := capi.OperationGetNumOperands(op.raw)
	operands := make([]Value, 0, n)
	for i := 0; i < n; i++ {
		operands = append(operands, Value{raw: capi.OperationGetOperand(op.raw, i)})
	}
	return operands
}

// SetOperand rewrites a single operand in-place.
func (op Operation) SetOperand(pos int, value Value) error {
	if op.IsNull() {
		return fmt.Errorf("mlir: null operation")
	}
	if value.IsNull() {
		return fmt.Errorf("mlir: null operand")
	}
	capi.OperationSetOperand(op.raw, pos, value.raw)
	return nil
}

// SetOperands replaces the full operand list.
func (op Operation) SetOperands(values ...Value) error {
	if op.IsNull() {
		return fmt.Errorf("mlir: null operation")
	}
	rawValues := make([]capi.Value, len(values))
	for i, value := range values {
		rawValues[i] = value.raw
	}
	capi.OperationSetOperands(op.raw, rawValues)
	return nil
}

// Successors returns the successor blocks in source order.
func (op Operation) Successors() []Block {
	if op.IsNull() {
		return nil
	}

	n := capi.OperationGetNumSuccessors(op.raw)
	blocks := make([]Block, 0, n)
	for i := 0; i < n; i++ {
		blocks = append(blocks, Block{raw: capi.OperationGetSuccessor(op.raw, i)})
	}
	return blocks
}

// SetSuccessor rewrites a single successor in-place.
func (op Operation) SetSuccessor(pos int, block Block) error {
	if op.IsNull() {
		return fmt.Errorf("mlir: null operation")
	}
	if block.IsNull() {
		return fmt.Errorf("mlir: null successor")
	}
	capi.OperationSetSuccessor(op.raw, pos, block.raw)
	return nil
}

// Attributes returns the attached attributes in source order.
func (op Operation) Attributes() []NamedAttribute {
	if op.IsNull() {
		return nil
	}

	n := capi.OperationGetNumAttributes(op.raw)
	attrs := make([]NamedAttribute, 0, n)
	for i := 0; i < n; i++ {
		attrs = append(attrs, NamedAttribute{raw: capi.OperationGetAttribute(op.raw, i)})
	}
	return attrs
}

// Attribute returns the attribute attached to name, or null if missing.
func (op Operation) Attribute(name string) Attribute {
	if op.IsNull() {
		return Attribute{}
	}
	return Attribute{raw: capi.OperationGetAttributeByName(op.raw, name)}
}

// SetAttribute attaches or replaces an attribute by name.
func (op Operation) SetAttribute(name string, attr Attribute) error {
	if op.IsNull() {
		return fmt.Errorf("mlir: null operation")
	}
	if name == "" {
		return fmt.Errorf("mlir: attribute name is required")
	}
	if attr.IsNull() {
		return fmt.Errorf("mlir: null attribute")
	}
	capi.OperationSetAttributeByName(op.raw, name, attr.raw)
	return nil
}

// RemoveAttribute removes an attribute by name and reports whether it existed.
func (op Operation) RemoveAttribute(name string) (bool, error) {
	if op.IsNull() {
		return false, fmt.Errorf("mlir: null operation")
	}
	if name == "" {
		return false, fmt.Errorf("mlir: attribute name is required")
	}
	return capi.OperationRemoveAttributeByName(op.raw, name), nil
}
