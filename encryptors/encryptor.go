package encryptors

import (
	"github.com/takahawk/shadownet/common"
)

type Encryptor interface {
	common.Nameable
	Encrypt(key []byte, data []byte) ([]byte, error)
	Decrypt(key []byte, cipher []byte) ([]byte, error)
}