//go:build !cgo

package mlir

import "fmt"

func IntegerAttribute(Type, int64) (Attribute, error) {
	return Attribute{}, fmt.Errorf("mlir: cgo is required")
}
func BoolAttribute(*Context, bool) (Attribute, error) {
	return Attribute{}, fmt.Errorf("mlir: cgo is required")
}
func StringAttribute(*Context, string) (Attribute, error) {
	return Attribute{}, fmt.Errorf("mlir: cgo is required")
}
func FlatSymbolRefAttribute(*Context, string) (Attribute, error) {
	return Attribute{}, fmt.Errorf("mlir: cgo is required")
}
func TypeAttribute(Type) (Attribute, error) { return Attribute{}, fmt.Errorf("mlir: cgo is required") }
func UnitAttribute(*Context) (Attribute, error) {
	return Attribute{}, fmt.Errorf("mlir: cgo is required")
}

func (a Attribute) IsInteger() bool                  { return false }
func (a Attribute) Int64Value() (int64, error)       { return 0, fmt.Errorf("mlir: cgo is required") }
func (a Attribute) IsBool() bool                     { return false }
func (a Attribute) BoolValue() (bool, error)         { return false, fmt.Errorf("mlir: cgo is required") }
func (a Attribute) IsString() bool                   { return false }
func (a Attribute) StringValue() (string, error)     { return "", fmt.Errorf("mlir: cgo is required") }
func (a Attribute) IsFlatSymbolRef() bool            { return false }
func (a Attribute) FlatSymbolValue() (string, error) { return "", fmt.Errorf("mlir: cgo is required") }
func (a Attribute) IsTypeAttribute() bool            { return false }
func (a Attribute) TypeValue() (Type, error)         { return Type{}, fmt.Errorf("mlir: cgo is required") }
func (a Attribute) IsUnit() bool                     { return false }
