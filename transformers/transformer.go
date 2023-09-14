package transformers

import (
	"github.com/takahawk/shadownet/common"
)

// Transformer is ShadowNet components that takes one sequence of bytes and
// make some operations to get the modified one. Sure it is also should provide
// a way for reverse transformation
type Transformer interface {
	common.Component
	// ForwardTransform given the original version of data returns the
	// transformed one
	ForwardTransform(data []byte) ([]byte, error)
	// ReverseTransform given the transformed version of data returns the
	// original one
	ReverseTransform(data []byte) ([]byte, error)
}
