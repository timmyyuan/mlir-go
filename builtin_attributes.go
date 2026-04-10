//go:build cgo

package mlir

import (
	"fmt"

	"github.com/timmyyuan/mlir-go/internal/capi"
)

func IntegerAttribute(typ Type, value int64) (Attribute, error) {
	if typ.IsNull() {
		return Attribute{}, fmt.Errorf("mlir: type is required")
	}
	attr := capi.IntegerAttrGet(typ.raw, value)
	if capi.AttributeIsNull(attr) {
		return Attribute{}, fmt.Errorf("mlir: failed to construct integer attribute")
	}
	return Attribute{raw: attr}, nil
}

func BoolAttribute(ctx *Context, value bool) (Attribute, error) {
	raw, err := ctx.requireRaw()
	if err != nil {
		return Attribute{}, err
	}
	attr := capi.BoolAttrGet(raw, value)
	if capi.AttributeIsNull(attr) {
		return Attribute{}, fmt.Errorf("mlir: failed to construct bool attribute")
	}
	return Attribute{raw: attr}, nil
}

func StringAttribute(ctx *Context, value string) (Attribute, error) {
	raw, err := ctx.requireRaw()
	if err != nil {
		return Attribute{}, err
	}
	attr := capi.StringAttrGet(raw, value)
	if capi.AttributeIsNull(attr) {
		return Attribute{}, fmt.Errorf("mlir: failed to construct string attribute")
	}
	return Attribute{raw: attr}, nil
}

func FlatSymbolRefAttribute(ctx *Context, symbol string) (Attribute, error) {
	raw, err := ctx.requireRaw()
	if err != nil {
		return Attribute{}, err
	}
	attr := capi.FlatSymbolRefAttrGet(raw, symbol)
	if capi.AttributeIsNull(attr) {
		return Attribute{}, fmt.Errorf("mlir: failed to construct flat symbol ref attribute")
	}
	return Attribute{raw: attr}, nil
}

func TypeAttribute(typ Type) (Attribute, error) {
	if typ.IsNull() {
		return Attribute{}, fmt.Errorf("mlir: type is required")
	}
	attr := capi.TypeAttrGet(typ.raw)
	if capi.AttributeIsNull(attr) {
		return Attribute{}, fmt.Errorf("mlir: failed to construct type attribute")
	}
	return Attribute{raw: attr}, nil
}

func UnitAttribute(ctx *Context) (Attribute, error) {
	raw, err := ctx.requireRaw()
	if err != nil {
		return Attribute{}, err
	}
	attr := capi.UnitAttrGet(raw)
	if capi.AttributeIsNull(attr) {
		return Attribute{}, fmt.Errorf("mlir: failed to construct unit attribute")
	}
	return Attribute{raw: attr}, nil
}

func (a Attribute) IsInteger() bool {
	return !a.IsNull() && capi.AttributeIsAInteger(a.raw)
}

func (a Attribute) Int64Value() (int64, error) {
	if !a.IsInteger() {
		return 0, fmt.Errorf("mlir: attribute is not an integer")
	}
	return capi.IntegerAttrGetValueInt(a.raw), nil
}

func (a Attribute) IsBool() bool {
	return !a.IsNull() && capi.AttributeIsABool(a.raw)
}

func (a Attribute) BoolValue() (bool, error) {
	if !a.IsBool() {
		return false, fmt.Errorf("mlir: attribute is not a bool")
	}
	return capi.BoolAttrGetValue(a.raw), nil
}

func (a Attribute) IsString() bool {
	return !a.IsNull() && capi.AttributeIsAString(a.raw)
}

func (a Attribute) StringValue() (string, error) {
	if !a.IsString() {
		return "", fmt.Errorf("mlir: attribute is not a string")
	}
	return capi.StringAttrGetValue(a.raw), nil
}

func (a Attribute) IsFlatSymbolRef() bool {
	return !a.IsNull() && capi.AttributeIsAFlatSymbolRef(a.raw)
}

func (a Attribute) FlatSymbolValue() (string, error) {
	if !a.IsFlatSymbolRef() {
		return "", fmt.Errorf("mlir: attribute is not a flat symbol ref")
	}
	return capi.FlatSymbolRefAttrGetValue(a.raw), nil
}

func (a Attribute) IsTypeAttribute() bool {
	return !a.IsNull() && capi.AttributeIsAType(a.raw)
}

func (a Attribute) TypeValue() (Type, error) {
	if !a.IsTypeAttribute() {
		return Type{}, fmt.Errorf("mlir: attribute is not a type attribute")
	}
	return Type{raw: capi.TypeAttrGetValue(a.raw)}, nil
}

func (a Attribute) IsUnit() bool {
	return !a.IsNull() && capi.AttributeIsAUnit(a.raw)
}
