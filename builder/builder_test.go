//go:build cgo

package builder_test

import (
	"strings"
	"testing"

	mlir "github.com/timmyyuan/mlir-go"
	"github.com/timmyyuan/mlir-go/builder"
	builderarith "github.com/timmyyuan/mlir-go/builder/arith"
	buildercf "github.com/timmyyuan/mlir-go/builder/cf"
	builderfunc "github.com/timmyyuan/mlir-go/builder/func"
)

func TestBuilderBuildsSimpleFunction(t *testing.T) {
	ctx, err := mlir.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer func() {
		if err := ctx.Close(); err != nil {
			t.Fatalf("ctx.Close() error = %v", err)
		}
	}()

	if err := ctx.RegisterAllDialects(); err != nil {
		t.Fatalf("RegisterAllDialects() error = %v", err)
	}

	b, err := builder.New(ctx)
	if err != nil {
		t.Fatalf("builder.New() error = %v", err)
	}
	loc, err := mlir.FileLineColLocation(ctx, "builder_test.mlir", 12, 3)
	if err != nil {
		t.Fatalf("FileLineColLocation() error = %v", err)
	}
	if err := b.SetLocation(loc); err != nil {
		t.Fatalf("b.SetLocation() error = %v", err)
	}

	i32, err := mlir.SignlessIntegerType(ctx, 32)
	if err != nil {
		t.Fatalf("SignlessIntegerType() error = %v", err)
	}
	fnType, err := mlir.FunctionType(ctx, []mlir.Type{i32}, []mlir.Type{i32})
	if err != nil {
		t.Fatalf("FunctionType() error = %v", err)
	}

	mod, err := b.BuildModule(func(b *builder.Builder, mod *mlir.Module) error {
		_, err := b.BuildFunction("increment", fnType, func(b *builder.Builder, entry mlir.Block) error {
			arg0 := entry.Arguments()[0]

			fiveAttr, err := mlir.IntegerAttribute(i32, 5)
			if err != nil {
				return err
			}
			c5, err := builderarith.Constant(b, i32, fiveAttr)
			if err != nil {
				return err
			}
			sum, err := builderarith.AddI(b, arg0, c5.Results()[0])
			if err != nil {
				return err
			}
			_, err = builderfunc.Return(b, sum.Results()[0])
			return err
		})
		return err
	})
	if err != nil {
		t.Fatalf("b.BuildModule() error = %v", err)
	}
	defer func() {
		if err := mod.Close(); err != nil {
			t.Fatalf("mod.Close() error = %v", err)
		}
	}()

	if err := mod.Verify(); err != nil {
		t.Fatalf("mod.Verify() error = %v", err)
	}
	ops := mod.Body().Operations()[0].Regions()[0].Blocks()[0].Operations()
	for i, op := range ops {
		if !strings.Contains(op.Location().String(), "builder_test.mlir") {
			t.Fatalf("op %d location = %q, want propagated builder filename", i, op.Location().String())
		}
	}
	runFileCheck(t, mod.String(), "builder_simple_module.mlir")
}

func TestBuilderBuildsMultiBlockControlFlow(t *testing.T) {
	ctx, err := mlir.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer func() {
		if err := ctx.Close(); err != nil {
			t.Fatalf("ctx.Close() error = %v", err)
		}
	}()

	if err := ctx.RegisterAllDialects(); err != nil {
		t.Fatalf("RegisterAllDialects() error = %v", err)
	}

	b, err := builder.New(ctx)
	if err != nil {
		t.Fatalf("builder.New() error = %v", err)
	}
	loc, err := mlir.FileLineColLocation(ctx, "builder_cf.mlir", 4, 7)
	if err != nil {
		t.Fatalf("FileLineColLocation() error = %v", err)
	}
	if err := b.SetLocation(loc); err != nil {
		t.Fatalf("b.SetLocation() error = %v", err)
	}

	i1, err := mlir.SignlessIntegerType(ctx, 1)
	if err != nil {
		t.Fatalf("SignlessIntegerType(i1) error = %v", err)
	}
	i32, err := mlir.SignlessIntegerType(ctx, 32)
	if err != nil {
		t.Fatalf("SignlessIntegerType(i32) error = %v", err)
	}
	fnType, err := mlir.FunctionType(ctx, []mlir.Type{i1}, []mlir.Type{i32})
	if err != nil {
		t.Fatalf("FunctionType() error = %v", err)
	}

	mod, err := b.BuildModule(func(b *builder.Builder, mod *mlir.Module) error {
		_, err := b.BuildFunction("select_const", fnType, func(b *builder.Builder, entry mlir.Block) error {
			thenBlock, err := b.AppendBlock()
			if err != nil {
				return err
			}
			elseBlock, err := b.AppendBlock()
			if err != nil {
				return err
			}

			cond := entry.Arguments()[0]
			if _, err := buildercf.CondBranch(b, cond, thenBlock, nil, elseBlock, nil); err != nil {
				return err
			}

			if err := b.PositionAtEnd(thenBlock); err != nil {
				return err
			}
			oneAttr, err := mlir.IntegerAttribute(i32, 1)
			if err != nil {
				return err
			}
			one, err := builderarith.Constant(b, i32, oneAttr)
			if err != nil {
				return err
			}
			if _, err := builderfunc.Return(b, one.Results()[0]); err != nil {
				return err
			}

			if err := b.PositionAtEnd(elseBlock); err != nil {
				return err
			}
			zeroAttr, err := mlir.IntegerAttribute(i32, 0)
			if err != nil {
				return err
			}
			zero, err := builderarith.Constant(b, i32, zeroAttr)
			if err != nil {
				return err
			}
			if _, err := builderfunc.Return(b, zero.Results()[0]); err != nil {
				return err
			}
			return nil
		})
		return err
	})
	if err != nil {
		t.Fatalf("b.BuildModule() error = %v", err)
	}
	defer func() {
		if err := mod.Close(); err != nil {
			t.Fatalf("mod.Close() error = %v", err)
		}
	}()

	if err := mod.Verify(); err != nil {
		t.Fatalf("mod.Verify() error = %v", err)
	}
	funcOp := mod.Body().Operations()[0]
	if got := len(funcOp.Regions()[0].Blocks()); got != 3 {
		t.Fatalf("len(function blocks) = %d, want 3", got)
	}
	runFileCheck(t, mod.String(), "builder_cf_module.mlir")
}

func TestBuilderRequiresTerminators(t *testing.T) {
	ctx, err := mlir.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer func() {
		if err := ctx.Close(); err != nil {
			t.Fatalf("ctx.Close() error = %v", err)
		}
	}()

	if err := ctx.RegisterAllDialects(); err != nil {
		t.Fatalf("RegisterAllDialects() error = %v", err)
	}

	b, err := builder.New(ctx)
	if err != nil {
		t.Fatalf("builder.New() error = %v", err)
	}
	i32, err := mlir.SignlessIntegerType(ctx, 32)
	if err != nil {
		t.Fatalf("SignlessIntegerType() error = %v", err)
	}
	fnType, err := mlir.FunctionType(ctx, nil, []mlir.Type{i32})
	if err != nil {
		t.Fatalf("FunctionType() error = %v", err)
	}

	_, err = b.BuildModule(func(b *builder.Builder, mod *mlir.Module) error {
		_, err := b.BuildFunction("broken", fnType, func(b *builder.Builder, entry mlir.Block) error {
			_, err := b.AppendBlock()
			return err
		})
		return err
	})
	if err == nil {
		t.Fatalf("b.BuildModule() error = nil, want missing terminator error")
	}
	if !strings.Contains(err.Error(), "terminator") {
		t.Fatalf("error = %q, want substring %q", err, "terminator")
	}
}

func TestBuilderDialectCallWrapper(t *testing.T) {
	ctx, err := mlir.NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer func() {
		if err := ctx.Close(); err != nil {
			t.Fatalf("ctx.Close() error = %v", err)
		}
	}()

	if err := ctx.RegisterAllDialects(); err != nil {
		t.Fatalf("RegisterAllDialects() error = %v", err)
	}

	b, err := builder.New(ctx)
	if err != nil {
		t.Fatalf("builder.New() error = %v", err)
	}
	i32, err := mlir.SignlessIntegerType(ctx, 32)
	if err != nil {
		t.Fatalf("SignlessIntegerType() error = %v", err)
	}
	fnType, err := mlir.FunctionType(ctx, nil, []mlir.Type{i32})
	if err != nil {
		t.Fatalf("FunctionType() error = %v", err)
	}

	mod, err := b.BuildModule(func(b *builder.Builder, mod *mlir.Module) error {
		if _, err := b.BuildFunction("callee", fnType, func(b *builder.Builder, entry mlir.Block) error {
			c7Attr, err := mlir.IntegerAttribute(i32, 7)
			if err != nil {
				return err
			}
			c7, err := builderarith.Constant(b, i32, c7Attr)
			if err != nil {
				return err
			}
			_, err = builderfunc.Return(b, c7.Results()[0])
			return err
		}); err != nil {
			return err
		}

		_, err := b.BuildFunction("caller", fnType, func(b *builder.Builder, entry mlir.Block) error {
			call, err := builderfunc.Call(b, "callee", []mlir.Type{i32})
			if err != nil {
				return err
			}
			_, err = builderfunc.Return(b, call.Results()[0])
			return err
		})
		return err
	})
	if err != nil {
		t.Fatalf("b.BuildModule() error = %v", err)
	}
	defer func() {
		if err := mod.Close(); err != nil {
			t.Fatalf("mod.Close() error = %v", err)
		}
	}()

	if err := mod.Verify(); err != nil {
		t.Fatalf("mod.Verify() error = %v", err)
	}
	runFileCheck(t, mod.String(), "builder_call_module.mlir")
}
