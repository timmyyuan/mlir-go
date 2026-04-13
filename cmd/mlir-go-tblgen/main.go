package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/timmyyuan/mlir-go/internal/gen"
)

func main() {
	mode := flag.String("mode", "validate-manifest", "one of: validate-manifest, emit-skip-report, emit-dialect")
	manifestPath := flag.String("manifest", "dialect_manifest.json", "path to the dialect manifest JSON file")
	includeRoot := flag.String("mlir-include-root", "", "MLIR include root; defaults to `llvm-config --includedir`")
	outputPath := flag.String("o", "", "output path for emit-skip-report; defaults to stdout")
	dialectName := flag.String("dialect", "", "dialect name for emit-dialect")
	outputRoot := flag.String("output-root", ".", "output root for emit-dialect")
	flag.Parse()

	root, err := gen.ResolveIncludeRoot(*includeRoot)
	if err != nil {
		fatal(err)
	}

	manifest, err := gen.LoadManifest(*manifestPath)
	if err != nil {
		fatal(err)
	}
	if err := manifest.Validate(root); err != nil {
		fatal(err)
	}

	switch *mode {
	case "validate-manifest":
		fmt.Fprintf(os.Stdout, "manifest OK: %d dialects\n", len(manifest.Dialects))
	case "emit-skip-report":
		report, err := gen.BuildSkipReport(manifest, root)
		if err != nil {
			fatal(err)
		}
		data, err := report.JSON()
		if err != nil {
			fatal(err)
		}
		if *outputPath == "" {
			if _, err := os.Stdout.Write(data); err != nil {
				fatal(err)
			}
			if _, err := os.Stdout.Write([]byte("\n")); err != nil {
				fatal(err)
			}
			return
		}
		if err := os.MkdirAll(filepath.Dir(*outputPath), 0o755); err != nil {
			fatal(err)
		}
		if err := os.WriteFile(*outputPath, data, 0o644); err != nil {
			fatal(err)
		}
	case "emit-dialect":
		if *dialectName == "" {
			fatal(fmt.Errorf("-dialect is required for emit-dialect"))
		}
		dialect, ok := manifest.FindDialect(*dialectName)
		if !ok {
			fatal(fmt.Errorf("dialect %q is not declared in manifest", *dialectName))
		}
		files, err := gen.EmitDialect(root, dialect)
		if err != nil {
			fatal(err)
		}
		for _, file := range files {
			path := filepath.Join(*outputRoot, filepath.FromSlash(file.Path))
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				fatal(err)
			}
			if err := os.WriteFile(path, file.Content, 0o644); err != nil {
				fatal(err)
			}
			fmt.Fprintln(os.Stdout, path)
		}
	default:
		fatal(fmt.Errorf("unsupported mode %q", *mode))
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
