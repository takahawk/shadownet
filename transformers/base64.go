package transformers

import (
	"encoding/base64"
	"errors"

	"github.com/takahawk/shadownet/logger"
)

// Base64TransformerName is component name of base64 transformer
const Base64TransformerName = "base64"

type base64Transformer struct {
	logger logger.Logger
}

// NewBase64Transformer returns transformer that provides encoding/decoding
// data to/from Base64
func NewBase64Transformer(logger logger.Logger) Transformer {
	return &base64Transformer{
		logger: logger,
	}
}

// NewBase64Transformer is convinience function to call base64 constructor
// with signature that contains params
func NewBase64TransformerWithParams(logger logger.Logger, params ...[]byte) (Transformer, error) {
	if len(params) != 0 {
		return nil, errors.New("base64 transformer doesn't accept any params")
	}

	return NewBase64Transformer(logger), nil
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
