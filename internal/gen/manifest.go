package gen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Manifest declares the dialects supported by the generator.
type Manifest struct {
	Dialects []DialectManifest `json:"dialects"`
}

// DialectManifest describes one dialect generation target.
type DialectManifest struct {
	Name          string   `json:"name"`
	TableGen      []string `json:"tablegen"`
	ImportPath    string   `json:"import_path"`
	Package       string   `json:"package"`
	StripPrefixes []string `json:"strip_prefixes"`
	GenerateTests bool     `json:"generate_tests"`
}

func LoadManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("decode manifest %s: %w", path, err)
	}
	if err := manifest.Validate(""); err != nil {
		return nil, err
	}
	return &manifest, nil
}

func (m *Manifest) Validate(includeRoot string) error {
	if m == nil {
		return fmt.Errorf("nil manifest")
	}
	if len(m.Dialects) == 0 {
		return fmt.Errorf("manifest must declare at least one dialect")
	}

	seenNames := make(map[string]struct{}, len(m.Dialects))
	seenImports := make(map[string]struct{}, len(m.Dialects))
	for i, dialect := range m.Dialects {
		prefix := fmt.Sprintf("dialects[%d]", i)
		if err := dialect.Validate(prefix, includeRoot); err != nil {
			return err
		}
		if _, ok := seenNames[dialect.Name]; ok {
			return fmt.Errorf("%s.name %q is duplicated", prefix, dialect.Name)
		}
		if _, ok := seenImports[dialect.ImportPath]; ok {
			return fmt.Errorf("%s.import_path %q is duplicated", prefix, dialect.ImportPath)
		}
		seenNames[dialect.Name] = struct{}{}
		seenImports[dialect.ImportPath] = struct{}{}
	}
	return nil
}

func (d DialectManifest) Validate(prefix string, includeRoot string) error {
	if prefix == "" {
		prefix = "dialect"
	}
	if d.Name == "" {
		return fmt.Errorf("%s.name is required", prefix)
	}
	if !isLowerIdentifier(d.Name) {
		return fmt.Errorf("%s.name %q must be lowercase and use [a-z0-9_]", prefix, d.Name)
	}
	if len(d.TableGen) == 0 {
		return fmt.Errorf("%s.tablegen must contain at least one entry", prefix)
	}
	for i, source := range d.TableGen {
		if source == "" {
			return fmt.Errorf("%s.tablegen[%d] is empty", prefix, i)
		}
		if filepath.IsAbs(source) {
			return fmt.Errorf("%s.tablegen[%d] must be relative, got %q", prefix, i, source)
		}
		clean := filepath.Clean(source)
		if clean == "." || strings.HasPrefix(clean, "..") {
			return fmt.Errorf("%s.tablegen[%d] %q escapes the include root", prefix, i, source)
		}
		if includeRoot != "" {
			if _, err := os.Stat(filepath.Join(includeRoot, clean)); err != nil {
				return fmt.Errorf("%s.tablegen[%d] %q does not exist under %q", prefix, i, source, includeRoot)
			}
		}
	}
	if d.ImportPath == "" {
		return fmt.Errorf("%s.import_path is required", prefix)
	}
	if d.Package == "" {
		return fmt.Errorf("%s.package is required", prefix)
	}
	return nil
}

func isLowerIdentifier(v string) bool {
	for _, r := range v {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= '0' && r <= '9':
		case r == '_':
		default:
			return false
		}
	}
	return true
}
