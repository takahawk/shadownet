package transformers

type Transformer interface {
	ForwardTransform(data []byte) ([]byte, error)
	ReverseTransform(data []byte) ([]byte, error)
}