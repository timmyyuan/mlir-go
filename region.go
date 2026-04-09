//go:build cgo

package mlir

import "github.com/timmyyuan/mlir-go/internal/capi"

// Region is a borrowed handle to an MLIR region.
type Region struct {
	raw capi.Region
}

// IsNull reports whether the region handle is null.
func (r Region) IsNull() bool {
	return capi.RegionIsNull(r.raw)
}

// Blocks returns the blocks contained in the region in source order.
func (r Region) Blocks() []Block {
	if r.IsNull() {
		return nil
	}

	var blocks []Block
	for raw := capi.RegionGetFirstBlock(r.raw); !capi.BlockIsNull(raw); raw = capi.BlockGetNextInRegion(raw) {
		blocks = append(blocks, Block{raw: raw})
	}
	return blocks
}
