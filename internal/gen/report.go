package gen

import (
	"encoding/json"
)

// SkipReport makes unsupported generated surface visible and reviewable.
type SkipReport struct {
	Dialects []DialectSkipReport `json:"dialects"`
}

type DialectSkipReport struct {
	Name       string             `json:"name"`
	ImportPath string             `json:"import_path"`
	Package    string             `json:"package"`
	Sources    []string           `json:"sources"`
	SkippedOps []SkippedOperation `json:"skipped_ops"`
}

type SkippedOperation struct {
	DefName    string `json:"def_name"`
	Mnemonic   string `json:"mnemonic"`
	Reason     string `json:"reason"`
	Limitation string `json:"limitation"`
	Source     string `json:"source"`
}

func BuildSkipReport(manifest *Manifest, includeRoot string) (*SkipReport, error) {
	report := &SkipReport{Dialects: make([]DialectSkipReport, 0, len(manifest.Dialects))}
	for _, dialect := range manifest.Dialects {
		ops, err := ScanDialect(includeRoot, dialect)
		if err != nil {
			return nil, err
		}
		entry := DialectSkipReport{
			Name:       dialect.Name,
			ImportPath: dialect.ImportPath,
			Package:    dialect.Package,
			Sources:    append([]string(nil), dialect.TableGen...),
			SkippedOps: make([]SkippedOperation, 0, len(ops)),
		}
		for _, op := range ops {
			if SupportedOp(dialect.Name, op.DefName) {
				continue
			}
			entry.SkippedOps = append(entry.SkippedOps, SkippedOperation{
				DefName:    op.DefName,
				Mnemonic:   op.Mnemonic,
				Reason:     "wrapper emission not implemented yet",
				Limitation: "generator",
				Source:     op.Source,
			})
		}
		report.Dialects = append(report.Dialects, entry)
	}
	return report, nil
}

func (r *SkipReport) JSON() ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}
