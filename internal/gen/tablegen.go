package gen

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// OperationDef is a lightweight view of an op definition in TableGen.
type OperationDef struct {
	DefName  string `json:"def_name"`
	Mnemonic string `json:"mnemonic"`
	Source   string `json:"source"`
}

var opHeaderPattern = regexp.MustCompile(`^def\s+([A-Za-z0-9_]+)\s*:\s*.*<[^>]*"([^"]+)"`)

func ScanDialect(includeRoot string, dialect DialectManifest) ([]OperationDef, error) {
	if includeRoot == "" {
		return nil, fmt.Errorf("include root is required")
	}
	if err := dialect.Validate("dialect", includeRoot); err != nil {
		return nil, err
	}

	var ops []OperationDef
	for _, source := range dialect.TableGen {
		found, err := ScanTableGenFile(filepath.Join(includeRoot, source), source)
		if err != nil {
			return nil, err
		}
		ops = append(ops, found...)
	}
	return ops, nil
}

func ScanTableGenFile(path string, relativeSource string) ([]OperationDef, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var ops []OperationDef
	scanner := bufio.NewScanner(file)
	var (
		inDef  bool
		header strings.Builder
	)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !inDef {
			if !strings.HasPrefix(line, "def ") {
				continue
			}
			inDef = true
			header.Reset()
		} else if line == "" || strings.HasPrefix(line, "//") {
			continue
		}
		if header.Len() > 0 {
			header.WriteByte(' ')
		}
		header.WriteString(line)

		combined := header.String()
		if matches := opHeaderPattern.FindStringSubmatch(combined); len(matches) == 3 {
			if !strings.HasSuffix(matches[1], "Op") {
				inDef = false
				header.Reset()
				continue
			}
			ops = append(ops, OperationDef{
				DefName:  matches[1],
				Mnemonic: matches[2],
				Source:   relativeSource,
			})
			inDef = false
			header.Reset()
			continue
		}
		if strings.Contains(line, "{") || strings.Contains(line, ";") {
			inDef = false
			header.Reset()
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan %s: %w", path, err)
	}
	return ops, nil
}
