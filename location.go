//go:build cgo

package mlir

import "github.com/timmyyuan/mlir-go/internal/capi"

// Location is a borrowed handle to an MLIR location.
type Location struct {
	raw capi.Location
}

// UnknownLocation returns an unknown location owned by the provided context.
func UnknownLocation(ctx *Context) (Location, error) {
	raw, err := ctx.requireRaw()
	if err != nil {
		return Location{}, err
	}
	return Location{raw: capi.LocationUnknownGet(raw)}, nil
}

// FileLineColLocation returns a file/line/column location owned by the context.
func FileLineColLocation(ctx *Context, filename string, line, col uint) (Location, error) {
	raw, err := ctx.requireRaw()
	if err != nil {
		return Location{}, err
	}
	return Location{raw: capi.LocationFileLineColGet(raw, filename, line, col)}, nil
}

// IsNull reports whether the location handle is null.
func (loc Location) IsNull() bool {
	return capi.LocationIsNull(loc.raw)
}

// String returns the textual assembly form of the location.
func (loc Location) String() string {
	if loc.IsNull() {
		return ""
	}
	return capi.LocationToString(loc.raw)
}
