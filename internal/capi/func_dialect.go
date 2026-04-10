//go:build cgo

package capi

/*
#include <stdlib.h>

#include "mlir-c/Dialect/Func.h"
*/
import "C"
import "unsafe"

func FuncSetArgAttr(op Operation, pos int, name string, attr Attribute) {
	cstr := C.CString(name)
	defer C.free(unsafe.Pointer(cstr))
	C.mlirFuncSetArgAttr(op.c, C.intptr_t(pos), C.mlirStringRefCreate(cstr, C.size_t(len(name))), attr.c)
}
