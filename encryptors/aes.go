package encryptors

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

const EncryptorName = "aes"

func NewAESEncryptor() Encryptor {
	return &aesEncryptor {}
}

func (ae *aesEncryptor) Name() string {
	return EncryptorName
}

// The key should be 48 bytes.
// First 32 bytes in key are actual key, the last 16 bytes are initialization vector
func (ae *aesEncryptor) Encrypt(key []byte, data []byte) ([]byte, error) {
	if len(key) != 48 {
		return nil, errors.New("key length should be 48 bytes (32 - key itself, then 16 - initialization vector)")
	}
	iv := key[32:]
	key = key[:32]
	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	// padding
	padding := block.BlockSize() - len(data) % block.BlockSize()
	padbytes := bytes.Repeat([]byte{byte(padding)}, padding)
	data = append(data, padbytes...)

	encrypted := make([]byte, len(data))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encrypted, data)


	return encrypted, nil
}

// The key should be 48 bytes.
// First 32 bytes in key are actual key, the last 16 bytes are initialization vector
func (ae *aesEncryptor) Decrypt(key []byte, data []byte) ([]byte, error) {
	if len(key) != 48 {
		return nil, errors.New("key length should be 48 bytes (32 - key itself, then 16 - initialization vector)")
	}
	iv := key[32:]
	key = key[:32]
	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	decrypted := make([]byte, len(data))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decrypted, data)
	// unpadding
	unpadding := int(decrypted[len(decrypted) - 1])
	decrypted = decrypted[:len(decrypted) - unpadding]

	return decrypted, nil
}