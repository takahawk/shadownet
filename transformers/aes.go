package transformers

import (
	"bytes"
	"errors"

	"crypto/aes"
	"crypto/cipher"
)

type aesEncryptor struct {
	key []byte
	iv []byte
}

const AESEncryptorName = "aes"

func NewAESEncryptor(key []byte, iv []byte) (Transformer, error) {
	if len(key) != 32 {
		return nil, errors.New("key length should be 32 bytes")
	}
	if len(iv) != 16 {
		return nil, errors.New("key length should be 16 bytes")
	}
	return &aesEncryptor {
		key: key,
		iv: iv,
	}, nil
}

func NewAESEncryptorWithParams(params... []byte) (Transformer, error) {
	if len(params) != 2 {
		return nil, errors.New("there should be 2 parameters: key and iv")
	}
	key := params[0]
	iv := params[1]

	return NewAESEncryptor(key, iv)
}

func (ae *aesEncryptor) Name() string {
	return AESEncryptorName
}

func (ae *aesEncryptor) Params() [][]byte {
	return [][]byte{ ae.key, ae.iv }
}

func (ae *aesEncryptor) ForwardTransform(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(ae.key)

	if err != nil {
		return nil, err
	}

	// padding
	padding := block.BlockSize() - len(data) % block.BlockSize()
	padbytes := bytes.Repeat([]byte{byte(padding)}, padding)
	data = append(data, padbytes...)

	encrypted := make([]byte, len(data))
	mode := cipher.NewCBCEncrypter(block, ae.iv)
	mode.CryptBlocks(encrypted, data)


	return encrypted, nil
}

// The key should be 48 bytes.
// First 32 bytes in key are actual key, the last 16 bytes are initialization vector
func (ae *aesEncryptor) ReverseTransform(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(ae.key)

	if err != nil {
		return nil, err
	}

	decrypted := make([]byte, len(data))
	mode := cipher.NewCBCDecrypter(block, ae.iv)
	mode.CryptBlocks(decrypted, data)
	// unpadding
	unpadding := int(decrypted[len(decrypted) - 1])
	decrypted = decrypted[:len(decrypted) - unpadding]

	return decrypted, nil
}