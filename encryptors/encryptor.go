package encryptors

type Encryptor interface {
	// TODO: change type to more general byte arrays
	Encrypt(data string) (string, error)
	Decrypt(cipher string) (string, error)
}