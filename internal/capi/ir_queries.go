//go:build cgo

package capi

/*
#include <stdlib.h>

#include "mlir-c/IR.h"
*/
import "C"
import "unsafe"

func OperationGetBlock(op Operation) Block {
	return Block{c: C.mlirOperationGetBlock(op.c)}
}

func OperationGetParentOperation(op Operation) Operation {
	return Operation{c: C.mlirOperationGetParentOperation(op.c)}
}

func OperationGetNumOperands(op Operation) int {
	return int(C.mlirOperationGetNumOperands(op.c))
}

func OperationGetOperand(op Operation, pos int) Value {
	return Value{c: C.mlirOperationGetOperand(op.c, C.intptr_t(pos))}
}

func OperationSetOperand(op Operation, pos int, value Value) {
	C.mlirOperationSetOperand(op.c, C.intptr_t(pos), value.c)
}

func OperationSetOperands(op Operation, operands []Value) {
	if len(operands) == 0 {
		C.mlirOperationSetOperands(op.c, 0, nil)
		return
	}
	cOperands := make([]C.MlirValue, len(operands))
	for i, operand := range operands {
		cOperands[i] = operand.c
	}
	C.mlirOperationSetOperands(op.c, C.intptr_t(len(cOperands)), &cOperands[0])
}

func OperationGetNumSuccessors(op Operation) int {
	return int(C.mlirOperationGetNumSuccessors(op.c))
}

func OperationGetSuccessor(op Operation, pos int) Block {
	return Block{c: C.mlirOperationGetSuccessor(op.c, C.intptr_t(pos))}
}

func OperationSetSuccessor(op Operation, pos int, block Block) {
	C.mlirOperationSetSuccessor(op.c, C.intptr_t(pos), block.c)
}

func OperationGetNumAttributes(op Operation) int {
	return int(C.mlirOperationGetNumAttributes(op.c))
}

func OperationGetAttribute(op Operation, pos int) NamedAttribute {
	return NamedAttribute{c: C.mlirOperationGetAttribute(op.c, C.intptr_t(pos))}
}

func OperationGetAttributeByName(op Operation, name string) Attribute {
	cstr := C.CString(name)
	defer C.free(unsafe.Pointer(cstr))
	return Attribute{
		c: C.mlirOperationGetAttributeByName(op.c, C.mlirStringRefCreate(cstr, C.size_t(len(name)))),
	}
}

func OperationSetAttributeByName(op Operation, name string, attr Attribute) {
	cstr := C.CString(name)
	defer C.free(unsafe.Pointer(cstr))
	C.mlirOperationSetAttributeByName(op.c, C.mlirStringRefCreate(cstr, C.size_t(len(name))), attr.c)
}

func OperationRemoveAttributeByName(op Operation, name string) bool {
	cstr := C.CString(name)
	defer C.free(unsafe.Pointer(cstr))
	return bool(C.mlirOperationRemoveAttributeByName(op.c, C.mlirStringRefCreate(cstr, C.size_t(len(name)))))
}

func BlockGetParentOperation(block Block) Operation {
	return Operation{c: C.mlirBlockGetParentOperation(block.c)}
}

func BlockGetParentRegion(block Block) Region {
	return Region{c: C.mlirBlockGetParentRegion(block.c)}
}

func BlockGetTerminator(block Block) Operation {
	return Operation{c: C.mlirBlockGetTerminator(block.c)}
}

func ValueEqual(lhs, rhs Value) bool {
	return bool(C.mlirValueEqual(lhs.c, rhs.c))
}

func BlockArgumentGetOwner(value Value) Block {
	return Block{c: C.mlirBlockArgumentGetOwner(value.c)}
}

func BlockArgumentGetArgNumber(value Value) int {
	return int(C.mlirBlockArgumentGetArgNumber(value.c))
}

func BlockArgumentSetType(value Value, typ Type) {
	C.mlirBlockArgumentSetType(value.c, typ.c)
}

func OpResultGetOwner(value Value) Operation {
	return Operation{c: C.mlirOpResultGetOwner(value.c)}
}

func OpResultGetResultNumber(value Value) int {
	return int(C.mlirOpResultGetResultNumber(value.c))
}

func ValueSetType(value Value, typ Type) {
	C.mlirValueSetType(value.c, typ.c)
}
