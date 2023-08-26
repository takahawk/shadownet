package transformers

import (
	"github.com/takahawk/shadownet/common"
)

type Transformer interface {
	common.Nameable
	ForwardTransform(data []byte) ([]byte, error)
	ReverseTransform(data []byte) ([]byte, error)
}