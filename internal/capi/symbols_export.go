//go:build cgo

package capi

/*
#include <stdbool.h>
#include <stdint.h>

#include "mlir-c/IR.h"
*/
import "C"
import (
	"runtime/cgo"
	"unsafe"
)

//export mlirGoSymbolTableWalkCallback
func mlirGoSymbolTableWalkCallback(op C.MlirOperation, allSymUsesVisible C.bool, userData unsafe.Pointer) {
	handle := cgo.Handle(*(*C.uintptr_t)(userData))
	callback := handle.Value().(func(Operation, bool))
	callback(Operation{c: op}, bool(allSymUsesVisible))
}
