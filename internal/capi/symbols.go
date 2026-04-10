//go:build cgo

package capi

/*
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

#include "mlir-c/IR.h"

extern void mlirGoSymbolTableWalkCallback(MlirOperation operation, bool allSymUsesVisible, void *userData);

static void mlirGoSymbolTableWalk(
    MlirOperation from, bool allSymUsesVisible, void *userData) {
	mlirSymbolTableWalkSymbolTables(
	    from, allSymUsesVisible, mlirGoSymbolTableWalkCallback, userData);
}
*/
import "C"
import (
	"runtime/cgo"
	"unsafe"
)

func NullSymbolTable() SymbolTable {
	return SymbolTable{}
}

func SymbolTableIsNull(symbolTable SymbolTable) bool {
	return bool(C.mlirSymbolTableIsNull(symbolTable.c))
}

func SymbolTableGetSymbolAttributeName() string {
	str := C.mlirSymbolTableGetSymbolAttributeName()
	if str.data == nil || str.length == 0 {
		return ""
	}
	return C.GoStringN(str.data, C.int(str.length))
}

func SymbolTableGetVisibilityAttributeName() string {
	str := C.mlirSymbolTableGetVisibilityAttributeName()
	if str.data == nil || str.length == 0 {
		return ""
	}
	return C.GoStringN(str.data, C.int(str.length))
}

func SymbolTableCreate(op Operation) SymbolTable {
	return SymbolTable{c: C.mlirSymbolTableCreate(op.c)}
}

func SymbolTableDestroy(symbolTable SymbolTable) {
	C.mlirSymbolTableDestroy(symbolTable.c)
}

func SymbolTableLookup(symbolTable SymbolTable, name string) Operation {
	cstr := C.CString(name)
	defer C.free(unsafe.Pointer(cstr))
	return Operation{
		c: C.mlirSymbolTableLookup(symbolTable.c, C.mlirStringRefCreate(cstr, C.size_t(len(name)))),
	}
}

func SymbolTableInsert(symbolTable SymbolTable, op Operation) Attribute {
	return Attribute{c: C.mlirSymbolTableInsert(symbolTable.c, op.c)}
}

func SymbolTableErase(symbolTable SymbolTable, op Operation) {
	C.mlirSymbolTableErase(symbolTable.c, op.c)
}

func SymbolTableReplaceAllSymbolUses(oldSymbol, newSymbol string, from Operation) bool {
	cOld := C.CString(oldSymbol)
	defer C.free(unsafe.Pointer(cOld))
	cNew := C.CString(newSymbol)
	defer C.free(unsafe.Pointer(cNew))
	return bool(C.mlirLogicalResultIsSuccess(C.mlirSymbolTableReplaceAllSymbolUses(
		C.mlirStringRefCreate(cOld, C.size_t(len(oldSymbol))),
		C.mlirStringRefCreate(cNew, C.size_t(len(newSymbol))),
		from.c,
	)))
}

func SymbolTableWalkSymbolTables(from Operation, allSymUsesVisible bool, callback func(Operation, bool)) {
	handle := cgo.NewHandle(callback)
	defer handle.Delete()

	userData := C.malloc(C.size_t(unsafe.Sizeof(C.uintptr_t(0))))
	defer C.free(userData)
	*(*C.uintptr_t)(userData) = C.uintptr_t(handle)

	C.mlirGoSymbolTableWalk(from.c, C.bool(allSymUsesVisible), userData)
}
