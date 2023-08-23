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

func NewAESEncryptor(key []byte, iv []byte) (Encryptor, error) {
	if len(key) != 32 {
		return nil, errors.New("key length should be 32 bytes")
	}

	if len(iv) != 16 {
		return nil, errors.New("initialization vector should be 16 bytes")
	}
	
	return &aesEncryptor {
		key: key,
		iv: iv,
	}, nil
}

func (ae *aesEncryptor) Encrypt(data []byte) ([]byte, error) {
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

func (ae *aesEncryptor) Decrypt(data []byte) ([]byte, error) {
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