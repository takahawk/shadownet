package encryptors

import (
	"crypto/aes"

	"encoding/hex"
)

type aesEncryptor struct {
	key string
}

func NewAESEncryptor(key string) Encryptor {
	// FIXME: check for key length (32 byte)
	// TODO: mb make key part of interface instead of storing it as fields?
	return &aesEncryptor {
		key: key,
	}
}

func (ae *aesEncryptor) Encrypt(data string) (string, error) {
	// FIXME: it encrypts only one block
	c, err := aes.NewCipher([]byte(ae.key))

	if err != nil {
		return "", err
	}

	out := make([]byte, len(data))
	c.Encrypt(out, []byte(data))

	// TODO: make hex (or base64) transformation a separate step.
	return hex.EncodeToString(out), nil
}

func (ae *aesEncryptor) Decrypt(data string) (string, error) {
	ciphertext, err := hex.DecodeString(data)

	if err != nil {
		return "", err
	}

	c, err := aes.NewCipher([]byte(ae.key))
	if err != nil {
		return "", err
	}

	out := make([]byte, len(ciphertext))
	c.Decrypt(out, ciphertext)

	return string(out[:]), nil
}