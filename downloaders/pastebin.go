package downloaders

import (
	"fmt"
	"io"
	"net/http"
)

const PastebinRawPrefix = "https://pastebin.com/raw"
const PastebinPostUrl = "https://pastebin.com/api/api_post.php"

type pastebinDownloader struct {}
const PastebinDownloaderName = "pastebin"

func NewPastebinDownloader() Downloader {
	return &pastebinDownloader{}
}

func (pd *pastebinDownloader) Name() string {
	return PastebinDownloaderName
}

func (pd *pastebinDownloader) Download(id string) (string, error) {
	res, err := http.Get(fmt.Sprintf("%s/%s", PastebinRawPrefix, id))
	if err != nil {
		// TODO: error handling (wrap etc.)?
		return "", err
	}

	content, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", err
	}

	return string(content), nil
}
