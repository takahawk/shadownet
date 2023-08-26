package transformers

import (
	"encoding/base64"
)

type base64Transformer struct {}
const TransformerName = "base64"

func NewBase64Transformer() Transformer {
	return &base64Transformer{}
}

func (b64t *base64Transformer) Name() string {
	return TransformerName
}

func (b64t *base64Transformer) ForwardTransform(data []byte) ([]byte, error) {
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(encoded, data)
	return encoded, nil
}

func (b64t *base64Transformer) ReverseTransform(data []byte) ([]byte, error) {
	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(decoded, data)
	if err != nil {
		return nil, err
	}

	return decoded[:n], nil
}