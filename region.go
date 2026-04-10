//go:build cgo

package mlir

import (
	"fmt"

	"github.com/timmyyuan/mlir-go/internal/capi"
)

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

// AppendOwnedBlock transfers ownership of a detached block into r.
func (r Region) AppendOwnedBlock(block *OwnedBlock) (Block, error) {
	if r.IsNull() {
		return Block{}, fmt.Errorf("mlir: null region")
	}
	rawBlock, err := block.takeRaw()
	if err != nil {
		return Block{}, err
	}
	capi.RegionAppendOwnedBlock(r.raw, rawBlock)
	return Block{raw: rawBlock}, nil
}
