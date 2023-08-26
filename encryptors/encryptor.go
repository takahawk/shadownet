package encryptors

import (
	"github.com/takahawk/shadownet/common"
)

type Encryptor interface {
	common.Nameable
	Encrypt(data []byte) ([]byte, error)
	Decrypt(cipher []byte) ([]byte, error)
}