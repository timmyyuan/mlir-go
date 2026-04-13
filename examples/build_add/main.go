package main

import (
	"fmt"
	"log"

	mlir "github.com/timmyyuan/mlir-go"
	"github.com/timmyyuan/mlir-go/builder"
	builderarith "github.com/timmyyuan/mlir-go/builder/arith"
	builderfunc "github.com/timmyyuan/mlir-go/builder/func"
)

func main() {
	ctx, err := mlir.NewContext()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := ctx.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := ctx.RegisterAllDialects(); err != nil {
		log.Fatal(err)
	}

	b, err := builder.New(ctx)
	if err != nil {
		log.Fatal(err)
	}

	i32, err := mlir.SignlessIntegerType(ctx, 32)
	if err != nil {
		log.Fatal(err)
	}
	fnType, err := mlir.FunctionType(ctx, []mlir.Type{i32}, []mlir.Type{i32})
	if err != nil {
		log.Fatal(err)
	}

	mod, err := b.BuildModule(func(b *builder.Builder, mod *mlir.Module) error {
		_, err := b.BuildFunction("increment", fnType, func(b *builder.Builder, entry mlir.Block) error {
			arg0 := entry.Arguments()[0]

			five, err := mlir.IntegerAttribute(i32, 5)
			if err != nil {
				return err
			}
			c5, err := builderarith.Constant(b, i32, five)
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
		log.Fatal(err)
	}
	defer func() {
		if err := mod.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := mod.Verify(); err != nil {
		log.Fatal(err)
	}
	fmt.Print(mod.String())
}
