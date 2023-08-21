package storages

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

const PastebinRawPrefix = "https://pastebin.com/raw"
const PastebinPostUrl = "https://pastebin.com/api/api_post.php"

type pastebinDownloader struct {}
type pastebinUploader struct {
	// TODO: add option to use authorized API instance
}

func NewPastebinDownloader() Downloader {
	return &pastebinDownloader{}
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

func NewPastebinUploader() Uploader {
	return &pastebinUploader{}
}

func (pd *pastebinUploader) Upload(id string, content string) error {
	// TODO: implement
	return errors.New("Unimplemented")
}