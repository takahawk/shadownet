package encryptors

type Encryptor interface {
	// TODO: change type to more general byte arrays
	Encrypt(data []byte) ([]byte, error)
	Decrypt(cipher []byte) ([]byte, error)
}