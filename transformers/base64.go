package transformers

import (
	"encoding/base64"
)

type base64Transformer struct{}

// Base64TransformerName is component name of base64 transformer
const Base64TransformerName = "base64"

// NewBase64Transformer returns transformer that provides encoding/decoding
// data to/from Base64
func NewBase64Transformer() Transformer {
	return &base64Transformer{}
}

// Name returns component name of base64 transformer.
// It is always Base64TransformerName
func (b64t *base64Transformer) Name() string {
	return Base64TransformerName
}

// Params returns empty slice just to be with accordance with general interface
func (b64t *base64Transformer) Params() [][]byte {
	return nil
}

// ForwardTransform encodes byte array using base64 encoding
func (b64t *base64Transformer) ForwardTransform(data []byte) ([]byte, error) {
	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(encoded, data)
	return encoded, nil
}

// ReverseTransform decodes base64'd byte array into original byte sequence
func (b64t *base64Transformer) ReverseTransform(data []byte) ([]byte, error) {
	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(decoded, data)
	if err != nil {
		return nil, err
	}

	return decoded[:n], nil
}
