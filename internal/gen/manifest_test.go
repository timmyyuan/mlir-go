package gen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadManifestAndValidate(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "mlir", "Dialect", "Sample", "IR", "SampleOps.td")
	if err := os.MkdirAll(filepath.Dir(source), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(source, []byte("def Sample_AddOp : Sample_Op<\"add\">;\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	manifestPath := filepath.Join(root, "dialect_manifest.json")
	manifestJSON := `{
  "dialects": [
    {
      "name": "sample",
      "tablegen": ["mlir/Dialect/Sample/IR/SampleOps.td"],
      "import_path": "github.com/timmyyuan/mlir-go/dialect/sample",
      "package": "sample",
      "strip_prefixes": ["Sample_"],
      "generate_tests": true
    }
  ]
}`
	if err := os.WriteFile(manifestPath, []byte(manifestJSON), 0o644); err != nil {
		t.Fatalf("WriteFile(manifest) error = %v", err)
	}

	manifest, err := LoadManifest(manifestPath)
	if err != nil {
		t.Fatalf("LoadManifest() error = %v", err)
	}
	if err := manifest.Validate(root); err != nil {
		t.Fatalf("manifest.Validate() error = %v", err)
	}
	if got := manifest.Dialects[0].Name; got != "sample" {
		t.Fatalf("dialect name = %q, want %q", got, "sample")
	}
}

func TestManifestValidateRejectsMissingSource(t *testing.T) {
	manifest := &Manifest{
		Dialects: []DialectManifest{{
			Name:       "sample",
			TableGen:   []string{"mlir/Dialect/Sample/IR/Missing.td"},
			ImportPath: "github.com/timmyyuan/mlir-go/dialect/sample",
			Package:    "sample",
		}},
	}
	err := manifest.Validate(t.TempDir())
	if err == nil {
		t.Fatalf("manifest.Validate() error = nil, want missing source error")
	}
	if !strings.Contains(err.Error(), "does not exist") {
		t.Fatalf("error = %q, want missing source detail", err)
	}
}
