package transformers

import (
	"bytes"
	"errors"

	"crypto/aes"
	"crypto/cipher"
)

type aesEncryptor struct {
	key []byte
	iv  []byte
}

// AESEncryptorName is component name for AES encryptor
const AESEncryptorName = "aes"

// NewAESEncryptor creates new transformer that allows to encrypt and decrypt
// data using AES algorithm.
// Key and initialization vector must be provided.
// Key should be 32 bytes long.
// And initialization vector should be 16 bytes long.
func NewAESEncryptor(key []byte, iv []byte) (Transformer, error) {
	if len(key) != 32 {
		return nil, errors.New("key length should be 32 bytes")
	}
	if len(iv) != 16 {
		return nil, errors.New("key length should be 16 bytes")
	}
	return &aesEncryptor{
		key: key,
		iv:  iv,
	}, nil
}

// NewAESEncryptorWithParams is convenience function that call NewAESEncryptor
// with parameters (that is, first is the key and second is initialization
// vector) packed into slice.
func NewAESEncryptorWithParams(params ...[]byte) (Transformer, error) {
	if len(params) != 2 {
		return nil, errors.New("there should be 2 parameters: key and iv")
	}
	key := params[0]
	iv := params[1]

	return NewAESEncryptor(key, iv)
}

// Name returns component name of AES encryptor. It is always AESEncryptorName
func (ae *aesEncryptor) Name() string {
	return AESEncryptorName
}

// Params returns key and initialization vector packed into slice
func (ae *aesEncryptor) Params() [][]byte {
	return [][]byte{ae.key, ae.iv}
}

// ForwardTransform returns byte sequence encoded with AES algorithm using
// key and initialization vector acquired during AES transformer creation
func (ae *aesEncryptor) ForwardTransform(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(ae.key)

	if err != nil {
		return nil, err
	}

	// padding
	padding := block.BlockSize() - len(data)%block.BlockSize()
	padbytes := bytes.Repeat([]byte{byte(padding)}, padding)
	data = append(data, padbytes...)

	encrypted := make([]byte, len(data))
	mode := cipher.NewCBCEncrypter(block, ae.iv)
	mode.CryptBlocks(encrypted, data)

	return encrypted, nil
}

// ReverseTransform returns byte sequence decoded with AES algorithm using
// key and initialization vector acquired during AES transformer creation
func (ae *aesEncryptor) ReverseTransform(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(ae.key)

	if err != nil {
		return nil, err
	}

	decrypted := make([]byte, len(data))
	mode := cipher.NewCBCDecrypter(block, ae.iv)
	mode.CryptBlocks(decrypted, data)
	// unpadding
	unpadding := int(decrypted[len(decrypted)-1])
	decrypted = decrypted[:len(decrypted)-unpadding]

	return decrypted, nil
}
