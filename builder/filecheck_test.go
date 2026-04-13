//go:build cgo

package builder_test

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
