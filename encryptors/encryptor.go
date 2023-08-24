package encryptors

type Encryptor interface {
	// TODO: change type to more general byte arrays
	Encrypt(key []byte, data []byte) ([]byte, error)
	Decrypt(key []byte, cipher []byte) ([]byte, error)
}