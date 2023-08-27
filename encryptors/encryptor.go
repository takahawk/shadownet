package encryptors

import (
	"github.com/takahawk/shadownet/common"
)

// TODO: mb, should be removed because signatures are the same for Transformer interface?
type Encryptor interface {
	common.Nameable
	Encrypt(data []byte) ([]byte, error)
	Decrypt(cipher []byte) ([]byte, error)
}