//go:build cgo

package builder

import (
	"fmt"

	mlir "github.com/timmyyuan/mlir-go"
	funcdialect "github.com/timmyyuan/mlir-go/dialect/func"
)

// Builder keeps a current location and insertion point while constructing IR.
//
// It intentionally stays close to MLIR structure instead of introducing a
// shadow AST.
type Builder struct {
	ctx           *mlir.Context
	loc           mlir.Location
	module        *mlir.Module
	currentBlock  mlir.Block
	currentRegion *mlir.OwnedRegion
}

// New creates a builder with an unknown default location.
func New(ctx *mlir.Context) (*Builder, error) {
	if ctx == nil {
		return nil, fmt.Errorf("mlir: context is required")
	}
	loc, err := mlir.UnknownLocation(ctx)
	if err != nil {
		return nil, err
	}
	return &Builder{ctx: ctx, loc: loc}, nil
}

// Context returns the MLIR context used by the builder.
func (b *Builder) Context() *mlir.Context {
	if b == nil {
		return nil
	}
	return b.ctx
}

// Location returns the default location used by builder helpers.
func (b *Builder) Location() mlir.Location {
	if b == nil {
		return mlir.Location{}
	}
	return b.loc
}

// SetLocation updates the default location used by builder helpers.
func (b *Builder) SetLocation(loc mlir.Location) error {
	if b == nil {
		return fmt.Errorf("mlir: builder is required")
	}
	if loc.IsNull() {
		return fmt.Errorf("mlir: null location")
	}
	b.loc = loc
	return nil
}

// SetUnknownLocation resets the default location to unknown in the current context.
func (b *Builder) SetUnknownLocation() error {
	if b == nil {
		return fmt.Errorf("mlir: builder is required")
	}
	loc, err := mlir.UnknownLocation(b.ctx)
	if err != nil {
		return err
	}
	b.loc = loc
	return nil
}

// Module returns the module currently under construction, if any.
func (b *Builder) Module() *mlir.Module {
	if b == nil {
		return nil
	}
	return b.module
}

// CurrentBlock returns the current insertion block.
func (b *Builder) CurrentBlock() mlir.Block {
	if b == nil {
		return mlir.Block{}
	}
	return b.currentBlock
}

// PositionAtEnd changes the insertion point to the end of block.
func (b *Builder) PositionAtEnd(block mlir.Block) error {
	if b == nil {
		return fmt.Errorf("mlir: builder is required")
	}
	if block.IsNull() {
		return fmt.Errorf("mlir: null block")
	}
	b.currentBlock = block
	return nil
}

// BuildModule creates an empty module and runs fn with the module body as the
// current insertion block.
func (b *Builder) BuildModule(fn func(*Builder, *mlir.Module) error) (*mlir.Module, error) {
	if b == nil {
		return nil, fmt.Errorf("mlir: builder is required")
	}
	if fn == nil {
		return nil, fmt.Errorf("mlir: module body callback is required")
	}

	mod, err := mlir.CreateEmptyModule(b.loc)
	if err != nil {
		return nil, err
	}

	prevModule := b.module
	prevBlock := b.currentBlock
	prevRegion := b.currentRegion
	b.module = mod
	b.currentBlock = mod.Body()
	b.currentRegion = nil

	runErr := fn(b, mod)

	b.module = prevModule
	b.currentBlock = prevBlock
	b.currentRegion = prevRegion

	if runErr != nil {
		_ = mod.Close()
		return nil, runErr
	}
	return mod, nil
}

// BuildFunction constructs a function operation with entry block arguments
// derived from functionType and inserts it at the current insertion point.
func (b *Builder) BuildFunction(name string, functionType mlir.Type, fn func(*Builder, mlir.Block) error) (mlir.Operation, error) {
	return b.BuildFunctionWithArgLocations(name, functionType, nil, fn)
}

// BuildFunctionWithArgLocations is like BuildFunction but lets the caller
// override entry block argument locations.
func (b *Builder) BuildFunctionWithArgLocations(name string, functionType mlir.Type, argLocs []mlir.Location, fn func(*Builder, mlir.Block) error) (mlir.Operation, error) {
	if b == nil {
		return mlir.Operation{}, fmt.Errorf("mlir: builder is required")
	}
	if fn == nil {
		return mlir.Operation{}, fmt.Errorf("mlir: function body callback is required")
	}
	if b.currentBlock.IsNull() {
		return mlir.Operation{}, fmt.Errorf("mlir: no insertion block")
	}
	if b.currentRegion != nil {
		return mlir.Operation{}, fmt.Errorf("mlir: nested function construction is not supported")
	}
	if name == "" {
		return mlir.Operation{}, fmt.Errorf("mlir: function name is required")
	}
	if !functionType.IsFunction() {
		return mlir.Operation{}, fmt.Errorf("mlir: function type is required")
	}

	inputs := functionType.FunctionInputs()
	effectiveArgLocs, err := b.defaultArgLocations(inputs, argLocs)
	if err != nil {
		return mlir.Operation{}, err
	}

	body, err := mlir.NewOwnedRegion()
	if err != nil {
		return mlir.Operation{}, err
	}
	entry, err := mlir.NewOwnedBlock(inputs, effectiveArgLocs)
	if err != nil {
		_ = body.Close()
		return mlir.Operation{}, err
	}
	entryBlock, err := body.AppendBlock(entry)
	if err != nil {
		_ = body.Close()
		return mlir.Operation{}, err
	}

	parentBlock := b.currentBlock
	prevRegion := b.currentRegion
	prevBlock := b.currentBlock
	b.currentRegion = body
	b.currentBlock = entryBlock

	runErr := fn(b, entryBlock)
	if runErr == nil {
		runErr = requireRegionTerminators(body.Region(), name)
	}

	b.currentRegion = prevRegion
	b.currentBlock = prevBlock

	if runErr != nil {
		_ = body.Close()
		return mlir.Operation{}, runErr
	}

	ownedFunc, err := funcdialect.Func(b.ctx, b.loc, name, functionType, body)
	if err != nil {
		return mlir.Operation{}, err
	}
	funcOp, err := parentBlock.AppendOwnedOperation(ownedFunc)
	if err != nil {
		_ = ownedFunc.Close()
		return mlir.Operation{}, err
	}
	return funcOp, nil
}

// AppendBlock appends a new block to the current region, using the builder's
// default location for all block arguments.
func (b *Builder) AppendBlock(argTypes ...mlir.Type) (mlir.Block, error) {
	return b.AppendBlockWithLocations(argTypes, nil)
}

// AppendBlockWithLocations appends a new block to the current region.
func (b *Builder) AppendBlockWithLocations(argTypes []mlir.Type, argLocs []mlir.Location) (mlir.Block, error) {
	if b == nil {
		return mlir.Block{}, fmt.Errorf("mlir: builder is required")
	}
	if b.currentRegion == nil {
		return mlir.Block{}, fmt.Errorf("mlir: no active region")
	}
	effectiveArgLocs, err := b.defaultArgLocations(argTypes, argLocs)
	if err != nil {
		return mlir.Block{}, err
	}
	block, err := mlir.NewOwnedBlock(argTypes, effectiveArgLocs)
	if err != nil {
		return mlir.Block{}, err
	}
	return b.currentRegion.AppendBlock(block)
}

// Emit appends a detached operation to the current insertion block.
func (b *Builder) Emit(op *mlir.OwnedOperation) (mlir.Operation, error) {
	if b == nil {
		return mlir.Operation{}, fmt.Errorf("mlir: builder is required")
	}
	if b.currentBlock.IsNull() {
		return mlir.Operation{}, fmt.Errorf("mlir: no insertion block")
	}
	return b.currentBlock.AppendOwnedOperation(op)
}

// EmitState constructs an operation from state and appends it to the current block.
func (b *Builder) EmitState(state *mlir.OperationState) (mlir.Operation, error) {
	op, err := mlir.CreateOperation(state)
	if err != nil {
		return mlir.Operation{}, err
	}
	return b.Emit(op)
}

// EmitWithLocation uses the builder's current location to create an operation
// and appends it to the current block.
func (b *Builder) EmitWithLocation(factory func(mlir.Location) (*mlir.OwnedOperation, error)) (mlir.Operation, error) {
	if b == nil {
		return mlir.Operation{}, fmt.Errorf("mlir: builder is required")
	}
	if factory == nil {
		return mlir.Operation{}, fmt.Errorf("mlir: operation factory is required")
	}
	return b.EmitFactory(func() (*mlir.OwnedOperation, error) {
		return factory(b.loc)
	})
}

// EmitFactory creates an operation and appends it to the current block.
func (b *Builder) EmitFactory(factory func() (*mlir.OwnedOperation, error)) (mlir.Operation, error) {
	if b == nil {
		return mlir.Operation{}, fmt.Errorf("mlir: builder is required")
	}
	if factory == nil {
		return mlir.Operation{}, fmt.Errorf("mlir: operation factory is required")
	}
	op, err := factory()
	if err != nil {
		return mlir.Operation{}, err
	}
	return b.Emit(op)
}

func (b *Builder) defaultArgLocations(argTypes []mlir.Type, argLocs []mlir.Location) ([]mlir.Location, error) {
	if len(argLocs) != 0 && len(argLocs) != len(argTypes) {
		return nil, fmt.Errorf("mlir: argument type/location length mismatch")
	}
	if len(argTypes) == 0 {
		return nil, nil
	}
	if len(argLocs) != 0 {
		return append([]mlir.Location(nil), argLocs...), nil
	}
	locs := make([]mlir.Location, len(argTypes))
	for i := range locs {
		locs[i] = b.loc
	}
	return locs, nil
}

func requireRegionTerminators(region mlir.Region, functionName string) error {
	blocks := region.Blocks()
	for i, block := range blocks {
		if block.Terminator().IsNull() {
			return fmt.Errorf("mlir: block %d in function %q is missing a terminator", i, functionName)
		}
	}
	return nil
}
