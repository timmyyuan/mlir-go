package gen

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestEmitDialectMatchesCheckedInFiles(t *testing.T) {
	repoRoot := filepath.Clean(filepath.Join("..", ".."))

	manifest, err := LoadManifest(filepath.Join(repoRoot, "dialect_manifest.json"))
	if err != nil {
		t.Fatalf("LoadManifest() error = %v", err)
	}
	includeRoot, err := ResolveIncludeRoot("")
	if err != nil {
		t.Fatalf("ResolveIncludeRoot() error = %v", err)
	}

	for _, dialectName := range []string{"arith", "func", "cf", "memref", "tensor", "vector", "linalg"} {
		dialect, ok := manifest.FindDialect(dialectName)
		if !ok {
			t.Fatalf("manifest.FindDialect(%s) = false, want true", dialectName)
		}
		files, err := EmitDialect(includeRoot, dialect)
		if err != nil {
			t.Fatalf("EmitDialect(%s) error = %v", dialectName, err)
		}
		if len(files) != 2 {
			t.Fatalf("%s: len(files) = %d, want 2", dialectName, len(files))
		}
		for _, file := range files {
			checkedIn, err := os.ReadFile(filepath.Join(repoRoot, filepath.FromSlash(file.Path)))
			if err != nil {
				t.Fatalf("ReadFile(%s) error = %v", file.Path, err)
			}
			if !bytes.Equal(checkedIn, file.Content) {
				t.Fatalf("generated file %s is stale; rerun cmd/mlir-go-tblgen", file.Path)
			}
		}
	}
}
