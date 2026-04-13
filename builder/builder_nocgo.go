//go:build !cgo

package builder

import (
	"fmt"

	mlir "github.com/timmyyuan/mlir-go"
)

type Builder struct{}

func New(*mlir.Context) (*Builder, error) {
	return nil, fmt.Errorf("mlir: cgo is required")
}

func (b *Builder) Context() *mlir.Context {
	return nil
}

func (b *Builder) Location() mlir.Location {
	return mlir.Location{}
}

func (b *Builder) SetLocation(mlir.Location) error {
	return fmt.Errorf("mlir: cgo is required")
}

func (b *Builder) SetUnknownLocation() error {
	return fmt.Errorf("mlir: cgo is required")
}

func (b *Builder) Module() *mlir.Module {
	return nil
}

func (b *Builder) CurrentBlock() mlir.Block {
	return mlir.Block{}
}

func (b *Builder) PositionAtEnd(mlir.Block) error {
	return fmt.Errorf("mlir: cgo is required")
}

func (b *Builder) BuildModule(func(*Builder, *mlir.Module) error) (*mlir.Module, error) {
	return nil, fmt.Errorf("mlir: cgo is required")
}

func (b *Builder) BuildFunction(string, mlir.Type, func(*Builder, mlir.Block) error) (mlir.Operation, error) {
	return mlir.Operation{}, fmt.Errorf("mlir: cgo is required")
}

func (b *Builder) BuildFunctionWithArgLocations(string, mlir.Type, []mlir.Location, func(*Builder, mlir.Block) error) (mlir.Operation, error) {
	return mlir.Operation{}, fmt.Errorf("mlir: cgo is required")
}

func (b *Builder) AppendBlock(...mlir.Type) (mlir.Block, error) {
	return mlir.Block{}, fmt.Errorf("mlir: cgo is required")
}

func (b *Builder) AppendBlockWithLocations([]mlir.Type, []mlir.Location) (mlir.Block, error) {
	return mlir.Block{}, fmt.Errorf("mlir: cgo is required")
}

func (b *Builder) Emit(*mlir.OwnedOperation) (mlir.Operation, error) {
	return mlir.Operation{}, fmt.Errorf("mlir: cgo is required")
}

func (b *Builder) EmitState(*mlir.OperationState) (mlir.Operation, error) {
	return mlir.Operation{}, fmt.Errorf("mlir: cgo is required")
}

func (b *Builder) EmitWithLocation(func(mlir.Location) (*mlir.OwnedOperation, error)) (mlir.Operation, error) {
	return mlir.Operation{}, fmt.Errorf("mlir: cgo is required")
}

func (b *Builder) EmitFactory(func() (*mlir.OwnedOperation, error)) (mlir.Operation, error) {
	return mlir.Operation{}, fmt.Errorf("mlir: cgo is required")
}
