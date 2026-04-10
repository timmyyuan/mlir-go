//go:build cgo

package mlir

import (
	"fmt"
	"sync"

	"github.com/timmyyuan/mlir-go/internal/capi"
)

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

func (d Diagnostic) String() string {
	if d.Location == "" {
		return fmt.Sprintf("%s: %s", d.Severity, d.Message)
	}
	return fmt.Sprintf("%s at %s: %s", d.Severity, d.Location, d.Message)
}

func (c *Context) CaptureDiagnostics(fn func() error) ([]Diagnostic, error) {
	if fn == nil {
		return nil, fmt.Errorf("mlir: diagnostic callback target is required")
	}
	raw, err := c.requireRaw()
	if err != nil {
		return nil, err
	}

	var (
		mu    sync.Mutex
		diags []Diagnostic
	)

	id := capi.ContextAttachDiagnosticCallback(raw, func(diag capi.Diagnostic) {
		mu.Lock()
		defer mu.Unlock()
		diags = append(diags, convertDiagnostic(diag))
	})

	runErr := fn()

	if _, err := c.requireRaw(); err == nil {
		capi.ContextDetachDiagnosticHandler(raw, id)
	}

	mu.Lock()
	defer mu.Unlock()
	return append([]Diagnostic(nil), diags...), runErr
}

func convertDiagnostic(diag capi.Diagnostic) Diagnostic {
	out := Diagnostic{
		Severity: convertDiagnosticSeverity(capi.DiagnosticGetSeverity(diag)),
		Message:  capi.DiagnosticToString(diag),
		Location: capi.LocationToString(capi.DiagnosticGetLocation(diag)),
	}
	n := capi.DiagnosticGetNumNotes(diag)
	if n == 0 {
		return out
	}
	out.Notes = make([]Diagnostic, 0, n)
	for i := 0; i < n; i++ {
		out.Notes = append(out.Notes, convertDiagnostic(capi.DiagnosticGetNote(diag, i)))
	}
	return out
}

func convertDiagnosticSeverity(severity capi.DiagnosticSeverity) DiagnosticSeverity {
	switch severity {
	case capi.DiagnosticSeverityWarning:
		return DiagnosticWarning
	case capi.DiagnosticSeverityNote:
		return DiagnosticNote
	case capi.DiagnosticSeverityRemark:
		return DiagnosticRemark
	default:
		return DiagnosticError
	}
}
