//go:build !cgo

package mlir

import "fmt"

func SignlessIntegerType(*Context, uint) (Type, error) {
	return Type{}, fmt.Errorf("mlir: cgo is required")
}
func SignedIntegerType(*Context, uint) (Type, error) {
	return Type{}, fmt.Errorf("mlir: cgo is required")
}
func UnsignedIntegerType(*Context, uint) (Type, error) {
	return Type{}, fmt.Errorf("mlir: cgo is required")
}
func IndexType(*Context) (Type, error) { return Type{}, fmt.Errorf("mlir: cgo is required") }
func F32Type(*Context) (Type, error)   { return Type{}, fmt.Errorf("mlir: cgo is required") }
func F64Type(*Context) (Type, error)   { return Type{}, fmt.Errorf("mlir: cgo is required") }
func NoneType(*Context) (Type, error)  { return Type{}, fmt.Errorf("mlir: cgo is required") }
func RankedTensorType([]int64, Type) (Type, error) {
	return Type{}, fmt.Errorf("mlir: cgo is required")
}
func RankedTensorTypeWithEncoding([]int64, Type, Attribute) (Type, error) {
	return Type{}, fmt.Errorf("mlir: cgo is required")
}
func CheckedRankedTensorType(Location, []int64, Type, Attribute) (Type, error) {
	return Type{}, fmt.Errorf("mlir: cgo is required")
}
func MemRefType([]int64, Type) (Type, error) {
	return Type{}, fmt.Errorf("mlir: cgo is required")
}
func MemRefTypeWithMemorySpace([]int64, Type, Attribute) (Type, error) {
	return Type{}, fmt.Errorf("mlir: cgo is required")
}
func CheckedMemRefType(Location, []int64, Type, Attribute) (Type, error) {
	return Type{}, fmt.Errorf("mlir: cgo is required")
}
func UnrankedMemRefType(Type) (Type, error) {
	return Type{}, fmt.Errorf("mlir: cgo is required")
}
func UnrankedMemRefTypeWithMemorySpace(Type, Attribute) (Type, error) {
	return Type{}, fmt.Errorf("mlir: cgo is required")
}
func CheckedUnrankedMemRefType(Location, Type, Attribute) (Type, error) {
	return Type{}, fmt.Errorf("mlir: cgo is required")
}
func FunctionType(*Context, []Type, []Type) (Type, error) {
	return Type{}, fmt.Errorf("mlir: cgo is required")
}

func (t Type) IsInteger() bool           { return false }
func (t Type) IntegerWidth() int         { return 0 }
func (t Type) IsSignlessInteger() bool   { return false }
func (t Type) IsSignedInteger() bool     { return false }
func (t Type) IsUnsignedInteger() bool   { return false }
func (t Type) IsIndex() bool             { return false }
func (t Type) IsFloat() bool             { return false }
func (t Type) FloatWidth() int           { return 0 }
func (t Type) IsF32() bool               { return false }
func (t Type) IsF64() bool               { return false }
func (t Type) IsNone() bool              { return false }
func (t Type) IsShaped() bool            { return false }
func (t Type) ElementType() Type         { return Type{} }
func (t Type) HasRank() bool             { return false }
func (t Type) Rank() int                 { return 0 }
func (t Type) HasStaticShape() bool      { return false }
func (t Type) IsDynamicDim(int) bool     { return false }
func (t Type) DimSize(int) int64         { return 0 }
func (t Type) Shape() []int64            { return nil }
func DynamicSize() int64                 { return 0 }
func IsDynamicSize(int64) bool           { return false }
func (t Type) IsRankedTensor() bool      { return false }
func (t Type) TensorEncoding() Attribute { return Attribute{} }
func (t Type) IsMemRef() bool            { return false }
func (t Type) IsUnrankedMemRef() bool    { return false }
func (t Type) MemRefLayout() Attribute   { return Attribute{} }
func (t Type) MemRefMemorySpace() Attribute {
	return Attribute{}
}
func (t Type) IsFunction() bool        { return false }
func (t Type) FunctionInputs() []Type  { return nil }
func (t Type) FunctionResults() []Type { return nil }
