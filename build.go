//go:build cgo

package mlir

import (
	"fmt"
	"sync"

	"github.com/timmyyuan/mlir-go/internal/capi"
)

// OwnedRegion owns a detached MLIR region until it is transferred.
type OwnedRegion struct {
	mu     sync.Mutex
	raw    capi.Region
	closed bool
}

func NewOwnedRegion() (*OwnedRegion, error) {
	raw := capi.RegionCreate()
	if capi.RegionIsNull(raw) {
		return nil, fmt.Errorf("mlir: failed to create region")
	}
	return &OwnedRegion{raw: raw}, nil
}

func (r *OwnedRegion) Close() error {
	if r == nil {
		return nil
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed {
		return nil
	}
	capi.RegionDestroy(r.raw)
	r.closed = true
	return nil
}

func (r *OwnedRegion) Region() Region {
	if r == nil {
		return Region{}
	}
	return Region{raw: r.raw}
}

func (r *OwnedRegion) AppendBlock(block *OwnedBlock) (Block, error) {
	rawRegion, err := r.requireRaw()
	if err != nil {
		return Block{}, err
	}
	rawBlock, err := block.takeRaw()
	if err != nil {
		return Block{}, err
	}
	capi.RegionAppendOwnedBlock(rawRegion, rawBlock)
	return Block{raw: rawBlock}, nil
}

func (r *OwnedRegion) requireRaw() (capi.Region, error) {
	if r == nil {
		return capi.Region{}, ErrClosed
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed || capi.RegionIsNull(r.raw) {
		return capi.Region{}, ErrClosed
	}
	return r.raw, nil
}

func (r *OwnedRegion) takeRaw() (capi.Region, error) {
	if r == nil {
		return capi.Region{}, ErrClosed
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.closed || capi.RegionIsNull(r.raw) {
		return capi.Region{}, ErrClosed
	}
	raw := r.raw
	r.raw = capi.Region{}
	r.closed = true
	return raw, nil
}

// OwnedBlock owns a detached MLIR block until it is transferred.
type OwnedBlock struct {
	mu     sync.Mutex
	raw    capi.Block
	closed bool
}

func NewOwnedBlock(argTypes []Type, argLocs []Location) (*OwnedBlock, error) {
	if len(argLocs) != 0 && len(argLocs) != len(argTypes) {
		return nil, fmt.Errorf("mlir: argument type/location length mismatch")
	}
	rawTypes := make([]capi.Type, len(argTypes))
	for i, v := range argTypes {
		rawTypes[i] = v.raw
	}
	rawLocs := make([]capi.Location, len(argLocs))
	for i, v := range argLocs {
		rawLocs[i] = v.raw
	}
	raw := capi.BlockCreate(rawTypes, rawLocs)
	if capi.BlockIsNull(raw) {
		return nil, fmt.Errorf("mlir: failed to create block")
	}
	return &OwnedBlock{raw: raw}, nil
}

func (b *OwnedBlock) Close() error {
	if b == nil {
		return nil
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return nil
	}
	capi.BlockDestroy(b.raw)
	b.closed = true
	return nil
}

func (b *OwnedBlock) Block() Block {
	if b == nil {
		return Block{}
	}
	return Block{raw: b.raw}
}

func (b *OwnedBlock) AppendOperation(op *OwnedOperation) (Operation, error) {
	rawBlock, err := b.requireRaw()
	if err != nil {
		return Operation{}, err
	}
	rawOp, err := op.takeRaw()
	if err != nil {
		return Operation{}, err
	}
	capi.BlockAppendOwnedOperation(rawBlock, rawOp)
	return Operation{raw: rawOp}, nil
}

func (b *OwnedBlock) InsertOperationBefore(ref Operation, op *OwnedOperation) (Operation, error) {
	rawBlock, err := b.requireRaw()
	if err != nil {
		return Operation{}, err
	}
	rawOp, err := op.takeRaw()
	if err != nil {
		return Operation{}, err
	}
	capi.BlockInsertOwnedOperationBefore(rawBlock, ref.raw, rawOp)
	return Operation{raw: rawOp}, nil
}

func (b *OwnedBlock) AddArgument(typ Type, loc Location) (Value, error) {
	rawBlock, err := b.requireRaw()
	if err != nil {
		return Value{}, err
	}
	return Value{raw: capi.BlockAddArgument(rawBlock, typ.raw, loc.raw)}, nil
}

func (b *OwnedBlock) requireRaw() (capi.Block, error) {
	if b == nil {
		return capi.Block{}, ErrClosed
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed || capi.BlockIsNull(b.raw) {
		return capi.Block{}, ErrClosed
	}
	return b.raw, nil
}

func (b *OwnedBlock) takeRaw() (capi.Block, error) {
	if b == nil {
		return capi.Block{}, ErrClosed
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed || capi.BlockIsNull(b.raw) {
		return capi.Block{}, ErrClosed
	}
	raw := b.raw
	b.raw = capi.Block{}
	b.closed = true
	return raw, nil
}

// OwnedOperation owns a detached MLIR operation until it is transferred.
type OwnedOperation struct {
	mu     sync.Mutex
	raw    capi.Operation
	closed bool
}

func (o *OwnedOperation) Close() error {
	if o == nil {
		return nil
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.closed {
		return nil
	}
	capi.OperationDestroy(o.raw)
	o.closed = true
	return nil
}

func (o *OwnedOperation) Operation() Operation {
	if o == nil {
		return Operation{}
	}
	return Operation{raw: o.raw}
}

func (o *OwnedOperation) takeRaw() (capi.Operation, error) {
	if o == nil {
		return capi.Operation{}, ErrClosed
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.closed || capi.OperationIsNull(o.raw) {
		return capi.Operation{}, ErrClosed
	}
	raw := o.raw
	o.raw = capi.Operation{}
	o.closed = true
	return raw, nil
}

// OperationState describes a generic operation under construction.
type OperationState struct {
	Name             string
	Location         Location
	Results          []Type
	Operands         []Value
	Successors       []Block
	Attributes       []NamedAttribute
	OwnedRegions     []*OwnedRegion
	InferResultTypes bool
}

func NewOperationState(name string, loc Location) *OperationState {
	return &OperationState{Name: name, Location: loc}
}

func (s *OperationState) AddResults(results ...Type) {
	s.Results = append(s.Results, results...)
}

func (s *OperationState) AddOperands(values ...Value) {
	s.Operands = append(s.Operands, values...)
}

func (s *OperationState) AddSuccessors(blocks ...Block) {
	s.Successors = append(s.Successors, blocks...)
}

func (s *OperationState) AddAttributes(attrs ...NamedAttribute) {
	s.Attributes = append(s.Attributes, attrs...)
}

func (s *OperationState) AddOwnedRegions(regions ...*OwnedRegion) {
	s.OwnedRegions = append(s.OwnedRegions, regions...)
}

func (s *OperationState) EnableResultTypeInference() {
	s.InferResultTypes = true
}

func CreateOperation(state *OperationState) (*OwnedOperation, error) {
	if state == nil {
		return nil, fmt.Errorf("mlir: nil operation state")
	}
	if state.Name == "" {
		return nil, fmt.Errorf("mlir: operation name is required")
	}
	if state.Location.IsNull() {
		return nil, fmt.Errorf("mlir: operation location is required")
	}

	rawResults := make([]capi.Type, len(state.Results))
	for i, v := range state.Results {
		rawResults[i] = v.raw
	}
	rawOperands := make([]capi.Value, len(state.Operands))
	for i, v := range state.Operands {
		rawOperands[i] = v.raw
	}
	rawSuccessors := make([]capi.Block, len(state.Successors))
	for i, v := range state.Successors {
		rawSuccessors[i] = v.raw
	}
	rawAttrs := make([]capi.NamedAttribute, len(state.Attributes))
	for i, v := range state.Attributes {
		rawAttrs[i] = v.raw
	}
	rawRegions := make([]capi.Region, 0, len(state.OwnedRegions))
	for _, r := range state.OwnedRegions {
		raw, err := r.takeRaw()
		if err != nil {
			return nil, err
		}
		rawRegions = append(rawRegions, raw)
	}

	op := capi.OperationCreate(state.Name, state.Location.raw, rawResults, rawOperands, rawRegions, rawSuccessors, rawAttrs, state.InferResultTypes)
	if capi.OperationIsNull(op) {
		return nil, fmt.Errorf("mlir: failed to create operation %q", state.Name)
	}
	return &OwnedOperation{raw: op}, nil
}

// ReplaceAllUsesWith rewrites every use of v to point at other.
func (v Value) ReplaceAllUsesWith(other Value) error {
	if v.IsNull() || other.IsNull() {
		return fmt.Errorf("mlir: null value")
	}
	capi.ValueReplaceAllUsesOfWith(v.raw, other.raw)
	return nil
}
