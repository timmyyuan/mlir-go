//go:build cgo

package capi

/*
#include <stdint.h>
#include <stdlib.h>

#include "mlir-c/Diagnostics.h"
*/
import "C"
import (
	"runtime/cgo"
	"unsafe"
)

//export mlirGoDiagnosticHandler
func mlirGoDiagnosticHandler(diag C.MlirDiagnostic, userData unsafe.Pointer) C.MlirLogicalResult {
	handle := cgo.Handle(*(*C.uintptr_t)(userData))
	callback := handle.Value().(func(Diagnostic))
	callback(Diagnostic{c: diag})
	return C.mlirLogicalResultSuccess()
}

//export mlirGoDiagnosticUserDataDestroy
func mlirGoDiagnosticUserDataDestroy(userData unsafe.Pointer) {
	handle := cgo.Handle(*(*C.uintptr_t)(userData))
	handle.Delete()
	C.free(userData)
}
