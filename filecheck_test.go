//go:build cgo

package mlir

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func fileCheckPath(t *testing.T) string {
	t.Helper()

	if p := os.Getenv("FILECHECK"); p != "" {
		return p
	}
	if p, err := exec.LookPath("FileCheck"); err == nil {
		return p
	}

	for _, envName := range []string{"MLIRGO_LLVM_CONFIG", "LLVM_CONFIG"} {
		llvmConfig := os.Getenv(envName)
		if llvmConfig == "" {
			continue
		}
		out, err := exec.Command(llvmConfig, "--bindir").Output()
		if err != nil {
			continue
		}
		p := filepath.Join(strings.TrimSpace(string(out)), "FileCheck")
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	t.Skip("FileCheck not available")
	return ""
}

func runFileCheck(t *testing.T, input string, fixture string) {
	t.Helper()

	checkFile := filepath.Join("testdata", "filecheck", fixture)
	cmd := exec.Command(fileCheckPath(t), checkFile)
	cmd.Stdin = strings.NewReader(input)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("FileCheck failed for %s: %v\n%s\ninput:\n%s", fixture, err, out.String(), input)
	}
}

func TestFileCheckEmptyModule(t *testing.T) {
	ctx, err := NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer func() {
		if err := ctx.Close(); err != nil {
			t.Fatalf("Close() error = %v", err)
		}
	}()

	mod, err := ParseModule(ctx, "module {\n}\n")
	if err != nil {
		t.Fatalf("ParseModule() error = %v", err)
	}
	defer func() {
		if err := mod.Close(); err != nil {
			t.Fatalf("mod.Close() error = %v", err)
		}
	}()

	runFileCheck(t, mod.String(), "empty_module.mlir")
}

func TestFileCheckAddModule(t *testing.T) {
	ctx, err := NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer func() {
		if err := ctx.Close(); err != nil {
			t.Fatalf("Close() error = %v", err)
		}
	}()

	if err := ctx.RegisterAllDialects(); err != nil {
		t.Fatalf("RegisterAllDialects() error = %v", err)
	}

	const input = `module {
  func.func @add(%arg0: i32) -> i32 {
    %0 = arith.addi %arg0, %arg0 : i32
    return %0 : i32
  }
}
`

	mod, err := ParseModule(ctx, input)
	if err != nil {
		t.Fatalf("ParseModule() error = %v", err)
	}
	defer func() {
		if err := mod.Close(); err != nil {
			t.Fatalf("mod.Close() error = %v", err)
		}
	}()

	runFileCheck(t, mod.String(), "add_module.mlir")
	runFileCheck(t, mod.Body().Operations()[0].Regions()[0].Blocks()[0].String(), "func_block.mlir")
}

func TestFileCheckModuleLocation(t *testing.T) {
	ctx, err := NewContext()
	if err != nil {
		t.Fatalf("NewContext() error = %v", err)
	}
	defer func() {
		if err := ctx.Close(); err != nil {
			t.Fatalf("Close() error = %v", err)
		}
	}()

	loc, err := FileLineColLocation(ctx, "test.mlir", 3, 9)
	if err != nil {
		t.Fatalf("FileLineColLocation() error = %v", err)
	}

	mod, err := CreateEmptyModule(loc)
	if err != nil {
		t.Fatalf("CreateEmptyModule() error = %v", err)
	}
	defer func() {
		if err := mod.Close(); err != nil {
			t.Fatalf("mod.Close() error = %v", err)
		}
	}()

	runFileCheck(t, mod.Operation().Location().String(), "location_string.txt")
}
