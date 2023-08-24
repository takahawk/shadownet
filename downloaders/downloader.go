package downloaders

type Downloader interface {
	Download(id string) (string, error)
}

// TODO: move to separate package
type Uploader interface {
	Upload(id string, content string) error
}