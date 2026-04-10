//go:build cgo

package mlir

import (
	"fmt"

	"github.com/timmyyuan/mlir-go/internal/capi"
)

func SignlessIntegerType(ctx *Context, width uint) (Type, error) {
	raw, err := ctx.requireRaw()
	if err != nil {
		return Type{}, err
	}
	typ := capi.IntegerTypeGet(raw, width)
	if capi.TypeIsNull(typ) {
		return Type{}, fmt.Errorf("mlir: failed to construct signless integer type")
	}
	return Type{raw: typ}, nil
}

func SignedIntegerType(ctx *Context, width uint) (Type, error) {
	raw, err := ctx.requireRaw()
	if err != nil {
		return Type{}, err
	}
	typ := capi.IntegerTypeSignedGet(raw, width)
	if capi.TypeIsNull(typ) {
		return Type{}, fmt.Errorf("mlir: failed to construct signed integer type")
	}
	return Type{raw: typ}, nil
}

func UnsignedIntegerType(ctx *Context, width uint) (Type, error) {
	raw, err := ctx.requireRaw()
	if err != nil {
		return Type{}, err
	}
	typ := capi.IntegerTypeUnsignedGet(raw, width)
	if capi.TypeIsNull(typ) {
		return Type{}, fmt.Errorf("mlir: failed to construct unsigned integer type")
	}
	return Type{raw: typ}, nil
}

func IndexType(ctx *Context) (Type, error) {
	raw, err := ctx.requireRaw()
	if err != nil {
		return Type{}, err
	}
	typ := capi.IndexTypeGet(raw)
	if capi.TypeIsNull(typ) {
		return Type{}, fmt.Errorf("mlir: failed to construct index type")
	}
	return Type{raw: typ}, nil
}

func F32Type(ctx *Context) (Type, error) {
	raw, err := ctx.requireRaw()
	if err != nil {
		return Type{}, err
	}
	typ := capi.F32TypeGet(raw)
	if capi.TypeIsNull(typ) {
		return Type{}, fmt.Errorf("mlir: failed to construct f32 type")
	}
	return Type{raw: typ}, nil
}

func F64Type(ctx *Context) (Type, error) {
	raw, err := ctx.requireRaw()
	if err != nil {
		return Type{}, err
	}
	typ := capi.F64TypeGet(raw)
	if capi.TypeIsNull(typ) {
		return Type{}, fmt.Errorf("mlir: failed to construct f64 type")
	}
	return Type{raw: typ}, nil
}

func NoneType(ctx *Context) (Type, error) {
	raw, err := ctx.requireRaw()
	if err != nil {
		return Type{}, err
	}
	typ := capi.NoneTypeGet(raw)
	if capi.TypeIsNull(typ) {
		return Type{}, fmt.Errorf("mlir: failed to construct none type")
	}
	return Type{raw: typ}, nil
}

func RankedTensorType(shape []int64, elementType Type) (Type, error) {
	return RankedTensorTypeWithEncoding(shape, elementType, Attribute{})
}

func RankedTensorTypeWithEncoding(shape []int64, elementType Type, encoding Attribute) (Type, error) {
	if elementType.IsNull() {
		return Type{}, fmt.Errorf("mlir: element type is required")
	}
	typ := capi.RankedTensorTypeGet(shape, elementType.raw, encoding.raw)
	if capi.TypeIsNull(typ) {
		return Type{}, fmt.Errorf("mlir: failed to construct ranked tensor type")
	}
	return Type{raw: typ}, nil
}

func CheckedRankedTensorType(loc Location, shape []int64, elementType Type, encoding Attribute) (Type, error) {
	if loc.IsNull() {
		return Type{}, fmt.Errorf("mlir: location is required")
	}
	if elementType.IsNull() {
		return Type{}, fmt.Errorf("mlir: element type is required")
	}
	typ := capi.RankedTensorTypeGetChecked(loc.raw, shape, elementType.raw, encoding.raw)
	if capi.TypeIsNull(typ) {
		return Type{}, fmt.Errorf("mlir: failed to construct ranked tensor type")
	}
	return Type{raw: typ}, nil
}

func FunctionType(ctx *Context, inputs []Type, results []Type) (Type, error) {
	raw, err := ctx.requireRaw()
	if err != nil {
		return Type{}, err
	}
	rawInputs := make([]capi.Type, len(inputs))
	for i, input := range inputs {
		rawInputs[i] = input.raw
	}
	rawResults := make([]capi.Type, len(results))
	for i, result := range results {
		rawResults[i] = result.raw
	}
	typ := capi.FunctionTypeGet(raw, rawInputs, rawResults)
	if capi.TypeIsNull(typ) {
		return Type{}, fmt.Errorf("mlir: failed to construct function type")
	}
	return Type{raw: typ}, nil
}

func (t Type) IsInteger() bool {
	return !t.IsNull() && capi.TypeIsAInteger(t.raw)
}

func (t Type) IntegerWidth() int {
	if !t.IsInteger() {
		return 0
	}
	return capi.IntegerTypeGetWidth(t.raw)
}

func (t Type) IsSignlessInteger() bool {
	return t.IsInteger() && capi.IntegerTypeIsSignless(t.raw)
}

func (t Type) IsSignedInteger() bool {
	return t.IsInteger() && capi.IntegerTypeIsSigned(t.raw)
}

func (t Type) IsUnsignedInteger() bool {
	return t.IsInteger() && capi.IntegerTypeIsUnsigned(t.raw)
}

func (t Type) IsIndex() bool {
	return !t.IsNull() && capi.TypeIsAIndex(t.raw)
}

func (t Type) IsFloat() bool {
	return !t.IsNull() && capi.TypeIsAFloat(t.raw)
}

func (t Type) FloatWidth() int {
	if !t.IsFloat() {
		return 0
	}
	return capi.FloatTypeGetWidth(t.raw)
}

func (t Type) IsF32() bool {
	return !t.IsNull() && capi.TypeIsAF32(t.raw)
}

func (t Type) IsF64() bool {
	return !t.IsNull() && capi.TypeIsAF64(t.raw)
}

func (t Type) IsNone() bool {
	return !t.IsNull() && capi.TypeIsANone(t.raw)
}

func (t Type) IsRankedTensor() bool {
	return !t.IsNull() && capi.TypeIsARankedTensor(t.raw)
}

func (t Type) TensorEncoding() Attribute {
	if !t.IsRankedTensor() {
		return Attribute{}
	}
	return Attribute{raw: capi.RankedTensorTypeGetEncoding(t.raw)}
}

func MemRefType(shape []int64, elementType Type) (Type, error) {
	return MemRefTypeWithMemorySpace(shape, elementType, Attribute{})
}

func MemRefTypeWithMemorySpace(shape []int64, elementType Type, memorySpace Attribute) (Type, error) {
	if elementType.IsNull() {
		return Type{}, fmt.Errorf("mlir: element type is required")
	}
	typ := capi.MemRefTypeContiguousGet(elementType.raw, shape, memorySpace.raw)
	if capi.TypeIsNull(typ) {
		return Type{}, fmt.Errorf("mlir: failed to construct memref type")
	}
	return Type{raw: typ}, nil
}

func CheckedMemRefType(loc Location, shape []int64, elementType Type, memorySpace Attribute) (Type, error) {
	if loc.IsNull() {
		return Type{}, fmt.Errorf("mlir: location is required")
	}
	if elementType.IsNull() {
		return Type{}, fmt.Errorf("mlir: element type is required")
	}
	typ := capi.MemRefTypeContiguousGetChecked(loc.raw, elementType.raw, shape, memorySpace.raw)
	if capi.TypeIsNull(typ) {
		return Type{}, fmt.Errorf("mlir: failed to construct memref type")
	}
	return Type{raw: typ}, nil
}

func UnrankedMemRefType(elementType Type) (Type, error) {
	return UnrankedMemRefTypeWithMemorySpace(elementType, Attribute{})
}

func UnrankedMemRefTypeWithMemorySpace(elementType Type, memorySpace Attribute) (Type, error) {
	if elementType.IsNull() {
		return Type{}, fmt.Errorf("mlir: element type is required")
	}
	typ := capi.UnrankedMemRefTypeGet(elementType.raw, memorySpace.raw)
	if capi.TypeIsNull(typ) {
		return Type{}, fmt.Errorf("mlir: failed to construct unranked memref type")
	}
	return Type{raw: typ}, nil
}

func CheckedUnrankedMemRefType(loc Location, elementType Type, memorySpace Attribute) (Type, error) {
	if loc.IsNull() {
		return Type{}, fmt.Errorf("mlir: location is required")
	}
	if elementType.IsNull() {
		return Type{}, fmt.Errorf("mlir: element type is required")
	}
	typ := capi.UnrankedMemRefTypeGetChecked(loc.raw, elementType.raw, memorySpace.raw)
	if capi.TypeIsNull(typ) {
		return Type{}, fmt.Errorf("mlir: failed to construct unranked memref type")
	}
	return Type{raw: typ}, nil
}

func (t Type) IsShaped() bool {
	return !t.IsNull() && capi.TypeIsAShaped(t.raw)
}

func (t Type) ElementType() Type {
	if !t.IsShaped() {
		return Type{}
	}
	return Type{raw: capi.ShapedTypeGetElementType(t.raw)}
}

func (t Type) HasRank() bool {
	return t.IsShaped() && capi.ShapedTypeHasRank(t.raw)
}

func (t Type) Rank() int {
	if !t.HasRank() {
		return 0
	}
	return capi.ShapedTypeGetRank(t.raw)
}

func (t Type) HasStaticShape() bool {
	return t.HasRank() && capi.ShapedTypeHasStaticShape(t.raw)
}

func (t Type) IsDynamicDim(dim int) bool {
	if !t.HasRank() {
		return false
	}
	return capi.ShapedTypeIsDynamicDim(t.raw, dim)
}

func (t Type) DimSize(dim int) int64 {
	if !t.HasRank() {
		return 0
	}
	return capi.ShapedTypeGetDimSize(t.raw, dim)
}

func (t Type) Shape() []int64 {
	if !t.HasRank() {
		return nil
	}
	shape := make([]int64, 0, t.Rank())
	for i := 0; i < t.Rank(); i++ {
		shape = append(shape, t.DimSize(i))
	}
	return shape
}

func DynamicSize() int64 {
	return capi.ShapedTypeGetDynamicSize()
}

func IsDynamicSize(size int64) bool {
	return capi.ShapedTypeIsDynamicSize(size)
}

func (t Type) IsMemRef() bool {
	return !t.IsNull() && capi.TypeIsAMemRef(t.raw)
}

func (t Type) IsUnrankedMemRef() bool {
	return !t.IsNull() && capi.TypeIsAUnrankedMemRef(t.raw)
}

func (t Type) MemRefLayout() Attribute {
	if !t.IsMemRef() {
		return Attribute{}
	}
	return Attribute{raw: capi.MemRefTypeGetLayout(t.raw)}
}

func (t Type) MemRefMemorySpace() Attribute {
	if !t.IsMemRef() && !t.IsUnrankedMemRef() {
		return Attribute{}
	}
	return Attribute{raw: capi.MemRefTypeGetMemorySpace(t.raw)}
}

func (t Type) IsFunction() bool {
	return !t.IsNull() && capi.TypeIsAFunction(t.raw)
}

func (t Type) FunctionInputs() []Type {
	if !t.IsFunction() {
		return nil
	}
	n := capi.FunctionTypeGetNumInputs(t.raw)
	inputs := make([]Type, 0, n)
	for i := 0; i < n; i++ {
		inputs = append(inputs, Type{raw: capi.FunctionTypeGetInput(t.raw, i)})
	}
	return inputs
}

func (t Type) FunctionResults() []Type {
	if !t.IsFunction() {
		return nil
	}
	n := capi.FunctionTypeGetNumResults(t.raw)
	results := make([]Type, 0, n)
	for i := 0; i < n; i++ {
		results = append(results, Type{raw: capi.FunctionTypeGetResult(t.raw, i)})
	}
	return results
}
