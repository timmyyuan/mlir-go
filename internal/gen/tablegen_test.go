package gen

import (
	"path/filepath"
	"testing"
)

func TestScanTableGenFile(t *testing.T) {
	ops, err := ScanTableGenFile(filepath.Join("testdata", "sample_ops.td"), "testdata/sample_ops.td")
	if err != nil {
		t.Fatalf("ScanTableGenFile() error = %v", err)
	}
	if len(ops) != 2 {
		t.Fatalf("len(ops) = %d, want 2 parsed ops", len(ops))
	}
	if ops[0].DefName != "Sample_AddOp" || ops[0].Mnemonic != "add" {
		t.Fatalf("first op = %+v, want Sample_AddOp/add", ops[0])
	}
	if ops[1].DefName != "Sample_SubOp" || ops[1].Mnemonic != "sub" {
		t.Fatalf("second op = %+v, want Sample_SubOp/sub", ops[1])
	}
}

func TestBuildSkipReport(t *testing.T) {
	manifest := &Manifest{
		Dialects: []DialectManifest{{
			Name:       "sample",
			TableGen:   []string{"testdata/sample_ops.td"},
			ImportPath: "github.com/timmyyuan/mlir-go/dialect/sample",
			Package:    "sample",
		}},
	}
	report, err := BuildSkipReport(manifest, ".")
	if err != nil {
		t.Fatalf("BuildSkipReport() error = %v", err)
	}
	if len(report.Dialects) != 1 {
		t.Fatalf("len(report.Dialects) = %d, want 1", len(report.Dialects))
	}
	if len(report.Dialects[0].SkippedOps) != 2 {
		t.Fatalf("len(report.Dialects[0].SkippedOps) = %d, want 2", len(report.Dialects[0].SkippedOps))
	}
	if report.Dialects[0].SkippedOps[0].Reason == "" {
		t.Fatalf("skip reason is empty")
	}
}
