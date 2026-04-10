//go:build cgo

package capi

/*
#include <stdlib.h>
#include <string.h>

#include "mlir-c/IR.h"
#include "mlir-c/Pass.h"
#include "mlir-c/ExecutionEngine.h"
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

static struct MlirGoStringBuffer mlirGoAttributeToString(MlirAttribute attr) {
	struct MlirGoStringBuffer buf = {0};
	mlirAttributePrint(attr, mlirGoStringCallback, &buf);
	return buf;
}

struct MlirGoParsePassPipelineResult {
	MlirLogicalResult result;
	struct MlirGoStringBuffer buffer;
};

static struct MlirGoParsePassPipelineResult mlirGoParsePassPipeline(
    MlirOpPassManager passManager, MlirStringRef pipeline) {
	struct MlirGoParsePassPipelineResult out = {0};
	out.result = mlirParsePassPipeline(passManager, pipeline, mlirGoStringCallback, &out.buffer);
	return out;
}

static struct MlirGoStringBuffer mlirGoPrintPassPipeline(MlirOpPassManager passManager) {
	struct MlirGoStringBuffer buf = {0};
	mlirPrintPassPipeline(passManager, mlirGoStringCallback, &buf);
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

type Attribute struct {
	c C.MlirAttribute
}

type NamedAttribute struct {
	c C.MlirNamedAttribute
}

type SymbolTable struct {
	c C.MlirSymbolTable
}

type PassManager struct {
	c C.MlirPassManager
}

type OpPassManager struct {
	c C.MlirOpPassManager
}

type ExecutionEngine struct {
	c C.MlirExecutionEngine
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

func RegisterAllPasses() {
	C.mlirRegisterAllPasses()
}

func RegisterAllLLVMTranslations(ctx Context) {
	C.mlirRegisterAllLLVMTranslations(ctx.c)
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

func OperationDestroy(op Operation) {
	C.mlirOperationDestroy(op.c)
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

func IdentifierGet(ctx Context, name string) Identifier {
	cstr := C.CString(name)
	defer C.free(unsafe.Pointer(cstr))
	return Identifier{
		c: C.mlirIdentifierGet(ctx.c, C.mlirStringRefCreate(cstr, C.size_t(len(name)))),
	}
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

func NullAttribute() Attribute {
	return Attribute{}
}

func AttributeParseGet(ctx Context, asm string) Attribute {
	cstr := C.CString(asm)
	defer C.free(unsafe.Pointer(cstr))

	return Attribute{
		c: C.mlirAttributeParseGet(ctx.c, C.mlirStringRefCreate(cstr, C.size_t(len(asm)))),
	}
}

func AttributeIsNull(attr Attribute) bool {
	return bool(C.mlirAttributeIsNull(attr.c))
}

func AttributeToString(attr Attribute) string {
	buf := C.mlirGoAttributeToString(attr.c)
	defer C.mlirGoStringBufferDestroy(buf)
	if buf.data == nil || buf.length == 0 {
		return ""
	}
	return C.GoStringN(buf.data, C.int(buf.length))
}

func NamedAttributeGet(name Identifier, attr Attribute) NamedAttribute {
	return NamedAttribute{c: C.mlirNamedAttributeGet(name.c, attr.c)}
}

func NamedAttributeName(attr NamedAttribute) Identifier {
	return Identifier{c: attr.c.name}
}

func NamedAttributeAttribute(attr NamedAttribute) Attribute {
	return Attribute{c: attr.c.attribute}
}

func RegionCreate() Region {
	return Region{c: C.mlirRegionCreate()}
}

func RegionDestroy(region Region) {
	C.mlirRegionDestroy(region.c)
}

func RegionAppendOwnedBlock(region Region, block Block) {
	C.mlirRegionAppendOwnedBlock(region.c, block.c)
}

func BlockCreate(argTypes []Type, argLocs []Location) Block {
	var typesPtr *C.MlirType
	var locsPtr *C.MlirLocation

	if len(argTypes) > 0 {
		cTypes := make([]C.MlirType, len(argTypes))
		for i, v := range argTypes {
			cTypes[i] = v.c
		}
		typesPtr = &cTypes[0]
		if len(argLocs) > 0 {
			cLocs := make([]C.MlirLocation, len(argLocs))
			for i, v := range argLocs {
				cLocs[i] = v.c
			}
			locsPtr = &cLocs[0]
			return Block{c: C.mlirBlockCreate(C.intptr_t(len(cTypes)), typesPtr, locsPtr)}
		}
		return Block{c: C.mlirBlockCreate(C.intptr_t(len(cTypes)), typesPtr, nil)}
	}

	return Block{c: C.mlirBlockCreate(0, nil, nil)}
}

func BlockDestroy(block Block) {
	C.mlirBlockDestroy(block.c)
}

func BlockAppendOwnedOperation(block Block, op Operation) {
	C.mlirBlockAppendOwnedOperation(block.c, op.c)
}

func BlockInsertOwnedOperationBefore(block Block, ref Operation, op Operation) {
	C.mlirBlockInsertOwnedOperationBefore(block.c, ref.c, op.c)
}

func BlockAddArgument(block Block, typ Type, loc Location) Value {
	return Value{c: C.mlirBlockAddArgument(block.c, typ.c, loc.c)}
}

func ValueReplaceAllUsesOfWith(from, to Value) {
	C.mlirValueReplaceAllUsesOfWith(from.c, to.c)
}

func OperationCreate(name string, loc Location, results []Type, operands []Value, regions []Region, successors []Block, attrs []NamedAttribute, inferResults bool) Operation {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	state := C.mlirOperationStateGet(C.mlirStringRefCreate(cname, C.size_t(len(name))), loc.c)

	if len(results) > 0 {
		cResults := make([]C.MlirType, len(results))
		for i, v := range results {
			cResults[i] = v.c
		}
		C.mlirOperationStateAddResults(&state, C.intptr_t(len(cResults)), &cResults[0])
	}
	if len(operands) > 0 {
		cOperands := make([]C.MlirValue, len(operands))
		for i, v := range operands {
			cOperands[i] = v.c
		}
		C.mlirOperationStateAddOperands(&state, C.intptr_t(len(cOperands)), &cOperands[0])
	}
	if len(regions) > 0 {
		cRegions := make([]C.MlirRegion, len(regions))
		for i, v := range regions {
			cRegions[i] = v.c
		}
		C.mlirOperationStateAddOwnedRegions(&state, C.intptr_t(len(cRegions)), &cRegions[0])
	}
	if len(successors) > 0 {
		cSuccessors := make([]C.MlirBlock, len(successors))
		for i, v := range successors {
			cSuccessors[i] = v.c
		}
		C.mlirOperationStateAddSuccessors(&state, C.intptr_t(len(cSuccessors)), &cSuccessors[0])
	}
	if len(attrs) > 0 {
		cAttrs := make([]C.MlirNamedAttribute, len(attrs))
		for i, v := range attrs {
			cAttrs[i] = v.c
		}
		C.mlirOperationStateAddAttributes(&state, C.intptr_t(len(cAttrs)), &cAttrs[0])
	}
	if inferResults {
		C.mlirOperationStateEnableResultTypeInference(&state)
	}

	return Operation{c: C.mlirOperationCreate(&state)}
}

func PassManagerCreate(ctx Context) PassManager {
	return PassManager{c: C.mlirPassManagerCreate(ctx.c)}
}

func PassManagerDestroy(pm PassManager) {
	C.mlirPassManagerDestroy(pm.c)
}

func PassManagerIsNull(pm PassManager) bool {
	return bool(C.mlirPassManagerIsNull(pm.c))
}

func PassManagerGetAsOpPassManager(pm PassManager) OpPassManager {
	return OpPassManager{c: C.mlirPassManagerGetAsOpPassManager(pm.c)}
}

func PassManagerRunOnOp(pm PassManager, op Operation) bool {
	return bool(C.mlirLogicalResultIsSuccess(C.mlirPassManagerRunOnOp(pm.c, op.c)))
}

func ParsePassPipeline(opm OpPassManager, pipeline string) (bool, string) {
	cstr := C.CString(pipeline)
	defer C.free(unsafe.Pointer(cstr))
	out := C.mlirGoParsePassPipeline(opm.c, C.mlirStringRefCreate(cstr, C.size_t(len(pipeline))))
	defer C.mlirGoStringBufferDestroy(out.buffer)
	msg := ""
	if out.buffer.data != nil && out.buffer.length > 0 {
		msg = C.GoStringN(out.buffer.data, C.int(out.buffer.length))
	}
	return bool(C.mlirLogicalResultIsSuccess(out.result)), msg
}

func PrintPassPipeline(opm OpPassManager) string {
	buf := C.mlirGoPrintPassPipeline(opm.c)
	defer C.mlirGoStringBufferDestroy(buf)
	if buf.data == nil || buf.length == 0 {
		return ""
	}
	return C.GoStringN(buf.data, C.int(buf.length))
}

func ExecutionEngineCreate(module Module, optLevel int) ExecutionEngine {
	return ExecutionEngine{
		c: C.mlirExecutionEngineCreate(module.c, C.int(optLevel), 0, nil, C.bool(false)),
	}
}

func ExecutionEngineDestroy(engine ExecutionEngine) {
	C.mlirExecutionEngineDestroy(engine.c)
}

func ExecutionEngineIsNull(engine ExecutionEngine) bool {
	return bool(C.mlirExecutionEngineIsNull(engine.c))
}

func ExecutionEngineInvokePacked(engine ExecutionEngine, name string, args []unsafe.Pointer) bool {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var cArgs unsafe.Pointer
	if len(args) > 0 {
		size := C.size_t(len(args)) * C.size_t(unsafe.Sizeof(uintptr(0)))
		cArgs = C.malloc(size)
		defer C.free(cArgs)
		slice := unsafe.Slice((*unsafe.Pointer)(cArgs), len(args))
		copy(slice, args)
	}

	result := C.mlirExecutionEngineInvokePacked(
		engine.c,
		C.mlirStringRefCreate(cname, C.size_t(len(name))),
		(*unsafe.Pointer)(cArgs),
	)
	return bool(C.mlirLogicalResultIsSuccess(result))
}
