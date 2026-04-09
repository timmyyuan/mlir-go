//go:build !cgo

package mlir

type Value struct{}

func (v Value) IsNull() bool {
	return true
}

func (v Value) IsBlockArgument() bool {
	return false
}

func (v Value) IsOpResult() bool {
	return false
}

func (v Value) Type() Type {
	return Type{}
}

func (v Value) String() string {
	return ""
}
