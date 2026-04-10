//go:build cgo

package capi

/*
#include <stdlib.h>
#include <stdint.h>
#include <string.h>

#include "mlir-c/Diagnostics.h"

extern MlirLogicalResult mlirGoDiagnosticHandler(MlirDiagnostic diagnostic, void *userData);
extern void mlirGoDiagnosticUserDataDestroy(void *userData);

struct MlirGoDiagnosticStringBuffer {
	char *data;
	size_t length;
};

static void mlirGoDiagnosticStringCallback(MlirStringRef str, void *userData) {
	struct MlirGoDiagnosticStringBuffer *buf = (struct MlirGoDiagnosticStringBuffer *)userData;
	size_t newLength = buf->length + str.length;
	if (newLength == 0) {
		return;
	}
	buf->data = (char *)realloc(buf->data, newLength);
	memcpy(buf->data + buf->length, str.data, str.length);
	buf->length = newLength;
}

static struct MlirGoDiagnosticStringBuffer mlirGoDiagnosticToString(MlirDiagnostic diagnostic) {
	struct MlirGoDiagnosticStringBuffer buf = {0};
	mlirDiagnosticPrint(diagnostic, mlirGoDiagnosticStringCallback, &buf);
	return buf;
}

static void mlirGoDiagnosticStringBufferDestroy(struct MlirGoDiagnosticStringBuffer buf) {
	free(buf.data);
}

static MlirDiagnosticHandlerID mlirGoContextAttachDiagnosticHandler(MlirContext context, void *userData) {
	return mlirContextAttachDiagnosticHandler(context, mlirGoDiagnosticHandler, userData, mlirGoDiagnosticUserDataDestroy);
}
*/
import "C"
import (
	"runtime/cgo"
	"unsafe"
)

type Diagnostic struct {
	c C.MlirDiagnostic
}

type DiagnosticSeverity int

const (
	DiagnosticSeverityError   DiagnosticSeverity = C.MlirDiagnosticError
	DiagnosticSeverityWarning DiagnosticSeverity = C.MlirDiagnosticWarning
	DiagnosticSeverityNote    DiagnosticSeverity = C.MlirDiagnosticNote
	DiagnosticSeverityRemark  DiagnosticSeverity = C.MlirDiagnosticRemark
)

type DiagnosticHandlerID uint64

func DiagnosticToString(diag Diagnostic) string {
	buf := C.mlirGoDiagnosticToString(diag.c)
	defer C.mlirGoDiagnosticStringBufferDestroy(buf)

	if buf.data == nil || buf.length == 0 {
		return ""
	}
	return C.GoStringN(buf.data, C.int(buf.length))
}

func DiagnosticGetLocation(diag Diagnostic) Location {
	return Location{c: C.mlirDiagnosticGetLocation(diag.c)}
}

func DiagnosticGetSeverity(diag Diagnostic) DiagnosticSeverity {
	switch C.mlirDiagnosticGetSeverity(diag.c) {
	case C.MlirDiagnosticWarning:
		return DiagnosticSeverityWarning
	case C.MlirDiagnosticNote:
		return DiagnosticSeverityNote
	case C.MlirDiagnosticRemark:
		return DiagnosticSeverityRemark
	default:
		return DiagnosticSeverityError
	}
}

func DiagnosticGetNumNotes(diag Diagnostic) int {
	return int(C.mlirDiagnosticGetNumNotes(diag.c))
}

func DiagnosticGetNote(diag Diagnostic, pos int) Diagnostic {
	return Diagnostic{c: C.mlirDiagnosticGetNote(diag.c, C.intptr_t(pos))}
}

func ContextAttachDiagnosticCallback(ctx Context, callback func(Diagnostic)) DiagnosticHandlerID {
	handle := cgo.NewHandle(callback)
	userData := C.malloc(C.size_t(unsafe.Sizeof(C.uintptr_t(0))))
	*(*C.uintptr_t)(userData) = C.uintptr_t(handle)
	return DiagnosticHandlerID(C.mlirGoContextAttachDiagnosticHandler(ctx.c, userData))
}

func ContextDetachDiagnosticHandler(ctx Context, id DiagnosticHandlerID) {
	C.mlirContextDetachDiagnosticHandler(ctx.c, C.MlirDiagnosticHandlerID(id))
}

func EmitError(loc Location, message string) {
	cstr := C.CString(message)
	defer C.free(unsafe.Pointer(cstr))
	C.mlirEmitError(loc.c, cstr)
}
