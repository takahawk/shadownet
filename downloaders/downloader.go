package downloaders

type Downloader interface {
	Download(id string) (string, error)
}
