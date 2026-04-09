//go:build cgo

package capi

/*
#include <stdlib.h>
#include <string.h>

#include "mlir-c/IR.h"
#include "mlir-c/RegisterEverything.h"

struct MlirGoStringBuffer {
	char *data;
	size_t length;
};

static void mlirGoStringCallback(MlirStringRef str, void *userData) {
	struct MlirGoStringBuffer *buf = (struct MlirGoStringBuffer *)userData;
	size_t newLength = buf->length + str.length;
	if (newLength == 0) {
		return;
	}
	buf->data = (char *)realloc(buf->data, newLength);
	memcpy(buf->data + buf->length, str.data, str.length);
	buf->length = newLength;
}

static struct MlirGoStringBuffer mlirGoOperationToString(MlirOperation op) {
	struct MlirGoStringBuffer buf = {0};
	mlirOperationPrint(op, mlirGoStringCallback, &buf);
	return buf;
}

static struct MlirGoStringBuffer mlirGoLocationToString(MlirLocation loc) {
	struct MlirGoStringBuffer buf = {0};
	mlirLocationPrint(loc, mlirGoStringCallback, &buf);
	return buf;
}

static struct MlirGoStringBuffer mlirGoBlockToString(MlirBlock block) {
	struct MlirGoStringBuffer buf = {0};
	mlirBlockPrint(block, mlirGoStringCallback, &buf);
	return buf;
}

static struct MlirGoStringBuffer mlirGoValueToString(MlirValue value) {
	struct MlirGoStringBuffer buf = {0};
	mlirValuePrint(value, mlirGoStringCallback, &buf);
	return buf;
}

static struct MlirGoStringBuffer mlirGoTypeToString(MlirType type) {
	struct MlirGoStringBuffer buf = {0};
	mlirTypePrint(type, mlirGoStringCallback, &buf);
	return buf;
}

static void mlirGoStringBufferDestroy(struct MlirGoStringBuffer buf) {
	free(buf.data);
}
*/
import "C"
import "unsafe"

type Context struct {
	c C.MlirContext
}

type DialectRegistry struct {
	c C.MlirDialectRegistry
}

type Module struct {
	c C.MlirModule
}

type Operation struct {
	c C.MlirOperation
}

type Location struct {
	c C.MlirLocation
}

type Region struct {
	c C.MlirRegion
}

type Block struct {
	c C.MlirBlock
}

type Identifier struct {
	c C.MlirIdentifier
}

type Value struct {
	c C.MlirValue
}

type Type struct {
	c C.MlirType
}

func NullContext() Context {
	return Context{}
}

func ContextCreate() Context {
	return Context{c: C.mlirContextCreate()}
}

func ContextDestroy(ctx Context) {
	C.mlirContextDestroy(ctx.c)
}

func ContextIsNull(ctx Context) bool {
	return bool(C.mlirContextIsNull(ctx.c))
}

func ContextSetAllowUnregisteredDialects(ctx Context, allow bool) {
	C.mlirContextSetAllowUnregisteredDialects(ctx.c, C.bool(allow))
}

func ContextAppendDialectRegistry(ctx Context, registry DialectRegistry) {
	C.mlirContextAppendDialectRegistry(ctx.c, registry.c)
}

func ContextLoadAllAvailableDialects(ctx Context) {
	C.mlirContextLoadAllAvailableDialects(ctx.c)
}

func ContextGetNumLoadedDialects(ctx Context) int {
	return int(C.mlirContextGetNumLoadedDialects(ctx.c))
}

func DialectRegistryCreate() DialectRegistry {
	return DialectRegistry{c: C.mlirDialectRegistryCreate()}
}

func DialectRegistryDestroy(registry DialectRegistry) {
	C.mlirDialectRegistryDestroy(registry.c)
}

func DialectRegistryIsNull(registry DialectRegistry) bool {
	return registry.c.ptr == nil
}

func RegisterAllDialects(registry DialectRegistry) {
	C.mlirRegisterAllDialects(registry.c)
}

func NullModule() Module {
	return Module{}
}

func ModuleCreateParse(ctx Context, asm string) Module {
	cstr := C.CString(asm)
	defer C.free(unsafe.Pointer(cstr))

	return Module{
		c: C.mlirModuleCreateParse(ctx.c, C.mlirStringRefCreate(cstr, C.size_t(len(asm)))),
	}
}

func ModuleCreateEmpty(loc Location) Module {
	return Module{c: C.mlirModuleCreateEmpty(loc.c)}
}

func ModuleIsNull(module Module) bool {
	return bool(C.mlirModuleIsNull(module.c))
}

func ModuleDestroy(module Module) {
	C.mlirModuleDestroy(module.c)
}

func ModuleGetOperation(module Module) Operation {
	return Operation{c: C.mlirModuleGetOperation(module.c)}
}

func ModuleGetBody(module Module) Block {
	return Block{c: C.mlirModuleGetBody(module.c)}
}

func OperationIsNull(op Operation) bool {
	return bool(C.mlirOperationIsNull(op.c))
}

func OperationToString(op Operation) string {
	buf := C.mlirGoOperationToString(op.c)
	defer C.mlirGoStringBufferDestroy(buf)

	if buf.data == nil || buf.length == 0 {
		return ""
	}
	return C.GoStringN(buf.data, C.int(buf.length))
}

func OperationVerify(op Operation) bool {
	return bool(C.mlirOperationVerify(op.c))
}

func OperationGetName(op Operation) Identifier {
	return Identifier{c: C.mlirOperationGetName(op.c)}
}

func OperationGetLocation(op Operation) Location {
	return Location{c: C.mlirOperationGetLocation(op.c)}
}

func OperationGetNumRegions(op Operation) int {
	return int(C.mlirOperationGetNumRegions(op.c))
}

func OperationGetRegion(op Operation, pos int) Region {
	return Region{c: C.mlirOperationGetRegion(op.c, C.intptr_t(pos))}
}

func OperationGetNextInBlock(op Operation) Operation {
	return Operation{c: C.mlirOperationGetNextInBlock(op.c)}
}

func NullLocation() Location {
	return Location{}
}

func LocationUnknownGet(ctx Context) Location {
	return Location{c: C.mlirLocationUnknownGet(ctx.c)}
}

func LocationFileLineColGet(ctx Context, filename string, line, col uint) Location {
	cstr := C.CString(filename)
	defer C.free(unsafe.Pointer(cstr))
	return Location{
		c: C.mlirLocationFileLineColGet(
			ctx.c,
			C.mlirStringRefCreate(cstr, C.size_t(len(filename))),
			C.uint(line),
			C.uint(col),
		),
	}
}

func LocationIsNull(loc Location) bool {
	return bool(C.mlirLocationIsNull(loc.c))
}

func LocationToString(loc Location) string {
	buf := C.mlirGoLocationToString(loc.c)
	defer C.mlirGoStringBufferDestroy(buf)

	if buf.data == nil || buf.length == 0 {
		return ""
	}
	return C.GoStringN(buf.data, C.int(buf.length))
}

func RegionIsNull(region Region) bool {
	return bool(C.mlirRegionIsNull(region.c))
}

func RegionGetFirstBlock(region Region) Block {
	return Block{c: C.mlirRegionGetFirstBlock(region.c)}
}

func RegionGetNextInOperation(region Region) Region {
	return Region{c: C.mlirRegionGetNextInOperation(region.c)}
}

func BlockIsNull(block Block) bool {
	return bool(C.mlirBlockIsNull(block.c))
}

func BlockGetFirstOperation(block Block) Operation {
	return Operation{c: C.mlirBlockGetFirstOperation(block.c)}
}

func BlockGetNextInRegion(block Block) Block {
	return Block{c: C.mlirBlockGetNextInRegion(block.c)}
}

func BlockToString(block Block) string {
	buf := C.mlirGoBlockToString(block.c)
	defer C.mlirGoStringBufferDestroy(buf)

	if buf.data == nil || buf.length == 0 {
		return ""
	}
	return C.GoStringN(buf.data, C.int(buf.length))
}

func IdentifierToString(id Identifier) string {
	str := C.mlirIdentifierStr(id.c)
	if str.data == nil || str.length == 0 {
		return ""
	}
	return C.GoStringN(str.data, C.int(str.length))
}

func BlockGetNumArguments(block Block) int {
	return int(C.mlirBlockGetNumArguments(block.c))
}

func BlockGetArgument(block Block, pos int) Value {
	return Value{c: C.mlirBlockGetArgument(block.c, C.intptr_t(pos))}
}

func NullValue() Value {
	return Value{}
}

func ValueIsNull(value Value) bool {
	return bool(C.mlirValueIsNull(value.c))
}

func ValueIsBlockArgument(value Value) bool {
	return bool(C.mlirValueIsABlockArgument(value.c))
}

func ValueIsOpResult(value Value) bool {
	return bool(C.mlirValueIsAOpResult(value.c))
}

func ValueGetType(value Value) Type {
	return Type{c: C.mlirValueGetType(value.c)}
}

func ValueToString(value Value) string {
	buf := C.mlirGoValueToString(value.c)
	defer C.mlirGoStringBufferDestroy(buf)

	if buf.data == nil || buf.length == 0 {
		return ""
	}
	return C.GoStringN(buf.data, C.int(buf.length))
}

func OperationGetNumResults(op Operation) int {
	return int(C.mlirOperationGetNumResults(op.c))
}

func OperationGetResult(op Operation, pos int) Value {
	return Value{c: C.mlirOperationGetResult(op.c, C.intptr_t(pos))}
}

func NullType() Type {
	return Type{}
}

func TypeParseGet(ctx Context, asm string) Type {
	cstr := C.CString(asm)
	defer C.free(unsafe.Pointer(cstr))

	return Type{
		c: C.mlirTypeParseGet(ctx.c, C.mlirStringRefCreate(cstr, C.size_t(len(asm)))),
	}
}

func TypeIsNull(typ Type) bool {
	return bool(C.mlirTypeIsNull(typ.c))
}

func TypeEqual(lhs, rhs Type) bool {
	return bool(C.mlirTypeEqual(lhs.c, rhs.c))
}

func TypeToString(typ Type) string {
	buf := C.mlirGoTypeToString(typ.c)
	defer C.mlirGoStringBufferDestroy(buf)

	if buf.data == nil || buf.length == 0 {
		return ""
	}
	return C.GoStringN(buf.data, C.int(buf.length))
}
