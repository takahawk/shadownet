package uploaders

type Uploader interface {
	// TODO: mb add some generic `params` to gain more control on specific upload?
	Upload(content []byte) (id string, err error)
}