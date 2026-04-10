//go:build cgo

package capi

/*
#include <stdlib.h>

#include "mlir-c/BuiltinAttributes.h"
#include "mlir-c/BuiltinTypes.h"
*/
import "C"
import "unsafe"

func IntegerTypeGet(ctx Context, width uint) Type {
	return Type{c: C.mlirIntegerTypeGet(ctx.c, C.uint(width))}
}

func IntegerTypeSignedGet(ctx Context, width uint) Type {
	return Type{c: C.mlirIntegerTypeSignedGet(ctx.c, C.uint(width))}
}

func IntegerTypeUnsignedGet(ctx Context, width uint) Type {
	return Type{c: C.mlirIntegerTypeUnsignedGet(ctx.c, C.uint(width))}
}

func TypeIsAInteger(typ Type) bool {
	return bool(C.mlirTypeIsAInteger(typ.c))
}

func IntegerTypeGetWidth(typ Type) int {
	return int(C.mlirIntegerTypeGetWidth(typ.c))
}

func IntegerTypeIsSignless(typ Type) bool {
	return bool(C.mlirIntegerTypeIsSignless(typ.c))
}

func IntegerTypeIsSigned(typ Type) bool {
	return bool(C.mlirIntegerTypeIsSigned(typ.c))
}

func IntegerTypeIsUnsigned(typ Type) bool {
	return bool(C.mlirIntegerTypeIsUnsigned(typ.c))
}

func IndexTypeGet(ctx Context) Type {
	return Type{c: C.mlirIndexTypeGet(ctx.c)}
}

func TypeIsAIndex(typ Type) bool {
	return bool(C.mlirTypeIsAIndex(typ.c))
}

func F32TypeGet(ctx Context) Type {
	return Type{c: C.mlirF32TypeGet(ctx.c)}
}

func F64TypeGet(ctx Context) Type {
	return Type{c: C.mlirF64TypeGet(ctx.c)}
}

func TypeIsAFloat(typ Type) bool {
	return bool(C.mlirTypeIsAFloat(typ.c))
}

func FloatTypeGetWidth(typ Type) int {
	return int(C.mlirFloatTypeGetWidth(typ.c))
}

func TypeIsAF32(typ Type) bool {
	return bool(C.mlirTypeIsAF32(typ.c))
}

func TypeIsAF64(typ Type) bool {
	return bool(C.mlirTypeIsAF64(typ.c))
}

func NoneTypeGet(ctx Context) Type {
	return Type{c: C.mlirNoneTypeGet(ctx.c)}
}

func TypeIsANone(typ Type) bool {
	return bool(C.mlirTypeIsANone(typ.c))
}

func RankedTensorTypeGet(shape []int64, elementType Type, encoding Attribute) Type {
	var shapePtr *C.int64_t
	if len(shape) > 0 {
		cShape := make([]C.int64_t, len(shape))
		for i, dim := range shape {
			cShape[i] = C.int64_t(dim)
		}
		shapePtr = &cShape[0]
		return Type{
			c: C.mlirRankedTensorTypeGet(C.intptr_t(len(cShape)), shapePtr, elementType.c, encoding.c),
		}
	}

	return Type{c: C.mlirRankedTensorTypeGet(0, nil, elementType.c, encoding.c)}
}

func RankedTensorTypeGetChecked(loc Location, shape []int64, elementType Type, encoding Attribute) Type {
	var shapePtr *C.int64_t
	if len(shape) > 0 {
		cShape := make([]C.int64_t, len(shape))
		for i, dim := range shape {
			cShape[i] = C.int64_t(dim)
		}
		shapePtr = &cShape[0]
		return Type{
			c: C.mlirRankedTensorTypeGetChecked(loc.c, C.intptr_t(len(cShape)), shapePtr, elementType.c, encoding.c),
		}
	}

	return Type{c: C.mlirRankedTensorTypeGetChecked(loc.c, 0, nil, elementType.c, encoding.c)}
}

func TypeIsAShaped(typ Type) bool {
	return bool(C.mlirTypeIsAShaped(typ.c))
}

func ShapedTypeGetElementType(typ Type) Type {
	return Type{c: C.mlirShapedTypeGetElementType(typ.c)}
}

func ShapedTypeHasRank(typ Type) bool {
	return bool(C.mlirShapedTypeHasRank(typ.c))
}

func ShapedTypeGetRank(typ Type) int {
	return int(C.mlirShapedTypeGetRank(typ.c))
}

func ShapedTypeHasStaticShape(typ Type) bool {
	return bool(C.mlirShapedTypeHasStaticShape(typ.c))
}

func ShapedTypeIsDynamicDim(typ Type, dim int) bool {
	return bool(C.mlirShapedTypeIsDynamicDim(typ.c, C.intptr_t(dim)))
}

func ShapedTypeGetDimSize(typ Type, dim int) int64 {
	return int64(C.mlirShapedTypeGetDimSize(typ.c, C.intptr_t(dim)))
}

func ShapedTypeIsDynamicSize(size int64) bool {
	return bool(C.mlirShapedTypeIsDynamicSize(C.int64_t(size)))
}

func ShapedTypeGetDynamicSize() int64 {
	return int64(C.mlirShapedTypeGetDynamicSize())
}

func TypeIsARankedTensor(typ Type) bool {
	return bool(C.mlirTypeIsARankedTensor(typ.c))
}

func RankedTensorTypeGetEncoding(typ Type) Attribute {
	return Attribute{c: C.mlirRankedTensorTypeGetEncoding(typ.c)}
}

func FunctionTypeGet(ctx Context, inputs, results []Type) Type {
	var inputsPtr *C.MlirType
	var resultsPtr *C.MlirType

	if len(inputs) > 0 {
		cInputs := make([]C.MlirType, len(inputs))
		for i, input := range inputs {
			cInputs[i] = input.c
		}
		inputsPtr = &cInputs[0]
		if len(results) > 0 {
			cResults := make([]C.MlirType, len(results))
			for i, result := range results {
				cResults[i] = result.c
			}
			resultsPtr = &cResults[0]
			return Type{
				c: C.mlirFunctionTypeGet(ctx.c, C.intptr_t(len(cInputs)), inputsPtr, C.intptr_t(len(cResults)), resultsPtr),
			}
		}
		return Type{
			c: C.mlirFunctionTypeGet(ctx.c, C.intptr_t(len(cInputs)), inputsPtr, 0, nil),
		}
	}

	if len(results) > 0 {
		cResults := make([]C.MlirType, len(results))
		for i, result := range results {
			cResults[i] = result.c
		}
		resultsPtr = &cResults[0]
	}
	return Type{
		c: C.mlirFunctionTypeGet(ctx.c, 0, nil, C.intptr_t(len(results)), resultsPtr),
	}
}

func TypeIsAFunction(typ Type) bool {
	return bool(C.mlirTypeIsAFunction(typ.c))
}

func FunctionTypeGetNumInputs(typ Type) int {
	return int(C.mlirFunctionTypeGetNumInputs(typ.c))
}

func FunctionTypeGetNumResults(typ Type) int {
	return int(C.mlirFunctionTypeGetNumResults(typ.c))
}

func FunctionTypeGetInput(typ Type, pos int) Type {
	return Type{c: C.mlirFunctionTypeGetInput(typ.c, C.intptr_t(pos))}
}

func FunctionTypeGetResult(typ Type, pos int) Type {
	return Type{c: C.mlirFunctionTypeGetResult(typ.c, C.intptr_t(pos))}
}

func MemRefTypeContiguousGet(elementType Type, shape []int64, memorySpace Attribute) Type {
	var shapePtr *C.int64_t
	if len(shape) > 0 {
		cShape := make([]C.int64_t, len(shape))
		for i, dim := range shape {
			cShape[i] = C.int64_t(dim)
		}
		shapePtr = &cShape[0]
		return Type{
			c: C.mlirMemRefTypeContiguousGet(elementType.c, C.intptr_t(len(cShape)), shapePtr, memorySpace.c),
		}
	}
	return Type{c: C.mlirMemRefTypeContiguousGet(elementType.c, 0, nil, memorySpace.c)}
}

func MemRefTypeContiguousGetChecked(loc Location, elementType Type, shape []int64, memorySpace Attribute) Type {
	var shapePtr *C.int64_t
	if len(shape) > 0 {
		cShape := make([]C.int64_t, len(shape))
		for i, dim := range shape {
			cShape[i] = C.int64_t(dim)
		}
		shapePtr = &cShape[0]
		return Type{
			c: C.mlirMemRefTypeContiguousGetChecked(loc.c, elementType.c, C.intptr_t(len(cShape)), shapePtr, memorySpace.c),
		}
	}
	return Type{c: C.mlirMemRefTypeContiguousGetChecked(loc.c, elementType.c, 0, nil, memorySpace.c)}
}

func UnrankedMemRefTypeGet(elementType Type, memorySpace Attribute) Type {
	return Type{c: C.mlirUnrankedMemRefTypeGet(elementType.c, memorySpace.c)}
}

func UnrankedMemRefTypeGetChecked(loc Location, elementType Type, memorySpace Attribute) Type {
	return Type{c: C.mlirUnrankedMemRefTypeGetChecked(loc.c, elementType.c, memorySpace.c)}
}

func TypeIsAMemRef(typ Type) bool {
	return bool(C.mlirTypeIsAMemRef(typ.c))
}

func TypeIsAUnrankedMemRef(typ Type) bool {
	return bool(C.mlirTypeIsAUnrankedMemRef(typ.c))
}

func MemRefTypeGetLayout(typ Type) Attribute {
	return Attribute{c: C.mlirMemRefTypeGetLayout(typ.c)}
}

func MemRefTypeGetMemorySpace(typ Type) Attribute {
	return Attribute{c: C.mlirMemRefTypeGetMemorySpace(typ.c)}
}

func AttributeIsAInteger(attr Attribute) bool {
	return bool(C.mlirAttributeIsAInteger(attr.c))
}

func IntegerAttrGet(typ Type, value int64) Attribute {
	return Attribute{c: C.mlirIntegerAttrGet(typ.c, C.int64_t(value))}
}

func IntegerAttrGetValueInt(attr Attribute) int64 {
	return int64(C.mlirIntegerAttrGetValueInt(attr.c))
}

func AttributeIsABool(attr Attribute) bool {
	return bool(C.mlirAttributeIsABool(attr.c))
}

func BoolAttrGet(ctx Context, value bool) Attribute {
	flag := 0
	if value {
		flag = 1
	}
	return Attribute{c: C.mlirBoolAttrGet(ctx.c, C.int(flag))}
}

func BoolAttrGetValue(attr Attribute) bool {
	return bool(C.mlirBoolAttrGetValue(attr.c))
}

func AttributeIsAString(attr Attribute) bool {
	return bool(C.mlirAttributeIsAString(attr.c))
}

func StringAttrGet(ctx Context, value string) Attribute {
	cstr := C.CString(value)
	defer C.free(unsafe.Pointer(cstr))
	return Attribute{
		c: C.mlirStringAttrGet(ctx.c, C.mlirStringRefCreate(cstr, C.size_t(len(value)))),
	}
}

func StringAttrGetValue(attr Attribute) string {
	str := C.mlirStringAttrGetValue(attr.c)
	if str.data == nil || str.length == 0 {
		return ""
	}
	return C.GoStringN(str.data, C.int(str.length))
}

func AttributeIsAFlatSymbolRef(attr Attribute) bool {
	return bool(C.mlirAttributeIsAFlatSymbolRef(attr.c))
}

func FlatSymbolRefAttrGet(ctx Context, symbol string) Attribute {
	cstr := C.CString(symbol)
	defer C.free(unsafe.Pointer(cstr))
	return Attribute{
		c: C.mlirFlatSymbolRefAttrGet(ctx.c, C.mlirStringRefCreate(cstr, C.size_t(len(symbol)))),
	}
}

func FlatSymbolRefAttrGetValue(attr Attribute) string {
	str := C.mlirFlatSymbolRefAttrGetValue(attr.c)
	if str.data == nil || str.length == 0 {
		return ""
	}
	return C.GoStringN(str.data, C.int(str.length))
}

func AttributeIsAType(attr Attribute) bool {
	return bool(C.mlirAttributeIsAType(attr.c))
}

func TypeAttrGet(typ Type) Attribute {
	return Attribute{c: C.mlirTypeAttrGet(typ.c)}
}

func TypeAttrGetValue(attr Attribute) Type {
	return Type{c: C.mlirTypeAttrGetValue(attr.c)}
}

func AttributeIsAUnit(attr Attribute) bool {
	return bool(C.mlirAttributeIsAUnit(attr.c))
}

func UnitAttrGet(ctx Context) Attribute {
	return Attribute{c: C.mlirUnitAttrGet(ctx.c)}
}
