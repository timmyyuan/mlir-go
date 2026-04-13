package gen

import (
	"bytes"
	"fmt"
	"go/format"
	"path/filepath"
	"slices"
)

// GeneratedFile is one checked-in output emitted by the generator.
type GeneratedFile struct {
	Path    string
	Content []byte
}

// SupportedOp reports whether the generator can currently emit code for defName.
func SupportedOp(dialectName, defName string) bool {
	switch dialectName {
	case "arith":
		return defName == "Arith_ConstantOp" ||
			defName == "Arith_AddIOp" ||
			defName == "Arith_SubIOp" ||
			defName == "Arith_MulIOp" ||
			defName == "Arith_ExtUIOp" ||
			defName == "Arith_ExtSIOp" ||
			defName == "Arith_TruncIOp" ||
			defName == "Arith_SIToFPOp" ||
			defName == "Arith_FPToSIOp" ||
			defName == "Arith_IndexCastOp" ||
			defName == "Arith_CmpIOp"
	case "func":
		return defName == "FuncOp" || defName == "CallOp" || defName == "ReturnOp"
	case "cf":
		return defName == "BranchOp" || defName == "CondBranchOp"
	case "memref":
		return defName == "MemRef_AllocaOp" || defName == "LoadOp" || defName == "MemRef_StoreOp"
	case "tensor":
		return defName == "Tensor_EmptyOp" || defName == "Tensor_DimOp" || defName == "Tensor_CastOp"
	case "vector":
		return defName == "Vector_BroadcastOp" || defName == "Vector_ExtractElementOp"
	case "linalg":
		return defName == "Linalg_YieldOp" || defName == "Linalg_IndexOp"
	default:
		return false
	}
}

// FindDialect returns the manifest entry for dialectName.
func (m *Manifest) FindDialect(dialectName string) (DialectManifest, bool) {
	if m == nil {
		return DialectManifest{}, false
	}
	for _, dialect := range m.Dialects {
		if dialect.Name == dialectName {
			return dialect, true
		}
	}
	return DialectManifest{}, false
}

// EmitDialect emits checked-in generated files for one supported dialect.
func EmitDialect(includeRoot string, dialect DialectManifest) ([]GeneratedFile, error) {
	ops, err := ScanDialect(includeRoot, dialect)
	if err != nil {
		return nil, err
	}

	switch dialect.Name {
	case "arith":
		return emitArithDialect(ops)
	case "func":
		return emitFuncDialect(ops)
	case "cf":
		return emitCFDialect(ops)
	case "memref":
		return emitMemRefDialect(ops)
	case "tensor":
		return emitTensorDialect(ops)
	case "vector":
		return emitVectorDialect(ops)
	case "linalg":
		return emitLinalgDialect(ops)
	default:
		return nil, fmt.Errorf("dialect %q emission is not implemented yet", dialect.Name)
	}
}

func emitArithDialect(ops []OperationDef) ([]GeneratedFile, error) {
	defNames := make([]string, 0, len(ops))
	for _, op := range ops {
		defNames = append(defNames, op.DefName)
	}
	for _, required := range []string{"Arith_ConstantOp", "Arith_AddIOp"} {
		if !slices.Contains(defNames, required) {
			return nil, fmt.Errorf("arith source is missing required op definition %q", required)
		}
	}
	for _, required := range []string{
		"Arith_SubIOp",
		"Arith_MulIOp",
		"Arith_ExtUIOp",
		"Arith_ExtSIOp",
		"Arith_TruncIOp",
		"Arith_SIToFPOp",
		"Arith_FPToSIOp",
		"Arith_IndexCastOp",
		"Arith_CmpIOp",
	} {
		if !slices.Contains(defNames, required) {
			return nil, fmt.Errorf("arith source is missing required op definition %q", required)
		}
	}

	cgoContent, err := format.Source([]byte(renderArithGenerated()))
	if err != nil {
		return nil, fmt.Errorf("format arith generated cgo file: %w", err)
	}
	nocgoContent, err := format.Source([]byte(renderArithGeneratedNoCgo()))
	if err != nil {
		return nil, fmt.Errorf("format arith generated nocgo file: %w", err)
	}

	return []GeneratedFile{
		{
			Path:    filepath.ToSlash(filepath.Join("dialect", "arith", "arith_gen.go")),
			Content: cgoContent,
		},
		{
			Path:    filepath.ToSlash(filepath.Join("dialect", "arith", "arith_gen_nocgo.go")),
			Content: nocgoContent,
		},
	}, nil
}

func emitMemRefDialect(ops []OperationDef) ([]GeneratedFile, error) {
	defNames := make([]string, 0, len(ops))
	for _, op := range ops {
		defNames = append(defNames, op.DefName)
	}
	for _, required := range []string{"MemRef_AllocaOp", "LoadOp", "MemRef_StoreOp"} {
		if !slices.Contains(defNames, required) {
			return nil, fmt.Errorf("memref source is missing required op definition %q", required)
		}
	}

	cgoContent, err := format.Source([]byte(renderMemRefGenerated()))
	if err != nil {
		return nil, fmt.Errorf("format memref generated cgo file: %w", err)
	}
	nocgoContent, err := format.Source([]byte(renderMemRefGeneratedNoCgo()))
	if err != nil {
		return nil, fmt.Errorf("format memref generated nocgo file: %w", err)
	}
	return []GeneratedFile{
		{
			Path:    filepath.ToSlash(filepath.Join("dialect", "memref", "memref_gen.go")),
			Content: cgoContent,
		},
		{
			Path:    filepath.ToSlash(filepath.Join("dialect", "memref", "memref_gen_nocgo.go")),
			Content: nocgoContent,
		},
	}, nil
}

func emitTensorDialect(ops []OperationDef) ([]GeneratedFile, error) {
	defNames := make([]string, 0, len(ops))
	for _, op := range ops {
		defNames = append(defNames, op.DefName)
	}
	for _, required := range []string{"Tensor_EmptyOp", "Tensor_DimOp", "Tensor_CastOp"} {
		if !slices.Contains(defNames, required) {
			return nil, fmt.Errorf("tensor source is missing required op definition %q", required)
		}
	}

	cgoContent, err := format.Source([]byte(renderTensorGenerated()))
	if err != nil {
		return nil, fmt.Errorf("format tensor generated cgo file: %w", err)
	}
	nocgoContent, err := format.Source([]byte(renderTensorGeneratedNoCgo()))
	if err != nil {
		return nil, fmt.Errorf("format tensor generated nocgo file: %w", err)
	}
	return []GeneratedFile{
		{
			Path:    filepath.ToSlash(filepath.Join("dialect", "tensor", "tensor_gen.go")),
			Content: cgoContent,
		},
		{
			Path:    filepath.ToSlash(filepath.Join("dialect", "tensor", "tensor_gen_nocgo.go")),
			Content: nocgoContent,
		},
	}, nil
}

func emitVectorDialect(ops []OperationDef) ([]GeneratedFile, error) {
	defNames := make([]string, 0, len(ops))
	for _, op := range ops {
		defNames = append(defNames, op.DefName)
	}
	for _, required := range []string{"Vector_BroadcastOp", "Vector_ExtractElementOp"} {
		if !slices.Contains(defNames, required) {
			return nil, fmt.Errorf("vector source is missing required op definition %q", required)
		}
	}

	cgoContent, err := format.Source([]byte(renderVectorGenerated()))
	if err != nil {
		return nil, fmt.Errorf("format vector generated cgo file: %w", err)
	}
	nocgoContent, err := format.Source([]byte(renderVectorGeneratedNoCgo()))
	if err != nil {
		return nil, fmt.Errorf("format vector generated nocgo file: %w", err)
	}
	return []GeneratedFile{
		{
			Path:    filepath.ToSlash(filepath.Join("dialect", "vector", "vector_gen.go")),
			Content: cgoContent,
		},
		{
			Path:    filepath.ToSlash(filepath.Join("dialect", "vector", "vector_gen_nocgo.go")),
			Content: nocgoContent,
		},
	}, nil
}

func emitLinalgDialect(ops []OperationDef) ([]GeneratedFile, error) {
	defNames := make([]string, 0, len(ops))
	for _, op := range ops {
		defNames = append(defNames, op.DefName)
	}
	for _, required := range []string{"Linalg_YieldOp", "Linalg_IndexOp"} {
		if !slices.Contains(defNames, required) {
			return nil, fmt.Errorf("linalg source is missing required op definition %q", required)
		}
	}

	cgoContent, err := format.Source([]byte(renderLinalgGenerated()))
	if err != nil {
		return nil, fmt.Errorf("format linalg generated cgo file: %w", err)
	}
	nocgoContent, err := format.Source([]byte(renderLinalgGeneratedNoCgo()))
	if err != nil {
		return nil, fmt.Errorf("format linalg generated nocgo file: %w", err)
	}
	return []GeneratedFile{
		{
			Path:    filepath.ToSlash(filepath.Join("dialect", "linalg", "linalg_gen.go")),
			Content: cgoContent,
		},
		{
			Path:    filepath.ToSlash(filepath.Join("dialect", "linalg", "linalg_gen_nocgo.go")),
			Content: nocgoContent,
		},
	}, nil
}

func emitFuncDialect(ops []OperationDef) ([]GeneratedFile, error) {
	defNames := make([]string, 0, len(ops))
	for _, op := range ops {
		defNames = append(defNames, op.DefName)
	}
	for _, required := range []string{"FuncOp", "CallOp", "ReturnOp"} {
		if !slices.Contains(defNames, required) {
			return nil, fmt.Errorf("func source is missing required op definition %q", required)
		}
	}

	cgoContent, err := format.Source([]byte(renderFuncGenerated()))
	if err != nil {
		return nil, fmt.Errorf("format func generated cgo file: %w", err)
	}
	nocgoContent, err := format.Source([]byte(renderFuncGeneratedNoCgo()))
	if err != nil {
		return nil, fmt.Errorf("format func generated nocgo file: %w", err)
	}

	return []GeneratedFile{
		{
			Path:    filepath.ToSlash(filepath.Join("dialect", "func", "func_gen.go")),
			Content: cgoContent,
		},
		{
			Path:    filepath.ToSlash(filepath.Join("dialect", "func", "func_gen_nocgo.go")),
			Content: nocgoContent,
		},
	}, nil
}

func emitCFDialect(ops []OperationDef) ([]GeneratedFile, error) {
	defNames := make([]string, 0, len(ops))
	for _, op := range ops {
		defNames = append(defNames, op.DefName)
	}
	for _, required := range []string{"BranchOp", "CondBranchOp"} {
		if !slices.Contains(defNames, required) {
			return nil, fmt.Errorf("cf source is missing required op definition %q", required)
		}
	}

	cgoContent, err := format.Source([]byte(renderCFGenerated()))
	if err != nil {
		return nil, fmt.Errorf("format cf generated cgo file: %w", err)
	}
	nocgoContent, err := format.Source([]byte(renderCFGeneratedNoCgo()))
	if err != nil {
		return nil, fmt.Errorf("format cf generated nocgo file: %w", err)
	}

	return []GeneratedFile{
		{
			Path:    filepath.ToSlash(filepath.Join("dialect", "cf", "cf_gen.go")),
			Content: cgoContent,
		},
		{
			Path:    filepath.ToSlash(filepath.Join("dialect", "cf", "cf_gen_nocgo.go")),
			Content: nocgoContent,
		},
	}, nil
}

func renderArithGenerated() string {
	var buf bytes.Buffer
	buf.WriteString("// Code generated by cmd/mlir-go-tblgen; DO NOT EDIT.\n")
	buf.WriteString("//go:build cgo\n\n")
	buf.WriteString("package arith\n\n")
	buf.WriteString("import (\n")
	buf.WriteString("\t\"fmt\"\n\n")
	buf.WriteString("\tmlir \"github.com/timmyyuan/mlir-go\"\n")
	buf.WriteString(")\n\n")

	buf.WriteString("func ConstantOp(ctx *mlir.Context, loc mlir.Location, resultType mlir.Type, value mlir.Attribute) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\tif resultType.IsNull() {\n")
	buf.WriteString("\t\treturn nil, fmt.Errorf(\"mlir: result type is required\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tif value.IsNull() {\n")
	buf.WriteString("\t\treturn nil, fmt.Errorf(\"mlir: value attribute is required\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tvalueAttr, err := mlir.NamedAttributeByName(ctx, \"value\", value)\n")
	buf.WriteString("\tif err != nil {\n")
	buf.WriteString("\t\treturn nil, err\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tstate := mlir.NewOperationState(\"arith.constant\", loc)\n")
	buf.WriteString("\tstate.AddResults(resultType)\n")
	buf.WriteString("\tstate.AddAttributes(valueAttr)\n")
	buf.WriteString("\treturn mlir.CreateOperation(state)\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func AddIOp(loc mlir.Location, lhs, rhs mlir.Value) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn emitSameTypeBinaryOp(\"arith.addi\", loc, lhs, rhs)\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func SubIOp(loc mlir.Location, lhs, rhs mlir.Value) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn emitSameTypeBinaryOp(\"arith.subi\", loc, lhs, rhs)\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func MulIOp(loc mlir.Location, lhs, rhs mlir.Value) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn emitSameTypeBinaryOp(\"arith.muli\", loc, lhs, rhs)\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func ExtUIOp(loc mlir.Location, input mlir.Value, resultType mlir.Type) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn emitUnaryCastOp(\"arith.extui\", loc, input, resultType)\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func ExtSIOp(loc mlir.Location, input mlir.Value, resultType mlir.Type) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn emitUnaryCastOp(\"arith.extsi\", loc, input, resultType)\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func TruncIOp(loc mlir.Location, input mlir.Value, resultType mlir.Type) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn emitUnaryCastOp(\"arith.trunci\", loc, input, resultType)\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func SIToFPOp(loc mlir.Location, input mlir.Value, resultType mlir.Type) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn emitUnaryCastOp(\"arith.sitofp\", loc, input, resultType)\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func FPToSIOp(loc mlir.Location, input mlir.Value, resultType mlir.Type) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn emitUnaryCastOp(\"arith.fptosi\", loc, input, resultType)\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func IndexCastOp(loc mlir.Location, input mlir.Value, resultType mlir.Type) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn emitUnaryCastOp(\"arith.index_cast\", loc, input, resultType)\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func CmpIOp(ctx *mlir.Context, loc mlir.Location, predicate string, lhs, rhs mlir.Value) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\tif lhs.IsNull() || rhs.IsNull() {\n")
	buf.WriteString("\t\treturn nil, fmt.Errorf(\"mlir: operands are required\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tresultType, err := cmpIResultType(ctx, lhs.Type(), rhs.Type())\n")
	buf.WriteString("\tif err != nil {\n")
	buf.WriteString("\t\treturn nil, err\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tpredicateAttr, err := cmpIPredicateAttr(ctx, predicate)\n")
	buf.WriteString("\tif err != nil {\n")
	buf.WriteString("\t\treturn nil, err\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tpredicateNamedAttr, err := mlir.NamedAttributeByName(ctx, \"predicate\", predicateAttr)\n")
	buf.WriteString("\tif err != nil {\n")
	buf.WriteString("\t\treturn nil, err\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tstate := mlir.NewOperationState(\"arith.cmpi\", loc)\n")
	buf.WriteString("\tstate.AddResults(resultType)\n")
	buf.WriteString("\tstate.AddOperands(lhs, rhs)\n")
	buf.WriteString("\tstate.AddAttributes(predicateNamedAttr)\n")
	buf.WriteString("\treturn mlir.CreateOperation(state)\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func emitSameTypeBinaryOp(name string, loc mlir.Location, lhs, rhs mlir.Value) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\tif lhs.IsNull() || rhs.IsNull() {\n")
	buf.WriteString("\t\treturn nil, fmt.Errorf(\"mlir: operands are required\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tresultType := lhs.Type()\n")
	buf.WriteString("\tif resultType.IsNull() || !resultType.Equal(rhs.Type()) {\n")
	buf.WriteString("\t\treturn nil, fmt.Errorf(\"mlir: operands must have the same non-null type\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tstate := mlir.NewOperationState(name, loc)\n")
	buf.WriteString("\tstate.AddResults(resultType)\n")
	buf.WriteString("\tstate.AddOperands(lhs, rhs)\n")
	buf.WriteString("\treturn mlir.CreateOperation(state)\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func emitUnaryCastOp(name string, loc mlir.Location, input mlir.Value, resultType mlir.Type) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\tif input.IsNull() {\n")
	buf.WriteString("\t\treturn nil, fmt.Errorf(\"mlir: input operand is required\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tif resultType.IsNull() {\n")
	buf.WriteString("\t\treturn nil, fmt.Errorf(\"mlir: result type is required\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tstate := mlir.NewOperationState(name, loc)\n")
	buf.WriteString("\tstate.AddResults(resultType)\n")
	buf.WriteString("\tstate.AddOperands(input)\n")
	buf.WriteString("\treturn mlir.CreateOperation(state)\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func cmpIPredicateAttr(ctx *mlir.Context, predicate string) (mlir.Attribute, error) {\n")
	buf.WriteString("\ti64, err := mlir.SignlessIntegerType(ctx, 64)\n")
	buf.WriteString("\tif err != nil {\n")
	buf.WriteString("\t\treturn mlir.Attribute{}, err\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tcode, ok := map[string]int64{\"eq\": 0, \"ne\": 1, \"slt\": 2, \"sle\": 3, \"sgt\": 4, \"sge\": 5, \"ult\": 6, \"ule\": 7, \"ugt\": 8, \"uge\": 9}[predicate]\n")
	buf.WriteString("\tif !ok {\n")
	buf.WriteString("\t\treturn mlir.Attribute{}, fmt.Errorf(\"mlir: unsupported cmpi predicate %q\", predicate)\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\treturn mlir.IntegerAttribute(i64, code)\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func cmpIResultType(ctx *mlir.Context, lhsType, rhsType mlir.Type) (mlir.Type, error) {\n")
	buf.WriteString("\tif lhsType.IsNull() || rhsType.IsNull() || !lhsType.Equal(rhsType) {\n")
	buf.WriteString("\t\treturn mlir.Type{}, fmt.Errorf(\"mlir: cmpi operands must have the same non-null type\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tif lhsType.IsInteger() || lhsType.IsIndex() {\n")
	buf.WriteString("\t\treturn mlir.SignlessIntegerType(ctx, 1)\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tif lhsType.IsRankedTensor() {\n")
	buf.WriteString("\t\ti1, err := mlir.SignlessIntegerType(ctx, 1)\n")
	buf.WriteString("\t\tif err != nil {\n")
	buf.WriteString("\t\t\treturn mlir.Type{}, err\n")
	buf.WriteString("\t\t}\n")
	buf.WriteString("\t\treturn mlir.RankedTensorType(lhsType.Shape(), i1)\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\treturn mlir.Type{}, fmt.Errorf(\"mlir: cmpi currently supports scalar integers, index values, and ranked tensors\")\n")
	buf.WriteString("}\n")
	return buf.String()
}

func renderArithGeneratedNoCgo() string {
	var buf bytes.Buffer
	buf.WriteString("// Code generated by cmd/mlir-go-tblgen; DO NOT EDIT.\n")
	buf.WriteString("//go:build !cgo\n\n")
	buf.WriteString("package arith\n\n")
	buf.WriteString("import (\n")
	buf.WriteString("\t\"fmt\"\n\n")
	buf.WriteString("\tmlir \"github.com/timmyyuan/mlir-go\"\n")
	buf.WriteString(")\n\n")

	buf.WriteString("func ConstantOp(*mlir.Context, mlir.Location, mlir.Type, mlir.Attribute) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn nil, fmt.Errorf(\"mlir: cgo is required\")\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func AddIOp(mlir.Location, mlir.Value, mlir.Value) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn nil, fmt.Errorf(\"mlir: cgo is required\")\n")
	buf.WriteString("}\n\n")
	buf.WriteString("func SubIOp(mlir.Location, mlir.Value, mlir.Value) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn nil, fmt.Errorf(\"mlir: cgo is required\")\n")
	buf.WriteString("}\n\n")
	buf.WriteString("func MulIOp(mlir.Location, mlir.Value, mlir.Value) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn nil, fmt.Errorf(\"mlir: cgo is required\")\n")
	buf.WriteString("}\n\n")
	buf.WriteString("func ExtUIOp(mlir.Location, mlir.Value, mlir.Type) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn nil, fmt.Errorf(\"mlir: cgo is required\")\n")
	buf.WriteString("}\n\n")
	buf.WriteString("func ExtSIOp(mlir.Location, mlir.Value, mlir.Type) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn nil, fmt.Errorf(\"mlir: cgo is required\")\n")
	buf.WriteString("}\n\n")
	buf.WriteString("func TruncIOp(mlir.Location, mlir.Value, mlir.Type) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn nil, fmt.Errorf(\"mlir: cgo is required\")\n")
	buf.WriteString("}\n\n")
	buf.WriteString("func SIToFPOp(mlir.Location, mlir.Value, mlir.Type) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn nil, fmt.Errorf(\"mlir: cgo is required\")\n")
	buf.WriteString("}\n\n")
	buf.WriteString("func FPToSIOp(mlir.Location, mlir.Value, mlir.Type) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn nil, fmt.Errorf(\"mlir: cgo is required\")\n")
	buf.WriteString("}\n\n")
	buf.WriteString("func IndexCastOp(mlir.Location, mlir.Value, mlir.Type) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn nil, fmt.Errorf(\"mlir: cgo is required\")\n")
	buf.WriteString("}\n\n")
	buf.WriteString("func CmpIOp(*mlir.Context, mlir.Location, string, mlir.Value, mlir.Value) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn nil, fmt.Errorf(\"mlir: cgo is required\")\n")
	buf.WriteString("}\n")
	return buf.String()
}

func renderFuncGenerated() string {
	var buf bytes.Buffer
	buf.WriteString("// Code generated by cmd/mlir-go-tblgen; DO NOT EDIT.\n")
	buf.WriteString("//go:build cgo\n\n")
	buf.WriteString("package funcdialect\n\n")
	buf.WriteString("import (\n")
	buf.WriteString("\t\"fmt\"\n\n")
	buf.WriteString("\tmlir \"github.com/timmyyuan/mlir-go\"\n")
	buf.WriteString(")\n\n")

	buf.WriteString("func FuncOp(ctx *mlir.Context, loc mlir.Location, name string, functionType mlir.Type, body *mlir.OwnedRegion) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\tif name == \"\" {\n")
	buf.WriteString("\t\treturn nil, fmt.Errorf(\"mlir: function name is required\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tif functionType.IsNull() {\n")
	buf.WriteString("\t\treturn nil, fmt.Errorf(\"mlir: function type is required\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tif body == nil {\n")
	buf.WriteString("\t\treturn nil, fmt.Errorf(\"mlir: function body region is required\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tsymName, err := mlir.StringAttribute(ctx, name)\n")
	buf.WriteString("\tif err != nil {\n")
	buf.WriteString("\t\treturn nil, err\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tfunctionTypeAttr, err := mlir.TypeAttribute(functionType)\n")
	buf.WriteString("\tif err != nil {\n")
	buf.WriteString("\t\treturn nil, err\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tnameAttr, err := mlir.NamedAttributeByName(ctx, \"sym_name\", symName)\n")
	buf.WriteString("\tif err != nil {\n")
	buf.WriteString("\t\treturn nil, err\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\ttypeAttr, err := mlir.NamedAttributeByName(ctx, \"function_type\", functionTypeAttr)\n")
	buf.WriteString("\tif err != nil {\n")
	buf.WriteString("\t\treturn nil, err\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tstate := mlir.NewOperationState(\"func.func\", loc)\n")
	buf.WriteString("\tstate.AddAttributes(nameAttr, typeAttr)\n")
	buf.WriteString("\tstate.AddOwnedRegions(body)\n")
	buf.WriteString("\treturn mlir.CreateOperation(state)\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func CallOp(ctx *mlir.Context, loc mlir.Location, callee string, results []mlir.Type, operands ...mlir.Value) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\tif callee == \"\" {\n")
	buf.WriteString("\t\treturn nil, fmt.Errorf(\"mlir: callee is required\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tcalleeAttr, err := mlir.FlatSymbolRefAttribute(ctx, callee)\n")
	buf.WriteString("\tif err != nil {\n")
	buf.WriteString("\t\treturn nil, err\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tnamedCallee, err := mlir.NamedAttributeByName(ctx, \"callee\", calleeAttr)\n")
	buf.WriteString("\tif err != nil {\n")
	buf.WriteString("\t\treturn nil, err\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tstate := mlir.NewOperationState(\"func.call\", loc)\n")
	buf.WriteString("\tstate.AddResults(results...)\n")
	buf.WriteString("\tstate.AddOperands(operands...)\n")
	buf.WriteString("\tstate.AddAttributes(namedCallee)\n")
	buf.WriteString("\treturn mlir.CreateOperation(state)\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func ReturnOp(loc mlir.Location, values ...mlir.Value) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\tstate := mlir.NewOperationState(\"func.return\", loc)\n")
	buf.WriteString("\tstate.AddOperands(values...)\n")
	buf.WriteString("\treturn mlir.CreateOperation(state)\n")
	buf.WriteString("}\n")
	return buf.String()
}

func renderFuncGeneratedNoCgo() string {
	var buf bytes.Buffer
	buf.WriteString("// Code generated by cmd/mlir-go-tblgen; DO NOT EDIT.\n")
	buf.WriteString("//go:build !cgo\n\n")
	buf.WriteString("package funcdialect\n\n")
	buf.WriteString("import (\n")
	buf.WriteString("\t\"fmt\"\n\n")
	buf.WriteString("\tmlir \"github.com/timmyyuan/mlir-go\"\n")
	buf.WriteString(")\n\n")

	buf.WriteString("func FuncOp(*mlir.Context, mlir.Location, string, mlir.Type, *mlir.OwnedRegion) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn nil, fmt.Errorf(\"mlir: cgo is required\")\n")
	buf.WriteString("}\n\n")
	buf.WriteString("func CallOp(*mlir.Context, mlir.Location, string, []mlir.Type, ...mlir.Value) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn nil, fmt.Errorf(\"mlir: cgo is required\")\n")
	buf.WriteString("}\n\n")
	buf.WriteString("func ReturnOp(mlir.Location, ...mlir.Value) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn nil, fmt.Errorf(\"mlir: cgo is required\")\n")
	buf.WriteString("}\n")
	return buf.String()
}

func renderCFGenerated() string {
	var buf bytes.Buffer
	buf.WriteString("// Code generated by cmd/mlir-go-tblgen; DO NOT EDIT.\n")
	buf.WriteString("//go:build cgo\n\n")
	buf.WriteString("package cf\n\n")
	buf.WriteString("import (\n")
	buf.WriteString("\t\"fmt\"\n")
	buf.WriteString("\t\"strconv\"\n\n")
	buf.WriteString("\tmlir \"github.com/timmyyuan/mlir-go\"\n")
	buf.WriteString(")\n\n")

	buf.WriteString("func BranchOp(loc mlir.Location, dest mlir.Block, operands ...mlir.Value) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\tif dest.IsNull() {\n")
	buf.WriteString("\t\treturn nil, fmt.Errorf(\"mlir: branch destination is required\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tstate := mlir.NewOperationState(\"cf.br\", loc)\n")
	buf.WriteString("\tstate.AddOperands(operands...)\n")
	buf.WriteString("\tstate.AddSuccessors(dest)\n")
	buf.WriteString("\treturn mlir.CreateOperation(state)\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func CondBranchOp(ctx *mlir.Context, loc mlir.Location, cond mlir.Value, trueDest mlir.Block, trueOperands []mlir.Value, falseDest mlir.Block, falseOperands []mlir.Value) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\tif cond.IsNull() {\n")
	buf.WriteString("\t\treturn nil, fmt.Errorf(\"mlir: branch condition is required\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tif trueDest.IsNull() || falseDest.IsNull() {\n")
	buf.WriteString("\t\treturn nil, fmt.Errorf(\"mlir: both branch destinations are required\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tsegmentAttr, err := mlir.ParseAttribute(ctx, \"array<i32: 1, \"+strconv.Itoa(len(trueOperands))+\", \"+strconv.Itoa(len(falseOperands))+\">\")\n")
	buf.WriteString("\tif err != nil {\n")
	buf.WriteString("\t\treturn nil, err\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tsegmentNamedAttr, err := mlir.NamedAttributeByName(ctx, \"operandSegmentSizes\", segmentAttr)\n")
	buf.WriteString("\tif err != nil {\n")
	buf.WriteString("\t\treturn nil, err\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tstate := mlir.NewOperationState(\"cf.cond_br\", loc)\n")
	buf.WriteString("\tstate.AddOperands(cond)\n")
	buf.WriteString("\tstate.AddOperands(trueOperands...)\n")
	buf.WriteString("\tstate.AddOperands(falseOperands...)\n")
	buf.WriteString("\tstate.AddSuccessors(trueDest, falseDest)\n")
	buf.WriteString("\tstate.AddAttributes(segmentNamedAttr)\n")
	buf.WriteString("\treturn mlir.CreateOperation(state)\n")
	buf.WriteString("}\n")
	return buf.String()
}

func renderCFGeneratedNoCgo() string {
	var buf bytes.Buffer
	buf.WriteString("// Code generated by cmd/mlir-go-tblgen; DO NOT EDIT.\n")
	buf.WriteString("//go:build !cgo\n\n")
	buf.WriteString("package cf\n\n")
	buf.WriteString("import (\n")
	buf.WriteString("\t\"fmt\"\n\n")
	buf.WriteString("\tmlir \"github.com/timmyyuan/mlir-go\"\n")
	buf.WriteString(")\n\n")

	buf.WriteString("func BranchOp(mlir.Location, mlir.Block, ...mlir.Value) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn nil, fmt.Errorf(\"mlir: cgo is required\")\n")
	buf.WriteString("}\n\n")
	buf.WriteString("func CondBranchOp(*mlir.Context, mlir.Location, mlir.Value, mlir.Block, []mlir.Value, mlir.Block, []mlir.Value) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn nil, fmt.Errorf(\"mlir: cgo is required\")\n")
	buf.WriteString("}\n")
	return buf.String()
}

func renderMemRefGenerated() string {
	var buf bytes.Buffer
	buf.WriteString("// Code generated by cmd/mlir-go-tblgen; DO NOT EDIT.\n")
	buf.WriteString("//go:build cgo\n\n")
	buf.WriteString("package memref\n\n")
	buf.WriteString("import (\n")
	buf.WriteString("\t\"fmt\"\n")
	buf.WriteString("\t\"strconv\"\n\n")
	buf.WriteString("\tmlir \"github.com/timmyyuan/mlir-go\"\n")
	buf.WriteString(")\n\n")
	buf.WriteString("func AllocaOp(ctx *mlir.Context, loc mlir.Location, resultType mlir.Type, dynamicSizes ...mlir.Value) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\tif resultType.IsNull() {\n")
	buf.WriteString("\t\treturn nil, fmt.Errorf(\"mlir: result type is required\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tsegmentAttr, err := mlir.ParseAttribute(ctx, \"array<i32: \"+strconv.Itoa(len(dynamicSizes))+\", 0>\")\n")
	buf.WriteString("\tif err != nil {\n")
	buf.WriteString("\t\treturn nil, err\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tsegmentNamedAttr, err := mlir.NamedAttributeByName(ctx, \"operandSegmentSizes\", segmentAttr)\n")
	buf.WriteString("\tif err != nil {\n")
	buf.WriteString("\t\treturn nil, err\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tstate := mlir.NewOperationState(\"memref.alloca\", loc)\n")
	buf.WriteString("\tstate.AddResults(resultType)\n")
	buf.WriteString("\tstate.AddOperands(dynamicSizes...)\n")
	buf.WriteString("\tstate.AddAttributes(segmentNamedAttr)\n")
	buf.WriteString("\treturn mlir.CreateOperation(state)\n")
	buf.WriteString("}\n\n")
	buf.WriteString("func LoadOp(loc mlir.Location, memref mlir.Value, indices ...mlir.Value) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\tif memref.IsNull() {\n")
	buf.WriteString("\t\treturn nil, fmt.Errorf(\"mlir: memref operand is required\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\telementType := memref.Type().ElementType()\n")
	buf.WriteString("\tif elementType.IsNull() {\n")
	buf.WriteString("\t\treturn nil, fmt.Errorf(\"mlir: memref operand must have shaped type\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tstate := mlir.NewOperationState(\"memref.load\", loc)\n")
	buf.WriteString("\tstate.AddResults(elementType)\n")
	buf.WriteString("\tstate.AddOperands(memref)\n")
	buf.WriteString("\tstate.AddOperands(indices...)\n")
	buf.WriteString("\treturn mlir.CreateOperation(state)\n")
	buf.WriteString("}\n\n")
	buf.WriteString("func StoreOp(loc mlir.Location, value mlir.Value, memref mlir.Value, indices ...mlir.Value) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\tif value.IsNull() || memref.IsNull() {\n")
	buf.WriteString("\t\treturn nil, fmt.Errorf(\"mlir: store operands are required\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tstate := mlir.NewOperationState(\"memref.store\", loc)\n")
	buf.WriteString("\tstate.AddOperands(value, memref)\n")
	buf.WriteString("\tstate.AddOperands(indices...)\n")
	buf.WriteString("\treturn mlir.CreateOperation(state)\n")
	buf.WriteString("}\n")
	return buf.String()
}

func renderMemRefGeneratedNoCgo() string {
	var buf bytes.Buffer
	buf.WriteString("// Code generated by cmd/mlir-go-tblgen; DO NOT EDIT.\n")
	buf.WriteString("//go:build !cgo\n\n")
	buf.WriteString("package memref\n\n")
	buf.WriteString("import (\n")
	buf.WriteString("\t\"fmt\"\n\n")
	buf.WriteString("\tmlir \"github.com/timmyyuan/mlir-go\"\n")
	buf.WriteString(")\n\n")
	buf.WriteString("func AllocaOp(*mlir.Context, mlir.Location, mlir.Type, ...mlir.Value) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn nil, fmt.Errorf(\"mlir: cgo is required\")\n")
	buf.WriteString("}\n\n")
	buf.WriteString("func LoadOp(mlir.Location, mlir.Value, ...mlir.Value) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn nil, fmt.Errorf(\"mlir: cgo is required\")\n")
	buf.WriteString("}\n\n")
	buf.WriteString("func StoreOp(mlir.Location, mlir.Value, mlir.Value, ...mlir.Value) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn nil, fmt.Errorf(\"mlir: cgo is required\")\n")
	buf.WriteString("}\n")
	return buf.String()
}

func renderTensorGenerated() string {
	var buf bytes.Buffer
	buf.WriteString("// Code generated by cmd/mlir-go-tblgen; DO NOT EDIT.\n")
	buf.WriteString("//go:build cgo\n\n")
	buf.WriteString("package tensor\n\n")
	buf.WriteString("import (\n")
	buf.WriteString("\t\"fmt\"\n\n")
	buf.WriteString("\tmlir \"github.com/timmyyuan/mlir-go\"\n")
	buf.WriteString(")\n\n")
	buf.WriteString("func EmptyOp(loc mlir.Location, resultType mlir.Type, dynamicSizes ...mlir.Value) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\tif resultType.IsNull() {\n")
	buf.WriteString("\t\treturn nil, fmt.Errorf(\"mlir: result type is required\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tif !resultType.IsRankedTensor() {\n")
	buf.WriteString("\t\treturn nil, fmt.Errorf(\"mlir: tensor.empty expects a ranked tensor result type\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\texpectedDynamic := 0\n")
	buf.WriteString("\tfor _, dim := range resultType.Shape() {\n")
	buf.WriteString("\t\tif mlir.IsDynamicSize(dim) {\n")
	buf.WriteString("\t\t\texpectedDynamic++\n")
	buf.WriteString("\t\t}\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tif len(dynamicSizes) != expectedDynamic {\n")
	buf.WriteString("\t\treturn nil, fmt.Errorf(\"mlir: tensor.empty expects %d dynamic sizes, got %d\", expectedDynamic, len(dynamicSizes))\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tstate := mlir.NewOperationState(\"tensor.empty\", loc)\n")
	buf.WriteString("\tstate.AddResults(resultType)\n")
	buf.WriteString("\tstate.AddOperands(dynamicSizes...)\n")
	buf.WriteString("\treturn mlir.CreateOperation(state)\n")
	buf.WriteString("}\n\n")
	buf.WriteString("func DimOp(loc mlir.Location, source mlir.Value, index mlir.Value) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\tif source.IsNull() || index.IsNull() {\n")
	buf.WriteString("\t\treturn nil, fmt.Errorf(\"mlir: source and index operands are required\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tif !index.Type().IsIndex() {\n")
	buf.WriteString("\t\treturn nil, fmt.Errorf(\"mlir: tensor.dim index operand must have index type\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tstate := mlir.NewOperationState(\"tensor.dim\", loc)\n")
	buf.WriteString("\tstate.AddResults(index.Type())\n")
	buf.WriteString("\tstate.AddOperands(source, index)\n")
	buf.WriteString("\treturn mlir.CreateOperation(state)\n")
	buf.WriteString("}\n\n")
	buf.WriteString("func CastOp(loc mlir.Location, source mlir.Value, resultType mlir.Type) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\tif source.IsNull() {\n")
	buf.WriteString("\t\treturn nil, fmt.Errorf(\"mlir: source operand is required\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tif resultType.IsNull() {\n")
	buf.WriteString("\t\treturn nil, fmt.Errorf(\"mlir: result type is required\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tstate := mlir.NewOperationState(\"tensor.cast\", loc)\n")
	buf.WriteString("\tstate.AddResults(resultType)\n")
	buf.WriteString("\tstate.AddOperands(source)\n")
	buf.WriteString("\treturn mlir.CreateOperation(state)\n")
	buf.WriteString("}\n")
	return buf.String()
}

func renderTensorGeneratedNoCgo() string {
	var buf bytes.Buffer
	buf.WriteString("// Code generated by cmd/mlir-go-tblgen; DO NOT EDIT.\n")
	buf.WriteString("//go:build !cgo\n\n")
	buf.WriteString("package tensor\n\n")
	buf.WriteString("import (\n")
	buf.WriteString("\t\"fmt\"\n\n")
	buf.WriteString("\tmlir \"github.com/timmyyuan/mlir-go\"\n")
	buf.WriteString(")\n\n")
	buf.WriteString("func EmptyOp(mlir.Location, mlir.Type, ...mlir.Value) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn nil, fmt.Errorf(\"mlir: cgo is required\")\n")
	buf.WriteString("}\n\n")
	buf.WriteString("func DimOp(mlir.Location, mlir.Value, mlir.Value) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn nil, fmt.Errorf(\"mlir: cgo is required\")\n")
	buf.WriteString("}\n\n")
	buf.WriteString("func CastOp(mlir.Location, mlir.Value, mlir.Type) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn nil, fmt.Errorf(\"mlir: cgo is required\")\n")
	buf.WriteString("}\n")
	return buf.String()
}

func renderVectorGenerated() string {
	var buf bytes.Buffer
	buf.WriteString("// Code generated by cmd/mlir-go-tblgen; DO NOT EDIT.\n")
	buf.WriteString("//go:build cgo\n\n")
	buf.WriteString("package vector\n\n")
	buf.WriteString("import (\n")
	buf.WriteString("\t\"fmt\"\n\n")
	buf.WriteString("\tmlir \"github.com/timmyyuan/mlir-go\"\n")
	buf.WriteString(")\n\n")
	buf.WriteString("func BroadcastOp(loc mlir.Location, source mlir.Value, resultType mlir.Type) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\tif source.IsNull() {\n")
	buf.WriteString("\t\treturn nil, fmt.Errorf(\"mlir: source operand is required\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tif resultType.IsNull() {\n")
	buf.WriteString("\t\treturn nil, fmt.Errorf(\"mlir: result type is required\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tstate := mlir.NewOperationState(\"vector.broadcast\", loc)\n")
	buf.WriteString("\tstate.AddResults(resultType)\n")
	buf.WriteString("\tstate.AddOperands(source)\n")
	buf.WriteString("\treturn mlir.CreateOperation(state)\n")
	buf.WriteString("}\n\n")
	buf.WriteString("func ExtractElementOp(loc mlir.Location, vector mlir.Value, resultType mlir.Type, position ...mlir.Value) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\tif vector.IsNull() {\n")
	buf.WriteString("\t\treturn nil, fmt.Errorf(\"mlir: vector operand is required\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tif resultType.IsNull() {\n")
	buf.WriteString("\t\treturn nil, fmt.Errorf(\"mlir: result type is required\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tif len(position) > 1 {\n")
	buf.WriteString("\t\treturn nil, fmt.Errorf(\"mlir: vector.extractelement accepts at most one dynamic position\")\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tstate := mlir.NewOperationState(\"vector.extractelement\", loc)\n")
	buf.WriteString("\tstate.AddResults(resultType)\n")
	buf.WriteString("\tstate.AddOperands(vector)\n")
	buf.WriteString("\tstate.AddOperands(position...)\n")
	buf.WriteString("\treturn mlir.CreateOperation(state)\n")
	buf.WriteString("}\n")
	return buf.String()
}

func renderVectorGeneratedNoCgo() string {
	var buf bytes.Buffer
	buf.WriteString("// Code generated by cmd/mlir-go-tblgen; DO NOT EDIT.\n")
	buf.WriteString("//go:build !cgo\n\n")
	buf.WriteString("package vector\n\n")
	buf.WriteString("import (\n")
	buf.WriteString("\t\"fmt\"\n\n")
	buf.WriteString("\tmlir \"github.com/timmyyuan/mlir-go\"\n")
	buf.WriteString(")\n\n")
	buf.WriteString("func BroadcastOp(mlir.Location, mlir.Value, mlir.Type) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn nil, fmt.Errorf(\"mlir: cgo is required\")\n")
	buf.WriteString("}\n\n")
	buf.WriteString("func ExtractElementOp(mlir.Location, mlir.Value, mlir.Type, ...mlir.Value) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn nil, fmt.Errorf(\"mlir: cgo is required\")\n")
	buf.WriteString("}\n")
	return buf.String()
}

func renderLinalgGenerated() string {
	var buf bytes.Buffer
	buf.WriteString("// Code generated by cmd/mlir-go-tblgen; DO NOT EDIT.\n")
	buf.WriteString("//go:build cgo\n\n")
	buf.WriteString("package linalg\n\n")
	buf.WriteString("import mlir \"github.com/timmyyuan/mlir-go\"\n\n")
	buf.WriteString("func IndexOp(ctx *mlir.Context, loc mlir.Location, dim int64) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\tindexType, err := mlir.IndexType(ctx)\n")
	buf.WriteString("\tif err != nil {\n")
	buf.WriteString("\t\treturn nil, err\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\ti64, err := mlir.SignlessIntegerType(ctx, 64)\n")
	buf.WriteString("\tif err != nil {\n")
	buf.WriteString("\t\treturn nil, err\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tdimAttr, err := mlir.IntegerAttribute(i64, dim)\n")
	buf.WriteString("\tif err != nil {\n")
	buf.WriteString("\t\treturn nil, err\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tdimNamedAttr, err := mlir.NamedAttributeByName(ctx, \"dim\", dimAttr)\n")
	buf.WriteString("\tif err != nil {\n")
	buf.WriteString("\t\treturn nil, err\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tstate := mlir.NewOperationState(\"linalg.index\", loc)\n")
	buf.WriteString("\tstate.AddResults(indexType)\n")
	buf.WriteString("\tstate.AddAttributes(dimNamedAttr)\n")
	buf.WriteString("\treturn mlir.CreateOperation(state)\n")
	buf.WriteString("}\n\n")
	buf.WriteString("func YieldOp(loc mlir.Location, values ...mlir.Value) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\tstate := mlir.NewOperationState(\"linalg.yield\", loc)\n")
	buf.WriteString("\tstate.AddOperands(values...)\n")
	buf.WriteString("\treturn mlir.CreateOperation(state)\n")
	buf.WriteString("}\n")
	return buf.String()
}

func renderLinalgGeneratedNoCgo() string {
	var buf bytes.Buffer
	buf.WriteString("// Code generated by cmd/mlir-go-tblgen; DO NOT EDIT.\n")
	buf.WriteString("//go:build !cgo\n\n")
	buf.WriteString("package linalg\n\n")
	buf.WriteString("import (\n")
	buf.WriteString("\t\"fmt\"\n\n")
	buf.WriteString("\tmlir \"github.com/timmyyuan/mlir-go\"\n")
	buf.WriteString(")\n\n")
	buf.WriteString("func IndexOp(*mlir.Context, mlir.Location, int64) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn nil, fmt.Errorf(\"mlir: cgo is required\")\n")
	buf.WriteString("}\n\n")
	buf.WriteString("func YieldOp(mlir.Location, ...mlir.Value) (*mlir.OwnedOperation, error) {\n")
	buf.WriteString("\treturn nil, fmt.Errorf(\"mlir: cgo is required\")\n")
	buf.WriteString("}\n")
	return buf.String()
}
