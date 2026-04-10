//go:build !cgo

package mlir

import "fmt"

type DiagnosticSeverity int

const (
	DiagnosticError DiagnosticSeverity = iota
	DiagnosticWarning
	DiagnosticNote
	DiagnosticRemark
)

func (s DiagnosticSeverity) String() string {
	switch s {
	case DiagnosticError:
		return "error"
	case DiagnosticWarning:
		return "warning"
	case DiagnosticNote:
		return "note"
	case DiagnosticRemark:
		return "remark"
	default:
		return "unknown"
	}
}

type Diagnostic struct {
	Severity DiagnosticSeverity
	Location string
	Message  string
	Notes    []Diagnostic
}

func (d Diagnostic) String() string { return d.Message }

func (c *Context) CaptureDiagnostics(func() error) ([]Diagnostic, error) {
	return nil, fmt.Errorf("mlir: cgo is required")
}
