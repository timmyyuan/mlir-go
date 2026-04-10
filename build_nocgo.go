//go:build !cgo

package mlir

import "fmt"

type OwnedRegion struct{}
type OwnedBlock struct{}
type OwnedOperation struct{}

func NewOwnedRegion() (*OwnedRegion, error) { return nil, fmt.Errorf("mlir: cgo is required") }
func (r *OwnedRegion) Close() error         { return fmt.Errorf("mlir: cgo is required") }
func (r *OwnedRegion) Region() Region       { return Region{} }
func (r *OwnedRegion) AppendBlock(*OwnedBlock) (Block, error) {
	return Block{}, fmt.Errorf("mlir: cgo is required")
}

func NewOwnedBlock([]Type, []Location) (*OwnedBlock, error) {
	return nil, fmt.Errorf("mlir: cgo is required")
}
func (b *OwnedBlock) Close() error { return fmt.Errorf("mlir: cgo is required") }
func (b *OwnedBlock) Block() Block { return Block{} }
func (b *OwnedBlock) AppendOperation(*OwnedOperation) (Operation, error) {
	return Operation{}, fmt.Errorf("mlir: cgo is required")
}
func (b *OwnedBlock) InsertOperationBefore(Operation, *OwnedOperation) (Operation, error) {
	return Operation{}, fmt.Errorf("mlir: cgo is required")
}
func (b *OwnedBlock) AddArgument(Type, Location) (Value, error) {
	return Value{}, fmt.Errorf("mlir: cgo is required")
}

func (o *OwnedOperation) Close() error         { return fmt.Errorf("mlir: cgo is required") }
func (o *OwnedOperation) Operation() Operation { return Operation{} }

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
func (s *OperationState) AddResults(results ...Type)  { s.Results = append(s.Results, results...) }
func (s *OperationState) AddOperands(values ...Value) { s.Operands = append(s.Operands, values...) }
func (s *OperationState) AddSuccessors(blocks ...Block) {
	s.Successors = append(s.Successors, blocks...)
}
func (s *OperationState) AddAttributes(attrs ...NamedAttribute) {
	s.Attributes = append(s.Attributes, attrs...)
}
func (s *OperationState) AddOwnedRegions(regions ...*OwnedRegion) {
	s.OwnedRegions = append(s.OwnedRegions, regions...)
}
func (s *OperationState) EnableResultTypeInference() { s.InferResultTypes = true }

func CreateOperation(*OperationState) (*OwnedOperation, error) {
	return nil, fmt.Errorf("mlir: cgo is required")
}

func (v Value) ReplaceAllUsesWith(Value) error {
	return fmt.Errorf("mlir: cgo is required")
}
