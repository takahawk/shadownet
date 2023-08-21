package storages

type Downloader interface {
	Download(id string) (string, error)
}

type Uploader interface {
	Upload(id string, content string) error
}